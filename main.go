package main

import (
  "subxgo/enumerators"
  "github.com/fatih/color"
  "github.com/deckarep/golang-set"
  "sync"
)

const TARGET = "grnet.gr"

func main() {
  printPreamble()

  chosenEnumerators := getEnumerators()
  out := merge(chosenEnumerators)
  // TODO: check if new subdomains are indeed new,
  // by converting a slice of strings into a set and diffing
  for {
    if subs, ok := <-out; ok{
      printNewSubdomains(subs)
    }else{
      break
    }
  }
}

func printNewSubdomains(NewSubdomains mapset.Set){
  sliced_subdomains := NewSubdomains.ToSlice()
  for _, subdomain := range sliced_subdomains {
    color.Magenta(subdomain.(string))
  }
}


func merge(enums map[string]Enumerator) (chan mapset.Set) {
  /*
    Takes pairs of strings and functions (Enumerators) and 
    returns a channel from which the caller will get receive
    the
  */
  var wg sync.WaitGroup
  out := make(chan mapset.Set)
  wg.Add(len(enums))

  for _, enumerator := range enums {
    go enumerator(TARGET, out, &wg)
  }

  go func (){
    wg.Wait()
    close(out)
  }()

  return out
}

func printPreamble() {
  color.Green(`
-------------------------------
 __       _       ___
/ _\_   _| |__   / _ \___
\ \| | | | '_ \ / /_\/ _ \
_\ \ |_| | |_) / /_\\ (_) |
\__/\__,_|_.__/\____/\___/

-------------------------------
  `)
  color.Green("Starting enumeration...")
}

type Enumerator func(string, chan<- mapset.Set, *sync.WaitGroup) error

func getEnumerators() (map[string]Enumerator){
  return map[string]Enumerator {
//    "Yahoo": enumerators.YahooEnum,
//    "Bing" : enumerators.BingEnum,
//    "Virustotal": enumerators.VirustotalEnum,
//    "Threatcrowd": enumerators.ThreatCrowdEnum,
    "Cert": enumerators.CrtEnum,
  }
}
