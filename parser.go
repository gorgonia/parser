package parser

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"

	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// Parser is a structure that parses an expression and returns the corresponding gorgonia node
type Parser struct {
	dico map[string]*G.Node
	g    *G.ExprGraph
}

// NewParser ...
func NewParser(g *G.ExprGraph) *Parser {
	return &Parser{
		dico: make(map[string]*G.Node),
		g:    g,
	}
}

// Get a pointer to the Node from its representation in the formulae
func (p *Parser) Get(ident string) *G.Node {
	return p.dico[ident]
}

// Set a value to the ident
func (p *Parser) Set(ident string, value *G.Node) {
	p.dico[ident] = value
}

// SetValue of the node identified by ident
func (p *Parser) SetValue(ident string, value G.Value) error {
	if _, ok := p.dico[ident]; !ok {
		return errors.New("Node not found")
	}
	_, err := G.Copy(p.dico[ident].Value(), value)
	return err
}

// Process the data.
// this function reads the data from `r` line by line.
// if replacer is non-nil, there is a call to `Replace` for every line.
// this function recognize assignations and declaration that match the corresponding regexp
// it return the last evaluated node
func (p *Parser) Process(r io.Reader, replacer *strings.Replacer) (*G.Node, error) {
	superScriptToNormal := strings.NewReplacer(
		`¹`, `1`,
		`²`, `2`,
		`³`, `3`,
		`⁴`, `4`,
		`⁵`, `5`,
		`⁶`, `6`,
		`⁷`, `7`,
		`⁸`, `8`,
		`⁹`, `9`,
		`⁰`, `0`,
	)

	scanner := bufio.NewScanner(r)
	l := &exprLex{}
	var node *G.Node
	for scanner.Scan() {
		var s string
		if replacer != nil {
			s = replacer.Replace(scanner.Text())
		} else {
			s = scanner.Text()
		}
		if s == "" {
			continue
		}
		switch {
		case isAssignation(s):
			assignation := regexp.MustCompile(AssignationRegexp)
			elements := assignation.FindAllStringSubmatch(s, -1)
			l.line = []byte(elements[0][2])
			l.dico = p.dico
			l.g = p.g
			gorgoniaParse(l)
			if l.err != nil {
				return nil, l.err
			}
			p.Set(elements[0][1], l.result)
			node = l.result
		case isDeclaration(s):
			declaration := regexp.MustCompile(DeclarationRegexp)
			elements := declaration.FindAllStringSubmatch(s, -1)
			ident := elements[0][1]
			var n *G.Node
			if elements[0][4] != "" {
				// It is a matrix
				h, err := strconv.Atoi(superScriptToNormal.Replace(elements[0][4]))
				if err != nil {
					return nil, errors.New("Bad assignment: " + s)
				}
				w, err := strconv.Atoi(superScriptToNormal.Replace(elements[0][2]))
				if err != nil {
					return nil, errors.New("Bad assignment: " + s)
				}
				n = G.NewMatrix(p.g, tensor.Float32, G.WithName(ident), G.WithShape(w, h))
			} else {
				// It is a vector
				w, err := strconv.Atoi(superScriptToNormal.Replace(elements[0][2]))
				if err != nil {
					return nil, errors.New("Bad assignment: " + s)
				}
				n = G.NewVector(p.g, tensor.Float32, G.WithName(ident), G.WithShape(w))
			}
			p.Set(ident, n)
			node = n
		default:
			l.line = []byte(s)
			l.dico = p.dico
			l.g = p.g
			gorgoniaParse(l)
			if l.err != nil {
				return nil, l.err
			}
			node = l.result
		}
	}
	return node, nil
}

// Parse a string and returns the node
func (p *Parser) Parse(s string) (*G.Node, error) {
	l := &exprLex{}
	l.line = []byte(s)
	l.dico = p.dico
	l.g = p.g
	gorgoniaParse(l)
	if l.err != nil {
		return nil, l.err
	}
	return l.result, nil
}

const (
	// AssignationRegexp that represents an assignation
	AssignationRegexp = `^([a-zA-Zĉ₋₊₀₁₂₃₄₅₆₇₈₉ₜᵢₒₖₓᵣₛₚₖₗₘₙₐ]+)=([^=.]*)$`
	// DeclarationRegexp that represents a declaration
	DeclarationRegexp = `^([a-zA-Zĉ₋₊₀₁₂₃₄₅₆₇₈₉ₜᵢₒₖₓᵣₛₚₖₗₘₙₐ]+)∈ℝ([⁰¹²³⁴⁵⁶⁷⁸⁹]+)(ˣ([⁰¹²³⁴⁵⁶⁷⁸⁹]+))?$`
)

func isAssignation(s string) bool {
	assignation := regexp.MustCompile(AssignationRegexp)
	if assignation.MatchString(s) {
		return true
	}
	return false
}

func isDeclaration(s string) bool {
	declaration := regexp.MustCompile(DeclarationRegexp)
	if declaration.MatchString(s) {
		return true
	}
	return false
}
