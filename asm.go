// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/buildtags"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func main() {
	Constraint(buildtags.Term("amd64"))
	zeros := GLOBL("zeros", RODATA|NOPTR)
	DATA(0, U64(0xcfcfcfcfcfcfcfd0))
	DATA(8, U64(0xcfcfcfcfcfcfcfd0))

	mul8 := GLOBL("mul8", RODATA|NOPTR)
	DATA(0, U32(0x010a010a))
	DATA(4, U32(0x010a010a))
	DATA(8, U32(0x010a010a))
	DATA(12, U32(0x010a010a))

	mul16 := GLOBL("mul16", RODATA|NOPTR)
	DATA(0, U32(0x00010064))
	DATA(4, U32(0x00010064))
	DATA(8, U32(0x00010064))
	DATA(12, U32(0x00010064))

	mul16Low := GLOBL("mul16Low", RODATA|NOPTR)
	DATA(0, U16(10000))
	DATA(2, U16(1))
	DATA(4, U16(10000))
	DATA(6, U16(1))
	DATA(8, U16(0))
	DATA(10, U16(0))
	DATA(12, U16(0))
	DATA(14, U16(0))

	TEXT("Parse16SSE", NOSPLIT, "func(s string) uint64")
	Load(Param("s").Base(), RCX)

	LDDQU(Mem{Base: RCX}, X0)
	PADDQ(zeros, X0)
	PMADDUBSW(mul8, X0)
	PMADDWL(mul16, X0)
	PACKUSDW(X0, X0)
	PMADDWL(mul16Low, X0)
	MOVQ(X0, RAX)
	MOVQ(RAX, RCX)
	ANDQ(U32(0x7fffffff), RCX)
	IMUL3Q(U32(100000000), RCX, RCX)
	SARQ(U8(32), RAX)
	ADDQ(RCX, RAX)
	Store(RAX, ReturnIndex(0))
	RET()

	TEXT("Parse16AVX", NOSPLIT, "func(s string) uint64")
	Load(Param("s").Base(), RCX)

	VLDDQU(Mem{Base: RCX}, X0)
	VPADDQ(zeros, X0, X0)
	VPMADDUBSW(mul8, X0, X0)
	VPMADDWD(mul16, X0, X0)
	VPACKUSDW(X0, X0, X0)
	VPMADDWD(mul16Low, X0, X0)
	VMOVQ(X0, RAX)
	MOVQ(RAX, RCX)
	ANDQ(U32(0x7FFFFFFF), RCX)
	IMUL3Q(U32(100000000), RCX, RCX)
	SARQ(U8(32), RAX)
	ADDQ(RCX, RAX)
	Store(RAX, ReturnIndex(0))
	RET()

	TEXT("Parse8SSE", NOSPLIT, "func(s string) uint64")
	Load(Param("s").Base(), RCX)

	MOVQ(Mem{Base: RCX, Scale: 8}, X0)
	PADDQ(zeros, X0)
	PSHUFD(U8(68), X0, X0)
	PMADDUBSW(mul8, X0)
	PMADDWL(mul16, X0)
	PACKUSDW(X0, X0)
	PMADDWL(mul16Low, X0)
	MOVQ(X0, RAX)
	SARQ(U8(32), RAX)
	Store(RAX, ReturnIndex(0))
	RET()

	TEXT("Parse8AVX", NOSPLIT, "func(s string) uint64")
	Load(Param("s").Base(), RCX)

	VMOVQ(Mem{Base: RCX, Scale: 8}, X0)
	VPADDQ(zeros, X0, X0)
	VPSHUFD(U8(68), X0, X0)
	VPMADDUBSW(mul8, X0, X0)
	VPMADDWD(mul16, X0, X0)
	VPACKUSDW(X0, X0, X0)
	VPMADDWD(mul16Low, X0, X0)
	VMOVQ(X0, RAX)
	SARQ(U8(32), RAX)
	Store(RAX, ReturnIndex(0))
	RET()

	Generate()
}
