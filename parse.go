package satol

import (
	"encoding/binary"
	"reflect"
	"unsafe"

	"github.com/klauspost/cpuid"
)

var hasAVX bool
var hasSSE42 bool

func init() {
	hasAVX = cpuid.CPU.AVX()
	hasSSE42 = cpuid.CPU.SSE42()
}

func Parse16(s string) uint64 {
	switch {
	case hasAVX:
		return Parse16AVX(s)
	case hasSSE42:
		return Parse16SSE(s)
	default:
		return Parse16GoBits(s)
	}
}

func Parse8(s string) uint64 {
	switch {
	case hasAVX:
		return Parse8AVX(s)
	default:
		return Parse8GoBits(s)
	}
}

func Parse16GoBits(s string) uint64 {
	b := s2b(s)[:16]
	high := binary.LittleEndian.Uint64(b[0:8])
	low := binary.LittleEndian.Uint64(b[8:16])
	return convert8Go(high)*100000000 + convert8Go(low)
}

func Parse16GoUnrolled(s string) uint64 {
	sum := uint64(0)
	sum += uint64(s[15]-'0') * 1
	sum += uint64(s[14]-'0') * 10
	sum += uint64(s[13]-'0') * 100
	sum += uint64(s[12]-'0') * 1000
	sum += uint64(s[11]-'0') * 10000
	sum += uint64(s[10]-'0') * 100000
	sum += uint64(s[9]-'0') * 1000000
	sum += uint64(s[8]-'0') * 10000000
	sum += uint64(s[7]-'0') * 100000000
	sum += uint64(s[6]-'0') * 1000000000
	sum += uint64(s[5]-'0') * 10000000000
	sum += uint64(s[4]-'0') * 100000000000
	sum += uint64(s[3]-'0') * 1000000000000
	sum += uint64(s[2]-'0') * 10000000000000
	sum += uint64(s[1]-'0') * 100000000000000
	sum += uint64(s[0]-'0') * 1000000000000000
	return sum
}

func Parse8GoBits(s string) uint64 {
	b := s2b(s)[:8]
	chunk := binary.LittleEndian.Uint64(b[:8])
	return convert8Go(chunk)
}

func Parse8GoUnrolled(s string) uint64 {
	sum := uint64(0)
	sum += uint64(s[7]-'0') * 1
	sum += uint64(s[6]-'0') * 10
	sum += uint64(s[5]-'0') * 100
	sum += uint64(s[4]-'0') * 1000
	sum += uint64(s[3]-'0') * 10000
	sum += uint64(s[2]-'0') * 100000
	sum += uint64(s[1]-'0') * 1000000
	sum += uint64(s[0]-'0') * 10000000
	return sum
}

func convert8Go(chunk uint64) uint64 {
	lowerDigits := (chunk & 0x0f000f000f000f00) >> 8
	upperDigits := (chunk & 0x000f000f000f000f) * 10
	chunk = lowerDigits + upperDigits

	lowerDigits = (chunk & 0x00ff000000ff0000) >> 16
	upperDigits = (chunk & 0x000000ff000000ff) * 100
	chunk = lowerDigits + upperDigits

	lowerDigits = (chunk & 0x0000ffff00000000) >> 32
	upperDigits = (chunk & 0x000000000000ffff) * 10000
	return lowerDigits + upperDigits
}

func Parse4GoBits(s string) uint64 {
	b := s2b(s)[:4]
	chunk := binary.LittleEndian.Uint32(b[0:4])
	return convert4Go(chunk)
}

func Parse4GoUnrolled(s string) uint64 {
	sum := uint64(0)
	sum += uint64(s[3]-'0') * 1
	sum += uint64(s[2]-'0') * 10
	sum += uint64(s[1]-'0') * 100
	sum += uint64(s[0]-'0') * 1000
	return sum
}

func convert4Go(chunk uint32) uint64 {
	lowerDigits := (chunk & 0x0f000f00) >> 8
	upperDigits := (chunk & 0x000f000f) * 10
	chunk = lowerDigits + upperDigits

	lowerDigits = (chunk & 0x00ff0000) >> 16
	upperDigits = (chunk & 0x000000ff) * 100
	return uint64(lowerDigits + upperDigits)
}

func s2b(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}
