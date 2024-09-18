package radix_tree

import (
	"bytes"
	"fmt"
	"hash/fnv"
)

const Term rune = -1

type Node struct {
	prefix string
	child  map[rune]*Node
}

func NewRadix() *Node {
	return &Node{
		child: make(map[rune]*Node),
	}
}

func (n *Node) hash(prefix string, buf *bytes.Buffer) {
	termS := ""
	if _, ok := n.child[Term]; ok {
		termS = " ($)"
	}

	buf.WriteString(
		fmt.Sprintf(
			"%s[%s]%s\n",
			prefix,
			n.prefix,
			termS,
		),
	)

	for r, child := range n.child {
		if r == Term {
			continue
		}
		child.hash(prefix+"  ", buf)
	}
}

func (n *Node) Hash() uint64 {
	buf := bytes.NewBufferString("")
	n.hash("", buf)

	h := fnv.New64a()
	_, _ = h.Write(buf.Bytes())
	return h.Sum64()
}

// ToDAWG converts a radix tree into a DAWG by replacing equivalent subtrees with references
func (n *Node) ToDAWG() *Node {
	subtrees := make(map[uint64]*Node)

	var processNode func(*Node) *Node
	processNode = func(node *Node) *Node {
		if node == nil {
			return nil
		}

		for r, child := range node.child {
			node.child[r] = processNode(child)
		}

		hash := node.Hash()

		if existing, ok := subtrees[hash]; ok {
			return existing
		}

		subtrees[hash] = node
		return node
	}

	return processNode(n)
}

func (n *Node) commonPrefixLength(s, prefix string) int {
	l := 0
	for l < len(s) && l < len(prefix) && s[l] == prefix[l] {
		l++
	}

	return l
}

func (n *Node) Add(s string) {
	if len(s) == 0 {
		return
	}

	l := n.commonPrefixLength(s, n.prefix)

	// If the common prefix is shorter than the node's prefix, split the node
	if l < len(n.prefix) {
		newChild := &Node{
			prefix: n.prefix[l:],
			child:  n.child,
		}
		n.child = map[rune]*Node{
			rune(n.prefix[l]): newChild,
		}
		n.prefix = n.prefix[:l]
	}

	// If there's remaining part of the string, add it as a new child
	if l < len(s) {
		r := rune(s[l])
		if child, exists := n.child[r]; exists {
			child.Add(s[l:])
		} else {
			n.child[r] = &Node{
				prefix: s[l:],
				child: map[rune]*Node{
					Term: {},
				},
			}
		}
	} else {
		n.child[Term] = &Node{}
	}
}

func (n *Node) string(prefix string, buf *bytes.Buffer, visited map[*Node]bool) {
	termS := ""
	if _, ok := n.child[Term]; ok {
		termS = " ($)"
	}

	refS := ""
	if visited[n] {
		refS = "*"
	}
	visited[n] = true

	buf.WriteString(
		fmt.Sprintf(
			"%p %s[%s]%s%s\n",
			n,
			prefix,
			n.prefix,
			refS,
			termS,
		),
	)

	for r, child := range n.child {
		if r == Term {
			continue
		}
		child.string(prefix+"  ", buf, visited)
	}
}

func (n *Node) String() string {
	buf := bytes.NewBufferString("")
	n.string("", buf, map[*Node]bool{})
	return buf.String()
}

func (n *Node) Find(s string) []string {
	var results []string
	n.find(s, 0, "", &results)
	return results
}

func (n *Node) find(s string, pos int, accPattern string, results *[]string) {
	var matchPrefix func(pattern string, pIdx int, s string, sIdx int) []int

	matchPrefix = func(pattern string, pIdx int, s string, sIdx int) []int {
		var endPositions []int

		if pIdx == len(pattern) {
			endPositions = append(endPositions, sIdx)
			return endPositions
		}

		if sIdx > len(s) {
			return endPositions
		}

		c := pattern[pIdx]

		if c == '*' {
			endPositions = append(endPositions, matchPrefix(pattern, pIdx+1, s, sIdx)...)
			if sIdx < len(s) {
				endPositions = append(endPositions, matchPrefix(pattern, pIdx, s, sIdx+1)...)
			}
		} else if c == '?' {
			if sIdx < len(s) {
				endPositions = append(endPositions, matchPrefix(pattern, pIdx+1, s, sIdx+1)...)
			}
		} else {
			if sIdx < len(s) && pattern[pIdx] == s[sIdx] {
				endPositions = append(endPositions, matchPrefix(pattern, pIdx+1, s, sIdx+1)...)
			}
		}

		return endPositions
	}

	endPositions := matchPrefix(n.prefix, 0, s, pos)

	for _, endPos := range endPositions {
		newAccPattern := accPattern + n.prefix

		if endPos == len(s) {
			if _, ok := n.child[Term]; ok {
				*results = append(*results, newAccPattern)
			}
		}

		for r, child := range n.child {
			if r == Term {
				continue
			}
			child.find(s, endPos, newAccPattern, results)
		}
	}
}
