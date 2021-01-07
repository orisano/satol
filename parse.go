// Copyright 2020-2021 Nao Yonashiro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package satol

import (
	"encoding/binary"
	"reflect"
	"unsafe"

	"github.com/klauspost/cpuid"
)

var hasAVX bool
var hasSSE42 bool

type UintParser interface {
	Parse(s string) uint64
}

var Parser UintParser = GoParser{}

func init() {
	hasAVX = cpuid.CPU.AVX()
	hasSSE42 = cpuid.CPU.SSE42()

	if hasAVX {
		Parser = AVXParser{}
	} else if hasSSE42 {
		Parser = SSEParser{}
	}
}

type AVXParser struct{}

func (AVXParser) Available() bool {
	return hasAVX
}

func (AVXParser) Parse(s string) uint64 {
	sum := uint64(0)
	if len(s) >= 16 {
		sum = Parse16AVX(s)
		s = s[16:]
	}
	if len(s) >= 8 {
		sum *= 100000000
		sum += Parse8AVX(s)
		s = s[8:]
	}
	if len(s) >= 4 {
		sum *= 10000
		sum += Parse4GoUnrolled(s)
		s = s[4:]
	}
	if len(s) == 0 {
		return sum
	} else {
		return parse3GoUnrolled(sum, s)
	}
}

type SSEParser struct{}

func (SSEParser) Available() bool {
	return hasSSE42
}

func (SSEParser) Parse(s string) uint64 {
	sum := uint64(0)
	if len(s) >= 16 {
		sum = Parse16SSE(s)
		s = s[16:]
	}
	if len(s) >= 8 {
		sum *= 100000000
		sum += Parse8SSE(s)
		s = s[8:]
	}
	if len(s) >= 4 {
		sum *= 10000
		sum += Parse4GoUnrolled(s)
		s = s[4:]
	}
	if len(s) == 0 {
		return sum
	} else {
		return parse3GoUnrolled(sum, s)
	}
}

type GoParser struct{}

func (GoParser) Available() bool {
	return true
}

func (GoParser) Parse(s string) uint64 {
	sum := uint64(0)
	if len(s) >= 16 {
		sum = Parse16GoBits(s)
		s = s[16:]
	}
	if len(s) >= 8 {
		sum *= 100000000
		sum += Parse8GoBits(s)
		s = s[8:]
	}
	if len(s) >= 4 {
		sum *= 10000
		sum += Parse4GoUnrolled(s)
		s = s[4:]
	}
	if len(s) == 0 {
		return sum
	} else {
		return parse3GoUnrolled(sum, s)
	}
}

func Parse(s string) uint64 {
	return Parser.Parse(s)
}

func ParseNaive(s string) uint64 {
	sum := uint64(0)
	for _, c := range s {
		sum = sum*10 + uint64(c-'0')
	}
	return sum
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

func parse3GoUnrolled(sum uint64, s string) uint64 {
	sum = sum*10 + uint64(s[0]-'0')
	if len(s) == 1 {
		return sum
	}
	sum = sum*10 + uint64(s[1]-'0')
	if len(s) == 2 {
		return sum
	}
	return sum*10 + uint64(s[2]-'0')
}

func s2b(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}
