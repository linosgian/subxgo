package enumerators

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/deckarep/golang-set"
)

func CrtEnum(domain string, outCh chan<- mapset.Set, errCh chan<- EnumError, wg *sync.WaitGroup) {
	client := http.Client{}
	defer wg.Done()

	url := fmt.Sprintf("https://crt.sh/?q=%%.%s", domain)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		errCh <- EnumError{err, "CrtEnum"}
		return
	}
	SetupRequestHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		errCh <- EnumError{err, "CrtEnum"}
		return
	}

	// Selector accepts table rows without anchor children and style attributes
	subs, err := GenericParser(resp, "td:not(:has(a)):not([style])")
	if err != nil {
		errCh <- EnumError{err, "CrtEnum"}
		return
	}

	outCh <- subs
}
