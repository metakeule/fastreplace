package benchmark

import (
	ŧ "fmt"
	. "github.com/metakeule/fastreplace"
	"regexp"
	"strings"
	"testing"
)

var Template = "a string with @@replacement1@@ and @@replacement2@@ that c@ntinues"
var ByteTemplate = []byte(Template)

var TemplateX = ""
var ByteTemplateX = []byte{}
var ExpectedX = ""

var MultiTemplate = ""
var MultiByteTemplate = []byte{}
var MultiExpected = ""
var MultiMap = map[string]string{}
var MultiByteMap = map[string][]byte{}

func PrepareMulti() {
	MultiMap = map[string]string{}
	MultiByteMap = map[string][]byte{}
	orig := []string{}
	exp := []string{}
	for i := 0; i < 5000; i++ {
		orig = append(orig, ŧ.Sprintf(`a string with @@replacement%v@@`, i))
		exp = append(exp, ŧ.Sprintf("a string with repl%v", i))
		key := ŧ.Sprintf("replacement%v", i)
		val := ŧ.Sprintf("repl%v", i)
		MultiMap["@@"+key+"@@"] = val
		MultiByteMap[key] = []byte(val)
	}
	MultiTemplate = strings.Join(orig, "")
	MultiExpected = strings.Join(exp, "")
	MultiByteTemplate = []byte(MultiTemplate)
}

func PrepareX() {
	orig := []string{}
	exp := []string{}
	for i := 0; i < 2500; i++ {
		orig = append(orig, Template)
		exp = append(exp, Expected)
	}
	TemplateX = strings.Join(orig, "")
	ExpectedX = strings.Join(exp, "")
	ByteTemplateX = []byte(TemplateX)
}

var Map = map[string]string{
	"@@replacement1@@": "repl1",
	"@@replacement2@@": "repl2",
}

var ByteMap = map[string][]byte{
	"replacement1": []byte("repl1"),
	"replacement2": []byte("repl2"),
}

var Expected = "a string with repl1 and repl2 that c@ntinues"

var mapperNaive = &Naive{}
var mapperReg = &Regexp{Regexp: regexp.MustCompile("(@@[^@]+@@)")}
var freplace = &FReplace{}

func TestReplace(t *testing.T) {
	mapperNaive.Map = Map
	mapperNaive.Template = Template
	if r := mapperNaive.Replace(); r != Expected {
		t.Errorf("unexpected result for %s: %#v", "mapperNaive", r)
	}

	mapperReg.Map = Map
	mapperReg.Template = Template
	mapperReg.Setup()
	if r := mapperReg.Replace(); r != Expected {
		t.Errorf("unexpected result for %s: %#v", "mapperReg", r)
	}

	freplace.Parse("@@", ByteTemplate)

	if r := freplace.Replace(ByteMap); string(r) != Expected {
		t.Errorf("unexpected result for %s: %#v", "freplace", string(r))
	}

	m := map[int][]byte{}

	for k, v := range ByteMap {
		pos := freplace.Pos(k)
		for _, p := range pos {
			m[p] = v
		}
	}

	if r := freplace.ReplacePos(m); string(r) != Expected {
		t.Errorf("unexpected result for %s: %#v", "freplace-ReplacePos", string(r))
	}
}

func TestReplaceX(t *testing.T) {
	PrepareX()
	mapperNaive.Map = Map
	mapperNaive.Template = TemplateX
	if r := mapperNaive.Replace(); r != ExpectedX {
		t.Errorf("unexpected result for %s: %#v", "mapperNaive", r)
	}

	mapperReg.Map = Map
	mapperReg.Template = TemplateX
	mapperReg.Setup()
	if r := mapperReg.Replace(); r != ExpectedX {
		t.Errorf("unexpected result for %s: %#v", "mapperReg", r)
	}

	freplace.Parse("@@", ByteTemplateX)

	if r := freplace.Replace(ByteMap); string(r) != ExpectedX {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "freplace", string(r), ExpectedX)
	}

	m := freplace.AllPos(ByteMap)

	if r := freplace.ReplacePos(m); string(r) != ExpectedX {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "freplace-ReplacePos", string(r), ExpectedX)
	}

}

func TestReplaceMulti(t *testing.T) {
	PrepareMulti()
	mapperNaive.Map = MultiMap
	mapperNaive.Template = MultiTemplate
	if r := mapperNaive.Replace(); r != MultiExpected {
		t.Errorf("unexpected result for %s: %#v", "mapperNaive", r)
	}

	mapperReg.Map = MultiMap
	mapperReg.Template = MultiTemplate
	mapperReg.Setup()
	if r := mapperReg.Replace(); r != MultiExpected {
		t.Errorf("unexpected result for %s: %#v", "mapperReg", r)
	}

	freplace.Parse("@@", MultiByteTemplate)

	if r := freplace.Replace(MultiByteMap); string(r) != MultiExpected {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "freplace", string(r), MultiExpected)
	}

	m := freplace.AllPos(MultiByteMap)

	if r := freplace.ReplacePos(m); string(r) != MultiExpected {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "freplace-ReplacePos", string(r), MultiExpected)
	}

}

func BenchmarkNaive(b *testing.B) {
	b.StopTimer()
	mapperNaive.Map = Map
	mapperNaive.Template = TemplateX
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mapperNaive.Replace()
	}
}

func BenchmarkReg(b *testing.B) {
	b.StopTimer()
	mapperReg.Map = Map
	mapperReg.Template = TemplateX
	mapperReg.Setup()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mapperReg.Replace()
	}
}

func BenchmarkFReplace(b *testing.B) {
	b.StopTimer()
	freplace.Parse("@@", ByteTemplateX)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		freplace.Replace(ByteMap)
	}
}

func BenchmarkFReplacePos(b *testing.B) {
	b.StopTimer()
	freplace.Parse("@@", ByteTemplateX)
	m := freplace.AllPos(ByteMap)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		freplace.ReplacePos(m)
	}
}

func BenchmarkNaiveM(b *testing.B) {
	b.StopTimer()
	PrepareMulti()
	mapperNaive.Map = MultiMap
	mapperNaive.Template = MultiTemplate
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mapperNaive.Replace()
	}
}

func BenchmarkRegM(b *testing.B) {
	b.StopTimer()
	PrepareMulti()
	mapperReg.Map = MultiMap
	mapperReg.Template = MultiTemplate
	mapperReg.Setup()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mapperReg.Replace()
	}
}

func BenchmarkFReplaceM(b *testing.B) {
	b.StopTimer()
	PrepareMulti()
	freplace.Parse("@@", MultiByteTemplate)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		freplace.Replace(MultiByteMap)
	}
}

func BenchmarkFReplacePosM(b *testing.B) {
	b.StopTimer()
	PrepareMulti()
	freplace.Parse("@@", MultiByteTemplate)
	m := freplace.AllPos(MultiByteMap)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		freplace.ReplacePos(m)
	}
}

func BenchmarkNaiveOneShot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mapperNaive.Map = Map
		mapperNaive.Template = TemplateX
		mapperNaive.Replace()
	}
}

func BenchmarkOneShotReg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mapperReg.Map = Map
		mapperReg.Template = TemplateX
		mapperReg.Setup()
		mapperReg.Replace()
	}
}

func BenchmarkFReplaceOneShot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		freplace.Parse("@@", ByteTemplateX)
		freplace.Replace(ByteMap)
	}
}
