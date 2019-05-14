package gann

import (
	"container/heap"
	"math"
	"sort"
)

func (idx *index) GetANNbyItemID(id int64, searchNum int, bucketScale float64) ([]int64, error) {
	it, ok := idx.itemIDToItem[itemId(id)]
	if !ok {
		return nil, errItemNotFoundOnGivenItemID
	}
	return idx.GetANNbyVector(it.vector, searchNum, bucketScale)
}

func (idx *index) GetANNbyVector(v []float64, searchNum int, bucketScale float64) ([]int64, error) {
	/*
		1. insert root nodes into the priority queue
		2. search all trees until len(`ann`) is enough.
		3. calculate actual distances to each elements in ann from v.
		4. sort `ann` by distances.
		5. Return the top `num` ones.
	*/

	if len(v) != idx.dim {
		return nil, errInvalidKeyVector
	}

	bucketSize := int(float64(searchNum) * bucketScale)
	annMap := make(map[itemId]struct{}, bucketSize)

	pq := priorityQueue{}

	// 1.
	for i, r := range idx.roots {
		n := &queueItem{
			value:    r.id,
			index:    i,
			priority: math.Inf(-1),
		}
		pq = append(pq, n)
	}

	heap.Init(&pq)

	// 2.
	for pq.Len() > 0 && len(annMap) < bucketSize {
		q, ok := heap.Pop(&pq).(*queueItem)
		d := q.priority
		n, ok := idx.nodeIDToNode[q.value]
		if !ok {
			return nil, errInvalidIndex
		}

		if len(n.leaf) > 0 {
			for _, id := range n.leaf {
				annMap[id] = struct{}{}
			}
			continue
		}

		dp := idx.metric.CalcDirectionPriority(n.vec, v)
		heap.Push(&pq, &queueItem{
			value:    n.children[left].id,
			priority: max(d, dp),
		})
		heap.Push(&pq, &queueItem{
			value:    n.children[right].id,
			priority: max(d, -dp),
		})
	}

	// 3.
	idToDist := make(map[int64]float64, len(annMap))
	ann := make([]int64, 0, len(annMap))
	for id := range annMap {
		iid := int64(id)
		ann = append(ann, iid)
		idToDist[iid] = idx.metric.CalcDistance(idx.itemIDToItem[id].vector, v)
	}

	// 4.
	sort.Slice(ann, func(i, j int) bool {
		return idToDist[ann[i]] < idToDist[ann[j]]
	})

	// 5.
	if len(ann) > searchNum {
		ann = ann[:searchNum]
	}
	return ann, nil
}

func max(a, b float64) float64 {
	if a < b {
		return b
	}
	return a
}
