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

package blsgo

import (
	"io/ioutil"
	"log"
	"testing"
)

// Test that the correct API endpoint is used.
// NOTE: dummy test
func TestURI(t *testing.T) {
	if uri != "https://api.bls.gov/publicAPI/v2/timeseries/data" {
		t.Error("Expected URI to be 'https://api.bls.gov/publicAPI/v2/timeseries/data', but got", uri)
	}
}

// Test that raw JSON is properly parsed into structs.
func TestParse(t *testing.T) {
	/*
		f, err := os.Open("testdata.json")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		dec := json.NewDecoder(f)
		var tr ResultData
		err = dec.Decode(&tr)
		if err != nil {
			panic(err)
		}
	*/

	f, err := ioutil.ReadFile("testdata/two-series.json")
	if err != nil {
		log.Fatalln("Error reading test data file:", err)
	}

	tr := ParseData(f)

	status := tr.Status
	series1 := tr.Results.Series[0].SeriesID
	catalog := tr.Results.Series[0].Catalog
	if status != "REQUEST_SUCCEEDED" {
		t.Error("Expected status to be 'REQUEST_SUCCEEDED', but got", status)
	}
	if series1 != "LEU0254555900" {
		t.Error("Expected first seriesID to be 'LEU0254555900', but got", series1)
	}
	if catalog != nil {
		t.Error("Expected a nil value for catalog data, but got", catalog)
	}
	if len(tr.Message) == 0 {
		t.Error("Expected to find 1 message, but found", len(tr.Message), "instead")
	}
}

// Test that arrays are properly reversed.
func TestReverse(t *testing.T) {
	// Create test periods
	// TODO: create these programmatically
	fn := new([]Footnote)
	calc := new(Calculation)
	p1 := Period{"1913", "M01", "January", "10.0", *fn, calc}
	p2 := Period{"1913", "M02", "February", "11.0", *fn, calc}
	p3 := Period{"1913", "M03", "March", "12.0", *fn, calc}

	pA := []Period{p1, p2, p3}
	pAR := []Period{p3, p2, p1}
	cases := []struct {
		in, want []Period
	}{
		{pA, pAR},
	}
	for _, c := range cases {
		got := Reverse(c.in)
		if got[0].Name != c.want[0].Name {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
