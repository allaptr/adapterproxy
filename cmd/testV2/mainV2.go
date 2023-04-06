package main

import (
	"flag"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var txtData = map[string]string{
	"621451430762XA": "{\"company_name\":\"ACME Incorporated\",\"tin\":\"2314-12-23T01:29:13Z\",\"dissolved_on\":\"2322-03-24T01:54:39Z\"}",
	"testv2": "{\"company_name\":\"laughing-lamarr-dirac\",\"tin\":\"2022-03-24T01:54:39Z\",\"dissolved_on\":null}",
	"440804282323AC": "{\"company_name\":\"elastic-aryabhata-kare\",\"tin\":\"1955-04-10T01:54:39Z\",\"dissolved_on\":\"2022-03-24T01:54:39Z\"}",
	"117441277681BD": "{\"company_name\":\"friendly-euler-ptolemy\",\"tin\":\"1977-04-04T01:54:39Z\",\"dissolved_on\":null}",
	"4164142148556E": "{\"company_name\":\"flamboyant-burnell-colden\",\"tin\":\"1745-05-30T01:54:39Z\",\"dissolved_on\":null}",
	"155674243726OF": "{\"company_name\":\"stoic-edison-feynman\",\"tin\":\"1786-05-20T01:54:39Z\",\"dissolved_on\":null}",
	"462059571630QG": "{\"company_name\":\"agitated-bell-easley\",\"tin\":\"2006-03-28T01:54:39Z\",\"dissolved_on\":null}",
	"854145342086DH": "{\"company_name\":\"unruffled-moore-hodgkin\",\"tin\":\"1777-05-22T01:54:39Z\",\"dissolved_on\":null}",
	"464415421689NI": "{\"company_name\":\"tender-swartz-murdock\",\"tin\":\"1771-05-24T01:54:39Z\",\"dissolved_on\":null}",
	"634374392121WJ": "{\"company_name\":\"serene-chaum-bhabha\",\"tin\":\"1869-04-30T01:54:39Z\",\"dissolved_on\":null}",
	"877991708633OK": "{\"company_name\":\"pensive-shockley-clarke\",\"tin\":\"1918-04-19T01:54:39Z\",\"dissolved_on\":null}",
	"6187217740522L": "{\"company_name\":\"elated-banach-goldstine\",\"tin\":\"1894-04-24T01:54:39Z\",\"dissolved_on\":null}",
	"918600199176TM": "{\"company_name\":\"reverent-chebyshev-engelbart\",\"tin\":\"1919-04-19T01:54:39Z\",\"dissolved_on\":null}",
	"369401320059GN": "{\"company_name\":\"festive-cray-bose\",\"tin\":\"1753-05-28T01:54:39Z\",\"dissolved_on\":null}",
	"525570322146UP": "{\"company_name\":\"mystifying-proskuriakova-mcclintock\",\"tin\":\"1759-05-27T01:54:39Z\",\"dissolved_on\":null}",
	"372581767233LQ": "{\"company_name\":\"dazzling-newton-cannon\",\"tin\":\"1925-04-17T01:54:39Z\",\"dissolved_on\":null}",
	"478803106995ZR": "{\"company_name\":\"pensive-bose-bhaskara\",\"tin\":\"1743-05-31T01:54:39Z\",\"dissolved_on\":null}",
	"655358414930YS": "{\"company_name\":\"tender-gauss-shirley\",\"tin\":\"1977-04-04T01:54:39Z\",\"dissolved_on\":null}",
	"156430540509YT": "{\"company_name\":\"confident-mcclintock-haibt\",\"tin\":\"1943-04-13T01:54:39Z\",\"dissolved_on\":\"1995-04-17T01:54:39Z\"}",
	"3768800898314U": "{\"company_name\":\"quirky-rhodes-carson\",\"tin\":\"1824-05-11T01:54:39Z\",\"dissolved_on\":null}",
	"3508905365551V": "{\"company_name\":\"epic-bose-brown\",\"tin\":\"1763-05-26T01:54:39Z\",\"dissolved_on\":null}",
	"599191407388UW": "{\"company_name\":\"agitated-cohen-lehmann\",\"tin\":\"1869-04-30T01:54:39Z\",\"dissolved_on\":null}",
	"150627018476X7": "{\"company_name\":\"gracious-tu-jepsen\",\"tin\":\"1864-05-01T01:54:39Z\",\"dissolved_on\":null}",
	"548431946394UY": "{\"company_name\":\"eager-jones-brahmagupta\",\"tin\":\"1869-04-30T01:54:39Z\",\"dissolved_on\":null}",
	"538767862374SZ": "{\"company_name\":\"baby-hertz-burnell\",\"tin\":\"1826-05-11T01:54:39Z\",\"dissolved_on\":null}",
	"920598541601S1": "{\"company_name\":\"pensive-austin-taussig\",\"tin\":\"1890-04-25T01:54:39Z\",\"dissolved_on\":null}",
	"311130734087N2": "{\"company_name\":\"heuristic-nash-darwin\",\"tin\":\"1898-04-23T01:54:39Z\",\"dissolved_on\":null}",
	"378339526656G3": "{\"company_name\":\"romantic-wu-blackwell\",\"tin\":\"1806-05-16T01:54:39Z\",\"dissolved_on\":null}",
	"552624153194S6": "{\"company_name\":\"intelligent-bassi-jones\",\"tin\":\"1795-05-18T01:54:39Z\",\"dissolved_on\":null}",
	"801922059877A7": "{\"company_name\":\"hardcore-satoshi-cohen\",\"tin\":\"1996-03-30T01:54:39Z\",\"dissolved_on\":null}",
	"272745309677HS": "{\"company_name\":\"beautiful-haslett-bohr\",\"tin\":\"1957-04-09T01:54:39Z\",\"dissolved_on\":null}",
	"208734339345LB": "{\"company_name\":\"sad-hawking-fermat\",\"tin\":\"1852-05-04T01:54:39Z\",\"dissolved_on\":null}",
	"138840155253JN": "{\"company_name\":\"friendly-babbage-pare\",\"tin\":\"1982-04-03T01:54:39Z\",\"dissolved_on\":null}",
	"642597231447KG": "{\"company_name\":\"sleepy-mcclintock-goldwasser\",\"tin\":\"1770-05-24T01:54:39Z\",\"dissolved_on\":null}",
	"961478751843NS": "{\"company_name\":\"heuristic-vaughan-lederberg\",\"tin\":\"1919-04-19T01:54:39Z\",\"dissolved_on\":null}",
	"739077213451MU": "{\"company_name\":\"priceless-swanson-sinoussi\",\"tin\":\"1870-04-30T01:54:39Z\",\"dissolved_on\":null}",
	"17937688759 7J": "{\"company_name\":\"zealous-germain-bartik\",\"tin\":\"1815-05-14T01:54:39Z\",\"dissolved_on\":null}",
	"66112741425T96": "{\"company_name\":\"condescending-shamir-chaplygin\",\"tin\":\"1843-05-07T01:54:39Z\",\"dissolved_on\":null}",
	"56705623572D7O": "{\"company_name\":\"amazing-proskuriakova-khorana\",\"tin\":\"1855-05-04T01:54:39Z\",\"dissolved_on\":null}",
	"29385698978P5G": "{\"company_name\":\"magical-almeida-napier\",\"tin\":\"1831-05-10T01:54:39Z\",\"dissolved_on\":null}",
	"78611535928K6E": "{\"company_name\":\"jovial-bartik-cannon\",\"tin\":\"1939-04-14T01:54:39Z\",\"dissolved_on\":null}",
	"8575350780T07Y": "{\"company_name\":\"gifted-mcnulty-chandrasekhar\",\"tin\":\"1993-03-31T01:54:39Z\",\"dissolved_on\":null}",
	"868649496F404U": "{\"company_name\":\"angry-williamson-hugle\",\"tin\":\"1838-05-08T01:54:39Z\",\"dissolved_on\":null}",
	"565058762G528T": "{\"company_name\":\"hopeful-chatelet-villani\",\"tin\":\"1874-04-29T01:54:39Z\",\"dissolved_on\":null}",
	"5216802H19801U": "{\"company_name\":\"quirky-banzai-einstein\",\"tin\":\"1743-05-31T01:54:39Z\",\"dissolved_on\":null}",
	"4118616J38379Q": "{\"company_name\":\"condescending-tharp-clarke\",\"tin\":\"1820-05-12T01:54:39Z\",\"dissolved_on\":null}",
	"51178474D1338D": "{\"company_name\":\"charming-jang-mccarthy\",\"tin\":\"1747-05-30T01:54:39Z\",\"dissolved_on\":null}",
	"1616114D78308D": "{\"company_name\":\"zen-lamport-curie\",\"tin\":\"1854-05-04T01:54:39Z\",\"dissolved_on\":null}",
	"6729D95902335D": "{\"company_name\":\"gallant-allen-northcutt\",\"tin\":\"1800-05-17T01:54:39Z\",\"dissolved_on\":null}",
	"4374990612D497": "{\"company_name\":\"vigorous-dirac-mclean\",\"tin\":\"1745-05-30T01:54:39Z\",\"dissolved_on\":null}",
	"18146D5231191A": "{\"company_name\":\"recursing-khorana-spence\",\"tin\":\"1867-05-01T01:54:39Z\",\"dissolved_on\":null}",
	"72588052F4592F": "{\"company_name\":\"festive-spence-burnell\",\"tin\":\"2010-03-27T01:54:39Z\",\"dissolved_on\":null}",
	"67418466F3833P": "{\"company_name\":\"angry-curran-boyd\",\"tin\":\"1988-04-01T01:54:39Z\",\"dissolved_on\":null}",
	"577F372190779C": "{\"company_name\":\"fervent-lalande-galileo\",\"tin\":\"1734-06-02T01:54:39Z\",\"dissolved_on\":null}",
	"57961072T94366": "{\"company_name\":\"sad-yonath-germain\",\"tin\":\"1929-04-16T01:54:39Z\",\"dissolved_on\":null}",
	"666898Y706315C": "{\"company_name\":\"flamboyant-poincare-galileo\",\"tin\":\"1816-05-13T01:54:39Z\",\"dissolved_on\":\"1925-04-17T01:54:39Z\"}",
	"8312921Y969034": "{\"company_name\":\"compassionate-wescoff-diffie\",\"tin\":\"1967-04-07T01:54:39Z\",\"dissolved_on\":null}",
	"40808U12759765": "{\"company_name\":\"intelligent-knuth-greider\",\"tin\":\"1764-05-25T01:54:39Z\",\"dissolved_on\":null}",
	"111I562191524H": "{\"company_name\":\"thirsty-yonath-jang\",\"tin\":\"1780-05-21T01:54:39Z\",\"dissolved_on\":null}",
	"7521I428441944": "{\"company_name\":\"festive-kirch-blackwell\",\"tin\":\"1783-05-21T01:54:39Z\",\"dissolved_on\":null}",
	"14189691O23973": "{\"company_name\":\"gifted-kirch-liskov\",\"tin\":\"1888-04-25T01:54:39Z\",\"dissolved_on\":null}",
	"60560163O0650G": "{\"company_name\":\"vigilant-driscoll-mccarthy\",\"tin\":\"1824-05-11T01:54:39Z\",\"dissolved_on\":null}",
	"2515842L13975A": "{\"company_name\":\"magical-lovelace-sinoussi\",\"tin\":\"1740-05-31T01:54:39Z\",\"dissolved_on\":null}",
	"90050247K8423A": "{\"company_name\":\"dreamy-cannon-bartik\",\"tin\":\"1901-04-23T01:54:39Z\",\"dissolved_on\":null}",
	"58993443H4900W": "{\"company_name\":\"zen-haibt-solomon\",\"tin\":\"1860-05-02T01:54:39Z\",\"dissolved_on\":null}",
	"758852J378501J": "{\"company_name\":\"suspicious-leavitt-lovelace\",\"tin\":\"1873-04-29T01:54:39Z\",\"dissolved_on\":null}",
	"28153912D5058U": "{\"company_name\":\"zealous-mccarthy-cartwright\",\"tin\":\"1920-04-18T01:54:39Z\",\"dissolved_on\":null}",
	"594163701S584O": "{\"company_name\":\"osom-merkle-margulis\",\"tin\":\"1805-05-16T01:54:39Z\",\"dissolved_on\":null}",
	"6695590F571003": "{\"company_name\":\"friendly-jackson-hermann\",\"tin\":\"1924-04-17T01:54:39Z\",\"dissolved_on\":null}",
	"17045668Z5224Z": "{\"company_name\":\"unruffled-nobel-dirac\",\"tin\":\"1842-05-07T01:54:39Z\",\"dissolved_on\":null}",
	"4511729Z044489": "{\"company_name\":\"tender-bartik-haslett\",\"tin\":\"1832-05-09T01:54:39Z\",\"dissolved_on\":null}",
	"457463C951243K": "{\"company_name\":\"jolly-kapitsa-bell\",\"tin\":\"1886-04-26T01:54:39Z\",\"dissolved_on\":null}",
	"32987353C3600A": "{\"company_name\":\"bold-lovelace-jennings\",\"tin\":\"1967-04-07T01:54:39Z\",\"dissolved_on\":null}",
	"1357V807041058": "{\"company_name\":\"friendly-lichterman-darwin\",\"tin\":\"1837-05-08T01:54:39Z\",\"dissolved_on\":\"2035-04-17T01:54:39Z\"}",
	"81011485C55277": "{\"company_name\":\"infallible-antonelli-cohen\",\"tin\":\"1988-04-01T01:54:39Z\",\"dissolved_on\":null}",
	"100009G4032731": "{\"company_name\":\"condescending-lamport-jang\",\"tin\":\"1796-05-17T01:54:39Z\",\"dissolved_on\":null}",
	"69939806H61241": "{\"company_name\":\"inspiring-mclean-beaver\",\"tin\":\"1795-05-18T01:54:39Z\",\"dissolved_on\":null}",
	"6012712F65201E": "{\"company_name\":\"stupefied-villani-cannon\",\"tin\":\"1969-04-06T01:54:39Z\",\"dissolved_on\":null}",
	"213093D596627C": "{\"company_name\":\"ACME Canary LTD\",\"tin\":\"1968-04-06T01:54:39Z\",\"dissolved_on\":null}",
	"23480123F7796A": "{\"company_name\":\"ACME LRU\",\"tin\":\"1972-04-05T01:54:39Z\",\"dissolved_on\":null}",
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}
// ru
func main() {
	upper := flag.Duration("delay", 250, "Upper limit of random delay in milliseconds")
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
		// add random delay
		upperlimit := int(*upper)
		log.Debugf("upperlimit %v", upperlimit)
		delay := rand.Intn(upperlimit)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		log.Debugf("Delay %v", delay)

		w.Header().Set("Content-Type", "application/x-company-v2")
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

	log.Info("Starting the test server us version2 on port 9003 ... ")
	log.Fatal(http.ListenAndServe(":9003", rtr))
}
