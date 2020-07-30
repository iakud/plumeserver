package rbtree

type color uint8

const (
	kRed   color = 0
	kBlack color = 1
)

type Node struct {
	left   *Node
	right  *Node
	parent *Node
	color  color

	Item
}

type Item interface {
	Less(than Item) bool
}

type Rbtree struct {
	root  *Node
	count uint
}

func New() *Rbtree {
	return &Rbtree{
		root:  nil,
		count: 0,
	}
}

func (t *Rbtree) Len() uint {
	return t.count
}

func (t *Rbtree) Min() Item {
	if x := t.min(t.root); x == nil {
		return x.Item
	}
	return nil
}

func (t *Rbtree) min(x *Node) *Node {
	if x == nil {
		return nil
	}
	for x.left != nil {
		x = x.left
	}
	return x
}

func (t *Rbtree) Max() Item {
	if x := t.max(t.root); x != nil {
		return x.Item
	}
	return nil
}

func (t *Rbtree) max(x *Node) *Node {
	if x == nil {
		return nil
	}
	for x.right != nil {
		x = x.right
	}
	return x
}

func (t *Rbtree) Get(item Item) Item {
	if item == nil {
		return nil
	}
	if x := t.search(item); x != nil {
		return x.Item
	}
	return nil
}

func (t *Rbtree) search(item Item) *Node {
	p := t.root
	for p != nil {
		if p.Item.Less(item) {
			p = p.right
		} else if item.Less(p.Item) {
			p = p.left
		} else {
			break
		}
	}
	return p
}

func (t *Rbtree) Insert(item Item) {
	if item == nil {
		return
	}

	x := t.root
	var y *Node

	for x != nil {
		y = x
		if item.Less(x.Item) {
			x = x.left
		} else if x.Item.Less(item) {
			x = x.right
		} else {
			return
		}
	}
	x = &Node{nil, nil, y, kRed, item}
	if y == nil {
		t.root = x
	} else if item.Less(y.Item) {
		y.left = x
	} else {
		y.right = x
	}

	t.count++
	t.insertFixup(x)
}

func (t *Rbtree) insertFixup(x *Node) {
	for x != t.root && x.parent.color == kRed {
		if x.parent == x.parent.parent.left {
			y := x.parent.parent.right
			if y != nil && y.color == kRed {
				x.parent.color = kBlack
				y.color = kBlack
				x.parent.parent.color = kRed
				x = x.parent.parent
			} else {
				if x == x.parent.right {
					x = x.parent
					t.leftRotate(x)
				}
				x.parent.color = kBlack
				x.parent.parent.color = kRed
				t.rightRotate(x.parent.parent)
			}
		} else {
			y := x.parent.parent.left
			if y != nil && y.color == kRed {
				x.parent.color = kBlack
				y.color = kBlack
				x.parent.parent.color = kRed
				x = x.parent.parent
			} else {
				if x == x.parent.left {
					x = x.parent
					t.rightRotate(x)
				}
				x.parent.color = kBlack
				x.parent.parent.color = kRed
				t.leftRotate(x.parent.parent)
			}
		}
	}
	t.root.color = kBlack
}

func (t *Rbtree) Delete(item Item) {
	if item == nil {
		return
	}

	z := t.search(item)

	if z == nil {
		return
	}

	var x, y *Node
	if z.left == nil || z.right == nil {
		y = z
	} else {
		y = t.successor(z)
	}

	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}

	if x != nil {
		x.parent = y.parent
	}

	if y.parent == nil {
		t.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	if y != z {
		z.Item = y.Item
	}

	if y.color == kBlack {
		t.deleteFixup(x, y.parent)
	}
	t.count--
}

func (t *Rbtree) deleteFixup(x, parent *Node) {
	for x != t.root && (x == nil || x.color == kBlack) {
		if x == parent.left {
			w := parent.right
			if w.color == kRed {
				w.color = kBlack
				parent.color = kRed
				t.leftRotate(parent)
				w = parent.right
			}
			if (w.left == nil || w.left.color == kBlack) && (w.right == nil || w.right.color == kBlack) {
				w.color = kRed
				x = parent
				parent = x.parent
			} else {
				if w.right == nil || w.right.color == kBlack {
					w.left.color = kBlack
					w.color = kRed
					t.rightRotate(w)
					w = parent.right
				}
				w.color = parent.color
				parent.color = kBlack
				w.right.color = kBlack
				t.leftRotate(parent)
				x = t.root
			}
		} else {
			w := parent.left
			if w != nil && w.color == kRed {
				w.color = kBlack
				parent.color = kRed
				t.rightRotate(parent)
				w = parent.left
			}
			if (w.left == nil || w.left.color == kBlack) && (w.right == nil || w.right.color == kBlack) {
				w.color = kRed
				x = parent
				parent = x.parent
			} else {
				if w.left == nil || w.left.color == kBlack {
					w.right.color = kBlack
					w.color = kRed
					t.leftRotate(w)
					w = parent.left
				}
				w.color = parent.color
				parent.color = kBlack
				w.left.color = kBlack
				t.rightRotate(parent)
				x = t.root
			}
		}
	}
	if x != nil {
		x.color = kBlack
	}
}

func (t *Rbtree) successor(x *Node) *Node {
	if x == nil {
		return nil
	}
	if x.right != nil {
		return t.min(x.right)
	}
	y := x.parent
	for y != nil && x == y.right {
		x = y
		y = y.parent
	}
	return y
}

func (t *Rbtree) leftRotate(x *Node) {
	if x.right == nil {
		return
	}

	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent

	if x.parent == nil {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.left = x
	x.parent = y
}

func (t *Rbtree) rightRotate(x *Node) {
	if x.left == nil {
		return
	}

	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	y.parent = x.parent

	if x.parent == nil {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.right = x
	x.parent = y
}
