package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
)

type listMapResults []map[string]string

var (
	totalResults int
	resultSearch *[]byte
	results      listMapResults
)

const urlToSearch string = "https://gist.githubusercontent.com/henriquecode/e6e4f907f1e2cd073705561d67b189fc/raw/d0dc9e98eed30e39605a8b454077be3ebf10dc28/page-search-go.html"

func main() {

	search()
	writeResultInFile()
	saveDataInFileJSON(extractResult())
}

func search() {

	response, err := http.Get(urlToSearch)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	resultSearch = &body
}

func writeResultInFile() {
	err := ioutil.WriteFile("result-google.txt", *resultSearch, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func extractResult() listMapResults {
	reItem := regexp.MustCompile(`class="item"`)
	reTitle := regexp.MustCompile(`<div class="title">([a-zA-Z0-9 \!\@\#\$\;':",.]+)<\/div>`)
	reContent := regexp.MustCompile(`<div class="content">([a-zA-Z0-9 \!\@\#\$\;':",.]+)<\/div>`)
	reURL := regexp.MustCompile(`<div class="url">([a-zA-Z0-9 \!\@\#\$\;':",.]+)<\/div>`)

	items := reItem.FindAll([]byte(*resultSearch), -1)
	totalItems := len(items)

	titles := reTitle.FindAllStringSubmatch(string(*resultSearch), -1)
	contents := reContent.FindAllStringSubmatch(string(*resultSearch), -1)
	urls := reURL.FindAllStringSubmatch(string(*resultSearch), -1)

	for totalItems > 0 {
		results = append(results, map[string]string{
			"id":      strconv.Itoa(totalItems),
			"title":   titles[totalItems-1][1],
			"content": contents[totalItems-1][1],
			"url":     urls[totalItems-1][1],
		})
		totalItems--
	}

	sort.Slice(results, func(i, j int) bool {
		vA, errA := strconv.Atoi(results[i]["id"])
		if errA != nil {
			panic(errA)
		}

		vB, errB := strconv.Atoi(results[j]["id"])
		if errB != nil {
			panic(errB)
		}

		return vA < vB
	})

	return results
}

func saveDataInFileJSON(value listMapResults) {
	response, errJSON := json.Marshal(value)

	if errJSON != nil {
		panic(errJSON)
	}

	errFile := ioutil.WriteFile("result-google.json", response, 0644)
	if errFile != nil {
		log.Fatal(errFile)
	}
}
