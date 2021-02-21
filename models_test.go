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
	r := "http://rdf.rhea-db.org/Compound_11132"
	fmt.Println(r)
	for _, a := range rhea.ReactiveParts {
		if a.CompoundReactionParticipantLink == r {
			fmt.Println(a)
		}
	}
}
