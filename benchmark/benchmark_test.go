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
var MultiByteMap2 = map[string][]byte{}

func PrepareMulti() {
	MultiMap = map[string]string{}
	MultiByteMap = map[string][]byte{}
	MultiByteMap2 = map[string][]byte{}
	orig := []string{}
	exp := []string{}
	for i := 0; i < 5000; i++ {
		orig = append(orig, ŧ.Sprintf(`a string with @@replacement%v@@`, i))
		exp = append(exp, ŧ.Sprintf("a string with repl%v", i))
		key := ŧ.Sprintf("replacement%v", i)
		val := ŧ.Sprintf("repl%v", i)
		MultiMap["@@"+key+"@@"] = val
		MultiByteMap[key] = []byte(val)
		MultiByteMap2["@@"+key+"@@"] = []byte(val)
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

var ByteMap2 = map[string][]byte{
	"@@replacement1@@": []byte("repl1"),
	"@@replacement2@@": []byte("repl2"),
}

var Expected = "a string with repl1 and repl2 that c@ntinues"

var mapperNaive = &Naive{}
var mapperReg = &Regexp{Regexp: regexp.MustCompile("(@@[^@]+@@)")}
var freplace = &FReplace{}
var byts = &Bytes{}

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

	byts.Map = ByteMap2
	byts.Parse(Template)
	if r := byts.Replace(); string(r) != Expected {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "byts", string(r), Expected)
	}

	freplace.ParseBytes([]byte("@@"), ByteTemplate)

	if r := freplace.ReplaceBytes(ByteMap); string(r) != Expected {
		t.Errorf("unexpected result for %s: %#v", "freplace", string(r))
	}

	m := freplace.AllPos(ByteMap)

	if r := freplace.ReplacePosBytes(m); string(r) != Expected {
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

	freplace.ParseBytes([]byte("@@"), ByteTemplateX)

	if r := freplace.ReplaceBytes(ByteMap); string(r) != ExpectedX {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "freplace", string(r), ExpectedX)
	}

	m := freplace.AllPos(ByteMap)

	if r := freplace.ReplacePosBytes(m); string(r) != ExpectedX {
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

	freplace.ParseBytes([]byte("@@"), MultiByteTemplate)

	if r := freplace.ReplaceBytes(MultiByteMap); string(r) != MultiExpected {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "freplace", string(r), MultiExpected)
	}

	m := freplace.AllPos(MultiByteMap)

	if r := freplace.ReplacePosBytes(m); string(r) != MultiExpected {
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

func BenchmarkByte(b *testing.B) {
	b.StopTimer()
	byts.Map = ByteMap2
	byts.Parse(TemplateX)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		byts.Replace()
	}
}

func BenchmarkFReplace(b *testing.B) {
	b.StopTimer()
	freplace.ParseBytes([]byte("@@"), ByteTemplateX)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		freplace.ReplaceBytes(ByteMap)
	}
}

func BenchmarkFReplacePos(b *testing.B) {
	b.StopTimer()
	freplace.ParseBytes([]byte("@@"), ByteTemplateX)
	m := freplace.AllPos(ByteMap)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		freplace.ReplacePosBytes(m)
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

func BenchmarkByteM(b *testing.B) {
	b.StopTimer()
	PrepareMulti()
	byts.Map = MultiByteMap2
	byts.Parse(MultiTemplate)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		byts.Replace()
	}
}

func BenchmarkFReplaceM(b *testing.B) {
	b.StopTimer()
	PrepareMulti()
	freplace.ParseBytes([]byte("@@"), MultiByteTemplate)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		freplace.ReplaceBytes(MultiByteMap)
	}
}

func BenchmarkFReplacePosM(b *testing.B) {
	b.StopTimer()
	PrepareMulti()
	freplace.ParseBytes([]byte("@@"), MultiByteTemplate)
	m := freplace.AllPos(MultiByteMap)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		freplace.ReplacePosBytes(m)
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
	mapperReg.Setup()
	for i := 0; i < b.N; i++ {
		mapperReg.Map = Map
		mapperReg.Template = TemplateX
		mapperReg.Replace()
	}
}

func BenchmarkOneShotByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byts.Parse(TemplateX)
		byts.Map = ByteMap2
		byts.Replace()
	}
}

func BenchmarkFReplaceOneShot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		freplace.ParseBytes([]byte("@@"), ByteTemplateX)
		freplace.ReplaceBytes(ByteMap)
	}
}
