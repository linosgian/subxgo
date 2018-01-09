package enumerators

import (
  "fmt"
  "net/http"
  "log"
  "os"
  "github.com/deckarep/golang-set"
  "sync"
)



func YahooEnum(domain string, ch chan<- mapset.Set, wg *sync.WaitGroup) (error) {
  prevSubdomains := mapset.NewSet()
  MaxDomains := 10
  noNewSubdomains := 0
  client := http.Client{}
  // Notify the master that the enumerator is done when this func is done
  defer wg.Done()
  subdomains := mapset.NewSet()
  
  page_num := 1
  for {
    if subdomains.Cardinality() >= MaxDomains {
      page_num += 10
    }

    final_uri := yahooBuildQuery(domain, page_num, &subdomains, MaxDomains)
    req, _ := http.NewRequest("GET", final_uri, nil)

    // Prepares HTTP Headers to bypass bot detection
    SetupRequestHeaders(req)

    resp, err := client.Do(req)
    if err != nil {
      // TODO: Handle error gracefully
      log.Fatal(err)
      os.Exit(1)
    }

    currSubdomains := GenericParser(resp, "span.fz-ms")
    resp.Body.Close()

    if currSubdomains == nil {
      return nil
    }
    diff := currSubdomains.Difference(subdomains)

    //TODO: Change this to a switch case??
    if currSubdomains.Equal(prevSubdomains){
      if noNewSubdomains == 3 {
        return nil
      }else {
        noNewSubdomains += 1
      }
    }else{
      ch <- diff
    }
    prevSubdomains = currSubdomains

    subdomains = subdomains.Union(currSubdomains)
  }
  fmt.Println("Yahoo finished...")
  return nil
}

func yahooBuildQuery(domain string, pageNum int, subdomains *mapset.Set, MaxDomains int) (string){
  baseUrl := "https://search.yahoo.com/search?b=%d&p=%s"
  baseQuery := fmt.Sprintf("site:%s+-domain:www.%s", domain, domain)

  finalQuery := fmt.Sprintf(baseUrl, pageNum, baseQuery)
  sliced_subdomains := (*subdomains).ToSlice()
  for index, subdomain := range sliced_subdomains{
    if subdomain == domain{
      continue
    }

    finalQuery = fmt.Sprintf("%s+-domain:%s", finalQuery, subdomain)

    if index >= MaxDomains - 1 {
      break
    }
  }
  return finalQuery
}
