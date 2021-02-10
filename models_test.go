package rhea

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	rheaBytes, err := ReadRhea("data/rhea.rdf.gz")
	if err != nil {
		t.Errorf("Failed to ReadRhea. Got error: %w", err)
	}
	rhea, err := ParseRhea(rheaBytes)
	if err != nil {
		t.Errorf("Failed to ParseRhea. Got error: %w", err)
	}
	fmt.Println(rhea.ReactiveParts[10:12])
}
