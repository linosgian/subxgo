package enumerators

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/deckarep/golang-set"
)

func VirustotalEnum(domain string, outCh chan<- mapset.Set, errCh chan<- EnumError, wg *sync.WaitGroup) {
	// TODO: Rate limiting should be enforced
	defer wg.Done()

	client := http.Client{}

	url := fmt.Sprintf("https://virustotal.com/en/domain/%s/information/", domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errCh <- EnumError{err, "VirustotalEnum"}
		return
	}

	SetupRequestHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		errCh <- EnumError{err, "VirustotalEnum"}
		return
	}

	subdomains, err := GenericParser(resp, "div#observed-subdomains a")
	resp.Body.Close()

	if err != nil {
		errCh <- EnumError{err, "VirustotalEnum"}
		return
	}
	outCh <- subdomains
}
