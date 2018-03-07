package enumerators

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/deckarep/golang-set"
	"github.com/goware/urlx"
)

func SetupRequestHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8")
}

func GenericParser(response *http.Response, selector string) (mapset.Set, error) {
	subdomains := mapset.NewSet()

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	selectedTags := doc.Find(selector)
	if selectedTags.Nodes == nil {
		return nil, fmt.Errorf("No nodes found by parser")
	}

	replacer := strings.NewReplacer(" ", "", "\n", "")
	for i := range selectedTags.Nodes {
		tag := selectedTags.Eq(i)
		trimmedUrl := replacer.Replace(tag.Text())
		parsedUrl, err := urlx.Parse(trimmedUrl)
		if err != nil {
			continue
		}
		subdomains.Add(parsedUrl.Host)
	}
	return subdomains, nil
}
