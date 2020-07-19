package satol

import (
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
			if got := Parse16(tt.args.s); got != tt.want {
				t.Errorf("Parse16() = %d(%x), want %d", got, got, tt.want)
			}
			if got := Parse16GoBits(tt.args.s); got != tt.want {
				t.Errorf("Parse16GoBits() = %d(%x), want %d", got, got, tt.want)
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
			if got := Parse8(tt.args.s); got != tt.want {
				t.Errorf("Parse8() = %d(%x), want %d", got, got, tt.want)
			}
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

func BenchmarkParse16(b *testing.B) {
	if hasAVX {
		b.Run("AVX", func(b *testing.B) {
			sum := uint64(0)
			for i := 0; i < b.N; i++ {
				sum += Parse16AVX("9999999999999999")
			}
		})
	}
	if hasSSE42 {
		b.Run("SSE42", func(b *testing.B) {
			sum := uint64(0)
			for i := 0; i < b.N; i++ {
				sum += Parse16SSE("9999999999999999")
			}
		})
	}
	b.Run("GoBits", func(b *testing.B) {
		sum := uint64(0)
		for i := 0; i < b.N; i++ {
			sum += Parse16GoBits("9999999999999999")
		}
	})
	b.Run("Auto", func(b *testing.B) {
		sum := uint64(0)
		for i := 0; i < b.N; i++ {
			sum += Parse16("9999999999999999")
		}
	})
	b.Run("ParseUint", func(b *testing.B) {
		sum := uint64(0)
		for i := 0; i < b.N; i++ {
			x, _ := strconv.ParseUint("9999999999999999", 10, 64)
			sum += x
		}
	})
}

func BenchmarkParse8(b *testing.B) {
	if hasAVX {
		b.Run("AVX", func(b *testing.B) {
			sum := uint64(0)
			for i := 0; i < b.N; i++ {
				sum += Parse8AVX("99999999")
			}
		})
	}
	if hasSSE42 {
		b.Run("SSE42", func(b *testing.B) {
			sum := uint64(0)
			for i := 0; i < b.N; i++ {
				sum += Parse8SSE("99999999")
			}
		})
	}
	b.Run("GoBits", func(b *testing.B) {
		sum := uint64(0)
		for i := 0; i < b.N; i++ {
			sum += Parse8GoBits("99999999")
		}
	})
	b.Run("GoUnrolled", func(b *testing.B) {
		sum := uint64(0)
		for i := 0; i < b.N; i++ {
			sum += Parse8GoUnrolled("99999999")
		}
	})
	b.Run("Auto", func(b *testing.B) {
		sum := uint64(0)
		for i := 0; i < b.N; i++ {
			sum += Parse8("99999999")
		}
	})
	b.Run("ParseUint", func(b *testing.B) {
		sum := uint64(0)
		for i := 0; i < b.N; i++ {
			x, _ := strconv.ParseUint("99999999", 10, 64)
			sum += x
		}
	})
}

func BenchmarkParse4(b *testing.B) {
	b.Run("GoBits", func(b *testing.B) {
		sum := uint64(0)
		for i := 0; i < b.N; i++ {
			sum += Parse4GoBits("9999")
		}
	})
	b.Run("GoUnrolled", func(b *testing.B) {
		sum := uint64(0)
		for i := 0; i < b.N; i++ {
			sum += Parse4GoUnrolled("9999")
		}
	})
	b.Run("ParseUint", func(b *testing.B) {
		sum := uint64(0)
		for i := 0; i < b.N; i++ {
			x, _ := strconv.ParseUint("9999", 10, 64)
			sum += x
		}
	})
}

