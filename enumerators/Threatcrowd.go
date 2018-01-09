package enumerators

import (
  "fmt"
  "io/ioutil"
  "encoding/json"
  "net/http"
  "log"
  "github.com/deckarep/golang-set"
  "sync"
)

type Subs struct {
  Subdomains []string
}

func ThreatCrowdEnum(domain string, ch chan<- mapset.Set, wg *sync.WaitGroup) (error) {
  defer wg.Done()
  client := http.Client{}

  url := fmt.Sprintf("https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=%s", domain)
  req, _ := http.NewRequest("GET", url, nil)
  SetupRequestHeaders(req)

  resp, err := client.Do(req)
  if err != nil {
    log.Fatal(err)
  }

  subdomains := JsonParser(resp)
  resp.Body.Close()

  if subdomains != nil {
    ch <- subdomains
  }
  fmt.Println("ThreatCrowd finished...")
  return nil
}

func JsonParser(resp *http.Response) (mapset.Set) {
  var subs Subs;
  subdomains := mapset.NewSet()

  body, _ := ioutil.ReadAll(resp.Body)
  err := json.Unmarshal(body, &subs)
  if err != nil{
    fmt.Println(err)
  }
  for _, subdomain := range subs.Subdomains {
    subdomains.Add(subdomain)
  }
  return subdomains
}
