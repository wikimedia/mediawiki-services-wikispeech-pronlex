package decompounder

import (
	"fmt"
	"unicode/utf8"
)

// TODO:
// remove
// compute freqs
// completion?
// bestGuess (heuristics)
// Handling 's' and other infixes between compound parts: 'handel+s+'trädgård'
// Handling tripple consonats clusters merged to two consonants in compounds ('natt+tåg' -> 'nattåg')

// TNode is kind of a trie-structure representing strings (words).
// A path trough the TNode that ends with leaf = true represents a
// string/word.
type TNode struct {
	// the current character in a string
	r rune
	// child nodes of this node
	sons map[rune]*TNode
	// true if this character ends an input string
	leaf bool
}

func NewTNode() *TNode {
	return &TNode{sons: make(map[rune]*TNode)}
}

// add inserts a strings into the TNode and builds up sub-nodes as it
// goes
func (t *TNode) add(s string) *TNode {

	if s == "" {
		return t
	}

	// Pick off first rune in string
	r, l := utf8.DecodeRuneInString(s)

	// This path so far already exists.
	// Recursively keep adding
	if son, ok := t.sons[r]; ok {
		if len(s) == 1 {
			son.leaf = true
			// This is where you could increment a frequency counter.
			// You'd want to add a frequency field to bot TNode and arc.
		}
		son.add(s[l:len(s)])

	} else { // new path
		son := NewTNode()
		son.r = r
		if len(s) == 1 {
			son.leaf = true
		}
		t.sons[r] = son
		son.add(s[l:len(s)])
	}

	return t
}

// arc represents a substring of a string, with a start and end index
// of the string.
type arc struct {
	start int
	end   int
}

// Returns the matching prefix substrings of s that exist in t in the
// form of arcs.
func (t *TNode) prefixes(s string) []arc {
	var res []arc

	sons := t.sons
	for i, r := range s {
		// path in tree
		if v, ok := sons[r]; ok {
			sons = v.sons
			if v.leaf {
				res = append(res, arc{end: i + 1})
			}
		} else { // not a path in tree
			return res
		}
	}

	return res
}

type PrefixTree struct {
	tree *TNode
}

func NewPrefixTree() PrefixTree {
	return PrefixTree{tree: NewTNode()}
}

func (t PrefixTree) Add(s string) {
	t.tree.add(s)
}

func (t PrefixTree) Prefixes(s string) []arc {
	return t.tree.prefixes(s)
}

type SuffixTree struct {
	tree *TNode
}

func NewSuffixTree() SuffixTree {
	return SuffixTree{tree: NewTNode()}
}

// Reverse returns its argument string reversed rune-wise left to right.
// Lifted from https://github.com/golang/example/blob/master/stringutil/reverse.go
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func (t SuffixTree) Add(s string) {
	r := reverse(s)
	t.tree.add(r)
}

func (t SuffixTree) Suffixes(s string) []arc {
	r := reverse(s)

	// the arcs going from right to left
	// (strings inverted)
	suffArcs := t.tree.prefixes(r)

	// invert arcs to go from left to right
	l := len(s)
	var res []arc
	for _, a := range suffArcs {
		res = append(res, arc{start: l - a.end, end: l - a.start})
	}

	return res
}

type Decompounder struct {
	prefixes PrefixTree
	suffixes SuffixTree
}

func NewDecompounder() Decompounder {
	return Decompounder{prefixes: NewPrefixTree(), suffixes: NewSuffixTree()}
}

func (d Decompounder) arcs(s string) []arc {
	var res []arc

	res = append(res, d.prefixes.Prefixes(s)...)
	res = append(res, d.suffixes.Suffixes(s)...)

	return res
}

func (d Decompounder) Decomp(s string) [][]string {
	var res [][]string

	arcs := d.arcs(s)
	paths := paths(arcs, 0, len(s))

	for _, p := range paths {
		res = append(res, pathToDecomp(p, s))
	}

	return res
}

func punk() {
	fmt.Println()
}

func paths(as []arc, from, to int) [][]arc {

	arcMap := make(map[int][]arc)
	for _, a := range as {
		v, _ := arcMap[a.start]
		arcMap[a.start] = append(v, a)
	}

	var path []arc
	var res [][]arc
	pathsAccu(arcMap, from, to, path, &res)
	return res
}

func pathsAccu(as map[int][]arc, from, to int, currPath []arc, paths *[][]arc) {

	arcs, ok := as[from]
	if !ok { // no path from 'from'
		return
	}

	for _, arc := range arcs {
		// Yeah! Complete path!
		if arc.end == to {
			path := currPath
			path = append(path, arc)
			*paths = append(*paths, path)
		}
		// Keep threading down the path
		newPath := currPath
		newPath = append(newPath, arc)
		pathsAccu(as, arc.end, to, newPath, paths)

	}
}

func pathToDecomp(p []arc, s string) []string {
	var res []string
	// TODO error checking
	for _, a := range p {
		s0 := s[a.start:a.end]
		res = append(res, s0)
	}
	return res
}
