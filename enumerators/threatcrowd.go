package enumerators

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/deckarep/golang-set"
)

// A struct in which we will unmarshal results to
type Subs struct {
	Subdomains []string
}

func ThreatCrowdEnum(domain string, outCh chan<- mapset.Set, errCh chan<- EnumError, wg *sync.WaitGroup) {
	defer wg.Done()
	client := http.Client{}

	url := fmt.Sprintf("https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=%s", domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errCh <- EnumError{err, "ThreatCrowdEnum"}
		return
	}

	SetupRequestHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		errCh <- EnumError{err, "ThreatCrowdEnum"}
		return
	}

	subdomains, err := JsonParser(resp)
	resp.Body.Close()

	if err != nil {
		errCh <- EnumError{err, "ThreatCrowdEnum"}
		return
	}
	outCh <- subdomains
}

func JsonParser(resp *http.Response) (mapset.Set, error) {
	var subs Subs
	extractedSubs := mapset.NewSet()

	body, err := ioutil.ReadAll(resp.Body) // We trust that the json returned won't be too big
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &subs); err != nil {
		return nil, err
	}

	for _, sub := range subs.Subdomains {
		extractedSubs.Add(sub)
	}
	return extractedSubs, nil
}
