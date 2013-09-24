package benchmark

import (
	"bytes"
	ŧ "fmt"
	. "github.com/metakeule/fastreplace"
	"github.com/metakeule/replacer"
	"regexp"
	"strings"
	"testing"
)

var Template = "a string with @@replacement1@@ and @@replacement2@@ that c@ntinues"
var TTemplate = "a string with {{.replacement1}} and {{.replacement2}} that c@ntinues"
var ByteTemplate = []byte(Template)

var TemplateX = ""
var TemplateTX = ""
var ByteTemplateX = []byte{}
var ExpectedX = ""

var MultiTemplate = ""
var MultiTemplateT = ""
var MultiByteTemplate = []byte{}
var MultiExpected = ""
var MultiMap = map[string]string{}
var MultiByteMap = map[string][]byte{}
var MultiByteMap2 = map[string][]byte{}
var T = NewTemplate()

//var MultiPlaceholderMap = map[fastreplace2.Placeholder]string{}
var MultiPlaceholderMap = map[string]string{}

func PrepareMulti() {
	MultiMap = map[string]string{}
	MultiByteMap = map[string][]byte{}
	MultiByteMap2 = map[string][]byte{}
	//MultiPlaceholderMap = map[fastreplace2.Placeholder]string{}
	MultiPlaceholderMap = map[string]string{}
	orig := []string{}
	exp := []string{}
	tx := []string{}
	for i := 0; i < 5000; i++ {
		orig = append(orig, ŧ.Sprintf(`a string with @@replacement%v@@`, i))
		tx = append(tx, ŧ.Sprintf(`a string with {{.replacement%v}}`, i))
		exp = append(exp, ŧ.Sprintf("a string with repl%v", i))
		key := ŧ.Sprintf("replacement%v", i)
		val := ŧ.Sprintf("repl%v", i)
		MultiMap["@@"+key+"@@"] = val
		MultiByteMap[key] = []byte(val)
		MultiByteMap2["@@"+key+"@@"] = []byte(val)
		//MultiPlaceholderMap[fastreplace2.NewPlaceholder(key)] = val
		MultiPlaceholderMap[key] = val
	}
	MultiTemplate = strings.Join(orig, "")
	MultiTemplateT = strings.Join(tx, "")
	MultiExpected = strings.Join(exp, "")
	MultiByteTemplate = []byte(MultiTemplate)
}

