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

	//for k, v := range rhea.ReactionParticipantToCompoundMap {
	//	fmt.Println(k)
	//	fmt.Println(v)
	//}
	for _, c := range rhea.Compounds {
		if c.Accession == "http://rdf.rhea-db.org/Compound_10594" {
			fmt.Println(c)
		}
	}
	//fmt.Println(rhea.ReactionParticipantToCompoundMap["http://rdf.rhea-db.org/Participant_10000_compound_1283"])
}
