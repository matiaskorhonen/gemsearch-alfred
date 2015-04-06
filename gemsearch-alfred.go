package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/matiaskorhonen/go-alfred"
)

func main() {
	queryTerms := os.Args[1:]
	query := strings.Join(queryTerms, " ")

	// optimize query terms for fuzzy matching
	alfred.InitTerms(queryTerms)

	// create a new alfred workflow response
	response := alfred.NewResponse()
	suggestions := getSuggestions(query)

	searchURL, _ := url.Parse("https://rubygems.org/search?")
	v := url.Values{}
	v.Set("query", query)
	searchURL.RawQuery = v.Encode()

	response.AddItem(&alfred.AlfredResponseItem{
		Valid:    true,
		Uid:      "000000-search-on-rubygems.org" + searchURL.String(),
		Title:    query,
		Subtitle: "Search RubyGems.org for '" + query + "'",
		Arg:      searchURL.String(),
		Icon:     "search.png",
	})

	for _, s := range suggestions {
		// it matched so add a new response item
		response.AddItem(&alfred.AlfredResponseItem{
			Valid:        true,
			Uid:          s.URL,
			Title:        s.Name,
			Autocomplete: s.Name,
			Subtitle:     s.URL,
			Arg:          s.URL,
		})
	}

	// finally print the resulting Alfred Workflow XML
	response.Print()
}

type suggestion struct {
	Name    string `json:"n"`
	Version string `json:"v"`
	URL     string `json:"u"`
}

func getSuggestions(query string) []suggestion {
	u, err := url.Parse("http://www.gemsear.ch/gems/suggestions")
	if err != nil {
		log.Fatal(err)
	}
	v := url.Values{}
	v.Set("q", query)
	v.Set("friendly", "true")

	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var suggestions []suggestion

	err = json.Unmarshal(body, &suggestions)
	if err != nil {
		panic(err)
	}

	return suggestions
}
