package enumerators

import (
  "fmt"
  "net/http"
  "log"
  "github.com/deckarep/golang-set"
  "sync"
)

func CrtEnum(domain string, ch chan<- mapset.Set, wg *sync.WaitGroup) (error) {
  client := http.Client{}
  defer wg.Done()

  url := fmt.Sprintf("https://crt.sh/?q=%%.%s", domain)
  req, _ := http.NewRequest("GET", url, nil)
  SetupRequestHeaders(req)

  resp, err := client.Do(req)
  if err != nil {
    log.Fatal(err)
  }

  // Selector accepts table rows without anchor children and style attributes
  subdomains := GenericParser(resp, "td:not(:has(a)):not([style])")
  if subdomains != nil {
    ch <- subdomains
  }
  fmt.Println("Crt finished...")
  return nil
}
