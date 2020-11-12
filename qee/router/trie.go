package router

import (
	"QeeWeb/qee/base"
	"fmt"
	"os"
	"strings"
)

type trieRouter struct {
	handlerMap map[string]func(ctx *base.Context) // pattern: handler
	root *trieNode
}

type trieNode struct {
	pattern string
	part string
	children []*trieNode
	isWild bool
}

func NewTrieRouter() *trieRouter {
	return &trieRouter{
		handlerMap: map[string]func(ctx *base.Context){},
		root: &trieNode{
			pattern: "",
			part: "/",
			children: []*trieNode{},
			isWild: false,
		},
	}
}

func (r *trieRouter) DisplayUrlPattern() {
	for k ,_ := range r.handlerMap {
		fmt.Println(k)
	}
}

func (node *trieNode) MatchChild(part string) *trieNode {
	for _, child := range node.children {
		if child.part == part && !child.isWild {
			return child
		}
	}
	return nil
}

func (node *trieNode) MatchChilds(part string) []*trieNode {
	resp := []*trieNode{}
	for _, child := range node.children {
		if child.part == part || child.isWild {
			resp = append(resp, child)
		}
	}
	return resp
}

func (node *trieNode) insertNode(pattern string, parts []string) {
	if len(parts) == 0 {
		return
	} else if len(parts) == 1 {
		child := &trieNode{
			pattern: pattern,
			part: parts[0],
			children: []*trieNode{},
			isWild: strings.HasPrefix(parts[0], ":") || strings.HasPrefix(parts[0], "*"),
		}
		node.children = append(node.children, child)
		return
	}

	child := node.MatchChild(parts[0])
	if child == nil {
		child = &trieNode{
			pattern: "",
			part: parts[0],
			children: []*trieNode{},
			isWild: strings.HasPrefix(parts[0], ":") || strings.HasPrefix(parts[0], "*"),
		}
		node.children = append(node.children, child)
	}
	child.insertNode(pattern, parts[1:])
}

func prioritySort(nodes []*trieNode) {
	absPatternNodes := []*trieNode{}
	halfPatternNodes := []*trieNode{}
	allPatternNodes := []*trieNode{}

	for _, node := range nodes {
		if node.pattern != "" && !node.isWild {
			absPatternNodes = append(absPatternNodes, node)
			continue
		}

		if strings.HasPrefix(node.part, ":") {
			halfPatternNodes = append(halfPatternNodes, node)
			continue
		}

		if strings.HasPrefix(node.part, "*") {
			allPatternNodes = append(allPatternNodes, node)
			continue
		}
	}
	nodes = append([]*trieNode{}, absPatternNodes...)
	nodes = append(nodes, halfPatternNodes...)
	nodes = append(nodes, allPatternNodes...)
}

func (node *trieNode) searchNode(parts []string, queryDict map[string]string) *trieNode {
	if len(parts) == 0 {
		return node
	}

	matchChilds := node.MatchChilds(parts[0])
	var rspNode *trieNode
	var firstAllPatternNode *trieNode
	var childNode *trieNode
	prioritySort(matchChilds) // abs pattern -> half pattern -> all pattern

	for _, child := range matchChilds {
		if strings.HasPrefix(child.part, "*") {
			firstAllPatternNode = child
			break
		}
		rspNode = child.searchNode(parts[1:], queryDict)
		if rspNode != nil {
			childNode = child
			break
		}
	}

	if rspNode == nil && firstAllPatternNode != nil {
		rspNode = firstAllPatternNode
		if strings.HasPrefix(rspNode.part, "*") {
			queryDict[rspNode.part[1:]] = strings.Join(parts, "/")
		}
	}

	if rspNode != nil && childNode != nil && strings.HasPrefix(childNode.part, ":") {
		if _, exist := queryDict[rspNode.part[1:]]; !exist {
			queryDict[childNode.part[1:]] = parts[0]
		}
	}
	return rspNode
}

// register pattern to trie router
func (r *trieRouter) RegisteredHandler(method string, pattern string, handler func(ctx *base.Context)) error {
	parts := strings.Split(pattern, "/")
	if len(parts) <= 1 {
		return fmt.Errorf("error pattern format: %s", pattern)
	}

	r.root.insertNode(pattern, parts)
	routerKey := method + "-" + pattern
	r.handlerMap[routerKey] = handler
	return nil
}

func (r *trieRouter) FindHandler(method string, path string) (func(ctx *base.Context), map[string]string) {
	parts := strings.Split(path, "/")
	if len(parts) <= 1 {
		fmt.Fprintf(os.Stderr, fmt.Errorf("error path format: %s", parts).Error())
		return nil, nil
	}
	patternMap := map[string]string{}
	node := r.root.searchNode(parts, patternMap)
	handler := r.handlerMap[method + "-" + node.pattern]
	return handler, patternMap
}