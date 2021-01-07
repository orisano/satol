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
	"math/rand"
	"strconv"
	"testing"
)

func TestParse16(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			args: args{
				s: "0000000000000000",
			},
			want: 0,
		},
		{
			args: args{
				s: "0000000000000001",
			},
			want: 1,
		},
		{
			args: args{
				s: "0000000000000002",
			},
			want: 2,
		},
		{
			args: args{
				s: "9999999999999999",
			},
			want: 9999999999999999,
		},
		{
			args: args{
				s: "8888888888888888",
			},
			want: 8888888888888888,
		},
		{
			args: args{
				s: "7777777777777777",
			},
			want: 7777777777777777,
		},
		{
			args: args{
				s: "1234567890123456",
			},
			want: 1234567890123456,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse16GoBits(tt.args.s); got != tt.want {
				t.Errorf("Parse16GoBits() = %d(%x), want %d", got, got, tt.want)
			}
			if got := Parse16GoUnrolled(tt.args.s); got != tt.want {
				t.Errorf("Parse16GoUnrolled() = %d(%x), want %d", got, got, tt.want)
			}
			if hasAVX {
				if got := Parse16AVX(tt.args.s); got != tt.want {
					t.Errorf("Parse16AVX() = %d(%x), want %d", got, got, tt.want)
				}
			}
			if hasSSE42 {
				if got := Parse16SSE(tt.args.s); got != tt.want {
					t.Errorf("Parse16SSE() = %d(%x), want %d", got, got, tt.want)
				}
			}
		})
	}
}

func TestParse8(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			args: args{
				s: "00000000",
			},
			want: 0,
		},
		{
			args: args{
				s: "00000001",
			},
			want: 1,
		},
		{
			args: args{
				s: "00000002",
			},
			want: 2,
		},
		{
			args: args{
				s: "99999999",
			},
			want: 99999999,
		},
		{
			args: args{
				s: "88888888",
			},
			want: 88888888,
		},
		{
			args: args{
				s: "77777777",
			},
			want: 77777777,
		},
		{
			args: args{
				s: "12345678",
			},
			want: 12345678,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse8GoBits(tt.args.s); got != tt.want {
				t.Errorf("Parse8GoBits() = %d(%x), want %d", got, got, tt.want)
			}
			if got := Parse8GoUnrolled(tt.args.s); got != tt.want {
				t.Errorf("Parse8GoUnrolled() = %d(%x), want %d", got, got, tt.want)
			}
			if hasAVX {
				if got := Parse8AVX(tt.args.s); got != tt.want {
					t.Errorf("Parse8AVX() = %d(%x), want %d", got, got, tt.want)
				}
			}
			if hasSSE42 {
				if got := Parse8SSE(tt.args.s); got != tt.want {
					t.Errorf("Parse8SSE() = %d(%x), want %d", got, got, tt.want)
				}
			}
		})
	}
}

func TestParse4(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			args: args{
				s: "0000",
			},
			want: 0,
		},
		{
			args: args{
				s: "0001",
			},
			want: 1,
		},
		{
			args: args{
				s: "0002",
			},
			want: 2,
		},
		{
			args: args{
				s: "9999",
			},
			want: 9999,
		},
		{
			args: args{
				s: "8888",
			},
			want: 8888,
		},
		{
			args: args{
				s: "7777",
			},
			want: 7777,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse4GoBits(tt.args.s); got != tt.want {
				t.Errorf("Parse4GoBits() = %d(%x), want %d", got, got, tt.want)
			}
			if got := Parse4GoUnrolled(tt.args.s); got != tt.want {
				t.Errorf("Parse4GoUnrolled() = %d(%x), want %d", got, got, tt.want)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	rng := rand.NewSource(42)
	data := make([]string, 1000000)
	for i := range data {
		rand.Uint64()
		data[i] = strconv.Itoa(int(rng.Int63()))
	}
	b.ResetTimer()
	result := uint64(0)
	b.Run("AVX", func(b *testing.B) {
		var p AVXParser
		for i := 0; i < b.N; i++ {
			for _, s := range data {
				result = p.Parse(s)
			}
		}
	})
	b.Run("SSE", func(b *testing.B) {
		var p SSEParser
		for i := 0; i < b.N; i++ {
			for _, s := range data {
				result = p.Parse(s)
			}
		}
	})
	b.Run("Go", func(b *testing.B) {
		var p GoParser
		for i := 0; i < b.N; i++ {
			for _, s := range data {
				result = p.Parse(s)
			}
		}
	})
	b.Run("Naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, s := range data {
				result = ParseNaive(s)
			}
		}
	})
	b.Run("ParseUint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, s := range data {
				result, _ = strconv.ParseUint(s, 10, 64)
			}
		}
	})
	result++
}

func BenchmarkParse16(b *testing.B) {
	result := uint64(0)
	if hasAVX {
		b.Run("AVX", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result = Parse16AVX("9999999999999999")
			}
		})
	}
	if hasSSE42 {
		b.Run("SSE42", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result = Parse16SSE("9999999999999999")
			}
		})
	}
	b.Run("GoBits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result = Parse16GoBits("9999999999999999")
		}
	})
	b.Run("GoUnrolled", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result = Parse16GoUnrolled("9999999999999999")
		}
	})
	b.Run("Naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result = ParseNaive("9999999999999999")
		}
	})
	b.Run("ParseUint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x, _ := strconv.ParseUint("9999999999999999", 10, 64)
			result = x
		}
	})
	result++
}

func BenchmarkParse8(b *testing.B) {
	result := uint64(0)
	if hasAVX {
		b.Run("AVX", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result = Parse8AVX("99999999")
			}
		})
	}
	if hasSSE42 {
		b.Run("SSE42", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result = Parse8SSE("99999999")
			}
		})
	}
	b.Run("GoBits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result = Parse8GoBits("99999999")
		}
	})
	b.Run("GoUnrolled", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result = Parse8GoUnrolled("99999999")
		}
	})
	b.Run("Naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result = ParseNaive("99999999")
		}
	})
	b.Run("ParseUint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x, _ := strconv.ParseUint("99999999", 10, 64)
			result = x
		}
	})
	result++
}

func BenchmarkParse4(b *testing.B) {
	result := uint64(0)
	b.Run("GoBits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result = Parse4GoBits("9999")
		}
	})
	b.Run("GoUnrolled", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result = Parse4GoUnrolled("9999")
		}
	})
	b.Run("Naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result = ParseNaive("9999")
		}
	})
	b.Run("ParseUint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x, _ := strconv.ParseUint("9999", 10, 64)
			result = x
		}
	})
	result++
}

func BenchmarkParse3(b *testing.B) {
	result := uint64(0)
	b.Run("GoUnrolled", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result = parse3GoUnrolled(0, "9")
			result = parse3GoUnrolled(0, "99")
			result = parse3GoUnrolled(0, "999")
		}
	})
	b.Run("Naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result = ParseNaive("9")
			result = ParseNaive("99")
			result = ParseNaive("999")
		}
	})
	b.Run("ParseUint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result, _ = strconv.ParseUint("9", 10, 64)
			result, _ = strconv.ParseUint("99", 10, 64)
			result, _ = strconv.ParseUint("999", 10, 64)
		}
	})
	result++
}
