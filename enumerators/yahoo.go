package enumerators

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/deckarep/golang-set"
)

func YahooEnum(domain string, outCh chan<- mapset.Set, errCh chan<- EnumError, wg *sync.WaitGroup) {
	const (
		maxDomains = 11 // Maximum number of excluded domains for Yahoo's -domain option
		maxRetries = 3
	)

	var (
		firstSubs []interface{} // Holds the first "maxDomains" subs
		currSubs  mapset.Set
	)

	prevSubs := mapset.NewSet() // Holds the previous request's extracted subs
	subs := mapset.NewSet()     // Holds all the subs found by the enumerator

	retries := 0
	pageNumber := 1

	client := http.Client{}

	// Notify the master that the enumerator is done when this func is done
	defer wg.Done()

	for {
		if subs.Cardinality() >= maxDomains {
			// When we reach the "maxDomains" amount of subdomains
			// Start iterating pages and dont change the requested URI anymore
			if len(firstSubs) != maxDomains { //
				firstSubs = (subs.ToSlice())[0:maxDomains]
			} else {
				pageNumber += 10 // Pages are changed every 10 search results
			}
		} else { // If we don't have enough subdomains yet, scrape a few pages more.
			firstSubs = subs.ToSlice()
			pageNumber += 10
		}

		finalUri := yahooBuildQuery(domain, pageNumber, &firstSubs)
		req, err := http.NewRequest("GET", finalUri, nil)

		if err != nil {
			errCh <- EnumError{err, "YahooEnum"}
			return
		}
		// Prepares HTTP Headers to bypass bot detection
		SetupRequestHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			errCh <- EnumError{err, "YahooEnum"}
			return
		}

		currSubs, err = GenericParser(resp, "span.fz-ms")
		resp.Body.Close()

		if err != nil {
			retries += 1 // If the parser didnt yield any results, maybe reached the end.
			return
		}

		if currSubs.Equal(prevSubs) {
			// If the same subdomains appear n times in a row, we might have reached the end.
			if retries == maxRetries {
				// log.Println("Yahoo finished after all pages...")
				outCh <- subs
				return
			}
			retries += 1
		} else {
			retries = 1
		}
		prevSubs = currSubs

		subs = subs.Union(currSubs)
	}
}

// Takes a target domain, the current page number and the first "maxDomains" subdomains
// and builds the final URL.
func yahooBuildQuery(domain string, pageNum int, firstSubs *[]interface{}) string {
	baseUrl := "https://search.yahoo.com/search?b=%d&p=%s"
	baseUri := fmt.Sprintf("site:%s+-domain:www.%s", domain, domain)

	finalQuery := fmt.Sprintf(baseUrl, pageNum, baseUri)
	for _, sub := range *firstSubs {
		if sub == domain { // if we exclude the target domain, the request will fail
			continue
		}
		finalQuery = fmt.Sprintf("%s+-domain:%s", finalQuery, sub)
	}
	return finalQuery
}
