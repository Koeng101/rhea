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
	//for _, b := range rhea.ReactionParticipants[100:110] {
	//	fmt.Println(b.Compound)
	//}

	for k, v := range rhea.ReactionParticipantToCompoundMap {
		fmt.Println(k)
		fmt.Println(v)
	}
	for _, c := range rhea.Compounds {
		if c.CompoundAccession == "GENERIC:10594" {
			fmt.Println(c)
		}
	}
}
