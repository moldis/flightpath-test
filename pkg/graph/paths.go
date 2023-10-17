package graph

import (
	"errors"
	"fmt"
)

var ErrTargetNotReachable = errors.New("target vertex not reachable from source")

// CreatesCycle determines whether adding an edge between the two given vertices
// would introduce a cycle in the graph. CreatesCycle will not create an edge.
//
// A potential edge would create a cycle if the target vertex is also a parent
// of the source vertex. In order to determine this, CreatesCycle runs a DFS.
func CreatesCycle[K comparable, T any](g Graph[K, T], source, target K) (bool, error) {
	if _, err := g.Vertex(source); err != nil {
		return false, fmt.Errorf("could not get vertex with hash %v: %w", source, err)
	}

	if _, err := g.Vertex(target); err != nil {
		return false, fmt.Errorf("could not get vertex with hash %v: %w", target, err)
	}

	if source == target {
		return true, nil
	}

	predecessorMap, err := g.PredecessorMap()
	if err != nil {
		return false, fmt.Errorf("failed to get predecessor map: %w", err)
	}

	stack := make([]K, 0)
	visited := make(map[K]bool)

	stack = append(stack, source)

	for len(stack) > 0 {
		currentHash := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if _, ok := visited[currentHash]; !ok {
			// If the adjacent vertex also is the target vertex, the target is a
			// parent of the source vertex. An edge would introduce a cycle.
			if currentHash == target {
				return true, nil
			}

			visited[currentHash] = true

			for adjacency := range predecessorMap[currentHash] {
				stack = append(stack, adjacency)
			}
		}
	}

	return false, nil
}
