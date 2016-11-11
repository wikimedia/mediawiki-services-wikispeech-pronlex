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
type tNode struct {
	// the current character in a string
	r rune
	// child nodes of this node
	sons map[rune]*tNode
	// true if this character ends an input string
	leaf bool
}

func NewtNode() *tNode {
	return &tNode{sons: make(map[rune]*tNode)}
}

// add inserts a strings into the tNode and builds up sub-nodes as it
// goes
func (t *tNode) add(s string) *tNode {

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
			// You'd want to add a frequency field to bot tNode and arc.
		}
		son.add(s[l:len(s)])

	} else { // new path
		son := NewtNode()
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
func (t *tNode) prefixes(s string) []arc {
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
	tree *tNode
}

func NewPrefixTree() PrefixTree {
	return PrefixTree{tree: NewtNode()}
}

func (t PrefixTree) Add(s string) {
	t.tree.add(s)
}

func (t PrefixTree) Prefixes(s string) []arc {
	return t.tree.prefixes(s)
}

func (t PrefixTree) RecursivePrefixes(s string) []arc {
	var res []arc
	// TODO the following call is probably broken
	t.recursivePrefixes(s, 0, len(s), &res)
	return res
}

// TODO Broke: it overgenerates, it seems.
// Yet to add test case.
func (t PrefixTree) recursivePrefixes(s string, from, to int, as *[]arc) {

	// TODO Where to look for infixes, like compounding 's'?
	// Probably somewhere around here

	newAs := t.Prefixes(s[from:])
	for _, a := range newAs {
		newArc := arc{start: a.start + from, end: a.end + from}
		if a.end < to {
			*as = append(*as, newArc)
			t.recursivePrefixes(s, from+a.end, to, as)
		}
	}
}

type SuffixTree struct {
	tree *tNode
}

func NewSuffixTree() SuffixTree {
	return SuffixTree{tree: NewtNode()}
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
	Prefixes PrefixTree
	Suffixes SuffixTree
}

func NewDecompounder() Decompounder {
	return Decompounder{Prefixes: NewPrefixTree(), Suffixes: NewSuffixTree()}
}

func (d Decompounder) arcs(s string) []arc {
	var res []arc

	res0 := append(res, d.Prefixes.RecursivePrefixes(s)...)
	res1 := append(res, d.Suffixes.Suffixes(s)...)

	// ensure no duplicate arcs
	found := make(map[arc]bool)
	for _, a := range res0 {
		if !found[a] {
			res = append(res, a)
			found[a] = true
		}
	}
	for _, a := range res1 {
		if !found[a] {
			res = append(res, a)
			found[a] = true
		}
	}

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

	// ensure that there are no duplicate arcs, since these will
	// generate multiple identical paths
	found := make(map[arc]bool)
	var uniqueAs []arc
	for _, a := range as {
		if !found[a] {
			uniqueAs = append(uniqueAs, a)
			found[a] = true
		}
	}

	arcMap := make(map[int][]arc)
	for _, a := range uniqueAs {
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