func PrepareX() {
	orig := []string{}
	exp := []string{}
	tx := []string{}
	for i := 0; i < 2500; i++ {
		orig = append(orig, Template)
		exp = append(exp, Expected)
		tx = append(tx, TTemplate)
	}
	TemplateTX = strings.Join(tx, "")
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

var PlaceholderMap = map[string]string{
	"replacement1": "repl1",
	"replacement2": "repl2",
}

var Expected = "a string with repl1 and repl2 that c@ntinues"

var mapperNaive = &Naive{}
var mapperReg = &Regexp{Regexp: regexp.MustCompile("(@@[^@]+@@)")}
var freplace = &FReplace{}
var byts = &Bytes{}
var freplace2 = replacer.New()

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

	T.Parse(TTemplate)
	var tbf bytes.Buffer
	if T.Replace(PlaceholderMap, &tbf); tbf.String() != Expected {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "template", tbf.String(), Expected)
	}

	freplace.ParseBytes([]byte("@@"), ByteTemplate)

	if r := freplace.ReplaceBytes(ByteMap); string(r) != Expected {
		t.Errorf("unexpected result for %s: %#v", "freplace", string(r))
	}

	m := freplace.AllPos(ByteMap)

	if r := freplace.ReplacePosBytes(m); string(r) != Expected {
		t.Errorf("unexpected result for %s: %#v", "freplace-ReplacePos", string(r))
	}

	efr2 := freplace2.Parse(ByteTemplate)
	if efr2 != nil {
		panic(efr2.Error())
	}

	var bf bytes.Buffer
	//if r := fr2.Replace(PlaceholderMap); string(r) != Expected {
	if freplace2.Replace(PlaceholderMap, &bf); bf.String() != Expected {
		t.Errorf("unexpected result for %s: %#v", "fastreplace2", bf.String())
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

	T.Parse(TemplateTX)
	var tbf bytes.Buffer
	if T.Replace(PlaceholderMap, &tbf); tbf.String() != ExpectedX {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "template", tbf.String(), ExpectedX)
	}

	freplace.ParseBytes([]byte("@@"), ByteTemplateX)

	if r := freplace.ReplaceBytes(ByteMap); string(r) != ExpectedX {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "freplace", string(r), ExpectedX)
	}

	m := freplace.AllPos(ByteMap)

	if r := freplace.ReplacePosBytes(m); string(r) != ExpectedX {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "freplace-ReplacePos", string(r), ExpectedX)
	}

	efr2 := freplace2.Parse(ByteTemplateX)

	if efr2 != nil {
		panic(efr2.Error())
	}

	var bf bytes.Buffer
	if freplace2.Replace(PlaceholderMap, &bf); bf.String() != ExpectedX {
		t.Errorf("unexpected result for %s: %#v", "fastreplace2", bf.String())
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

	T.Parse(MultiTemplateT)
	var tbf bytes.Buffer
	if T.Replace(MultiPlaceholderMap, &tbf); tbf.String() != MultiExpected {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "template", tbf.String(), MultiExpected)
	}

	freplace.ParseBytes([]byte("@@"), MultiByteTemplate)

	if r := freplace.ReplaceBytes(MultiByteMap); string(r) != MultiExpected {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "freplace", string(r), MultiExpected)
	}

	m := freplace.AllPos(MultiByteMap)

	if r := freplace.ReplacePosBytes(m); string(r) != MultiExpected {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "freplace-ReplacePos", string(r), MultiExpected)
	}

	efr2 := freplace2.Parse(MultiByteTemplate)
	if efr2 != nil {
		panic(efr2.Error())
	}

	var bf bytes.Buffer

	//if r := fr2.Replace(MultiPlaceholderMap); string(r) != MultiExpected {
	if freplace2.Replace(MultiPlaceholderMap, &bf); bf.String() != MultiExpected {
		t.Errorf("unexpected result for %s: %#v", "fastreplace2", bf.String())
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
func BenchmarkTemplate(b *testing.B) {
	b.StopTimer()
	T.Parse(TemplateTX)
	var tbf bytes.Buffer
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		T.Replace(PlaceholderMap, &tbf)
		tbf.Reset()
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

func BenchmarkFReplace2(b *testing.B) {
	b.StopTimer()
	freplace2.Parse(ByteTemplateX)
	var bf bytes.Buffer
	//freplace.ParseBytes([]byte("@@"), ByteTemplateX)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		//freplace.ReplaceBytes(ByteMap)
		freplace2.Replace(PlaceholderMap, &bf)
		bf.Reset()
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

func BenchmarkTemplateM(b *testing.B) {
	b.StopTimer()
	PrepareMulti()
	T.Parse(MultiTemplateT)
	var tbf bytes.Buffer
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		T.Replace(MultiPlaceholderMap, &tbf)
		tbf.Reset()
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

func BenchmarkFReplace2M(b *testing.B) {
	b.StopTimer()
	PrepareMulti()
	freplace2.Parse(MultiByteTemplate)
	var bf bytes.Buffer
	//freplace.ParseBytes([]byte("@@"), MultiByteTemplate)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		freplace2.Replace(MultiPlaceholderMap, &bf)
		bf.Reset()
		//freplace.ReplaceBytes(MultiByteMap)
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

func BenchmarkTemplateOneShot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		T.Parse(TemplateTX)
		var tbf bytes.Buffer
		T.Replace(PlaceholderMap, &tbf)
	}
}

func BenchmarkFReplaceOneShot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		freplace.ParseBytes([]byte("@@"), ByteTemplateX)
		freplace.ReplaceBytes(ByteMap)
	}
}

func BenchmarkFReplace2OneShot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//freplace.ParseBytes([]byte("@@"), ByteTemplateX)
		freplace2.Parse(ByteTemplateX)
		//freplace.ReplaceBytes(ByteMap)
		var bf bytes.Buffer
		freplace2.Replace(PlaceholderMap, &bf)
	}
}
