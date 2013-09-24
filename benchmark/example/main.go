package main

import (
	"bytes"
	"fmt"
	"github.com/metakeule/fastreplace/benchmark"
)

func main() {
	t := benchmark.NewTemplate()
	err := t.Parse(`hi {{.ho}} {{.hu}}`)

	if err != nil {
		panic(err.Error())
	}

	var buffer bytes.Buffer
	m := map[string]string{
		"ho": "holla",
		"hu": "hulloa",
	}
	if err = t.Replace(m, &buffer); err != nil {
		panic(err.Error())
	}

	fmt.Println(buffer.String())
}
