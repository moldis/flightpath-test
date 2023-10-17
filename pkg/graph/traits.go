package graph

// Traits represents a set of graph traits and types, such as directedness or acyclicness. These
// traits can be set when creating a graph by passing the corresponding functional options, for
// example:
//
//	g := graph.New(graph.IntHash, graph.Directed())
//
// This will set the IsDirected field to true.
type Traits struct {
	IsDirected    bool
	IsAcyclic     bool
	PreventCycles bool
}

// Directed creates a directed graph. This has implications on graph traversal and the order of
// arguments of the Edge and AddEdge functions.
func Directed() func(*Traits) {
	return func(t *Traits) {
		t.IsDirected = true
	}
}

// Acyclic creates an acyclic graph. Note that creating edges that form a cycle will still be
// possible. To prevent this explicitly, use PreventCycles.
func Acyclic() func(*Traits) {
	return func(t *Traits) {
		t.IsAcyclic = true
	}
}

// PreventCycles creates an acyclic graph that prevents and proactively prevents the creation of
// cycles. These cycle checks affect the performance and complexity of operations such as AddEdge.
func PreventCycles() func(*Traits) {
	return func(t *Traits) {
		Acyclic()(t)
		t.PreventCycles = true
	}
}
