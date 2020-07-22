package rbtree

type color uint8

const (
	kRBTreeRed   color = 1
	kRBTreeBlack color = 2
)

type Node struct {
	Left   *Node
	Right  *Node
	Parent *Node
	color  color

	// for use by client.
	Interface
}

type Interface interface {
	Less(than Interface) bool
}

type Rbtree struct {
	NIL   *Node
	root  *Node
	count uint
}

func New() *Rbtree {
	node := &Node{color: kRBTreeBlack}
	return &Rbtree{
		NIL:   node,
		root:  node,
		count: 0,
	}
}
