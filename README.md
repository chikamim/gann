# gann
[![CircleCI](https://circleci.com/gh/mathetake/gann.svg?style=svg)](https://circleci.com/gh/mathetake/gann)

gann (go-approximate-nearest-neighbor) is a library for approximate nearest neighbor search purely written in golang.

The implemented algorithm is truly inspired by Annoy (https://github.com/spotify/annoy).

# feature
1. __ONLY__ written in golang, no dependencies out of go world.
2. easy to tune with a bit of parameters
3. __ONLY support for cosine similarity search.__ (issue: https://github.com/mathetake/gann/issues/12)

# usage

```golang
import (
	"fmt"
	"github.com/mathetake/gann"
	"math/rand"
	"time"
)

func main() {

	dim := 3
	nTrees := 2
	k := 10

	rawItems := [][]float32{}
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 1000; i++ {
		rawItems = append(rawItems, []float32{
			rand.Float32(),
			rand.Float32(),
			rand.Float32(),
		})
	}

	// create index
	gIDx, err := gann.GetIndex(rawItems, dim, nTrees, k)
	if err != nil {
		panic(err)
	}
	// build index
	gIDx.Build()

	// do search
	q := []float32{0.1, 0.02, 0.001}
	ann, err := gIDx.GetANNbyVector(q, 5, 10)
	fmt.Println("result:", ann)
}
```
# interfaces

https://github.com/mathetake/gann/blob/master/gann.go#L32-L36
```golang
type GannIndex interface {
	Build() error
	GetANNbyItemID(id int64, num int, bucketScale float64) (ann []int64, err error)
	GetANNbyVector(v []float32, num int, bucketScale float64) (ann []int64, err error)
}
```

# parameters

To be explained

# references

- https://github.com/spotify/annoy
- https://en.wikipedia.org/wiki/Nearest_neighbor_search#Approximate_nearest_neighbor

# License

MIT