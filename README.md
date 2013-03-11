fastreplace
===========

[![Build Status](https://secure.travis-ci.org/metakeule/fastreplace.png)](http://travis-ci.org/metakeule/fastreplace)

The typical scenario is that your template never changes but the replacements for you placeholders. fastreplace
is faster than using strings.replace or regexp.ReplaceAllStringFunc in this situation.

Performance
-----------

Run benchmarks in benchmark directory.
On my laptop I get the following results:

replacing 2 placeholders that occur 2500x in the template

	BenchmarkNaive	         500	   3653172 ns/op      2,6x (strings.replace)
	BenchmarkReg	         100	  24978530 ns/op     17,6x (regexp.ReplaceAllStringFunc)
	BenchmarkFReplace	    1000	   2037904 ns/op      1,4x
	BenchmarkFReplacePos	1000	   1418960 ns/op      1,0x

replacing 5000 placeholders that occur 1x in the template

	BenchmarkNaiveM		       1	5522101000 ns/op   3847,2x (strings.replace)
	BenchmarkRegM		      50	  28921900 ns/op     20,1x (regexp.ReplaceAllStringFunc)
	BenchmarkFReplaceM	    1000	   2382313 ns/op      1,7x
	BenchmarkFReplacePosM	1000	   1435357 ns/op      1,0x

replacing 2 placeholders that occur 1x in the template, parse template each time (you should not do this)

	BenchmarkNaiveOneShot	        500	   3743120 ns/op  1,0x (strings.replace)
	BenchmarkOneShotReg	            100	  25118460 ns/op  6,7x (regexp.ReplaceAllStringFunc)
	BenchmarkFReplaceOneShot	    100	  15681430 ns/op  4,2x


Example
-------

```go
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
```

results in

```
IBM / Armonk  -  Microsoft / Redmond
Google / Mountain View  -  Apple / Cupertino
```



Documentation
-------------

see http://godoc.org/github.com/metakeule/fastreplace
