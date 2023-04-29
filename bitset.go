package timtab

import (
	"bytes"
	"strings"
)

/*
0 		-> 	0
1..8	->  1
9..16	->  2
*/

func NewBitset(size int) Bitset {
	if size <= 0 {
		return nil
	}
	return make(Bitset, (size-1)/8+1)
}

type Bitset []byte

func (bs Bitset) Get(i int) bool {
	return bs[i/8]&(byte(1)<<(i%8)) != 0
}

func (bs Bitset) Set(i int, v bool) {
	if v {
		bs[i/8] |= (byte(1) << (i % 8))
	} else {
		bs[i/8] &= ^(^bs[i/8] | (byte(1) << (i % 8)))
	}
}

func (bs Bitset) String() string {
	var sl []string
	for i := 0; i < len(bs)*8; i++ {
		if bs.Get(i) {
			sl = append(sl, "1")
		} else {
			sl = append(sl, "0")
		}
	}
	return strings.Join(sl, "")
}

func (bs Bitset) Clone() Bitset {
	return bytes.Clone(bs)
}
