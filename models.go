package rhea

import (
	"compress/gzip"
	"encoding/xml"
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

	- Chebi
	- SmallMolecule
	- GenericCompound
	- Polymer
	- ReactionSide
	- Reaction

Relationally, the entire structure of Rhea can simply be thought of as:

	- Human readable REACTIONS exist, associated with ec numbers (and Uniprot identifiers)
	- REACTIONS have two SIDES of their equation: a left side, and a right side.
	- Each SIDE of an equation contains a variable number of POLYMERS, GENERIC COMPOUNDS, OR SMALL MOLECULES
	- Each POLYMER, GENERIC COMPOUND, OR SMALL MOLECULE is associated with a CHEBI
	- A CHEBI represents an underlying chemical

******************************************************************************/

type Rhea struct {
	ReactiveParts []ReactivePart
	Compounds     []Compound
	ReactionSides []ReactionSide
	Reactions     []Reaction
}

type ReactivePart struct {
	// Small Molecule portions
	Id                  int
	Accession           string
	Position            string
	Name                string
	HtmlName            string
	Formula             string
	Charge              string
	Chebi               string
	SubclassOfChebi     string
	PolymerizationIndex string
}

type Compound struct {
	Id           int
	Accession    string
	Name         string
	HtmlName     string
	CompoundType string // SmallMolecule, Polymer, GenericPolypeptide, GenericPolynucleotide, GenericHeteropolysaccharide
	ReactivePart string
}

type ReactionSide struct {
	Accession string
	Contains  int
	ContainsN bool
	Minus     bool // Only set to true if ContainsN == true to handle Nminus1
	Plus      bool // Only set to true if ContainsN == true to handle Nplus1
	Compound  string
}

type Reaction struct {
	Id                   int
	Directional          bool
	Accession            string
	Status               string
	Comment              string
	Equation             string
	HtmlEquation         string
	IsChemicallyBalanced bool
	IsTransport          bool
	Ec                   string
	Citations            []string
	Substrates           []string
	Products             []string
	SubstrateOrProducts  []string
	Location             string
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

	for _, description := range rdf.Descriptions {
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
					Chebi:     description.Chebi.Resource}
				if compoundType == "Polymer" {
					newReactivePart.Chebi = description.UnderlyingChebi.Resource
				}
				// Add subclass Chebi
				for _, sc := range description.Subclass {
					if strings.Contains(sc.Resource, "CHEBI") {
						newReactivePart.SubclassOfChebi = sc.Resource
					}
				}
				newCompound := Compound{
					Id:           description.Id,
					Accession:    description.Accession,
					Name:         description.Name,
					HtmlName:     description.HtmlName,
					CompoundType: compoundType,
					ReactivePart: description.Accession}

				// Add new reactive parts and new compounds to rhea
				rhea.ReactiveParts = append(rhea.ReactiveParts, newReactivePart)
				rhea.Compounds = append(rhea.Compounds, newCompound)
			case "http://rdf.rhea-db.org/GenericPolypeptide", "http://rdf.rhea-db.org/GenericPolynucleotide", "http://rdf.rhea-db.org/GenericHeteropolysaccharide":
				compoundType := subclass.Resource[23:]
				newCompound := Compound{
					Id:           description.Id,
					Accession:    description.Accession,
					Name:         description.Name,
					HtmlName:     description.HtmlName,
					CompoundType: compoundType,
					ReactivePart: description.ReactivePartXml.Resource}
				rhea.Compounds = append(rhea.Compounds, newCompound)
			case "http://rdf.rhea-db.org/ReactivePart":
				newReactivePart := ReactivePart{
					Id:        description.Id,
					Accession: description.Accession,
					Position:  description.Position,
					Name:      description.Name,
					HtmlName:  description.HtmlName,
					Formula:   description.Formula,
					Charge:    description.Charge,
					Chebi:     description.Chebi.Resource}
				rhea.ReactiveParts = append(rhea.ReactiveParts, newReactivePart)
			}
		}
		for _, containsx := range description.ContainsX {
			if strings.Contains(containsx.XMLName.Local, "contains") {
				// Get reaction sides
				// gzip -d -k -c rhea.rdf.gz | grep -o -P '(?<=contains).*(?= rdf)' | tr ' ' '\n' | sort -u | tr '\n' ' '
				// The exceptions to numeric contains are 2n, N, Nminus1, and Nplus1
				var newReactionSide ReactionSide
				switch containsx.XMLName.Local {
				case "containsN":
					newReactionSide = ReactionSide{
						Accession: description.About,
						Contains:  1,
						ContainsN: true,
						Minus:     false,
						Plus:      false,
						Compound:  containsx.Content}
				case "contains2n":
					newReactionSide = ReactionSide{
						Accession: description.About,
						Contains:  2,
						ContainsN: true,
						Minus:     false,
						Plus:      false,
						Compound:  containsx.Content}
				case "containsNminus1":
					newReactionSide = ReactionSide{
						Accession: description.About,
						Contains:  1,
						ContainsN: true,
						Minus:     true,
						Plus:      false,
						Compound:  containsx.Content}
				case "containsNplus1":
					newReactionSide = ReactionSide{
						Accession: description.About,
						Contains:  1,
						ContainsN: true,
						Minus:     false,
						Plus:      true,
						Compound:  containsx.Content}
				default:
					i, err := strconv.Atoi(containsx.XMLName.Local[8:])
					if err != nil {
						return Rhea{}, err
					}
					newReactionSide = ReactionSide{
						Accession: description.About,
						Contains:  i,
						ContainsN: false,
						Minus:     false,
						Plus:      false,
						Compound:  containsx.Content}
				}
				rhea.ReactionSides = append(rhea.ReactionSides, newReactionSide)
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
