package main

type jsonObj interface {
	parseJSON() any
}

type cliCommand struct {
	name        string
	description string
	callBack    func(args ...string) error
}

type locationAreaJSON struct {
	Count   int    `json:"count"`
	Next    string `json:"next"`
	Prev    string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
