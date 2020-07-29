package rbtree

type color uint8

const (
	RED   color = 0
	BLACK color = 1
)

type Node struct {
	Left   *Node
	Right  *Node
	Parent *Node
	Color  color

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
	node := &Node{Color: BLACK}
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
	z := &Node{t.NIL, t.NIL, y, RED, item}
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
	for z.Parent.Color == RED {
		if z.Parent == z.Parent.Parent.Left {
			y := z.Parent.Parent.Right
			if y.Color == RED {
				z.Parent.Color = BLACK
				y.Color = BLACK
				z.Parent.Parent.Color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right {
					z = z.Parent
					t.leftRotate(z)
				}
				z.Parent.Color = BLACK
				z.Parent.Parent.Color = RED
				t.rightRotate(z.Parent.Parent)
			}
		} else {
			y := z.Parent.Parent.Left
			if y.Color == RED {
				z.Parent.Color = BLACK
				y.Color = BLACK
				z.Parent.Parent.Color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Left {
					z = z.Parent
					t.rightRotate(z)
				}
				z.Parent.Color = BLACK
				z.Parent.Parent.Color = RED
				t.leftRotate(z.Parent.Parent)
			}
		}
	}
	t.root.Color = BLACK
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

	if y.Color == BLACK {
		t.deleteFixup(x)
	}

	t.count--
}

func (t *Rbtree) deleteFixup(x *Node) {
	for x != t.root && x.Color == BLACK {
		if x == x.Parent.Left {
			w := x.Parent.Right
			if w.Color == RED {
				w.Color = BLACK
				x.Parent.Color = RED
				t.leftRotate(x.Parent)
				w = x.Parent.Right
			}
			if w.Left.Color == BLACK && w.Right.Color == BLACK {
				w.Color = RED
				x = x.Parent
			} else {
				if w.Right.Color == BLACK {
					w.Left.Color = BLACK
					w.Color = RED
					t.rightRotate(w)
					w = x.Parent.Right
				}
				w.Color = x.Parent.Color
				x.Parent.Color = BLACK
				w.Right.Color = BLACK
				t.leftRotate(x.Parent)
				x = t.root
			}
		} else {
			w := x.Parent.Left
			if w.Color == RED {
				w.Color = BLACK
				x.Parent.Color = RED
				t.rightRotate(x.Parent)
				w = x.Parent.Left
			}
			if w.Left.Color == BLACK && w.Right.Color == BLACK {
				w.Color = RED
				x = x.Parent
			} else {
				if w.Left.Color == BLACK {
					w.Right.Color = BLACK
					w.Color = RED
					t.leftRotate(w)
					w = x.Parent.Left
				}
				w.Color = x.Parent.Color
				x.Parent.Color = BLACK
				w.Left.Color = BLACK
				t.rightRotate(x.Parent)
				x = t.root
			}
		}
	}
	x.Color = BLACK
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
