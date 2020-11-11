package router

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type trieRouter struct {
	router map[string]func(ctx context.Context) // pattern: handler
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
		router: map[string]func(ctx context.Context){},
		root: &trieNode{
			pattern: "",
			part: "/",
			children: []*trieNode{},
			isWild: false,
		},
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

func (node *trieNode) searchNode(parts []string) *trieNode {
	if len(parts) == 0 {
		return nil
	}

	matchChilds := node.MatchChilds(parts[0])
	var firstAllPatternNode *trieNode
	prioritySort(matchChilds)

	fmt.Println("part: ", parts[0])
	for _, child := range matchChilds {
		fmt.Println(child.part)
	}

	if len(parts) == 1 {
		for _, child := range matchChilds {
			if firstAllPatternNode == nil && strings.HasPrefix(child.part, "*") {
				firstAllPatternNode = child
			}
			if len(child.children) != 0 {
				continue
			} else {
				return child
			}
		}
	} else {
		for _, child := range matchChilds {
			if firstAllPatternNode == nil && strings.HasPrefix(child.part, "*") {
				firstAllPatternNode = child
			}
			resultLeaf := child.searchNode(parts[1:])
			if resultLeaf != nil {
				return resultLeaf
			}
		}
	}

	return firstAllPatternNode
}

// register pattern to trie router
func (r *trieRouter) RegisteredHandler(pattern string) error {
	parts := strings.Split(pattern, "/")
	if len(parts) <= 1 {
		return fmt.Errorf("error pattern format: %s", pattern)
	}

	r.root.insertNode(pattern, parts)
	return nil
}

func (r *trieRouter) FindHandler(path string) *trieNode {
	parts := strings.Split(path, "/")
	if len(parts) <= 1 {
		fmt.Fprintf(os.Stderr, fmt.Errorf("error path format: %s", parts).Error())
		return nil
	}

	node := r.root.searchNode(parts)
	return node
}