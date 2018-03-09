# goxml

>xml parser for golang

**no othor dependencies**

file in and struct out

```
type Node struct {
	TagName string
	Content string
	Childs  []*Node
	Attrs   map[string]string
}
```
free style

use

```
go get github.com/wenlaizhou/goxml
```

```
package main

import (
	"github.com/wenlaizhou/goxml"
)

func main() {
	root, err := goxml.ParseFile("demo.xml")
	if err == nil {
		println(root.TagName)
	}
}

```
