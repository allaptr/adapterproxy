package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var txtData = map[string]string{
	"621451430762X": "{\"cn\":\"ACME Incorporated\",\"created_on\":\"2314-12-23T01:29:13Z\",\"closed_on\":null}",
	"282676033074S": "{\"cn\":\"laughing-lamarr-dirac\",\"created_on\":\"2022-03-24T01:54:39Z\",\"closed_on\":null}",
	"440804282323A": "{\"cn\":\"elastic-aryabhata-kare\",\"created_on\":\"1955-04-10T01:54:39Z\",\"closed_on\":null}",
	"testv1": "{\"cn\":\"friendly-euler-ptolemy\",\"created_on\":\"1977-04-04T01:54:39Z\",\"closed_on\":null}",
	"4164142148556": "{\"cn\":\"flamboyant-burnell-colden\",\"created_on\":\"1745-05-30T01:54:39Z\",\"closed_on\":null}",
	"155674243726O": "{\"cn\":\"stoic-edison-feynman\",\"created_on\":\"1786-05-20T01:54:39Z\",\"closed_on\":null}",
	"462059571630Q": "{\"cn\":\"agitated-bell-easley\",\"created_on\":\"2006-03-28T01:54:39Z\",\"closed_on\":null}",
	"854145342086D": "{\"cn\":\"unruffled-moore-hodgkin\",\"created_on\":\"1777-05-22T01:54:39Z\",\"closed_on\":null}",
	"464415421689N": "{\"cn\":\"tender-swartz-murdock\",\"created_on\":\"1771-05-24T01:54:39Z\",\"closed_on\":null}",
	"634374392121W": "{\"cn\":\"serene-chaum-bhabha\",\"created_on\":\"1869-04-30T01:54:39Z\",\"closed_on\":null}",
	"877991708633O": "{\"cn\":\"pensive-shockley-clarke\",\"created_on\":\"1918-04-19T01:54:39Z\",\"closed_on\":null}",
	"6187217740522": "{\"cn\":\"elated-banach-goldstine\",\"created_on\":\"1894-04-24T01:54:39Z\",\"closed_on\":null}",
	"918600199176T": "{\"cn\":\"reverent-chebyshev-engelbart\",\"created_on\":\"1919-04-19T01:54:39Z\",\"closed_on\":null}",
	"369401320059G": "{\"cn\":\"festive-cray-bose\",\"created_on\":\"1753-05-28T01:54:39Z\",\"closed_on\":null}",
	"525570322146U": "{\"cn\":\"mystifying-proskuriakova-mcclintock\",\"created_on\":\"1759-05-27T01:54:39Z\",\"closed_on\":null}",
	"372581767233L": "{\"cn\":\"dazzling-newton-cannon\",\"created_on\":\"1925-04-17T01:54:39Z\",\"closed_on\":null}",
	"478803106995Z": "{\"cn\":\"pensive-bose-bhaskara\",\"created_on\":\"1743-05-31T01:54:39Z\",\"closed_on\":null}",
	"655358414930Y": "{\"cn\":\"tender-gauss-shirley\",\"created_on\":\"1977-04-04T01:54:39Z\",\"closed_on\":null}",
	"156430540509Y": "{\"cn\":\"confident-mcclintock-haibt\",\"created_on\":\"1943-04-13T01:54:39Z\",\"closed_on\":null}",
	"3768800898314": "{\"cn\":\"quirky-rhodes-carson\",\"created_on\":\"1824-05-11T01:54:39Z\",\"closed_on\":null}",
	"3508905365551": "{\"cn\":\"epic-bose-brown\",\"created_on\":\"1763-05-26T01:54:39Z\",\"closed_on\":null}",
	"599191407388U": "{\"cn\":\"agitated-cohen-lehmann\",\"created_on\":\"1869-04-30T01:54:39Z\",\"closed_on\":null}",
	"1506270184767": "{\"cn\":\"gracious-tu-jepsen\",\"created_on\":\"1864-05-01T01:54:39Z\",\"closed_on\":null}",
	"548431946394U": "{\"cn\":\"eager-jones-brahmagupta\",\"created_on\":\"1869-04-30T01:54:39Z\",\"closed_on\":null}",
	"538767862374S": "{\"cn\":\"baby-hertz-burnell\",\"created_on\":\"1826-05-11T01:54:39Z\",\"closed_on\":null}",
	"920598541601S": "{\"cn\":\"pensive-austin-taussig\",\"created_on\":\"1890-04-25T01:54:39Z\",\"closed_on\":null}",
	"311130734087N": "{\"cn\":\"heuristic-nash-darwin\",\"created_on\":\"1898-04-23T01:54:39Z\",\"closed_on\":null}",
	"378339526656G": "{\"cn\":\"romantic-wu-blackwell\",\"created_on\":\"1806-05-16T01:54:39Z\",\"closed_on\":null}",
	"552624153194S": "{\"cn\":\"intelligent-bassi-jones\",\"created_on\":\"1795-05-18T01:54:39Z\",\"closed_on\":null}",
	"801922059877A": "{\"cn\":\"hardcore-satoshi-cohen\",\"created_on\":\"1996-03-30T01:54:39Z\",\"closed_on\":null}",
	"272745309677H": "{\"cn\":\"beautiful-haslett-bohr\",\"created_on\":\"1957-04-09T01:54:39Z\",\"closed_on\":null}",
	"208734339345L": "{\"cn\":\"sad-hawking-fermat\",\"created_on\":\"1852-05-04T01:54:39Z\",\"closed_on\":null}",
	"138840155253J": "{\"cn\":\"friendly-babbage-pare\",\"created_on\":\"1982-04-03T01:54:39Z\",\"closed_on\":null}",
	"642597231447K": "{\"cn\":\"sleepy-mcclintock-goldwasser\",\"created_on\":\"1770-05-24T01:54:39Z\",\"closed_on\":null}",
	"961478751843S": "{\"cn\":\"heuristic-vaughan-lederberg\",\"created_on\":\"1919-04-19T01:54:39Z\",\"closed_on\":null}",
	"739077213451U": "{\"cn\":\"priceless-swanson-sinoussi\",\"created_on\":\"1870-04-30T01:54:39Z\",\"closed_on\":null}",
	"179376887597J": "{\"cn\":\"zealous-germain-bartik\",\"created_on\":\"1815-05-14T01:54:39Z\",\"closed_on\":null}",
	"6611274142596": "{\"cn\":\"condescending-shamir-chaplygin\",\"created_on\":\"1843-05-07T01:54:39Z\",\"closed_on\":null}",
	"567056235727O": "{\"cn\":\"amazing-proskuriakova-khorana\",\"created_on\":\"1855-05-04T01:54:39Z\",\"closed_on\":null}",
	"293856989785G": "{\"cn\":\"magical-almeida-napier\",\"created_on\":\"1831-05-10T01:54:39Z\",\"closed_on\":null}",
	"786115359286E": "{\"cn\":\"jovial-bartik-cannon\",\"created_on\":\"1939-04-14T01:54:39Z\",\"closed_on\":null}",
	"857535078007Y": "{\"cn\":\"gifted-mcnulty-chandrasekhar\",\"created_on\":\"1993-03-31T01:54:39Z\",\"closed_on\":null}",
	"868649496404U": "{\"cn\":\"angry-williamson-hugle\",\"created_on\":\"1838-05-08T01:54:39Z\",\"closed_on\":null}",
	"565058762528T": "{\"cn\":\"hopeful-chatelet-villani\",\"created_on\":\"1874-04-29T01:54:39Z\",\"closed_on\":null}",
	"521680219801U": "{\"cn\":\"quirky-banzai-einstein\",\"created_on\":\"1743-05-31T01:54:39Z\",\"closed_on\":null}",
	"411861638379Q": "{\"cn\":\"condescending-tharp-clarke\",\"created_on\":\"1820-05-12T01:54:39Z\",\"closed_on\":null}",
	"511784741338D": "{\"cn\":\"charming-jang-mccarthy\",\"created_on\":\"1747-05-30T01:54:39Z\",\"closed_on\":null}",
	"161611478308D": "{\"cn\":\"zen-lamport-curie\",\"created_on\":\"1854-05-04T01:54:39Z\",\"closed_on\":null}",
	"672995902335D": "{\"cn\":\"gallant-allen-northcutt\",\"created_on\":\"1800-05-17T01:54:39Z\",\"closed_on\":null}",
	"4374990612497": "{\"cn\":\"vigorous-dirac-mclean\",\"created_on\":\"1745-05-30T01:54:39Z\",\"closed_on\":null}",
	"181465231191A": "{\"cn\":\"recursing-khorana-spence\",\"created_on\":\"1867-05-01T01:54:39Z\",\"closed_on\":null}",
	"725880524592F": "{\"cn\":\"festive-spence-burnell\",\"created_on\":\"2010-03-27T01:54:39Z\",\"closed_on\":null}",
	"674184663833P": "{\"cn\":\"angry-curran-boyd\",\"created_on\":\"1988-04-01T01:54:39Z\",\"closed_on\":null}",
	"577372190779C": "{\"cn\":\"fervent-lalande-galileo\",\"created_on\":\"1734-06-02T01:54:39Z\",\"closed_on\":null}",
	"5796107294366": "{\"cn\":\"sad-yonath-germain\",\"created_on\":\"1929-04-16T01:54:39Z\",\"closed_on\":null}",
	"666898706315C": "{\"cn\":\"flamboyant-poincare-galileo\",\"created_on\":\"1816-05-13T01:54:39Z\",\"closed_on\":null}",
	"8312921969034": "{\"cn\":\"compassionate-wescoff-diffie\",\"created_on\":\"1967-04-07T01:54:39Z\",\"closed_on\":null}",
	"4080812759765": "{\"cn\":\"intelligent-knuth-greider\",\"created_on\":\"1764-05-25T01:54:39Z\",\"closed_on\":null}",
	"111562191524H": "{\"cn\":\"thirsty-yonath-jang\",\"created_on\":\"1780-05-21T01:54:39Z\",\"closed_on\":null}",
	"7521428441944": "{\"cn\":\"festive-kirch-blackwell\",\"created_on\":\"1783-05-21T01:54:39Z\",\"closed_on\":null}",
	"1418969123973": "{\"cn\":\"gifted-kirch-liskov\",\"created_on\":\"1888-04-25T01:54:39Z\",\"closed_on\":null}",
	"605601630650G": "{\"cn\":\"vigilant-driscoll-mccarthy\",\"created_on\":\"1824-05-11T01:54:39Z\",\"closed_on\":null}",
	"251584213975A": "{\"cn\":\"magical-lovelace-sinoussi\",\"created_on\":\"1740-05-31T01:54:39Z\",\"closed_on\":null}",
	"900502478423A": "{\"cn\":\"dreamy-cannon-bartik\",\"created_on\":\"1901-04-23T01:54:39Z\",\"closed_on\":null}",
	"589934434900W": "{\"cn\":\"zen-haibt-solomon\",\"created_on\":\"1860-05-02T01:54:39Z\",\"closed_on\":null}",
	"758852378501J": "{\"cn\":\"suspicious-leavitt-lovelace\",\"created_on\":\"1873-04-29T01:54:39Z\",\"closed_on\":null}",
	"281539125058U": "{\"cn\":\"zealous-mccarthy-cartwright\",\"created_on\":\"1920-04-18T01:54:39Z\",\"closed_on\":null}",
	"594163701584O": "{\"cn\":\"osom-merkle-margulis\",\"created_on\":\"1805-05-16T01:54:39Z\",\"closed_on\":null}",
	"6695590571003": "{\"cn\":\"friendly-jackson-hermann\",\"created_on\":\"1924-04-17T01:54:39Z\",\"closed_on\":null}",
	"170456685224Z": "{\"cn\":\"unruffled-nobel-dirac\",\"created_on\":\"1842-05-07T01:54:39Z\",\"closed_on\":null}",
	"4511729044489": "{\"cn\":\"tender-bartik-haslett\",\"created_on\":\"1832-05-09T01:54:39Z\",\"closed_on\":null}",
	"457463951243K": "{\"cn\":\"jolly-kapitsa-bell\",\"created_on\":\"1886-04-26T01:54:39Z\",\"closed_on\":null}",
	"329873533600A": "{\"cn\":\"bold-lovelace-jennings\",\"created_on\":\"1967-04-07T01:54:39Z\",\"closed_on\":null}",
	"1357807041058": "{\"cn\":\"friendly-lichterman-darwin\",\"created_on\":\"1837-05-08T01:54:39Z\",\"closed_on\":null}",
	"8101148555277": "{\"cn\":\"infallible-antonelli-cohen\",\"created_on\":\"1988-04-01T01:54:39Z\",\"closed_on\":null}",
	"1000094032731": "{\"cn\":\"condescending-lamport-jang\",\"created_on\":\"1796-05-17T01:54:39Z\",\"closed_on\":null}",
	"6993980661241": "{\"cn\":\"inspiring-mclean-beaver\",\"created_on\":\"1795-05-18T01:54:39Z\",\"closed_on\":null}",
	"601271265201E": "{\"cn\":\"stupefied-villani-cannon\",\"created_on\":\"1969-04-06T01:54:39Z\",\"closed_on\":null}",
	"213093596627C": "{\"cn\":\"ACME Canary LTD\",\"created_on\":\"1968-04-06T01:54:39Z\",\"closed_on\":null}",
	"234801237796A": "{\"cn\":\"ACME LRU\",\"created_on\":\"1972-04-05T01:54:39Z\",\"closed_on\":null}",
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}
// us
func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/companies/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// var params map[string]string
		params := mux.Vars(r)
		id := params["id"]
		log.Debugf("Requested id %v", id)
		if _, ok := txtData[id]; !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/x-company-v1")
		w.WriteHeader(http.StatusOK)
		txt, err := w.Write([]byte(txtData[id]+"\n"))
		if err != nil {
			log.Debugf("Error writing Response %v ", err)
			return
		}
		if txt == 0 {
			log.Debug("Faied to write bytes to Response")
		}
	})

	log.Info("Starting the test server us version1 on port 9002 ... ")
	log.Fatal(http.ListenAndServe(":9002", rtr))
}
