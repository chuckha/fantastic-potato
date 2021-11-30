package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	contents, err := os.ReadFile("prefectures copy.json")
	if err != nil {
		panic(fmt.Sprintf("failed to readfile %+v", err))
	}
	out := &geojson{}
	if err := json.Unmarshal(contents, out); err != nil {
		panic(fmt.Sprintf("failed to unmarshal %+v", err))
	}
	data, err := json.Marshal(out)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal %+v", err))
	}
	if err := ioutil.WriteFile("test.json", data, 0o755); err != nil {
		panic(fmt.Sprintf("failed to write file %+v", err))
	}

}

type geojson struct {
	Type     string    `json:"type"`
	Name     string    `json:"name"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string          `json:"type"`
	Properties Properties      `json:"properties"`
	Geometry   json.RawMessage `json:"geometry"`
}

type Properties struct {
	Name         string `json:"name"`
	NameEnglish  string `json:"name_en"`
	NameJapanese string `json:"name_ja"`
}
