package bitset

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// BitSet bit set
type BitSet struct {
	set []byte
	mu  sync.RWMutex
}

// NewBitSet creates a bit set object.
func NewBitSet(init ...byte) *BitSet {
	return &BitSet{set: init}
}

// Set sets the bit bool value on the specified offset,
// and returns the value of before setting.
// Notes:
//  0 means the 1st bit, -1 means the bottom 1th bit, -2 means the bottom 2th bit and so on;
//  If offset>=len(b.set), automatically grow the bit set;
//  If the bit offset is out of the left range, returns error.
func (b *BitSet) Set(offset int, value bool) (bool, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	size := b.size()
	// 0 means the 1st bit, -1 means the bottom 1th bit,
	// -2 means the bottom 2th bit and so on.
	if offset < 0 {
		offset += size
	}
	if offset < 0 {
		return false, errors.New("the bit offset is out of the left range")
	}

	// the bit group index
	gi := offset / 8
	// the bit index of in the group
	bi := offset % 8

	// if the bit offset is out of the right range, automatically grow.
	if gi >= len(b.set) {
		newSet := make([]byte, gi+1)
		copy(newSet, b.set)
		b.set = newSet
	}

	gb := b.set[gi]
	rOff := byte(7 - bi)
	var mask byte = 1 << rOff
	oldVal := gb & mask >> rOff
	if (oldVal == 1) != value {
		if oldVal == 1 {
			b.set[gi] = gb &^ mask
		} else {
			b.set[gi] = gb | mask
		}
	}
	return oldVal == 1, nil
}

// Get gets the bit bool value on the specified offset.
// Notes:
//  0 means the 1st bit, -1 means the bottom 1th bit, -2 means the bottom 2th bit and so on;
//  If offset>=len(b.set), returns false.
func (b *BitSet) Get(offset int) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	size := b.size()
	// 0 means the 1st bit, -1 means the bottom 1th bit,
	// -2 means the bottom 2th bit and so on.
	if offset < 0 {
		offset += size
	}
	if offset < 0 || offset >= size {
		return false
	}
	return getBit(b.set[offset/8], byte(offset%8)) == 1
}

// Count counts the amount of bit set to 1 within the specified range of the bit set.
// Notes:
//  0 means the 1st bit, -1 means the bottom 1th bit, -2 means the bottom 2th bit and so on.
func (b *BitSet) Count(start, end int) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	sgi, sbi, egi, ebi, valid := b.validRange(start, end)
	if !valid {
		return 0
	}
	var n int
	gb := b.set[sgi]
	for i := byte(sbi); i < 8; i++ {
		if getBit(gb, i) == 1 {
			n++
		}
	}
	for h := sgi + 1; h < egi; h++ {
		gb = b.set[h]
		for i := byte(0); i < 8; i++ {
			if getBit(gb, i) == 1 {
				n++
			}
		}
	}
	gb = b.set[egi]
	for i := uint(0); i <= ebi; i++ {
		if getBit(gb, byte(i)) == 1 {
			n++
		}
	}
	return n
}

func getBit(gb, offset byte) byte {
	var rOff = 7 - offset
	var mask byte = 1 << rOff
	return gb & mask >> rOff
}

func (b *BitSet) validRange(start, end int) (sgi, sbi, egi, ebi uint, valid bool) {
	size := b.size()
	if start < 0 {
		start += size
	}
	if start >= size {
		return
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += size
	}
	if end >= size {
		end = size - 1
	}
	if start > end {
		return
	}
	valid = true
	sgi, sbi = uint(start/8), uint(start%8)
	egi, ebi = uint(end/8), uint(end%8)
	return
}

// Size returns the bits size.
func (b *BitSet) Size() int {
	b.mu.RLock()
	size := b.size()
	b.mu.RUnlock()
	return size
}

func (b *BitSet) size() int {
	size := len(b.set) * 8
	if size/8 != len(b.set) {
		panic("overflows when calculating the bit set size")
	}
	return size
}

// Bytes returns the bit set copy bytes.
func (b *BitSet) Bytes() []byte {
	set := make([]byte, len(b.set))
	copy(set, b.set)
	return set
}

// Binary returns the bit binary by hex type.
// Notes:
//  Paramter sep is the separator between bytes.
func (b *BitSet) Binary(sep string) string {
	if len(b.set) == 0 {
		return ""
	}
	var s strings.Builder
	for _, i := range b.set {
		s.WriteString(fmt.Sprintf("%s%08b", sep, i))
	}
	return strings.TrimPrefix(s.String(), sep)
}

// String returns the bit set by hex type.
func (b *BitSet) String() string {
	return hex.EncodeToString(b.set)
}

// Sub returns the bit subset within the specified range of the bit set.
// Notes:
//  0 means the 1st bit, -1 means the bottom 1th bit, -2 means the bottom 2th bit and so on.
func (b *BitSet) Sub(start, end int) *BitSet {
	newBitSet := &BitSet{
		set: make([]byte, 0, len(b.set)),
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	sgi, sbi, egi, ebi, valid := b.validRange(start, end)
	if !valid {
		return newBitSet
	}
	pre := b.set[sgi] << sbi
	for _, v := range b.set[sgi+1 : egi] {
		newBitSet.set = append(newBitSet.set, pre|v>>(7-sbi))
		pre = v << sbi
	}
	last := b.set[egi] >> (7 - ebi) << ebi
	newBitSet.set = append(newBitSet.set, pre|last>>(7-sbi), last<<sbi)
	return newBitSet
}
