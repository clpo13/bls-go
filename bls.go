// bls-go - a Go interface for BLS.gov
// Copyright 2018 Cody Logan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package blsgo is a Go interface for getting and parsing data from the BLS.gov API.
package blsgo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Default API endpoint
const uri = "https://api.bls.gov/publicAPI/v2/timeseries/data"

// Footnote represents a footnote sent with the received data.
type Footnote struct {
	Code string `json:"code"`
	Text string `json:"text"`
}

// Change represents calculations of change over 1, 3, 6, and 12 months.
type Change struct {
	OneMonth    string `json:"1"`
	ThreeMonth  string `json:"3"`
	SixMonth    string `json:"6"`
	TwelveMonth string `json:"12"`
}

// Calculation represents calculations sent with the received data.
type Calculation struct {
	NetChange Change `json:"net_changes"`
	PctChange Change `json:"pct_changes"`
}

// Period represents an individual period (usually a month) of data.
type Period struct {
	Year         string       `json:"year"`
	Num          string       `json:"period"`
	Name         string       `json:"periodName"`
	Value        string       `json:"value"`
	Footnotes    []Footnote   `json:"footnotes"`
	Calculations *Calculation `json:"calculations,omitempty"`
}

// Catalog represents the catalog data of a series.
type Catalog struct {
	Title    string `json:"series_title"`
	ID       string `json:"series_id"`
	Season   string `json:"seasonality"`
	Name     string `json:"survey_name"`
	Abbr     string `json:"survey_abbreviation"`
	DataType string `json:"measure_data_type"`
	Area     string `json:"area"`
	AreaType string `json:"area_type"`
}

// SeriesData represents data from a single series.
type SeriesData struct {
	SeriesID string   `json:"seriesID"`
	Catalog  *Catalog `json:"catalog,omitempty"`
	Data     []Period `json:"data"`
}

// Series represents a JSON object containing arrays of series.
type Series struct {
	Series []SeriesData `json:"series"`
}

// ResultData represents the top-level structure of the received data.
type ResultData struct {
	Status       string   `json:"status"`
	ResponseTime int      `json:"responseTime"`
	Message      []string `json:"message"`
	Results      Series   `json:"Results"`
}

// Payload defines the structure of data to be sent to an API endpoint.
type Payload struct {
	Start   string   `json:"startyear"`
	End     string   `json:"endyear"`
	Series  []string `json:"seriesid"`
	Catalog bool     `json:"catalog"`
	Calc    bool     `json:"calculations"`
	Avg     bool     `json:"annualaverage"`
	Key     string   `json:"registrationkey"`
}

// GetData takes a Payload, converts it to JSON, and sends it to the specified
// API endpoint. It returns a ResultData object that contains all the received
// data for the given Payload.
func GetData(payload Payload) ResultData {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln("There was an error parsing the payload:", err)
	}

	// Set up an HTTP GET request.
	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalln("There was an error getting the URI:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// Body returned as []byte, so convert to string to print it out.
	//fmt.Printf("Raw response as a string:\n%v\n", string(body))

	return parseData(body)
}

// Reverse takes an array of Periods and reverses it.
func Reverse(a []Period) []Period {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

// Map JSON response to Go structs.
func parseData(body []byte) ResultData {
	var rd ResultData
	err := json.Unmarshal(body, &rd)
	if err != nil {
		log.Fatalln("There was an error parsing the JSON response:", err)
	}
	return rd
}
