package decompounder

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

// TODO:
// remove
// compute freqs
// completion?
// bestGuess (heuristics)
// Handling tripple consonats clusters merged into two consonants in compounds ('natt+tåg' -> 'nattåg')

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
			// You'd want to add a frequency field to both tNode and arc.
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

// contains returns true iff s is a leaf
func (t *tNode) contains(s string) bool {

	if s == "" {
		return false
	}

	res := false
	sons := t.sons
	for _, r := range s {
		v, ok := sons[r]
		if !ok { // s not a path in t
			return false
		}
		res = v.leaf
		sons = v.sons
	}

	return res
}

// remove sets leaf = false, but does not actually remove the
// character sequence of the string from the tree (since it may be a
// substring of another string).
// If the string is not in t, nothing happens.
// Returns true if the string was 'removed' from the tree, false otherwise.
// TODO Purge paths that do not lead to a leaf = true.
func (t *tNode) remove(s string) bool {

	if s == "" {
		return false
	}

	sons := t.sons
	for i, r := range s {
		if v, ok := sons[r]; ok {
			if i == len(s)-1 { // last rune of s
				v.leaf = false
				return true
			}
			// keep following the path
			sons = v.sons

		} else {
			return false // s not a path in t
		}
	}

	return false
}

// list returns all paths ending in a leaf, that is all strings added
// to the top tNode (if not subsequently removed).
func (t *tNode) list() []string {
	var res []string

	for _, s := range t.sons {
		listAccu(s, "", &res)
	}

	return res
}

// listAccu is an accumulator helper function to list() above
func listAccu(t *tNode, soFar string, accu *[]string) {
	if t.leaf {
		*accu = append(*accu, soFar+string(t.r))
	}
	for _, k := range t.sons {
		listAccu(k, soFar+string(t.r), accu)
	}
}

// arc represents a substring of a string, with a start and end index
// of the string.
type arc struct {
	start int
	end   int
}

// Returns the matching prefix substrings of s that exist in t in the
// form of arcs. A prefix must be shorter than the input string.
func (t *tNode) prefixes(s string) []arc {
	var res []arc

	sons := t.sons
	for i, r := range s {
		// path in tree
		if v, ok := sons[r]; ok {
			sons = v.sons
			// '&& i < len(s)-1' ensures that the prefix is shorter than s
			if v.leaf && i < len(s)-1 {
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
	// infixes are gluing parts that may appear once after a prefix
	infixes *tNode
}

func NewPrefixTree() PrefixTree {
	return PrefixTree{tree: NewtNode(), infixes: NewtNode()}
}

func (t PrefixTree) Add(s string) {
	t.tree.add(s)
}
func (t PrefixTree) AddInfix(s string) {
	t.infixes.add(s)
}

func (t PrefixTree) Prefixes(s string) []arc {
	return t.tree.prefixes(s)
}

func (t PrefixTree) Infixes(s string) []arc {
	return t.infixes.prefixes(s)
}

func (t PrefixTree) RecursivePrefixes(s string) []arc {
	var res []arc
	t.recursivePrefixes(s, 0, len(s), &res)
	return res
}

func (t PrefixTree) recursivePrefixes(s string, from, to int, as *[]arc) {

	newAs := t.Prefixes(s[from:])

	for _, a := range newAs {
		newArc := arc{start: a.start + from, end: a.end + from}

		//fmt.Printf("newArc: %#v\n", newArc)
		//fmt.Printf("s: %s %s\n", s, s[newArc.start:newArc.end])

		if a.end < to {
			*as = append(*as, newArc)

			// We have found a prefix above.
			// Go looking for potential infixes, and add these to prefix list
			infixes := t.Infixes(s[newArc.end:])
			for _, in := range infixes {
				infix := arc{start: newArc.end, end: in.end + newArc.end}
				if infix.end < to {
					*as = append(*as, infix)
					// TODO Aouch... nested recursion. Fix this to have only one recursive call below.
					// I guess this might blow things up.
					// 'from' could be a list of arcs instead of a single int?
					t.recursivePrefixes(s, infix.end, to, as)
				}
			}

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

// Suffixes returns the arc for suffixes of s in t. A suffix may not
// span the complete s.
func (t SuffixTree) Suffixes(s string) []arc {
	r := reverse(s)

	// the arcs going from right to left
	// (strings inverted)
	suffArcs := t.tree.prefixes(r)

	// invert arcs to go from left to right
	l := len(r)
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

	res0 := append(res, d.prefixes.RecursivePrefixes(s)...)
	res1 := append(res, d.suffixes.Suffixes(s)...)

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

func (d Decompounder) AddPrefix(s string) {
	d.prefixes.Add(s)
}

func (d Decompounder) AddInfix(s string) {
	d.prefixes.AddInfix(s)
}

func (d Decompounder) AddSuffix(s string) {
	d.suffixes.Add(s)
}

// NewDecompounderFromFile initializes a Decompounder from a text file of the following format:
// # line starting with '#' is ignored
// '' empty lines are ignore
//<PREFIX|INFIX|SUFFIX>:<lower-case string>
func NewDecompounderFromFile(fileName string) (Decompounder, error) {
	res := NewDecompounder()
	fh, err := os.Open(fileName)
	if err != nil {
		return res, err
	}
	defer fh.Close()

	s := bufio.NewScanner(fh)
	for s.Scan() {
		l := strings.TrimSpace(s.Text())
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		// TODO parse string and report mismatching strings
		// 1) print failing line to STDERR
		// 2) count nummber of failing lines and if > 0 return error
	}

	return res, nil
}

// sorting [][]string according to length
type ByLen [][]string

func (b ByLen) Len() int {
	return len(b)
}
func (b ByLen) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByLen) Less(i, j int) bool {
	// TODO? add frequency to arcs, and sort them instead: first by
	// length, second by highest lowest freq.
	return len(b[i]) < len(b[j])
}

func (d Decompounder) Decomp(s string) [][]string {
	var res [][]string

	arcs := d.arcs(s)
	paths := paths(arcs, 0, len(s))

	for _, p := range paths {
		res = append(res, pathToDecomp(p, s))
	}

	sort.Sort(ByLen(res))
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
