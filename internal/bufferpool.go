package internal

import (
	"math"
	"sync"
)

const (
	minByte = uint32(1 << 7) //128 b
	capGrow = uint32(1)
	maxByte = 1 << 24 //16M
)

var BUFFERPOOL *BufferPool

func init() {
	rawDataCap := minByte
	bp := &BufferPool{
		capSlice: make([]uint32, 0),
		pool:     map[uint32]*sync.Pool{},
	}
	for rawDataCap <= maxByte {
		bp.capSlice = append(bp.capSlice, rawDataCap)
		tmpCap := rawDataCap
		bp.pool[tmpCap] = &sync.Pool{
			New: func() interface{} {
				slice := make([]byte, tmpCap)
				return &slice
			},
		}
		//fmt.Println(rawDataCap)
		rawDataCap <<= capGrow
	}
	BUFFERPOOL = bp
}

type BufferPool struct {
	capSlice []uint32
	pool     map[uint32]*sync.Pool
}

func (b *BufferPool) Put(bs []byte) {
	p := b.pool[uint32(cap(bs))]
	if p != nil {
		p.Put(&bs)
	}
}
func (b *BufferPool) Get(l uint32) *[]byte {
	offset := int(math.Log2(float64(b.capSlice[0])))

	index := int(math.Ceil(math.Log2(float64(l))))
	//fmt.Println(index)
	if index <= offset {
		index = 0
	} else {
		index = index - offset
	}
	if index > len(b.capSlice)-1 {
		index = len(b.capSlice) - 1
	}
	if index < 0 {
		index = 0
	}
	capSize := b.capSlice[index]
	p := b.pool[capSize]
	return p.Get().(*[]byte)
}
