package skiplist

import (
	"math/rand"
)

const (
	kSkipListMaxLevel = 32
	kSkipListP        = 0.25
)

type Interface interface {
	Less(other Interface) bool
}

type skiplistLevel struct {
	forward *Element
	span    int
}

type Element struct {
	Value    Interface
	backward *Element
	level    []skiplistLevel
}

func newElement(level int, v Interface) *Element {
	return &Element{
		Value:    v,
		backward: nil,
		level:    make([]skiplistLevel, level),
	}
}

func (e *Element) Next() *Element {
	return e.level[0].forward
}

func (e *Element) Prev() *Element {
	return e.backward
}

type SkipList struct {
	header *Element
	tail   *Element
	length int
	level  int
}

func New() *SkipList {
	sl := &SkipList{
		header: newElement(kSkipListMaxLevel, nil),
		tail:   nil,
		length: 0,
		level:  1,
	}
	return sl
}

func (sl *SkipList) Front() *Element {
	return sl.header.level[0].forward
}

func (sl *SkipList) Back() *Element {
	return sl.tail
}

func (sl *SkipList) Len() int {
	return sl.length
}

func randomLevel() int {
	level := 1
	for rand.Int()&0xFFFF < kSkipListP*(0xFFFF+1) {
		level++
	}
	if level < kSkipListMaxLevel {
		return level
	}
	return kSkipListMaxLevel
}

func (sl *SkipList) Insert(v Interface) *Element {
	var update [kSkipListMaxLevel]*Element
	var rank [kSkipListMaxLevel]int

	e := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		if i != sl.level-1 {
			rank[i] = rank[i+1]
		}
		for e.level[i].forward != nil && e.level[i].forward.Value.Less(v) {
			rank[i] += e.level[i].span
			e = e.level[i].forward
		}
		update[i] = e
	}

	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.header
			update[i].level[i].span = sl.length
		}
		sl.level = level
	}

	e = newElement(level, v)
	for i := 0; i < level; i++ {
		e.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = e

		e.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	for i := level; i < sl.level; i++ {
		update[i].level[i].span++
	}
	if update[0] != sl.header {
		e.backward = update[0]
	}
	if e.level[0].forward != nil {
		e.level[0].forward.backward = e
	} else {
		sl.tail = e
	}
	sl.length++
	return e
}

func (sl *SkipList) deleteElement(e *Element, update []*Element) {
	for i := 0; i < sl.level; i++ {
		if update[i].level[i].forward == e {
			update[i].level[i].span += e.level[i].span - 1
			update[i].level[i].forward = e.level[i].forward
		} else {
			update[i].level[i].span -= 1
		}
	}

	if e.level[0].forward != nil {
		e.level[0].forward.backward = e.backward
	} else {
		sl.tail = e.backward
	}

	for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
		sl.level--
	}
	sl.length--
}

func (sl *SkipList) Delete(v Interface) *Element {
	var update [kSkipListMaxLevel]*Element
	e := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for e.level[i].forward != nil && e.level[i].forward.Value.Less(v) {
			e = e.level[i].forward
		}
		update[i] = e
	}

	e = e.level[0].forward
	if e != nil && !v.Less(e.Value) {
		sl.deleteElement(e, update[:])
		return e
	}
	return nil // not found
}

func (sl *SkipList) GetRank(v Interface) int {
	e := sl.header
	rank := 0
	for i := sl.level - 1; i >= 0; i-- {
		for e.level[i].forward != nil && !v.Less(e.level[i].forward.Value) {
			rank += e.level[i].span
			e = e.level[i].forward
		}
		if e.Value != nil && !e.Value.Less(v) {
			return rank
		}
	}
	return 0
}

func (sl *SkipList) GetElementByRank(rank int) *Element {
	e := sl.header
	traversed := 0
	for i := sl.level - 1; i >= 0; i-- {
		for e.level[i].forward != nil && traversed+e.level[i].span <= rank {
			traversed += e.level[i].span
			e = e.level[i].forward
		}
		if traversed == rank {
			return e
		}
	}
	return nil
}
