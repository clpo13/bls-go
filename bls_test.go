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
	expURI := "https://api.bls.gov/publicAPI/v2/timeseries/data"
	if uri != expURI {
		t.Errorf("Expected URI to be '%v', but got '%v'", expURI, uri)
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

	tr, err := parseData(f)

	expStatus := "REQUEST_SUCCEEDED"
	expSeries := "LEU0254555900"

	status := tr.Status
	series1 := tr.Results.Series[0].SeriesID
	catalog := tr.Results.Series[0].Catalog
	if status != expStatus {
		t.Errorf("Expected status to be '%v', but got '%v'", expStatus, status)
	}
	if series1 != expSeries {
		t.Errorf("Expected first seriesID to be '%v', but got '%v'", expSeries, series1)
	}
	if catalog != nil {
		t.Error("Expected a nil value for catalog data, but got", catalog)
	}
	if len(tr.Message) != 1 {
		t.Errorf("Expected to find 1 message, but found %v messages instead", len(tr.Message))
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
			t.Errorf("Reverse(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}

// Test that a request with a bad series is handled properly.
func TestBadSeries(t *testing.T) {
	f, err := ioutil.ReadFile("testdata/bad-series.json")
	if err != nil {
		log.Fatalln("Error reading test data file:", err)
	}

	_, err = parseData(f)

	if err == nil {
		t.Error("Expected an error")
	}

	if aerr, ok := err.(*DataError); ok {
		expMsg := "An invalid series was requested"
		expDetail1 := "Invalid Series for Series FOO"
		expDetail2 := "Unable to get Catalog Data for series FOO"

		if aerr.Msg != expMsg {
			t.Errorf("Expected status to be '%v', but got '%v'", expMsg, aerr.Msg)
		}
		if aerr.Details[0] != expDetail1 {
			t.Errorf("Expected error message to be '%v', but got '%v'", expDetail1, aerr.Details[0])
		}
		if aerr.Details[1] != expDetail2 {
			t.Errorf("Expected error message to be '%v', but got '%v'", expDetail2, aerr.Details[1])
		}
	}
}

// Test that error messages sent by the server are stored in an error struct.
func TestErrorStatus(t *testing.T) {
	f, err := ioutil.ReadFile("testdata/error-status.json")
	if err != nil {
		log.Fatalln("Error reading test data file:", err)
	}

	_, err = parseData(f)

	if err == nil {
		t.Error("Expected an error")
	}

	if aerr, ok := err.(*DataError); ok {
		expMsg := "REQUEST_FAILED_INVALID_PARAMETERS"
		expDetail := "startyear: Value must be a four-digit number."

		if aerr.Msg != expMsg {
			t.Errorf("Expected status to be '%v', but got '%v'", expMsg, aerr.Msg)
		}
		if aerr.Details[0] != expDetail {
			t.Errorf("Expected error message to be '%v', but got '%v'", expDetail, aerr.Details[0])
		}
	}
}
