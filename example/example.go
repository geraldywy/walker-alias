package example

import (
	"fmt"
	"time"

	wa "github.com/geraldywy/walker-alias"
)

func main() {
	items := []string{"first", "second", "third"}
	pMap := map[int]float64{
		0: 3.5,
		1: 6.5,
		2: 10,
	}
	w := wa.NewWalkerAlias(pMap, time.Now().Unix())
	freq := make(map[string]int)
	for i := 0; i < 1000000; i++ {
		freq[items[w.Random()]]++
	}

	/*
		to verify, reflect back the rough corresponding probability eg:
		first  0.175
		second 0.325
		third  0.500
	*/
	for k, v := range freq {
		fmt.Println(k, float64(v)/float64(1000000))
	}
}
