package skiplist

import (
	"math/rand"
	"sort"
	"testing"
)

type Int int

func (i Int) Less(other Interface) bool {
	return i < other.(Int)
}

func TestInt(t *testing.T) {
	sl := New()
	if sl.Len() != 0 || sl.Front() != nil && sl.Back() != nil {
		t.Fatal()
	}

	testData := []Int{Int(1), Int(2), Int(3)}
	sl.Insert(testData[0])
	if sl.Len() != 1 || sl.Front().Value.(Int) != testData[0] || sl.Back().Value.(Int) != testData[0] {
		t.Fatal()
	}
	sl.Insert(testData[2])
	if sl.Len() != 2 || sl.Front().Value.(Int) != testData[0] || sl.Back().Value.(Int) != testData[2] {
		t.Fatal()
	}
	sl.Insert(testData[1])
	if sl.Len() != 3 || sl.Front().Value.(Int) != testData[0] || sl.Back().Value.(Int) != testData[2] {
		t.Fatal()
	}

	sl.Insert(Int(-999))
	sl.Insert(Int(-888))
	sl.Insert(Int(888))
	sl.Insert(Int(999))
	sl.Insert(Int(1000))

	expect := []Int{Int(-999), Int(-888), Int(1), Int(2), Int(3), Int(888), Int(999), Int(1000)}
	ret := make([]Int, 0)
	for e := sl.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(Int))
	}
	for i := 0; i < len(ret); i++ {
		if ret[i] != expect[i] {
			t.Fatal()
		}
	}

	sl.Delete(Int(2))
	sl.Delete(Int(888))
	sl.Delete(Int(1000))

	expect = []Int{Int(-999), Int(-888), Int(1), Int(3), Int(999)}
	ret = make([]Int, 0)

	for e := sl.Back(); e != nil; e = e.Prev() {
		ret = append(ret, e.Value.(Int))
	}

	for i := 0; i < len(ret); i++ {
		if ret[i] != expect[len(ret)-i-1] {
			t.Fatal()
		}
	}

	if sl.Front().Value.(Int) != -999 {
		t.Fatal()
	}

	sl.Delete(sl.Front().Value)
	if sl.Front().Value.(Int) != -888 || sl.Back().Value.(Int) != 999 {
		t.Fatal()
	}

	sl.Delete(sl.Back().Value)
	if sl.Front().Value.(Int) != -888 || sl.Back().Value.(Int) != 3 {
		t.Fatal()
	}

	if e := sl.Insert(Int(2)); e.Value.(Int) != 2 {
		t.Fatal()
	}
	sl.Delete(Int(-888))

	if r := sl.Delete(Int(123)); r != nil {
		t.Fatal()
	}

	if sl.Len() != 3 {
		t.Fatal()
	}
}

func TestRank(t *testing.T) {
	sl := New()

	for i := 1; i <= 10; i++ {
		sl.Insert(Int(i))
	}
	for i := 1; i <= 10; i++ {
		if sl.GetRank(Int(i)) != i {
			t.Fatal()
		}
	}
	for i := 1; i <= 10; i++ {
		if sl.GetElementByRank(i).Value != Int(i) {
			t.Fatal()
		}
	}

	if sl.GetRank(Int(0)) != 0 || sl.GetRank(Int(11)) != 0 {
		t.Fatal()
	}
	if sl.GetElementByRank(11) != nil || sl.GetElementByRank(12) != nil {
		t.Fatal()
	}

	expect := []Int{Int(7), Int(8), Int(9), Int(10)}
	for e, i := sl.GetElementByRank(7), 0; e != nil; e, i = e.Next(), i+1 {
		if e.Value != expect[i] {
			t.Fatal()
		}
	}

	sl = New()
	mark := make(map[int]bool)
	ss := make([]int, 0)

	for i := 1; i <= 20000; i++ {
		x := rand.Int()
		if !mark[x] {
			mark[x] = true
			sl.Insert(Int(x))
			ss = append(ss, x)
		}
	}
	sort.Ints(ss)

	for i := 0; i < len(ss); i++ {
		if sl.GetElementByRank(i+1).Value != Int(ss[i]) || sl.GetRank(Int(ss[i])) != i+1 {
			t.Fatal()
		}
	}
}

func BenchmarkIntInsertOrder(b *testing.B) {
	sl := New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(Int(i))
	}
}

func BenchmarkIntInsertRandom(b *testing.B) {
	sl := New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(Int(rand.Int()))
	}
}

func BenchmarkIntDeleteOrder(b *testing.B) {
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Insert(Int(i))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Delete(Int(i))
	}
}

func BenchmarkIntDeleteRandom(b *testing.B) {
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Insert(Int(rand.Int()))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Delete(Int(rand.Int()))
	}
}

func BenchmarkIntRankOrder(b *testing.B) {
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Insert(Int(i))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.GetRank(Int(i))
	}
}

func BenchmarkIntRankRandom(b *testing.B) {
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Insert(Int(rand.Int()))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.GetRank(Int(rand.Int()))
	}
}
