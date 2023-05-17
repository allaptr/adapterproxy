package main

import (
	"backendify/cache"
	"backendify/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cactus/go-statsd-client/v5/statsd"
	log "github.com/sirupsen/logrus"
)

const (
	v1 = "application/x-company-v1"
	v2 = "application/x-company-v2"
)

var (
	countryError  = []byte("{\"error\":\"country not found\"}")
	companyError  = []byte("{\"error\":\"company not found\"}")
	paramError    = []byte("{\"error\":\"missing request parameters\"}")
	backendError  = []byte("{\"error\":\"calling backend error\"}")
	timeoutError  = []byte("{\"error\":\"time out waiting for response from backend\"}")
	catchallError = []byte("{\"error\":\"temporary unknown server error\"}")
)

type backendCaller interface {
	callBackend(route, id string) (*http.Response, error)
}

type backender struct {
	timeout time.Duration
}

type apiHandler struct {
	// a list of backend endpoints is passed in on startup
	// ex: ru=http://localhost:9001 us=http://localhost:9002 and cached in the map
	routs map[string]string
	// Caches responses by the key created from country_iso and company id
	// Ex.: us-605601630650G
	cache cache.Cache
	b     backendCaller
	// StatsD metrics client
	metricClient statsd.Statter
}

func newApiHandler(m map[string]string, r cache.Cache, statsdClient statsd.Statter, timeout time.Duration) *apiHandler {
	return &apiHandler{
		routs: m,
		cache: r,
		b: &backender{
			timeout: timeout,
		},
		metricClient: statsdClient,
	}
}

func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	begin := time.Now().UTC().UnixMilli()
	a.collectMetric(2)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	country := r.URL.Query().Get("country_iso")
	if id == "" || country == "" {
		a.writeResponse(w, http.StatusBadRequest, paramError)
		return
	}

	// Use cached data
	if a.serveCached(w, country, id) {
		return
	}
	// Get counry endpoint
	var route string
	var ok bool
	if route, ok = a.routs[country]; !ok {
		a.writeResponse(w, http.StatusNotFound, countryError)
		return
	}
	callBegin := time.Now().UTC().UnixMilli()
	resp, err := a.b.callBackend(route, id)
	if err != nil && resp == nil {
		if errors.Is(err, context.DeadlineExceeded) {
			a.collectMetric(4)
			a.writeResponse(w, http.StatusGatewayTimeout, timeoutError)
			return
		}
		a.writeResponse(w, http.StatusGatewayTimeout, catchallError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		a.collectMetric(5)
		log.WithFields(log.Fields{"country_iso": country, "id": id}).Debugf("Status code %v", resp.StatusCode)
		if resp.StatusCode == http.StatusNotFound {
			a.writeResponse(w, http.StatusNotFound, companyError)
			return
		}
		a.writeResponse(w, http.StatusGatewayTimeout, backendError)
		return
	}
	callEnd := time.Now().UTC().UnixMilli()
	log.Debugf("Backend Call: %d", callEnd-callBegin)
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{"country_iso": country, "id": id}).Debugf("Error reading response body %v", err)
		a.writeResponse(w, http.StatusInternalServerError, catchallError)
		return
	}

	if len(body) == 0 {
		log.WithFields(log.Fields{"country_iso": country, "id": id}).Debug("Company not found")
		a.writeResponse(w, http.StatusNotFound, companyError)
		return
	}
	ct := resp.Header.Get("Content-Type")
	compInfo := a.processBackendResponse(id, ct, body)
	if compInfo == nil {
		a.writeResponse(w, http.StatusNotFound, companyError)
		return
	}

	// Create response body
	b, err := json.Marshal(compInfo)
	if err != nil {
		log.WithFields(log.Fields{"country_iso": country, "id": id}).Debugf("Error marshalling response %v", err)
		a.writeResponse(w, http.StatusInternalServerError, catchallError)
		return
	}
	a.writeResponse(w, http.StatusOK, b)
	a.cacheResponse(country, id, b)
	end := time.Now().UTC().UnixMilli()
	log.Debugf("Total Request Processing Time: %d", end-begin)
}

func (a *apiHandler) writeResponse(w http.ResponseWriter, statusCode int, msgJson []byte) {
	w.WriteHeader(statusCode)
	ret, err := w.Write(msgJson)
	if err != nil || ret == 0 {
		log.Errorf("response write error %v", err)
	}
}

func (b *backender) callBackend(route, id string) (*http.Response, error) {
	// Call the backend with timeout
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()
	reqStr := requestString(route, id)
	log.Debugf("Sending GET request %s ", reqStr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqStr, nil)
	if err != nil {
		log.Debugf("Error calling backend with ctx %v ", err)
		return nil, err
	}
	client := &http.Client{}
	return client.Do(req)
}

func (a *apiHandler) collectMetric(value int64) {
	var b strings.Builder
	fmt.Fprintf(&b, "metric.%d", value)
	mname := b.String()
	err := a.metricClient.Inc(mname, value, 1.0)
	if err != nil {
		log.Debugf("Metric error for %s", mname)
	}
}

func (a *apiHandler) cacheResponse(country, id string, data []byte) {
	key := cache.MakeKey(country, id)
	a.cache.Put(key, data)
}

func (a *apiHandler) serveCached(w http.ResponseWriter, country, id string) bool {
	key := cache.MakeKey(country, id)
	cachedResponse := a.cache.Get(key)
	if cachedResponse != nil {
		log.WithFields(log.Fields{"country_iso": country, "id": id}).Debug("Cached response")
		w.WriteHeader(http.StatusOK)
		cnt, err := w.Write(cachedResponse)
		if err != nil || cnt == 0 {
			log.Debugf("Error writing response %v", err)
			return false
		}
		a.collectMetric(3)
		return true
	}
	return false
}

func (a *apiHandler) processBackendResponse(id, ct string, body []byte) *model.CompanyInfo {
	//read backend response
	var companyName string
	var active bool
	var active_until string
	switch ct {
	case v1:
		companyName, active, active_until = processV1Response(body)
	case v2:
		companyName, active, active_until = processV2Response(body)
	default:
		return nil
	}

	compInfo := model.CompanyInfo{
		Id:           id,
		Name:         companyName,
		Active:       active,
		Active_until: active_until,
	}
	return &compInfo
}
func processV2Response(resp []byte) (string, bool, string) {
	var d model.CompanyLegacyV2
	if err := json.Unmarshal(resp, &d); err != nil {
		return "", false, ""
	}
	if d.Tin == "" {
		return d.Company_name, false, ""
	}
	if d.Dissolved_on == "" {
		return d.Company_name, true, ""
	}
	closedTime, _ := time.Parse(time.RFC3339, d.Dissolved_on)
	pastClosed := time.Now().UTC().After(closedTime)
	return d.Company_name, !pastClosed, d.Dissolved_on
}

// Returns comp name, isActive, active_until (RFC3339)
func processV1Response(body []byte) (string, bool, string) {
	var d model.CompanyLegacyV1
	err := json.NewDecoder(bytes.NewReader(body)).Decode(&d)
	if err != nil {
		return "", false, ""
	}
	if d.Created_on == "" {
		return d.Cn, false, ""
	}
	if d.Closed_on == "" {
		return d.Cn, true, ""
	}
	closedTime, _ := time.Parse(time.RFC3339, d.Closed_on)
	pastClosed := time.Now().UTC().After(closedTime)
	return d.Cn, !pastClosed, d.Closed_on
}
func requestString(address, id string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/companies/%s", address, id)
	return b.String()
}
