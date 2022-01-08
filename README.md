# walker alias method

[Walker's alias method](https://en.wikipedia.org/wiki/Alias_method) is an efficient algorithm to sample from a discrete probability distribution. This means given an arbitrary probability distribution like ```{A: 1, B: 2, C: 3, D: 4}```, the odds of sampling A is roughly 1 in 10 and the odds of sampling D is roughly 4 in 10.

The algorithm is able to generate a random sample in O(1) time and is ideal for cases where frequent queries are done against a fixed arbitrary weight distribution. <br>
The preprocessing step for probability table creation is O(n) time. For floating point rounding errors, a negligible error is to be expected (<0.05% for 10000000 iterations).

***
[![Go Report Card](https://goreportcard.com/badge/github.com/geraldywy/walker-alias)](https://goreportcard.com/report/github.com/geraldywy/walker-alias)
[![Coverage Status](https://coveralls.io/repos/github/geraldywy/walker-alias/badge.svg?branch=master)](https://coveralls.io/github/geraldywy/walker-alias?branch=master)

## Install with go mod
```
go get github.com/geraldywy/walker-alias
```

## Usage
```go
package main

import (
	"fmt"
	"time"

	wa "github.com/geraldywy/walker-alias"
)

func main() {
	// key: weight mapping, keys can be any valid int, weights can be any valid float64
	pMap := map[int]float64{
		0: 3.5,
		1: 6.5,
		2: 10,
	}
	w := wa.NewWalkerAlias(pMap, time.Now().Unix()) // init with weights and a random seed

	for i := 0; i < 5; i++ {
		randomKey := w.Random() // generates a random key in O(1)
		fmt.Println(randomKey)
	}
}
```

> For usage with a non int key, one workaround is to use a lookup table to map the key to a unique int and map back after sampling.
For a simple example with string keys, refer to ```example.go``` in the ```example``` folder.
Alternatively, you could implement it with generics and submit a pull request.


## Benchmarking

Walker Alias is much faster compared to alternatives like sampling via naive linear search or binary search via partitions.
Below shows the result of sampling from 10 million keys using all three methods, with Walker Alias clearly outperforming others.
The full implementation for the benchmarking can be found in main_test.go.

```
$ go test -bench=. -benchtime=60s -count=3  -run=^# -timeout 20m
goos: darwin
goarch: amd64
pkg: github.com/geraldywy/walker-alias
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkWalkerAlias_Random/naive_search-16                             20083                  3514372 ns/op
BenchmarkWalkerAlias_Random/naive_search-16                             20355                  3504450 ns/op
BenchmarkWalkerAlias_Random/naive_search-16                             20383                  3534106 ns/op
BenchmarkWalkerAlias_Random/binary_searching_partitions-16              130537809              558.2 ns/op
BenchmarkWalkerAlias_Random/binary_searching_partitions-16              127571498              569.2 ns/op
BenchmarkWalkerAlias_Random/binary_searching_partitions-16              126131037              562.6 ns/op
BenchmarkWalkerAlias_Random/walker_alias-16                             425654750              159.7 ns/op
BenchmarkWalkerAlias_Random/walker_alias-16                             412184373              159.7 ns/op
BenchmarkWalkerAlias_Random/walker_alias-16                             424909584              158.0 ns/op
PASS
ok      github.com/geraldywy/walker-alias       947.729s
```

## Licensing

MIT License