package rhea

import (
	"compress/gzip"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

/******************************************************************************

Lower level structs

These structs operate at the lowest level to parse the RheaRdf database dump
into structs that Golang can understand. Getting all of Rhea from Rdf->Golang
is quite verbose, and so most of the time you should not use these structs unless
you know what you are doing and it cannot be accomplished with higher level structs.

******************************************************************************/

type RheaRdf struct {
	XMLName      xml.Name      `xml:"RDF"`
	Descriptions []Description `xml:"Description"`
}

type Description struct {
	// Reaction
	XMLName              xml.Name             `xml:"Description"`
	About                string               `xml:"about,attr"`
	Id                   int                  `xml:"id"`
	Accession            string               `xml:"accession"`
	Equation             string               `xml:"equation"`
	HtmlEquation         string               `xml:"htmlEquation"`
	IsChemicallyBalanced bool                 `xml:"isChemicallyBalanced"`
	IsTransport          bool                 `xml:"isTransport"`
	Citations            []Citation           `xml:"citation"`
	Substrates           []Substrate          `xml:"substrates"`
	Products             []Product            `xml:"products"`
	SubstrateOrProducts  []SubstrateOrProduct `xml:"substratesOrProducts"`
	Subclass             []Subclass           `xml:"subClassOf"`
	Comment              string               `xml:"comment"`
	EC                   EC                   `xml:"ec"`
	Status               Status               `xml:"status"`
	Compound             CompoundXml          `xml:"compound"`

	// ReactionSide / Reaction Participant
	BidirectionalReactions []BidirectionalReaction `xml:"bidirectionalReaction"`
	DirectionalReactions   []DirectionalReaction   `xml:"directionalReaction"`
	Side                   Side                    `xml:"side"`
	SeeAlsos               SeeAlso                 `xml:"seeAlso"`
	TransformableTo        TransformableTo         `xml:"transformableTo"`
	CuratedOrder           int                     `xml:"curatedOrder"`
	Contains               Contains                `xml:"contains"`

	// ContainsX contains all other name-attribute pairs, with names like "contains1" in mind
	ContainsX []ContainsX `xml:",any"`

	// Small Molecule tags
	Name     string   `xml:"name"`
	HtmlName string   `xml:"htmlName"`
	Formula  string   `xml:"formula"`
	Charge   string   `xml:"charge"`
	Chebi    ChebiXml `xml:"chebi"`

	// Generic Compound
	ReactivePartXml ReactivePartXml `xml:"reactivePart"`

	// ReactivePart
	Position string `xml:"position"`

	// Polymer
	UnderlyingChebi     UnderlyingChebi `xml:"underlyingChebi"`
	PolymerizationIndex string          `xml:"polymerizationIndex"`

	// Transport
	Location Location `xml:"location"`
}

func (d *Description) CitationStrings() []string {
	var output []string
	for _, x := range d.Citations {
		output = append(output, x.Resource)
	}
	return output
}

func (d *Description) SubstrateStrings() []string {
	var output []string
	for _, x := range d.Substrates {
		output = append(output, x.Resource)
	}
	return output
}

func (d *Description) ProductStrings() []string {
	var output []string
	for _, x := range d.Products {
		output = append(output, x.Resource)
	}
	return output
}

func (d *Description) SubstrateOrProductStrings() []string {
	var output []string
	for _, x := range d.SubstrateOrProducts {
		output = append(output, x.Resource)
	}
	return output
}

type Citation struct {
	XMLName  xml.Name `xml:"citation"`
	Resource string   `xml:"resource,attr"`
}

type Substrate struct {
	XMLName  xml.Name `xml:"substrates"`
	Resource string   `xml:"resource,attr"`
}

type Product struct {
	XMLName  xml.Name `xml:"products"`
	Resource string   `xml:"resource,attr"`
}

type SubstrateOrProduct struct {
	XMLName  xml.Name `xml:"substratesOrProducts"`
	Resource string   `xml:"resource,attr"`
}

type Subclass struct {
	XMLName  xml.Name `xml:"subClassOf"`
	Resource string   `xml:"resource,attr"`
}

type EC struct {
	XMLName  xml.Name `xml:"ec"`
	Resource string   `xml:"resource,attr"`
}

type Status struct {
	XMLName  xml.Name `xml:"status"`
	Resource string   `xml:"resource,attr"`
}

type CompoundXml struct {
	XMLName  xml.Name `xml:"compound"`
	Resource string   `xml:"resource,attr"`
}

type BidirectionalReaction struct {
	XMLName  xml.Name `xml:"bidirectionalReaction"`
	Resource string   `xml:"resource,attr"`
}

type DirectionalReaction struct {
	XMLName  xml.Name `xml:"directionalReaction"`
	Resource string   `xml:"resource,attr"`
}

type Side struct {
	XMLName  xml.Name `xml:"side"`
	Resource string   `xml:"resource,attr"`
}

type SeeAlso struct {
	XMLName  xml.Name `xml:"seeAlso"`
	Resource string   `xml:"resource,attr"`
}

type TransformableTo struct {
	XMLName  xml.Name `xml:"transformableTo"`
	Resource string   `xml:"resource,attr"`
}

type Contains struct {
	XMLName  xml.Name `xml:"contains"`
	Resource string   `xml:"resource,attr"`
}

type ContainsX struct {
	XMLName xml.Name
	Content string `xml:"resource,attr"`
}

type ChebiXml struct {
	XMLName  xml.Name `xml:"chebi"`
	Resource string   `xml:"resource,attr"`
}

type UnderlyingChebi struct {
	XMLName  xml.Name `xml:"underlyingChebi"`
	Resource string   `xml:"resource,attr"`
}

type ReactivePartXml struct {
	XMLName  xml.Name `xml:"reactivePart"`
	Resource string   `xml:"resource,attr"`
}

type Location struct {
	XMLName  xml.Name `xml:"location"`
	Resource string   `xml:"resource,attr"`
}

/******************************************************************************

Higher level structs

These structs are what you would put into a database or directly use. In order to
create a tree or insert into a normalized database, you would insert in the following
order:

	- ReactivePart
	- ReactiveSide (derived from ReactionParticipant/Reaction)
	- ReactionParticipant
	- Reaction

The entire structure of Rhea can simply be thought of as:

	- There are Reactions. Those Reactions can have substrates and products, or substratesOrProducts
	  in the case that the reaction is bidirectional.
	- There are ReactionSides. ReactionSides can be thought of as a many-to-many table between Reactions
	  and ReactionParticipants. It only serves as an inbetween, saying "this Reaction has these
	  ReactionParticipants on this side"
	- There are ReactionParticipants. ReactionParticipants link ReactionSides with ReactiveParts and include
	  useful information like the number of ReactiveParts (or chemicals) needed to do a certain Reaction.
	- There are ReactiveParts. These are physical molecules represented by Chebis. "CompoundReactionParticipantLink"
	  links to the "Compound" of a ReactionParticipant

******************************************************************************/

type Rhea struct {
	ReactionParticipants []ReactionParticipant `json:"reactionParticipants"`
	ReactiveParts        []ReactivePart        `json:"reactiveParts"`
	Reactions            []Reaction            `json:"reactions"`
}

type ReactivePart struct {
	Id                              int    `json:"id" db:"id"`
	Accession                       string `json:"accession" db:"accession"`
	Position                        string `json:"position" db: "position"`
	Name                            string `json:"name" db:"name"`
	HtmlName                        string `json:"htmlName" db:"htmlname"`
	Formula                         string `json:"formula" db:"formula"`
	Charge                          string `json:"charge" db:"charge"`
	Chebi                           string `json:"chebi" db:"chebi"`
	SubclassOfChebi                 string `json:"subclassOfChebi"`
	PolymerizationIndex             string `json:"polymerizationIndex" db:"polymerizationindex"`
	CompoundReactionParticipantLink string `json:"reactionparticipantlink" db:"reactionparticipantlink"`
	CompoundId                      int    `json:"id" db:"compoundid"`
	CompoundAccession               string `json:"accession" db:"compoundaccession"`
	CompoundName                    string `json:"name" db:"compoundname"`
	CompoundHtmlName                string `json:"htmlName" db:"compoundhtmlname"`
	CompoundType                    string `json:"compoundType" db:"compoundtype"`
}

type ReactionParticipant struct {
	ReactionSide string `json:"reactionside" db:"reactionside"`
	Contains     int    `json:"contains" db:"contains"`
	ContainsN    bool   `json:"containsn" db:"containsn"`
	Minus        bool   `json:"minus" db:"minus"` // Only set to true if ContainsN == true to handle Nminus1
	Plus         bool   `json:"plus" db:"plus"`   // Only set to true if ContainsN == true to handle Nplus1
	Compound     string `json:"compound" db:"compound"`
}

type Reaction struct {
	Id                   int    `json:"id" db:"id"`
	Directional          bool   `json:"directional" db:"directional"`
	Accession            string `json:"accession" db:"accession"`
	Status               string `json:"status" db:"status"`
	Comment              string `json:"comment" db:"comment"`
	Equation             string `json:"equation" db:"equation"`
	HtmlEquation         string `json:"htmlequation" db:"htmlequation"`
	IsChemicallyBalanced bool   `json:"ischemicallybalanced" db:"ischemicallybalanced"`
	IsTransport          bool   `json:"istransport" db:"istransport"`
	Ec                   string `json:"ec" db:"ec"`
	Location             string `json:"location" db:"location"`
	Citations            []string
	Substrates           []string
	Products             []string
	SubstrateOrProducts  []string
}

/******************************************************************************

Parse functions

These functions take in the rhea.rdf.gz dump file and return a Rhea struct,
which contains all of the higher level structs

******************************************************************************/

func ParseRhea(rheaBytes []byte) (Rhea, error) {
	var err error
	// Read rheaBytes into a RheaRdf object
	var rdf RheaRdf
	err = xml.Unmarshal(rheaBytes, &rdf)
	if err != nil {
		return Rhea{}, err
	}

	// Initialize Rhea
	var rhea Rhea
	compoundMap := make(map[string]string)
	reactivePartMap := make(map[string]ReactivePart)

	for _, description := range rdf.Descriptions {
		// Handle the case of a single compound -> reactive part, such as
		// <rdf:Description rdf:about="http://rdf.rhea-db.org/Compound_10594">
		// 	<rh:reactivePart rdf:resource="http://rdf.rhea-db.org/Compound_10594_rp2"/>
		// </rdf:Description>
		if (len(description.Subclass) == 0) && (description.ReactivePartXml.Resource != "") {
			compoundMap[description.ReactivePartXml.Resource] = description.About
		}
		if description.Compound.Resource != "" {
			compoundMap[description.About] = description.Compound.Resource
		}

		for _, subclass := range description.Subclass {
			switch subclass.Resource {
			case "http://rdf.rhea-db.org/DirectionalReaction":
				newReaction := Reaction{
					Id:                   description.Id,
					Directional:          true,
					Accession:            description.Accession,
					Status:               description.Status.Resource,
					Comment:              description.Comment,
					Equation:             description.Equation,
					HtmlEquation:         description.HtmlEquation,
					IsChemicallyBalanced: description.IsChemicallyBalanced,
					IsTransport:          description.IsTransport,
					Ec:                   description.EC.Resource,
					Citations:            description.CitationStrings(),
					Substrates:           description.SubstrateStrings(),
					Products:             description.ProductStrings(),
					SubstrateOrProducts:  description.SubstrateOrProductStrings(),
					Location:             description.Location.Resource}
				rhea.Reactions = append(rhea.Reactions, newReaction)
			case "http://rdf.rhea-db.org/BidirectionalReaction":
				newReaction := Reaction{
					Id:                   description.Id,
					Directional:          false,
					Accession:            description.Accession,
					Status:               description.Status.Resource,
					Comment:              description.Comment,
					Equation:             description.Equation,
					HtmlEquation:         description.HtmlEquation,
					IsChemicallyBalanced: description.IsChemicallyBalanced,
					IsTransport:          description.IsTransport,
					Ec:                   description.EC.Resource,
					Citations:            description.CitationStrings(),
					Substrates:           description.SubstrateStrings(),
					Products:             description.ProductStrings(),
					SubstrateOrProducts:  description.SubstrateOrProductStrings(),
					Location:             description.Location.Resource}
				rhea.Reactions = append(rhea.Reactions, newReaction)
			case "http://rdf.rhea-db.org/SmallMolecule", "http://rdf.rhea-db.org/Polymer":
				compoundType := subclass.Resource[23:]
				newReactivePart := ReactivePart{
					Id:        description.Id,
					Accession: description.Accession,
					Position:  description.Position,
					Name:      description.Name,
					HtmlName:  description.HtmlName,
					Formula:   description.Formula,
					Charge:    description.Charge,
					Chebi:     description.Chebi.Resource,

					CompoundReactionParticipantLink: description.About,
					CompoundId:                      description.Id,
					CompoundAccession:               description.Accession,
					CompoundName:                    description.Name,
					CompoundHtmlName:                description.HtmlName,
					CompoundType:                    compoundType}
				if compoundType == "Polymer" {
					newReactivePart.Chebi = description.UnderlyingChebi.Resource
				}
				// Add subclass Chebi
				for _, sc := range description.Subclass {
					if strings.Contains(sc.Resource, "CHEBI") {
						newReactivePart.SubclassOfChebi = sc.Resource
					}
				}
				// Add new reactive parts and new compounds to rhea
				rhea.ReactiveParts = append(rhea.ReactiveParts, newReactivePart)
			case "http://rdf.rhea-db.org/GenericPolypeptide", "http://rdf.rhea-db.org/GenericPolynucleotide", "http://rdf.rhea-db.org/GenericHeteropolysaccharide":
				compoundType := subclass.Resource[23:]
				newReactivePart := ReactivePart{
					CompoundId:        description.Id,
					CompoundAccession: description.Accession,
					CompoundName:      description.Name,
					CompoundHtmlName:  description.HtmlName,
					CompoundType:      compoundType}
				reactivePartMap[description.About] = newReactivePart
				compoundMap[description.ReactivePartXml.Resource] = description.About
			}
		}
	}

	// Go back and get the ReactiveParts
	for _, description := range rdf.Descriptions {
		for _, containsx := range description.ContainsX {
			if strings.Contains(containsx.XMLName.Local, "contains") {
				// Get reaction sides
				// gzip -d -k -c rhea.rdf.gz | grep -o -P '(?<=contains).*(?= rdf)' | tr ' ' '\n' | sort -u | tr '\n' ' '
				// The exceptions to numeric contains are 2n, N, Nminus1, and Nplus1
				var newReactionParticipant ReactionParticipant
				switch containsx.XMLName.Local {
				case "containsN":
					newReactionParticipant = ReactionParticipant{
						ReactionSide: description.About,
						Contains:     1,
						ContainsN:    true,
						Minus:        false,
						Plus:         false,
						Compound:     compoundMap[containsx.Content]}
				case "contains2n":
					newReactionParticipant = ReactionParticipant{
						ReactionSide: description.About,
						Contains:     2,
						ContainsN:    true,
						Minus:        false,
						Plus:         false,
						Compound:     compoundMap[containsx.Content]}
				case "containsNminus1":
					newReactionParticipant = ReactionParticipant{
						ReactionSide: description.About,
						Contains:     1,
						ContainsN:    true,
						Minus:        true,
						Plus:         false,
						Compound:     compoundMap[containsx.Content]}
				case "containsNplus1":
					newReactionParticipant = ReactionParticipant{
						ReactionSide: description.About,
						Contains:     1,
						ContainsN:    true,
						Minus:        false,
						Plus:         true,
						Compound:     compoundMap[containsx.Content]}
				default:
					i, err := strconv.Atoi(containsx.XMLName.Local[8:])
					if err != nil {
						return Rhea{}, err
					}
					newReactionParticipant = ReactionParticipant{
						ReactionSide: description.About,
						Contains:     i,
						ContainsN:    false,
						Minus:        false,
						Plus:         false,
						Compound:     compoundMap[containsx.Content]}
				}
				rhea.ReactionParticipants = append(rhea.ReactionParticipants, newReactionParticipant)
			}
		}

		for _, subclass := range description.Subclass {
			switch subclass.Resource {
			case "http://rdf.rhea-db.org/ReactivePart":
				newReactivePart, ok := reactivePartMap[compoundMap[description.About]]
				if ok != true {
					return Rhea{}, errors.New("Could not find " + description.About)
				}
				newReactivePart.CompoundReactionParticipantLink = description.About
				newReactivePart.Id = description.Id
				newReactivePart.Accession = description.Accession
				newReactivePart.Position = description.Position
				newReactivePart.Name = description.Name
				newReactivePart.HtmlName = description.HtmlName
				newReactivePart.Formula = description.Formula
				newReactivePart.Charge = description.Charge
				newReactivePart.Chebi = description.Chebi.Resource
				rhea.ReactiveParts = append(rhea.ReactiveParts, newReactivePart)
			}
		}
	}
	return rhea, nil
}

func ReadRhea(gzipPath string) ([]byte, error) {
	// Get gz'd file bytes
	xmlFile, err := os.Open("data/rhea.rdf.gz")
	if err != nil {
		return []byte{}, err
	}

	// Decompress gz'd file
	r, err := gzip.NewReader(xmlFile)
	if err != nil {
		return []byte{}, err
	}

	// Read decompressed gz'd file
	rheaBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}
	return rheaBytes, nil
}
