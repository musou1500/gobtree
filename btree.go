package gobtree

import "golang.org/x/exp/slices"

type Item interface {
	Less(than Item) bool
}

type items []Item

type node struct {
	items    items
	children []*node
}

type BTree struct {
	root   *node
	degree int
}

/**
 * find index to insert item
 */
func (items items) find(item Item) (bool, int) {
	for i := 0; i < len(items); i++ {
		if !items[i].Less(item) {
			return !item.Less(items[i]), i
		}
	}
	return false, len(items)
}

func New(degree int) *BTree {
	if degree < 2 {
		panic("degree must be greater than 1")
	}
	return &BTree{degree: degree, root: &node{}}
}

func (n *node) split(middle int) (Item, *node) {
	item := n.items[middle]
	next := &node{}
	next.items = append(next.items, n.items[middle+1:]...)
	if len(n.children) > 0 {
		next.children = append(next.children, n.children[middle+1:]...)
	}
	n.items = n.items[:middle]
	return item, next
}

func (n *node) insert(item Item, maxItems int) {
	found, index := n.items.find(item)
	if found {
		n.items[index] = item
		return
	}

	if len(n.children) == 0 {
		// insert to leaf node
		n.items = slices.Insert(n.items, index, item)
		return
	}

	// need to split?
	if len(n.children[index].items) >= maxItems {
		newItem, newNode := n.children[index].split(maxItems / 2)
		n.items = slices.Insert(n.items, index, newItem)
		n.children = slices.Insert(n.children, index+1, newNode)
		if index+1 < len(n.children) && n.items[index].Less(item) {
			index++
		}
	}

	// insert to internal node
	n.children[index].insert(item, maxItems)
}

func (btree *BTree) InsertOrReplace(item Item) {
	if len(btree.root.items) >= btree.maxItems() {
		newItem, newNode := btree.root.split(btree.maxItems() / 2)
		newRoot := &node{}
		newRoot.items = append(newRoot.items, newItem)
		newRoot.children = append(newRoot.children, btree.root, newNode)
		btree.root = newRoot
	}
	btree.root.insert(item, btree.maxItems())
}

func (btree *BTree) maxItems() int {
	return btree.degree*2 - 1
}
