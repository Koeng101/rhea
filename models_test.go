package rhea

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	xmlFile, err := os.Open("data/rhea.rdf.gz")
	if err != nil {
		t.Errorf("%w", err)
	}
	r, err := gzip.NewReader(xmlFile)
	if err != nil {
		t.Errorf("%w", err)
	}
	rheaBytes, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("%w", err)
	}

	var rdf RheaRdf
	err = xml.Unmarshal(rheaBytes, &rdf)
	if err != nil {
		t.Errorf("%w", err)
	}
	fmt.Println(rdf.XMLName)
	for _, description := range rdf.Descriptions {
		if len(description.ContainsX) != 0 {
			fmt.Println(description.ContainsX)
		}

	}
}
