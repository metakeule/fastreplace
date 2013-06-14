package fastreplace

import (
	ŧ "fmt"
	"strings"
	"testing"
)

var Template = []byte{}
var Expected = ""
var Map = map[string][]byte{}

func Prepare() {
	Map = map[string][]byte{}
	orig := []string{}
	exp := []string{}
	for i := 0; i < 5; i++ {
		orig = append(orig, ŧ.Sprintf(`a string with @@replacement%v@@`, i))
		exp = append(exp, ŧ.Sprintf("a string with repl%v", i))
		Map[ŧ.Sprintf("replacement%v", i)] = []byte(ŧ.Sprintf("repl%v", i))
	}
	Expected = strings.Join(exp, "")
	Template = []byte(strings.Join(orig, ""))
}

var replace = &FReplace{}

func TestReplaceMulti(t *testing.T) {
	Prepare()
	replace.ParseBytes([]byte("@@"), Template)

	if r := replace.GetString(Map); string(r) != Expected {
		t.Errorf("unexpected result: %#v, expected: %#v", string(r), Expected)
	}

	m := replace.AllPos(Map)

	if r := replace.SetPosBytes(m); string(r) != Expected {
		t.Errorf("unexpected result for: %#v, expected: %#v", string(r), Expected)
	}
}

func TestReplaceSyntaxError(t *testing.T) {
	ſ := replace.ParseBytes([]byte("@@"), []byte("before @@one@@@@two@@ after"))
	// ŧ.Println(ſ)
	if ſ == nil {
		t.Errorf("expected syntax error for 2 placeholders side by side, got none")
	}
}
