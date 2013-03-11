package main

import (
	"fmt"
	"github.com/metakeule/fastreplace"
)

var r = fastreplace.NewString("@@", "@@name@@ / @@city@@")
var posCity = r.Pos("city")[0]
var posName = r.Pos("name")[0]

func main() {
	i := map[string][]byte{
		"city": []byte("Armonk"),
		"name": []byte("IBM"),
	}

	// more comfort
	m := r.Instance()
	m.AssignString("city", "Redmond")
	m.AssignString("name", "Microsoft")

	fmt.Println(r.ReplaceString(i), " - ", m)

	// to be even faster (FReplacePos in the benchmarks)
	// you could cache the positions of the needed placeholders (see merg func below)
	fmt.Println(
		string(merge([]byte("Google"), []byte("Mountain View"))),
		" - ",
		string(merge([]byte("Apple"), []byte("Cupertino"))),
	)
}

func merge(name []byte, city []byte) []byte {
	m := map[int][]byte{posCity: city, posName: name}
	return r.ReplacePos(m)
}
