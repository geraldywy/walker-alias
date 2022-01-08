package walkeralias

import (
	"math"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"
)

func TestWalkerAlias_Random(t *testing.T) {
	tests := []struct {
		name             string
		pMap             map[int]float64
		setupFunc        func(pMap map[int]float64) randomizer
		iterations       int
		allowedThreshold float64 // allowable probability variance, from [0-1]
		skipInCI bool // indicate whether to skip test in CI to prevent timeouts
	}{
		{
			name: "[WalkerAlias] no floating point rounding errors",
			pMap: map[int]float64{
				0: 3.5,
				1: 6.5,
				2: 10,
			},
			setupFunc: func(pMap map[int]float64) randomizer {
				w := NewWalkerAlias(pMap, time.Now().Unix())
				return w
			},
			iterations:       10000000,
			allowedThreshold: 0.0005, // 0.05%
			skipInCI: false,
		},
		{
			name: "[WalkerAlias] with floating point rounding errors",
			pMap: map[int]float64{
				100: 2,
				300: 1,
				500: 1,
				600: 1,
				102: 2,
				320: 1,
				111: 2,
			},
			setupFunc: func(pMap map[int]float64) randomizer {
				w := NewWalkerAlias(pMap, time.Now().Unix())
				return w
			},
			iterations:       10000000,
			allowedThreshold: 0.0005,
			skipInCI: false,
		},
		{
			name: "[WalkerAlias] more iterations should reflect a tighter threshold", // takes ~50s to run
			pMap: map[int]float64{
				0: 3.5,
				1: 6.5,
				2: 10,
			},
			setupFunc: func(pMap map[int]float64) randomizer {
				w := NewWalkerAlias(pMap, time.Now().Unix())
				return w
			},
			iterations:       1000000000,
			allowedThreshold: 0.00005, // 0.005%
			skipInCI: true,
		},
		{
			name: "[WalkerAlias] more iterations should reflect a tighter threshold again", // takes ~50s to run
			pMap: map[int]float64{
				100: 2,
				300: 1,
				500: 1,
				600: 1,
				102: 2,
				320: 1,
				111: 2,
			},
			setupFunc: func(pMap map[int]float64) randomizer {
				w := NewWalkerAlias(pMap, time.Now().Unix())
				return w
			},

			iterations:       1000000000,
			allowedThreshold: 0.00005,
			skipInCI: true,
		},
		{
			name: "[NaiveSearch] no floating point rounding errors",
			pMap: map[int]float64{
				0: 1,
				1: 2,
				2: 3,
				3: 4,
			},
			setupFunc: func(pMap map[int]float64) randomizer {
				n := newNaiveSearch(pMap, time.Now().Unix())
				return n
			},
			iterations:       10000000,
			allowedThreshold: 0.0005, // 0.05%
			skipInCI: false,
		},
		{
			name: "[BinarySearchPartitions] no floating point rounding errors",
			pMap: map[int]float64{
				0: 1,
				1: 2,
				2: 3,
				3: 4,
			},
			setupFunc: func(pMap map[int]float64) randomizer {
				n := newNaiveSearch(pMap, time.Now().Unix())
				return n
			},
			iterations:       10000000,
			allowedThreshold: 0.0005, // 0.05%
			skipInCI: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipInCI && os.Getenv("CI") != "" {
				t.Skip("Skipping testing in CI environment")
			}
			randomizer := tt.setupFunc(tt.pMap)
			actualPMap := make(map[int]float64)
			for i := 0; i < tt.iterations; i++ {
				actualPMap[randomizer.Random()] += float64(1) / float64(tt.iterations)
			}
			var sumWeights float64
			for _, v := range tt.pMap {
				sumWeights += v
			}
			for key, weight := range tt.pMap {
				expectedProb := weight / sumWeights
				if math.Abs(actualPMap[key]-expectedProb) > tt.allowedThreshold {
					t.Errorf("actual probability (%.5f%%) differed from expected prob (%.5f%%) by more than acceptable range (%.5f%%)",
						actualPMap[key], expectedProb, tt.allowedThreshold)
				}
			}
		})
	}
}

type randomizer interface {
	Random() int
}

// Benchmarking tests across naive randomizer O(N), binary searching partitions O(NlogN), alias walker O(1)
func BenchmarkWalkerAlias_Random(b *testing.B) {
	tests := []struct {
		name      string
		pMap      map[int]float64
		setupFunc func(pMap map[int]float64) randomizer
	}{
		{
			name: "naive search",
			setupFunc: func(pMap map[int]float64) randomizer {
				n := newNaiveSearch(pMap, time.Now().Unix())
				return n
			},
		},
		{
			name: "binary searching partitions",
			setupFunc: func(pMap map[int]float64) randomizer {
				b := newBSearchPartitions(pMap, time.Now().Unix())
				return b
			},
		},
		{
			name: "walker alias",
			setupFunc: func(pMap map[int]float64) randomizer {
				w := NewWalkerAlias(pMap, time.Now().Unix())
				return w
			},
		},
	}
	for _, tt := range tests {
		pMap := make(map[int]float64)
		for i := 1; i <= 10000000; i++ { // 10 million entries
			pMap[i] = float64(i)
		}
		b.Run(tt.name, func(b *testing.B) {
			r := tt.setupFunc(pMap)
			for i := 0; i < b.N; i++ {
				r.Random()
			}
		})
	}
}

/*
	Alternatives to Walker Alias for benchmarking purposes
*/

func newNaiveSearch(pMap map[int]float64, seed int64) *naiveSearch {
	keys := make([]int, 0)
	probs := make([]float64, 0)
	var sum float64
	for k, w := range pMap {
		sum += w
		keys = append(keys, k)
		probs = append(probs, sum)
	}

	return &naiveSearch{keys: keys, probs: probs, r: rand.New(rand.NewSource(seed))}
}

type naiveSearch struct {
	probs []float64
	keys  []int
	r     *rand.Rand
}

func (n *naiveSearch) Random() int {
	prob := n.r.Float64() * n.probs[len(n.probs)-1]
	for i, p := range n.probs {
		if p >= prob {
			return n.keys[i]
		}
	}

	return -1
}

func newBSearchPartitions(pMap map[int]float64, seed int64) *bSearchPartitions {
	keys := make([]int, 0)
	probs := make([]float64, 0)
	var cnt float64
	for k, w := range pMap {
		cnt += w
		keys = append(keys, k)
		probs = append(probs, cnt)
	}

	return &bSearchPartitions{keys: keys, probs: probs, r: rand.New(rand.NewSource(seed))}
}

type bSearchPartitions struct {
	probs []float64
	keys  []int
	sum   float64
	r     *rand.Rand
}

func (b *bSearchPartitions) Random() int {
	prob := b.r.Float64() * b.probs[len(b.probs)-1]
	idx := sort.SearchFloat64s(b.probs, prob)

	return b.keys[idx]
}
