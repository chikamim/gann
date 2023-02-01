package gann

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/chikamim/gann/metric"
)

func TestIndex_GetANNbyItemID(t *testing.T) {
	for i, c := range []struct {
		dim, num, nTree, k int
	}{
		{dim: 2, num: 1000, nTree: 10, k: 2},
		{dim: 10, num: 100, nTree: 5, k: 10},
		{dim: 1000, num: 10000, nTree: 5, k: 10},
	} {
		c := c
		t.Run(fmt.Sprintf("%d-th case", i), func(t *testing.T) {
			rawItems := make([][]float32, c.num)
			for i := range rawItems {
				v := make([]float32, c.dim)

				var norm float32
				for j := range v {
					cof := float32(rand.Float64() - 0.5)
					v[j] = cof
					norm += cof * cof
				}

				norm = float32(math.Sqrt(float64(norm)))
				for j := range v {
					v[j] /= norm
				}

				rawItems[i] = v
			}

			m, err := metric.NewCosineMetric(c.dim)
			if err != nil {
				t.Fatal(err)
			}

			idx, err := CreateNewIndex(rawItems, c.dim, c.nTree, c.k, m)
			if err != nil {
				t.Fatal(err)
			}

			if _, err = idx.GetANNbyItemID(0, 10, 2); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIndex_GetANNbyVector(t *testing.T) {
	for i, c := range []struct {
		dim, num, nTree, k int
	}{
		{dim: 2, num: 1000, nTree: 10, k: 2},
		{dim: 10, num: 100, nTree: 5, k: 10},
		{dim: 1000, num: 10000, nTree: 5, k: 10},
	} {
		c := c
		t.Run(fmt.Sprintf("%d-th case", i), func(t *testing.T) {
			rawItems := make([][]float32, c.num)
			for i := range rawItems {
				v := make([]float32, c.dim)

				var norm float32
				for j := range v {
					cof := float32(rand.Float64() - 0.5)
					v[j] = cof
					norm += cof * cof
				}

				norm = float32(math.Sqrt(float64(norm)))
				for j := range v {
					v[j] /= norm
				}

				rawItems[i] = v
			}

			m, err := metric.NewCosineMetric(c.dim)
			if err != nil {
				t.Fatal(err)
			}

			idx, err := CreateNewIndex(rawItems, c.dim, c.nTree, c.k, m)
			if err != nil {
				t.Fatal(err)
			}

			key := make([]float32, c.dim)
			for i := range key {
				key[i] = float32(rand.Float64() - 0.5)
			}

			if _, err = idx.GetANNbyVector(key, 10, 2); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// This unit test is made to verify if our algorithm can correctly find
// the `exact` neighbors. That is done by checking the ratio of exact
// neighbors in the result returned by `getANNbyVector` is less than
// the given threshold.
func TestAnnSearchAccuracy(t *testing.T) {
	for i, c := range []struct {
		k, dim, num, nTree, searchNum int
		threshold, bucketScale        float64
	}{
		{
			k:           2,
			dim:         20,
			num:         10000,
			nTree:       20,
			threshold:   0.90,
			searchNum:   200,
			bucketScale: 20,
		},
		{
			k:           2,
			dim:         20,
			num:         10000,
			nTree:       20,
			threshold:   0.8,
			searchNum:   20,
			bucketScale: 1000,
		},
	} {
		c := c
		t.Run(fmt.Sprintf("%d-th case", i), func(t *testing.T) {
			rawItems := make([][]float32, c.num)
			for i := range rawItems {
				v := make([]float32, c.dim)

				var norm float32
				for j := range v {
					cof := float32(rand.Float64() - 0.5)
					v[j] = cof
					norm += cof * cof
				}

				norm = float32(math.Sqrt(float64(norm)))
				for j := range v {
					v[j] /= norm
				}

				rawItems[i] = v
			}

			m, err := metric.NewCosineMetric(c.dim)
			if err != nil {
				t.Fatal(err)
			}

			idx, err := CreateNewIndex(rawItems, c.dim, c.nTree, c.k, m)
			if err != nil {
				t.Fatal(err)
			}

			rawIdx, ok := idx.(*index)
			if !ok {
				t.Fatal("assertion failed")
			}

			// query vector
			query := make([]float32, c.dim)
			query[0] = 0.1

			// exact neighbors
			aDist := map[int64]float32{}
			ids := make([]int64, len(rawItems))
			for i, v := range rawItems {
				ids[i] = int64(i)
				aDist[int64(i)] = rawIdx.metric.CalcDistance(v, query)
			}
			sort.Slice(ids, func(i, j int) bool {
				return aDist[ids[i]] < aDist[ids[j]]
			})

			expectedIDsMap := make(map[int64]struct{}, c.searchNum)
			for _, id := range ids[:c.searchNum] {
				expectedIDsMap[int64(id)] = struct{}{}
			}

			ass, err := idx.GetANNbyVector(query, c.searchNum, c.bucketScale)
			if err != nil {
				t.Fatal(err)
			}

			var count int
			for _, id := range ass {
				if _, ok := expectedIDsMap[id]; ok {
					count++
				}
			}

			if ratio := float64(count) / float64(c.searchNum); ratio < c.threshold {
				t.Fatalf("Too few exact neighbors found in approximated result: %d / %d = %f", count, c.searchNum, ratio)
			} else {
				t.Logf("ratio of exact neighbors in approximated result: %d / %d = %f", count, c.searchNum, ratio)
			}
		})
	}
}
