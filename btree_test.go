package gobtree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func dump(n *node, lv int) string {
	prefix := strings.Repeat("\t", lv)
	s := ""
	for _, item := range n.items {
		s += fmt.Sprintf("%s%v\n", prefix, item)
	}

	for i, child := range n.children {
		s += dump(child, lv+1)
		if i < len(n.children)-1 {
			s += "\n"
		}
	}

	return s
}

type Int int

func (i Int) Less(than Item) bool {
	return i < than.(Int)
}

func TestBTreeSplit(t *testing.T) {

	entries := []struct {
		items    []int
		expected string
	}{
		{items: []int{1, 2, 3}, expected: `1
2
3
`,
		},
		{items: []int{1, 2, 3, 4, 5, 6}, expected: `2
4
	1

	3

	5
	6
`,
		},
		{items: []int{1, 6, 2, 3}, expected: `2
	1

	3
	6
`,
		},
		{items: []int{5, 6, 7, 2, 3}, expected: `6
	2
	3
	5

	7
`,
		},
		{items: []int{1, 2, 3, 1, 2}, expected: `2
	1

	3
`,
		},
	}

	for _, entry := range entries {
		bt := New(2)
		for _, item := range entry.items {
			bt.InsertOrReplace(Int(item))
		}
		result := dump(bt.root, 0)
		fmt.Print(result)
		require.Equal(t, entry.expected, result)
	}
}
