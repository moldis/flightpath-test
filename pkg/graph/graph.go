package graph

import "errors"

var (
	ErrVertexNotFound      = errors.New("vertex not found")
	ErrVertexAlreadyExists = errors.New("vertex already exists")
	ErrEdgeNotFound        = errors.New("edge not found")
	ErrEdgeAlreadyExists   = errors.New("edge already exists")
	ErrEdgeCreatesCycle    = errors.New("edge would create a cycle")
	ErrVertexHasEdges      = errors.New("vertex has edges")
)

// Graph represents a generic graph data structure consisting of vertices of
// type T identified by a hash of type K.
type Graph[K comparable, T any] interface {
	// Traits returns the graph's traits. Those traits must be set when creating
	// a graph using New.
	Traits() *Traits

	// AddVertex creates a new vertex in the graph. If the vertex already exists
	// in the graph, ErrVertexAlreadyExists will be returned.
	//
	// AddVertex accepts a variety of functional options to set further edge
	// details such as the weight or an attribute:
	//
	//	_ = graph.AddVertex("A", "B")
	//
	AddVertex(value T) error

	// Vertex returns the vertex with the given hash or ErrVertexNotFound if it
	// doesn't exist.
	Vertex(hash K) (T, error)

	// RemoveVertex removes the vertex with the given hash value from the graph.
	//
	// The vertex is not allowed to have edges and thus must be disconnected.
	// Potential edges must be removed first. Otherwise, ErrVertexHasEdges will
	// be returned. If the vertex doesn't exist, ErrVertexNotFound is returned.
	RemoveVertex(hash K) error

	// AddEdge creates an edge between the source and the target vertex.
	//
	// If either vertex cannot be found, ErrVertexNotFound will be returned. If
	// the edge already exists, ErrEdgeAlreadyExists will be returned. If cycle
	// prevention has been activated using PreventCycles and if adding the edge
	// would create a cycle, ErrEdgeCreatesCycle will be returned.
	//
	//	_ = g.AddEdge("A", "B")
	//
	AddEdge(sourceHash, targetHash K) error

	// Edge returns the edge joining two given vertices or ErrEdgeNotFound if
	// the edge doesn't exist. In an undirected graph, an edge with swapped
	// source and target vertices does match.
	Edge(sourceHash, targetHash K) (Edge[T], error)

	// Edges returns a slice of all edges in the graph. These edges are of type
	// Edge[K] and hence will contain the vertex hashes, not the vertex values.
	Edges() ([]Edge[K], error)

	// RemoveEdge removes the edge between the given source and target vertices.
	// If the edge cannot be found, ErrEdgeNotFound will be returned.
	RemoveEdge(source, target K) error

	// AdjacencyMap computes an adjacency map with all vertices in the graph.
	//
	// There is an entry for each vertex. Each of those entries is another map
	// whose keys are the hash values of the adjacent vertices. The value is an
	// Edge instance that stores the source and target hash values along with
	// the edge metadata.
	//
	// For a directed graph with two edges AB and AC, AdjacencyMap would return
	// the following map:
	//
	//	map[string]map[string]Edge[string]{
	//		"A": map[string]Edge[string]{
	//			"B": {Source: "A", Target: "B"},
	//			"C": {Source: "A", Target: "C"},
	//		},
	//		"B": map[string]Edge[string]{},
	//		"C": map[string]Edge[string]{},
	//	}
	//
	// This design makes AdjacencyMap suitable for a wide variety of algorithms.
	AdjacencyMap() (map[K]map[K]Edge[K], error)

	// PredecessorMap computes a predecessor map with all vertices in the graph.
	//
	// It has the same map layout and does the same thing as AdjacencyMap, but
	// for ingoing instead of outgoing edges of each vertex.
	//
	// For a directed graph with two edges AB and AC, PredecessorMap would
	// return the following map:
	//
	//	map[string]map[string]Edge[string]{
	//		"A": map[string]Edge[string]{},
	//		"B": map[string]Edge[string]{
	//			"A": {Source: "A", Target: "B"},
	//		},
	//		"C": map[string]Edge[string]{
	//			"A": {Source: "A", Target: "C"},
	//		},
	//	}
	//
	// For an undirected graph, PredecessorMap is the same as AdjacencyMap. This
	// is because there is no distinction between "outgoing" and "ingoing" edges
	// in an undirected graph.
	PredecessorMap() (map[K]map[K]Edge[K], error)

	// Order returns the number of vertices in the graph.
	Order() (int, error)

	// Size returns the number of edges in the graph.
	Size() (int, error)
}

// Edge represents an edge that joins two vertices. Even though these edges are
// always referred to as source and target, whether the graph is directed or not
// is determined by its traits.
type Edge[T any] struct {
	Source T
	Target T
}

// Hash is a hashing function that takes a vertex of type T and returns a hash
// value of type K.
//
// Every graph has a hashing function and uses that function to retrieve the
// hash values of its vertices. You can either use one of the predefined hashing
// functions or provide your own one for custom data types:
//
//	cityHash := func(c City) string {
//		return c.Name
//	}
//
// The cityHash function returns the city name as a hash value. The types of T
// and K, in this case City and string, also define the types of the graph.
type Hash[K comparable, T any] func(T) K

// New creates a new graph with vertices of type T, identified by hash values of
// type K. These hash values will be obtained using the provided hash function.
//
// The graph will use the default in-memory store for persisting vertices and
// edges. To use a different [Store], use [NewWithStore].
func New[K comparable, T any](hash Hash[K, T], options ...func(*Traits)) Graph[K, T] {
	return NewWithStore(hash, newMemoryStore[K, T](), options...)
}

// NewWithStore creates a new graph same as [New] but uses the provided store
// instead of the default memory store.
func NewWithStore[K comparable, T any](hash Hash[K, T], store Store[K, T], options ...func(*Traits)) Graph[K, T] {
	var p Traits

	for _, option := range options {
		option(&p)
	}

	return newDirected(hash, &p, store)
}

// StringHash is a hashing function that accepts a string and uses that exact
// string as a hash value. Using it as Hash will yield a Graph[string, string].
func StringHash(v string) string {
	return v
}

// IntHash is a hashing function that accepts an integer and uses that exact
// integer as a hash value. Using it as Hash will yield a Graph[int, int].
func IntHash(v int) int {
	return v
}
