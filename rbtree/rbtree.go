package rbtree

type color uint8

const (
	kRed   color = 0
	kBlack color = 1
)

type Node struct {
	Left   *Node
	Right  *Node
	Parent *Node
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
	for x.Left != t.NIL {
		x = x.Left
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
	for x.Right != t.NIL {
		x = x.Right
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
			p = p.Right
		} else if item.Less(p.Item) {
			p = p.Left
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
			x = x.Left
		} else if x.Item.Less(item) {
			x = x.Right
		} else {
			return
		}
	}
	z := &Node{t.NIL, t.NIL, y, kRed, item}
	if y == t.NIL {
		t.root = z
	} else if item.Less(y.Item) {
		y.Left = z
	} else {
		y.Right = z
	}

	t.count++
	t.insertFixup(z)
}

func (t *Rbtree) insertFixup(z *Node) {
	for z.Parent.color == kRed {
		if z.Parent == z.Parent.Parent.Left {
			y := z.Parent.Parent.Right
			if y.color == kRed {
				z.Parent.color = kBlack
				y.color = kBlack
				z.Parent.Parent.color = kRed
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right {
					z = z.Parent
					t.leftRotate(z)
				}
				z.Parent.color = kBlack
				z.Parent.Parent.color = kRed
				t.rightRotate(z.Parent.Parent)
			}
		} else {
			y := z.Parent.Parent.Left
			if y.color == kRed {
				z.Parent.color = kBlack
				y.color = kBlack
				z.Parent.Parent.color = kRed
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Left {
					z = z.Parent
					t.rightRotate(z)
				}
				z.Parent.color = kBlack
				z.Parent.Parent.color = kRed
				t.leftRotate(z.Parent.Parent)
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
	if z.Left == t.NIL || z.Right == t.NIL {
		y = z
	} else {
		y = t.successor(z)
	}

	if y.Left != t.NIL {
		x = y.Left
	} else {
		x = y.Right
	}

	x.Parent = y.Parent

	if y.Parent == t.NIL {
		t.root = x
	} else if y == y.Parent.Left {
		y.Parent.Left = x
	} else {
		y.Parent.Right = x
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
		if x == x.Parent.Left {
			w := x.Parent.Right
			if w.color == kRed {
				w.color = kBlack
				x.Parent.color = kRed
				t.leftRotate(x.Parent)
				w = x.Parent.Right
			}
			if w.Left.color == kBlack && w.Right.color == kBlack {
				w.color = kRed
				x = x.Parent
			} else {
				if w.Right.color == kBlack {
					w.Left.color = kBlack
					w.color = kRed
					t.rightRotate(w)
					w = x.Parent.Right
				}
				w.color = x.Parent.color
				x.Parent.color = kBlack
				w.Right.color = kBlack
				t.leftRotate(x.Parent)
				x = t.root
			}
		} else {
			w := x.Parent.Left
			if w.color == kRed {
				w.color = kBlack
				x.Parent.color = kRed
				t.rightRotate(x.Parent)
				w = x.Parent.Left
			}
			if w.Left.color == kBlack && w.Right.color == kBlack {
				w.color = kRed
				x = x.Parent
			} else {
				if w.Left.color == kBlack {
					w.Right.color = kBlack
					w.color = kRed
					t.leftRotate(w)
					w = x.Parent.Left
				}
				w.color = x.Parent.color
				x.Parent.color = kBlack
				w.Left.color = kBlack
				t.rightRotate(x.Parent)
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

	if x.Right != t.NIL {
		return t.min(x.Right)
	}

	y := x.Parent
	for y != t.NIL && x == y.Right {
		x = y
		y = y.Parent
	}
	return y
}

func (t *Rbtree) leftRotate(x *Node) {
	if x.Right == t.NIL {
		return
	}

	y := x.Right
	x.Right = y.Left
	if y.Left != t.NIL {
		y.Left.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == t.NIL {
		t.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Left = x
	x.Parent = y
}

func (t *Rbtree) rightRotate(x *Node) {
	if x.Left == t.NIL {
		return
	}

	y := x.Left
	x.Left = y.Right
	if y.Right != t.NIL {
		y.Right.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == t.NIL {
		t.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Right = x
	x.Parent = y
}
