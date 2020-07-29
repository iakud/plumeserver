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
	NIL   *Node
	root  *Node
	count uint
}

func New() *Rbtree {
	node := &Node{color: kBlack}
	return &Rbtree{
		NIL:   node,
		root:  node,
		count: 0,
	}
}

func (t *Rbtree) Len() uint {
	return t.count
}

func (t *Rbtree) Min() Item {
	if x := t.min(t.root); x == t.NIL {
		return x.Item
	}
	return nil
}

func (t *Rbtree) min(x *Node) *Node {
	if x == t.NIL {
		return t.NIL
	}
	for x.left != t.NIL {
		x = x.left
	}
	return x
}

func (t *Rbtree) Max() Item {
	if x := t.max(t.root); x != t.NIL {
		return x.Item
	}
	return nil
}

func (t *Rbtree) max(x *Node) *Node {
	if x == t.NIL {
		return t.NIL
	}
	for x.right != t.NIL {
		x = x.right
	}
	return x
}

func (t *Rbtree) Get(item Item) Item {
	if item == nil {
		return nil
	}
	if x := t.search(item); x != t.NIL {
		return x.Item
	}
	return nil
}

func (t *Rbtree) search(item Item) *Node {
	p := t.root
	for p != t.NIL {
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
	t.insert(item)
}

func (t *Rbtree) insert(item Item) {
	x := t.root
	y := t.NIL

	for x != t.NIL {
		y = x
		if item.Less(x.Item) {
			x = x.left
		} else if x.Item.Less(item) {
			x = x.right
		} else {
			return
		}
	}
	z := &Node{t.NIL, t.NIL, y, kRed, item}
	if y == t.NIL {
		t.root = z
	} else if item.Less(y.Item) {
		y.left = z
	} else {
		y.right = z
	}

	t.count++
	t.insertFixup(z)
}

func (t *Rbtree) insertFixup(z *Node) {
	for z.parent.color == kRed {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y.color == kRed {
				z.parent.color = kBlack
				y.color = kBlack
				z.parent.parent.color = kRed
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					t.leftRotate(z)
				}
				z.parent.color = kBlack
				z.parent.parent.color = kRed
				t.rightRotate(z.parent.parent)
			}
		} else {
			y := z.parent.parent.left
			if y.color == kRed {
				z.parent.color = kBlack
				y.color = kBlack
				z.parent.parent.color = kRed
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					t.rightRotate(z)
				}
				z.parent.color = kBlack
				z.parent.parent.color = kRed
				t.leftRotate(z.parent.parent)
			}
		}
	}
	t.root.color = kBlack
}

func (t *Rbtree) Delete(item Item) {
	if item == nil {
		return
	}
	t.delete(item)
}

func (t *Rbtree) delete(item Item) {
	z := t.search(item)

	if z == t.NIL {
		return
	}

	var x, y *Node
	if z.left == t.NIL || z.right == t.NIL {
		y = z
	} else {
		y = t.successor(z)
	}

	if y.left != t.NIL {
		x = y.left
	} else {
		x = y.right
	}

	x.parent = y.parent

	if y.parent == t.NIL {
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
		t.deleteFixup(x)
	}

	t.count--
}

func (t *Rbtree) deleteFixup(x *Node) {
	for x != t.root && x.color == kBlack {
		if x == x.parent.left {
			w := x.parent.right
			if w.color == kRed {
				w.color = kBlack
				x.parent.color = kRed
				t.leftRotate(x.parent)
				w = x.parent.right
			}
			if w.left.color == kBlack && w.right.color == kBlack {
				w.color = kRed
				x = x.parent
			} else {
				if w.right.color == kBlack {
					w.left.color = kBlack
					w.color = kRed
					t.rightRotate(w)
					w = x.parent.right
				}
				w.color = x.parent.color
				x.parent.color = kBlack
				w.right.color = kBlack
				t.leftRotate(x.parent)
				x = t.root
			}
		} else {
			w := x.parent.left
			if w.color == kRed {
				w.color = kBlack
				x.parent.color = kRed
				t.rightRotate(x.parent)
				w = x.parent.left
			}
			if w.left.color == kBlack && w.right.color == kBlack {
				w.color = kRed
				x = x.parent
			} else {
				if w.left.color == kBlack {
					w.right.color = kBlack
					w.color = kRed
					t.leftRotate(w)
					w = x.parent.left
				}
				w.color = x.parent.color
				x.parent.color = kBlack
				w.left.color = kBlack
				t.rightRotate(x.parent)
				x = t.root
			}
		}
	}
	x.color = kBlack
}

func (t *Rbtree) successor(x *Node) *Node {
	if x == t.NIL {
		return t.NIL
	}

	if x.right != t.NIL {
		return t.min(x.right)
	}

	y := x.parent
	for y != t.NIL && x == y.right {
		x = y
		y = y.parent
	}
	return y
}

func (t *Rbtree) leftRotate(x *Node) {
	if x.right == t.NIL {
		return
	}

	y := x.right
	x.right = y.left
	if y.left != t.NIL {
		y.left.parent = x
	}
	y.parent = x.parent

	if x.parent == t.NIL {
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
	if x.left == t.NIL {
		return
	}

	y := x.left
	x.left = y.right
	if y.right != t.NIL {
		y.right.parent = x
	}
	y.parent = x.parent

	if x.parent == t.NIL {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.right = x
	x.parent = y
}
