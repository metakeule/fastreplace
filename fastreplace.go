package fastreplace

import (
	"sort"
)

/*
	Usage:

	see README.md
*/

type FReplace struct {
	original  []byte
	positions map[int]string // maps a replacement string/sequence of bytes to an array of the indices of the original before which it appears
	sortedPos []int
}

func New(delimiter string, input []byte) (ø *FReplace) {
	ø = &FReplace{}
	ø.Parse(delimiter, input)
	return
}

func NewString(delimiter string, input string) (ø *FReplace) {
	return New(delimiter, []byte(input))
}

// returns a map with all positions and their replacements for the placeholders given in the input map
func (ø *FReplace) AllPos(in map[string][]byte) (out map[int][]byte) {
	out = map[int][]byte{}
	for k, v := range in {
		pos := ø.Pos(k)
		for _, p := range pos {
			out[p] = v
		}
	}
	return
}

// returns all positions for the placeholder string
func (ø *FReplace) Pos(str string) (p []int) {
	for pos, s := range ø.positions {
		if s == str {
			p = append(p, pos)
		}
	}
	return
}

// replace based on positions
func (ø *FReplace) ReplacePos(m map[int][]byte) (res []byte) {
	last := 0
	for _, pos := range ø.sortedPos {
		res = append(res, ø.original[last:pos]...)
		if repl := m[pos]; repl != nil {
			// we have a replacement for the anchor
			res = append(res, repl...)
		}
		last = pos
	}
	res = append(res, ø.original[last:len(ø.original)]...)
	return
}

// like Replace but returns a string
func (ø *FReplace) ReplaceString(m map[string][]byte) (res string) {
	return string(ø.Replace(m))
}

// replace based on placeholders
func (ø *FReplace) Replace(m map[string][]byte) (res []byte) {
	last := 0
	for _, pos := range ø.sortedPos {
		res = append(res, ø.original[last:pos]...)
		posBt := ø.positions[pos]
		if repl := m[posBt]; repl != nil {
			// we have a replacement for the anchor
			res = append(res, repl...)
		}
		last = pos
	}
	res = append(res, ø.original[last:len(ø.original)]...)
	return
}

func (ø *FReplace) ParseString(delimiter string, s string) {
	ø.Parse(delimiter, []byte(s))
}

// parse the input for placeholders and caches the result
func (ø *FReplace) Parse(delimiter string, in []byte) {
	ø.positions = map[int]string{}
	ø.original = []byte{}
	ø.sortedPos = []int{}
	inAnchor := false
	startPos := 0
	var anchor []byte
	hs := []byte(delimiter)
	lenDel := len(hs)
	lenIn := len(in)
	h := hs[0]
	i := 0
	for ii := 0; ii < lenIn; ii++ {
		b := in[ii]
		if b == h && lenIn-ii >= lenDel && string(in[ii:ii+lenDel]) == delimiter {
			for jj := 1; jj < lenDel; jj++ {
				// fast forward to the end of the delimiter
				ii++
			}
			if inAnchor {
				inAnchor = false
				ø.positions[startPos] = string(anchor) // register the cached anchor
				continue
			} else {
				inAnchor = true
				startPos = i
				anchor = []byte{}
				continue
			}
		}
		if !inAnchor {
			ø.original = append(ø.original, b)
			i++
		} else {
			anchor = append(anchor, b)
		}
	}

	sorted := []int{}

	for pos, _ := range ø.positions {
		sorted = append(sorted, pos)
	}
	sort.Ints(sorted)
	ø.sortedPos = sorted
}

// returns an Instance that offers more comfort and caching of replacements
func (ø *FReplace) Instance() *Instance {
	return &Instance{replace: ø, replacePos: map[int][]byte{}}
}

type Instance struct {
	replace    *FReplace
	replacePos map[int][]byte
}

func (ø *Instance) String() string {
	return string(ø.Replace())
}

func (ø *Instance) Replace() []byte {
	return ø.replace.ReplacePos(ø.replacePos)
}

func (ø *Instance) Assign(key string, val []byte) {
	poses := ø.replace.Pos(key)
	for _, pos := range poses {
		ø.replacePos[pos] = val
	}
}

func (ø *Instance) AssignString(key string, val string) {
	ø.Assign(key, []byte(val))
}
