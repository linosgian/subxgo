package enumerators

import (
  "fmt"
  "net/http"
  "log"
  "github.com/deckarep/golang-set"
  "sync"
)

func VirustotalEnum(domain string, ch chan<- mapset.Set, wg *sync.WaitGroup) (error) {
  // TODO: Rate limiting should be enforced
  defer wg.Done()

  client := http.Client{}
  url := fmt.Sprintf("https://virustotal.com/en/domain/%s/information/", domain)
  req, _ := http.NewRequest("GET", url, nil)
  SetupRequestHeaders(req)

  resp, err := client.Do(req)
  if err != nil {
    log.Fatal(err)
  }

  subdomains := GenericParser(resp, "div#observed-subdomains a")
  resp.Body.Close()

  if subdomains != nil {
    ch <- subdomains
  }
  fmt.Println("Virustotal finished...")
  return nil
}
