package main

import (
	enums "subxgo/enumerators"
	"sync"

	"github.com/deckarep/golang-set"
	"github.com/fatih/color"
)

const TARGET = "grnet.gr"

func main() {
	// printPreamble()

	subs := mapset.NewSet()

	chosenEnumerators := getEnumerators()
	out, errs := merge(chosenEnumerators)
	// TODO: check if new subdomains are indeed new,
	// by converting a slice of strings into a set and diffing
	for {
		select {
		case newSubs, ok := <-out:
			if !ok {
				out = nil
				break // This only breaks out of the select/case statement
			}
			printNewSubdomains(&newSubs, &subs)
			subs = subs.Union(newSubs)
		case enumerr, ok := <-errs:
			if !ok {
				errs = nil
				break
			}
			color.Red(enumerr.Error())
		}

		if out == nil && errs == nil {
			color.Green("All done!")
			return
		}
	}
}

func merge(es map[string]enums.Enumerator) (chan mapset.Set, chan enums.EnumError) {
	/*
		Takes pairs of strings and functions (Enumerators) and
		returns a channel from which the caller will get receive
		the results
	*/
	var wg sync.WaitGroup
	out := make(chan mapset.Set)
	errs := make(chan enums.EnumError)
	wg.Add(len(es))

	for _, enumerator := range es {
		go enumerator(TARGET, out, errs, &wg)
	}

	go func() {
		wg.Wait()
		close(out)
		close(errs)
	}()

	return out, errs
}

func getEnumerators() map[string]enums.Enumerator {
	return map[string]enums.Enumerator{
		"Yahoo":       enums.YahooEnum,
		"Bing":        enums.BingEnum,
		"Virustotal":  enums.VirustotalEnum,
		"Threatcrowd": enums.ThreatCrowdEnum,
		"Cert":        enums.CrtEnum,
	}
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

func printNewSubdomains(subs, newSubs *mapset.Set) {
	uniqueNewSubs := (*subs).Difference(*newSubs)
	slicedSubs := uniqueNewSubs.ToSlice()
	for _, sub := range slicedSubs {
		color.Magenta(sub.(string))
	}
}
