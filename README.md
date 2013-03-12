fastreplace
===========

[![Build Status](https://secure.travis-ci.org/metakeule/fastreplace.png)](http://travis-ci.org/metakeule/fastreplace)

The typical scenario is that your template never changes but the replacements for you placeholders. fastreplace
is faster than using (strings|bytes).replace or regexp.ReplaceAllStringFunc in this situation.

Performance
-----------

Run benchmarks in benchmark directory.
I get the following results on my laptop:

replacing 2 placeholders that occur 2500x in the template

	BenchmarkNaive	         500	 3759998 ns/op    2,5x (strings.Replace)
	BenchmarkReg	         100	25826150 ns/op   17,1x (regexp.ReplaceAllStringFunc)
	BenchmarkByte	        1000	 2862196 ns/op    1,9x (bytes.Replace)
	BenchmarkFReplace	    1000	 2115185 ns/op    1,4x
	BenchmarkFReplacePos	1000	 1513158 ns/op    1,0x

replacing 5000 placeholders that occur 1x in the template

	BenchmarkNaiveM	           1  5494931000 ns/op 3652,1x (strings.Replace)
	BenchmarkRegM	          50	29896260 ns/op	 19,9x (regexp.ReplaceAllStringFunc)
	BenchmarkByteM	           1  3881286000 ns/op 2579,7x (bytes.Replace)
	BenchmarkFReplaceM	    1000     2368541 ns/op	  1,6x
	BenchmarkFReplacePosM	1000	 1504577 ns/op	  1,0x

replacing 2 placeholders that occur 1x in the template, parse template each time (you should not do this)

	BenchmarkNaiveOneShot	 500	 3756994 ns/op	  1,3x (strings.Replace)
	BenchmarkOneShotReg	     100    25417730 ns/op	  8,6x (regexp.ReplaceAllStringFunc)
	BenchmarkOneShotByte    1000	 2941726 ns/op	  1,0x (bytes.Replace)
	BenchmarkFReplaceOneShot 200	 7813775 ns/op	  2,7x


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
