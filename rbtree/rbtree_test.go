package rbtree

import (
	"math/rand"
	"testing"
)

type Int int

func (i Int) Less(other Item) bool {
	return i < other.(Int)
}

func BenchmarkIntInsertOrder(b *testing.B) {
	t := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Insert(Int(i))
	}
}

func BenchmarkIntInsertRandom(b *testing.B) {
	t := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Insert(Int(rand.Int()))
	}
}

func BenchmarkIntDeleteOrder(b *testing.B) {
	t := New()
	for i := 0; i < 1000000; i++ {
		t.Insert(Int(i))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t.Delete(Int(i))
	}
}

func BenchmarkIntDeleteRandome(b *testing.B) {
	t := New()
	for i := 0; i < 1000000; i++ {
		t.Insert(Int(rand.Int()))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t.Delete(Int(rand.Int()))
	}
}
