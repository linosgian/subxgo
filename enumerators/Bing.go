package enumerators
import (
  "fmt"
  "net/http"
  "log"
  "os"
  "github.com/deckarep/golang-set"
  "sync"
)


func BingEnum(domain string, ch chan<- mapset.Set, wg *sync.WaitGroup) (error){
  prevSubdomains := mapset.NewSet()
  MaxDomains := 30
  noNewSubdomains := 0
  client := http.Client{}
  defer wg.Done()
  subdomains := mapset.NewSet()
  
  page_num := 1
  for {
    if subdomains.Cardinality() >= MaxDomains {
      page_num += 10
    }
    // TODO: Check if Request objects can be reused
    final_uri := bingBuildQuery(domain, page_num, &subdomains, MaxDomains)
    fmt.Println(final_uri)
    req, _ := http.NewRequest("GET", final_uri, nil)
    SetupRequestHeaders(req)

    resp, err := client.Do(req)
    if err != nil {
      // TODO: Handle error gracefully
      log.Fatal(err)
      os.Exit(1)
    }

    currSubdomains := GenericParser(resp, "cite")
    resp.Body.Close()

    if currSubdomains == nil {
      return nil
    }
    fmt.Println(currSubdomains)
    diff := currSubdomains.Difference(subdomains)

    if currSubdomains.Equal(prevSubdomains) {
      fmt.Println("Same")
      if noNewSubdomains == 3{
        break
      }else{
        noNewSubdomains += 1
      }
    }else{
      ch <- diff
    }
    prevSubdomains = currSubdomains

    subdomains = subdomains.Union(currSubdomains)
  }
  fmt.Println("Bing finished...")
  return nil
}

func bingBuildQuery(domain string, pageNum int, subdomains *mapset.Set, MaxDomains int) (string) {
  baseUrl := "https://www.bing.com/search?go=Submit&first=%d&q=%s"
  baseQuery := fmt.Sprintf("domain:%s+-www.%s", domain, domain)

  finalQuery := fmt.Sprintf(baseUrl, pageNum, baseQuery)

  sliced_subdomains := (*subdomains).ToSlice()
  for index, subdomain := range sliced_subdomains{
    if subdomain == domain {
      continue
    }

    finalQuery = fmt.Sprintf("%s+-domain:%s", finalQuery, subdomain)

    if index >= MaxDomains - 1 {
      break
    }
  }
  return finalQuery
}
