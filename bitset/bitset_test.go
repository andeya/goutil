package bitset_test

import (
	"fmt"
	"testing"

	"github.com/henrylee2cn/goutil/bitset"
)

func TestBitSet(t *testing.T) {
	bitSet := bitset.New(1, 1, 1)
	t.Log(bitSet.Set(0, true))
	t.Log(bitSet.Set(1, true))
	t.Log(bitSet.Set(2, true))
	t.Log(bitSet.Set(3, true))
	t.Log(bitSet.Set(4, true))
	t.Log(bitSet.Set(5, true))
	t.Log(bitSet.Set(6, true))
	t.Log(bitSet.Set(7, true))
	t.Log(bitSet.Binary(" "))
	count := bitSet.Count(0, -1)
	if count != 10 {
		t.Fatalf("[0,-1] bit count: get %d, want %d", count, 10)
	}
	if !bitSet.Get(3) {
		t.Fatalf("bit in offset 3: get %d, want %d", 0, 1)
	}
	if !bitSet.Get(15) {
		t.Fatalf("bit in offset 15: get %d, want %d", 0, 1)
	}

	t.Log(bitSet.Set(500, true))
	count = bitSet.Count(0, 500)
	if count != 11 {
		t.Fatalf("1 bit count: get %d, want %d", count, 11)
	}
	if !bitSet.Get(500) {
		t.Fatalf("bit in offset 15: get %d, want %d", 0, 1)
	}
	t.Log(bitSet.Binary(" "))

	old, err := bitSet.Set(7, false)
	if err != nil {
		t.Fatal(err)
	}
	if !old {
		t.Fatalf("old bit in offset 7: get %d, want %d", 0, 1)
	}

	sub := bitSet.Sub(5, 20)
	t.Log(sub.Binary(" "))
	t.Log(sub)
	count = sub.Count(0, -1)
	if count != 3 {
		t.Fatalf("[0,-1] bit count: get %d, want %d", count, 3)
	}
	if !sub.Get(0) {
		t.Fatalf("bit in offset 0: get %d, want %d", 0, 1)
	}
	if sub.Get(2) {
		t.Fatalf("bit in offset 2: get %d, want %d", 1, 0)
	}

	if bitSet.Size() != (500/8+1)*8 {
		t.Fatalf("bitSet size: get %d, want %d", bitSet.Size(), (500/8+1)*8)
	}
	if sub.Size() != ((20-5)/8+1)*8 {
		t.Fatalf("sub size: get %d, want %d", sub.Size(), ((20-5)/8+1)*8)
	}

	sub.Clear()
	t.Log("cleared:", sub.Binary(" "))
	count = sub.Count(0, -1)
	if count != 0 {
		t.Fatalf("sub size: get %d, want %d", count, 0)
	}
}

func ExampleBitSet() {
	bs, _ := bitset.NewFromHex("c020")
	fmt.Println("Origin:", bs.Binary(" "))
	not := bs.Not()
	fmt.Println("Not:", not.Binary(" "))
	fmt.Println("AndNot:", not.AndNot(bitset.New(1, 1)).Binary(" "))
	fmt.Println("And:", not.And(bitset.New(1<<1, 1<<1)).Binary(" "))
	fmt.Println("Or:", not.Or(bitset.New(1<<7, 1<<7)).Binary(" "))
	fmt.Println("Xor:", not.Xor(bitset.New(1<<7, 1<<7)).Binary(" "))

	not.Range(func(k int, v bool) bool {
		fmt.Println(v)
		return true
	})

	// Output:
	// Origin: 11000000 00100000
	// Not: 00111111 11011111
	// AndNot: 00111110 11011110
	// And: 00000010 00000010
	// Or: 10111111 11011111
	// Xor: 10111111 01011111
	// false
	// false
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// false
	// true
	// true
	// true
	// true
	// true
}
