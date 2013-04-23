package fastreplace

import (
	"bytes"
	"fmt"
	"sort"
)

/*
	Usage:

	see README.md
*/

type Replacer interface {
	String() string
	Replace() []byte
	Assign(key string, val []byte)
	AssignString(key string, val string)
}

type FReplace struct {
	original  []byte
	positions map[int]string // maps a replacement string/sequence of bytes to an array of the indices of the original before which it appears
	sortedPos []int
}

func New(delimiter []byte, input []byte) (ø *FReplace, ſ error) {
	ø = &FReplace{}
	ſ = ø.Parse(delimiter, input)
	return
}

func NewString(delimiter string, input string) (ø *FReplace, ſ error) {
	return New([]byte(delimiter), []byte(input))
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
func (ø *FReplace) Pos(placeholder string) (p []int) {
	for pos, s := range ø.positions {
		if placeholder == s {
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
	ø.Parse([]byte(delimiter), []byte(s))
}

// parse the input for placeholders and caches the result
func (ø *FReplace) Parse(delimiter []byte, in []byte) error {
	ø.positions = map[int]string{}
	ø.original = []byte{}
	ø.sortedPos = []int{}
	lenDel := len(delimiter)
	lenIn := len(in)
	for i := 0; i < lenIn; i++ {
		found := bytes.Index(in[i:], delimiter)
		if found != -1 {
			if found == 0 && i != 0 {
				return fmt.Errorf("Syntax error: can't have 2 or more placeholders side by side: %#v\n", string(in[:i+lenDel]))
			}
			start := found + i
			ø.original = append(ø.original, in[i:start]...)
			startPlaceH := start + lenDel
			found = bytes.Index(in[startPlaceH:], delimiter)
			if found == -1 {
				// is not a delimiter
				ø.original = append(ø.original, in[startPlaceH:]...)
				break
			} else {
				end := found + start + lenDel
				pos := len(ø.original)
				ø.sortedPos = append(ø.sortedPos, pos)
				ø.positions[pos] = string(in[startPlaceH:end])
				i = end + 1
			}
		} else {
			ø.original = append(ø.original, in[i:]...)
			break
		}
	}
	sort.Ints(ø.sortedPos)
	return nil
}

// returns an Instance that offers more comfort and caching of replacements
func (ø *FReplace) Instance() Replacer {
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

func (ø *Instance) Append(key string, val []byte) {
	poses := ø.replace.Pos(key)
	for _, pos := range poses {
		ø.replacePos[pos] = append(ø.replacePos[pos], val...)
	}
}

func (ø *Instance) AppendString(key string, val string) {
	ø.Append(key, []byte(val))
}

func (ø *Instance) AssignString(key string, val string) {
	ø.Assign(key, []byte(val))
}
