# gann
[![CircleCI](https://circleci.com/gh/mathetake/gann.svg?style=shield&circle-token=9a6608c5baa7a400661a700127778a9ff8baeee3)](https://circleci.com/gh/mathetake/gann)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

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
	gIDx, err := gann.GetIndex(rawItems, dim, nTrees, k, true)
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

```golang
type GannIndex interface {
	Build() error
	GetANNbyItemID(id int64, num int, bucketScale float32) (ann []int64, err error)
	GetANNbyVector(v []float32, num int, bucketScale float32) (ann []int64, err error)
}
```

# parameters

See the blog post describing the parameters and algorithms in _gann_  :

https://mathetake.github.io/blogs/gann.html

# references

- https://github.com/spotify/annoy
- https://en.wikipedia.org/wiki/Nearest_neighbor_search#Approximate_nearest_neighbor

# License

MIT