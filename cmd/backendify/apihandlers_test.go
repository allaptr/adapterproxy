package main

import (
	"backendify/cache"
	"backendify/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cactus/go-statsd-client/v5/statsd"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_apiHandler(t *testing.T) {
	tds := []struct {
		name          string
		requestMethod string
		requestString string
		responseCode  int
		responseBody  []byte
	}{
		{"Unsupported http method", http.MethodPut, "/company?id=abc&country_iso=us", http.StatusMethodNotAllowed, []byte{}},
		{"Company id missing", http.MethodGet, "/company?country_iso=us", http.StatusBadRequest, paramError},
		{"Country_iso is missing", http.MethodGet, "/company?id=123", http.StatusBadRequest, paramError},
		{"Country not found", http.MethodGet, "/company?id=abc&country_iso=cc", http.StatusNotFound, countryError},
	}

	for _, td := range tds {
		t.Run(td.name, func(t *testing.T) {
			config := &statsd.ClientConfig{
				Address: "127.0.0.1:8125",
			}
			mclient, err := statsd.NewClientWithConfig(config)
			assert.NoError(t, err)

			req := httptest.NewRequest(td.requestMethod, td.requestString, nil)
			w := httptest.NewRecorder()
			m := map[string]string{"us": "http://localhost:9002"}
			c := cache.NewCache()

			ah := newApiHandler(m, c, mclient, 200*time.Millisecond)
			ah.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, td.responseBody, data)
			assert.Equal(t, td.responseCode, res.StatusCode)
		})
	}
}

func Test_serveCached(t *testing.T) {
	cacheVal := []byte("{\"cn\":\"elastic-aryabhata-kare\",\"created_on\":\"1955-04-10T01:54:39Z\",\"closed_on\":null}")
	tds := []struct {
		name         string
		country      string
		companyID    string
		success      bool
		responseBody []byte
	}{
		{"Cache miss", "us", "123", false, []byte{}},
		{"Cache hit", "us", "440804282323A", true, cacheVal},
	}
	c := cache.NewCache()
	c.Put(cache.MakeKey("us", "440804282323A"), cacheVal)

	for _, td := range tds {
		tf := func(t *testing.T) {
			t.Parallel()
			config := &statsd.ClientConfig{
				Address: "127.0.0.1:8125",
			}
			mclient, err := statsd.NewClientWithConfig(config)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			m := map[string]string{"us": "http://localhost:9002"}

			ah := newApiHandler(m, c, mclient, 200*time.Millisecond)
			success := ah.serveCached(w, td.country, td.companyID)
			assert.Equal(t, td.success, success)
			res := w.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.Equal(t, td.responseBody, data)
		}
		t.Run(td.name, tf)
	}
}

func Test_callBackend(t *testing.T) {
	expected := "{\"cn\":\"mystifying-proskuriakova-mcclintock\",\"created_on\":\"1759-05-27T01:54:39Z\",\"closed_on\":null}"
	tds := []struct {
		name         string
		delay        time.Duration
		err          error
		responseBody []byte
	}{
		{"Happy path", 25 * time.Millisecond, nil, []byte(expected)},
		{"Timeout", 250 * time.Millisecond, context.DeadlineExceeded, nil},
	}
	be := backender{
		timeout: 200 * time.Millisecond,
	}
	for _, td := range tds {
		t.Run(td.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(td.delay)
				w.WriteHeader(http.StatusOK)
				ret, err := w.Write([]byte(expected))
				if err != nil || ret == 0 {
					log.Errorf("Writing error %v", err)
					return
				}
			}))
			assert.NotNil(t, svr)
			defer svr.Close()
			fmt.Println(svr.URL)
			resp, err := be.callBackend(svr.URL, "634374392121W")
			if td.responseBody == nil {
				assert.Nil(t, resp)
				assert.True(t, errors.Is(err, td.err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				var body []byte
				body, err = io.ReadAll(resp.Body)
				assert.NoError(t, err)
				defer resp.Body.Close()
				assert.Equal(t, expected, string(body))
			}
		})
	}
}

func Test_processV1Response(t *testing.T) {

	td := []struct {
		name           string
		data           model.CompanyLegacyV1
		expName        string
		expActive      bool
		expActiveUntil string
	}{
		{"Active", model.CompanyLegacyV1{Cn: "foo", Created_on: "2005-04-12T23:20:50.52Z", Closed_on: ""}, "foo", true, ""},
		{"Active with active_until", model.CompanyLegacyV1{Cn: "foo", Created_on: "2005-04-12T23:20:50.52Z", Closed_on: "3005-04-12T23:20:50.52Z"}, "foo", true, "3005-04-12T23:20:50.52Z"},
		{"Inactive", model.CompanyLegacyV1{Cn: "foo", Created_on: "2005-04-12T23:20:50.52Z", Closed_on: "2015-04-01"}, "foo", false, "2015-04-01"},
		{"Active creation in the future", model.CompanyLegacyV1{Cn: "foo", Created_on: "2075-04-01T00:13:59Z", Closed_on: "2065-04-01T00:13:12Z"}, "foo", true, "2065-04-01T00:13:12Z"},
	}
	for _, td := range td {
		t.Run(td.name, func(t *testing.T) {
			b, err := json.Marshal(td.data)
			assert.NoError(t, err)
			name, active, active_until := processV1Response(b)
			assert.Equal(t, td.expName, name)
			assert.Equal(t, td.expActive, active)
			assert.Equal(t, td.expActiveUntil, active_until)
		})
	}
}

func Test_processV1RealData(t *testing.T) {

	td := []struct {
		name           string
		data           string
		expName        string
		expActive      bool
		expActiveUntil string
	}{
		{"Active 1", `{"cn":"osom-robinson-chaum","created_on":"1909-04-21T00:13:59Z","closed_on":null}`, "osom-robinson-chaum", true, ""},
		{"Active 2", `{"cn":"reverent-haslett-chebyshev","created_on":"1855-05-04T00:13:59Z","closed_on":null}`, "reverent-haslett-chebyshev", true, ""},
		{"Active 3", `{"cn":"infallible-mcnulty-mestorf","created_on":"1856-05-03T00:13:59Z","closed_on":null}`, "infallible-mcnulty-mestorf", true, ""},
		{"Active 4", `{"cn":"keen-lederberg-noether","created_on":"1983-04-03T00:13:59Z","closed_on":null}`, "keen-lederberg-noether", true, ""},
	}
	for _, td := range td {
		t.Run(td.name, func(t *testing.T) {
			name, active, active_until := processV1Response([]byte(td.data))
			assert.Equal(t, td.expName, name)
			assert.Equal(t, td.expActive, active)
			assert.Equal(t, td.expActiveUntil, active_until)
		})
	}
}

func Test_requestString(t *testing.T) {
	str := requestString("http://localhost:9002", "170456685224Z")
	assert.Equal(t, "http://localhost:9002/companies/170456685224Z", str)
}

func TestGo(t *testing.T) {
	var π = 22 / 7.0
	fmt.Println(π)

	var m map[string]int
	fmt.Println(m["Hello"])

	s := "a\tb"
	fmt.Println(s)
}
