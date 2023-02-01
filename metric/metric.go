package metric

// Metric is the interface of metrics which defines target search spaces.
type Metric interface {
	// CalcDistance ... calculates the distance between given vectors
	CalcDistance(v1, v2 []float32) float32
	// GetSplittingVector ... calculates the splitting vector which becomes a node's vector in the index
	GetSplittingVector(vs [][]float32) []float32
	// CalcDirectionPriority ... calculates the priority of the children nodes which can be used for determining
	// which way (right or left child) should go next traversal. The return values must be contained in [-1, 1].
	CalcDirectionPriority(base, target []float32) float32
}
