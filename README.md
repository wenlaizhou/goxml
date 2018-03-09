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
root, err := goxml.ParseFile(fileName)
println(root.TagName)
```
