# satol
![Test](https://github.com/orisano/satol/workflows/Test/badge.svg)

satol is a fast integer parser.
inspired by this article. https://kholdstare.github.io/technical/2020/05/26/faster-integer-parsing.html

## Installation
```bash
go get github.com/orisano/satol
```

## Benchmark
```
goos: darwin
goarch: amd64
pkg: github.com/orisano/satol
BenchmarkParse/AVX-8 	     147	   7831060 ns/op
BenchmarkParse/SSE-8 	     148	   7978048 ns/op
BenchmarkParse/Go-8  	      99	  11176608 ns/op
BenchmarkParse/Naive-8         	      67	  16171480 ns/op
BenchmarkParse/ParseUint-8     	      32	  35728000 ns/op
BenchmarkParse16/AVX-8         	536560414	         2.24 ns/op
BenchmarkParse16/SSE42-8       	449927114	         2.90 ns/op
BenchmarkParse16/GoBits-8      	215787661	         5.68 ns/op
BenchmarkParse16/GoUnrolled-8  	142361196	         8.34 ns/op
BenchmarkParse16/Naive-8       	100000000	        12.0 ns/op
BenchmarkParse16/ParseUint-8   	38155237	        30.8 ns/op
BenchmarkParse8/AVX-8          	564110046	         2.14 ns/op
BenchmarkParse8/SSE42-8        	432063027	         2.94 ns/op
BenchmarkParse8/GoBits-8       	332060721	         3.72 ns/op
BenchmarkParse8/GoUnrolled-8   	268983666	         4.47 ns/op
BenchmarkParse8/Naive-8        	174564190	         6.70 ns/op
BenchmarkParse8/ParseUint-8    	82321137	        14.8 ns/op
BenchmarkParse4/GoBits-8       	359085199	         3.36 ns/op
BenchmarkParse4/GoUnrolled-8   	1000000000	         0.795 ns/op
BenchmarkParse4/Naive-8        	281167344	         4.32 ns/op
BenchmarkParse4/ParseUint-8    	125097560	         9.49 ns/op
BenchmarkParse3/GoUnrolled-8   	1000000000	         0.666 ns/op
BenchmarkParse3/Naive-8        	120597782	        10.1 ns/op
BenchmarkParse3/ParseUint-8    	60491136	        20.3 ns/op
PASS
ok  	github.com/orisano/satol	37.800s
```

## Author
Nao Yonashiro (@orisano)

## License
Apache v2.0
