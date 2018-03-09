package goxml

import (
	"regexp"
	"strings"
	"io/ioutil"
	"errors"
	"fmt"
	"log"
)

type Node struct {
	TagName string
	Content string
	Childs  []*Node
	Attrs   map[string]string
}

//find first tagName
var startTag = regexp.MustCompile(`<\s*(\w+)\s*.*?>`)

//find whole tag
var oneTag = `(?ms:<\s*%s\s*(.*?)>(.*?)</\s*%s\s*>(.*))`

//parse xml file
//
//return rootNodePtr and error
func ParseFile(filePath string) (*Node, error) {
	if len(filePath) <= 0 {
		return nil, errors.New("filePath is empty")
	}
	fileBytes, err := ioutil.ReadFile(filePath)
	if processError(err) {
		return nil, errors.New("file not exist or can not read")
	}
	xmlStr := strings.TrimSpace(string(fileBytes))
	if len(xmlStr) <= 0 {
		return nil, errors.New("file size error")
	}
	xmlStr = processInclude(xmlStr)
	root, contentStr, _ := firstNode(nil, xmlStr)
	processXml(contentStr, root)
	return root, nil
}

//do process xml str
func processXml(xmlStr string, parent *Node) {
	nodePtr, content, nextStr := firstNode(parent, xmlStr)
	if nodePtr != nil && len(content) > 0 {
		processXml(content, nodePtr)
	}
	if parent != nil && len(nextStr) > 0 {
		processXml(nextStr, parent)
	}
}

//find first node
//
//return: node, contentStr, leftStr
func firstNode(parent *Node, xmlStr string) (*Node, string, string) {
	firstMatch := startTag.FindStringSubmatch(xmlStr)
	if len(firstMatch) < 2 { //no node
		return nil, "", ""
	}
	xmlStr = strings.TrimSpace(xmlStr)
	nodeTag := firstMatch[1]
	oneReg := fmt.Sprintf(oneTag, nodeTag, nodeTag)
	endMatch := regexp.MustCompile(oneReg).FindStringSubmatch(xmlStr)
	if len(endMatch) < 4 {
		return nil, "", ""
	}
	nodeAttrStr := strings.TrimSpace(endMatch[1])
	nodeContentText := endMatch[2]
	nextStr := strings.TrimSpace(endMatch[3])
	node := createNode(nodeTag, nodeContentText, nodeAttrStr)
	if parent != nil {
		parent.Childs = append(parent.Childs, node)
	}
	return node, nodeContentText, nextStr
}

func createNode(tag string, content string, attrs string) *Node {
	return &Node{
		TagName: strings.TrimSpace(tag),
		Content: content,
		Attrs:   buildNodeAttrs(attrs),
	}
}

// node attribute regex
var attrReg = regexp.MustCompile(`((\w+.*?)=("(.*?)"|'(.*?)'))`)

func buildNodeAttrs(attrsStr string) map[string]string {
	res := make(map[string]string)
	if len(attrsStr) <= 0 {
		return res
	}
	matches := attrReg.FindAllStringSubmatch(attrsStr, -1)
	if len(matches) <= 0 {
		return res
	}
	for _, match := range matches {
		if len(match) == 6 {
			value, match4, match5 := "", match[4], match[5]
			if len(match4) > 0 {
				value = match[4]
			}
			if len(match5) > 0 {
				value = match5
			}
			res[match[2]] = value
		}
	}
	return res
}

// read <include src="filePath" /> and replace string
func processInclude(data string) string {

	includeReg := regexp.MustCompile(`<\s*include\s*src\s*=\s*("(.*)"|'(.*)').*/>`)
	allSub := includeReg.FindAllStringSubmatch(string(data), -1)
	if len(allSub) <= 0 {
		return data
	}
	for _, subList := range allSub {
		var filePath string
		if len(subList[2]) > 0 {
			filePath = subList[2]
		} else {
			filePath = subList[3]
		}
		includeData, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}
		data = strings.Replace(data, subList[0],
			strings.TrimSpace(string(includeData)), 1)
	}
	return data
}

// error processor
//
// return true error is not nil
//
// false error is nil
func processError(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
