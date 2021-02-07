package rhea

import (
	"encoding/xml"
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
	Subclass             Subclass             `xml:"subclassOf"`
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
	ReactivePart ReactivePart `xml:"reactivePart"`

	// ReactivePart
	Position string `xml:"position"`

	// Polymer
	UnderlyingChebi     UnderlyingChebi `xml:"underlyingChebi"`
	PolymerizationIndex string          `xml:"polymerizationIndex"`

	// Transport
	Location Location `xml:"location"`
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
	XMLName  xml.Name `xml:"subclassOf"`
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

type ReactivePart struct {
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
}

type Chebi struct {
}

type SmallMolecule struct {
}

type GenericCompound struct {
}

type Polymer struct {
}

type ReactionSide struct {
}

type Reaction struct {
}
