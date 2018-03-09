# goxml

>xml parser for golang

file in and struct out
free style
```
type Node struct {
	TagName string
	Content string
	Childs  []*Node
	Attrs   map[string]string
}
```

use

```
root, err := goxml.ParseFile(fileName)
println(root.TagName)
```
