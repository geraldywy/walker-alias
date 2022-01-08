package walkeralias

import (
	"math/rand"
)

// walkerAlias holds an internal rand reference instead of sharing with global rand.
type walkerAlias struct {
	buckets []*bucket
	r       *rand.Rand
}

// NewWalkerAlias accepts a map {key: probability} and a seed to init a new rand for its own use.
// Returns a reference to a WalkerAlias object to use
// An example of a probabilityMap: {1: 2, 2: 3}, the ratio of selecting key 1 and key 2 is 0.4 and 0.6 respectively
// WalkerAlias involves an O(n) preprocessing step to generate a probability table.
// Subsequent sampling are all O(1).
func NewWalkerAlias(probabilityMap map[int]float64, seed int64) *walkerAlias {
	n := len(probabilityMap)
	buckets := make([]*bucket, 0)

	var sumWeights float64
	for _, w := range probabilityMap {
		sumWeights += w
	}
	for k, w := range probabilityMap {
		prob := w * float64(n) / sumWeights
		buckets = append(buckets, newBucket(k, prob))
	}

	underfull := make([]int, 0)
	overfull := make([]int, 0)
	for i, b := range buckets {
		if b.threshold < 1 {
			underfull = append(underfull, i)
		} else if b.threshold > 1 {
			overfull = append(overfull, i)
		}
	}

	for len(underfull) > 0 && len(overfull) > 0 {
		u, o := underfull[len(underfull)-1], overfull[len(overfull)-1]
		underfull = underfull[:len(underfull)-1]
		under, over := buckets[u], buckets[o]
		under.key2 = over.key1
		over.threshold -= 1 - under.threshold
		if over.threshold < 1 {
			underfull = append(underfull, o)
			overfull = overfull[:len(overfull)-1]
		}
	}

	return &walkerAlias{buckets: buckets, r: rand.New(rand.NewSource(seed))}
}

// Random returns a random key following the given probability
func (w *walkerAlias) Random() int {
	bucketIdx := rand.Intn(len(w.buckets))
	b := w.buckets[bucketIdx]
	prob := rand.Float64()
	if prob > b.threshold {
		return b.key2
	}

	return b.key1
}

// newBucket returns a ref to a bucket object with the given key
// and sets its initial threshold to the prob (probability) given
func newBucket(key int, prob float64) *bucket {
	return &bucket{threshold: prob, key1: key, key2: -1}
}

// bucket holds 2 keys at most,
// Returns Key1 below or equal to the threshold, Key2 strictly above the threshold
type bucket struct {
	threshold float64 // threshold point
	key1      int     // key below or equal to threshold
	key2      int     // key above threshold
}
