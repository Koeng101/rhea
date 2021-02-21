package rhea

import (
	"fmt"
	"strings"
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
	//	if k == "http://rdf.rhea-db.org/Participant_10008_compound_10594" {
	//		fmt.Println(k)
	//		fmt.Println(v)
	//	}
	//}
	for _, c := range rhea.Compounds {
		if strings.Contains(c.Accession, "http://rdf.rhea-db.org/Participant_10000_compound_1283") {
			fmt.Println(c)
		}
	}
}
