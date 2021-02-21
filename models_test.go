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
	//r := rhea.ReactionParticipants[100].Compound
	for _, a := range rhea.ReactiveParts {
		if a.CompoundReactionParticipantLink == "http://rdf.rhea-db.org/Compound_10594" {
			fmt.Println(a)
		}
	}
	//for _, b := range rhea.ReactionParticipants[100:110] {
	//	fmt.Println(b.Compound)
	//}

	for _, c := range rhea.ReactiveParts[100:105] {
		fmt.Println(c.CompoundReactionParticipantLink)
	}
}
