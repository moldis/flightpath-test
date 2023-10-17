package graph

import (
	"errors"
	"fmt"
)

type directed[K comparable, T any] struct {
	hash   Hash[K, T]
	traits *Traits
	store  Store[K, T]
}

func newDirected[K comparable, T any](hash Hash[K, T], traits *Traits, store Store[K, T]) *directed[K, T] {
	return &directed[K, T]{
		hash:   hash,
		traits: traits,
		store:  store,
	}
}

func (d *directed[K, T]) Traits() *Traits {
	return d.traits
}

func (d *directed[K, T]) AddVertex(value T) error {
	hash := d.hash(value)
	return d.store.AddVertex(hash, value)
}

func (d *directed[K, T]) Vertex(hash K) (T, error) {
	vertex, err := d.store.Vertex(hash)
	return vertex, err
}

func (d *directed[K, T]) RemoveVertex(hash K) error {
	return d.store.RemoveVertex(hash)
}

func (d *directed[K, T]) AddEdge(sourceHash, targetHash K) error {
	_, err := d.store.Vertex(sourceHash)
	if err != nil {
		return fmt.Errorf("source vertex %v: %w", sourceHash, err)
	}

	_, err = d.store.Vertex(targetHash)
	if err != nil {
		return fmt.Errorf("target vertex %v: %w", targetHash, err)
	}

	if _, err := d.Edge(sourceHash, targetHash); !errors.Is(err, ErrEdgeNotFound) {
		return ErrEdgeAlreadyExists
	}

	// If the user opted in to preventing cycles, run a cycle check.
	if d.traits.PreventCycles {
		createsCycle, err := d.createsCycle(sourceHash, targetHash)
		if err != nil {
			return fmt.Errorf("check for cycles: %w", err)
		}
		if createsCycle {
			return ErrEdgeCreatesCycle
		}
	}

	edge := Edge[K]{
		Source: sourceHash,
		Target: targetHash,
	}

	return d.addEdge(sourceHash, targetHash, edge)
}

func (d *directed[K, T]) Edge(sourceHash, targetHash K) (Edge[T], error) {
	_, err := d.store.Edge(sourceHash, targetHash)
	if err != nil {
		return Edge[T]{}, err
	}

	sourceVertex, err := d.store.Vertex(sourceHash)
	if err != nil {
		return Edge[T]{}, err
	}

	targetVertex, err := d.store.Vertex(targetHash)
	if err != nil {
		return Edge[T]{}, err
	}

	return Edge[T]{
		Source: sourceVertex,
		Target: targetVertex,
	}, nil
}

func (d *directed[K, T]) Edges() ([]Edge[K], error) {
	return d.store.ListEdges()
}

func (d *directed[K, T]) RemoveEdge(source, target K) error {
	if _, err := d.Edge(source, target); err != nil {
		return err
	}

	if err := d.store.RemoveEdge(source, target); err != nil {
		return fmt.Errorf("failed to remove edge from %v to %v: %w", source, target, err)
	}

	return nil
}

func (d *directed[K, T]) AdjacencyMap() (map[K]map[K]Edge[K], error) {
	vertices, err := d.store.ListVertices()
	if err != nil {
		return nil, fmt.Errorf("failed to list vertices: %w", err)
	}

	edges, err := d.store.ListEdges()
	if err != nil {
		return nil, fmt.Errorf("failed to list edges: %w", err)
	}

	m := make(map[K]map[K]Edge[K], len(vertices))

	for _, vertex := range vertices {
		m[vertex] = make(map[K]Edge[K])
	}

	for _, edge := range edges {
		m[edge.Source][edge.Target] = edge
	}

	return m, nil
}

func (d *directed[K, T]) PredecessorMap() (map[K]map[K]Edge[K], error) {
	vertices, err := d.store.ListVertices()
	if err != nil {
		return nil, fmt.Errorf("failed to list vertices: %w", err)
	}

	edges, err := d.store.ListEdges()
	if err != nil {
		return nil, fmt.Errorf("failed to list edges: %w", err)
	}

	m := make(map[K]map[K]Edge[K], len(vertices))

	for _, vertex := range vertices {
		m[vertex] = make(map[K]Edge[K])
	}

	for _, edge := range edges {
		if _, ok := m[edge.Target]; !ok {
			m[edge.Target] = make(map[K]Edge[K])
		}
		m[edge.Target][edge.Source] = edge
	}

	return m, nil
}

func (d *directed[K, T]) addEdge(sourceHash, targetHash K, edge Edge[K]) error {
	return d.store.AddEdge(sourceHash, targetHash, edge)
}

func (d *directed[K, T]) Order() (int, error) {
	return d.store.VertexCount()
}

func (d *directed[K, T]) Size() (int, error) {
	size := 0
	outEdges, err := d.AdjacencyMap()
	if err != nil {
		return 0, fmt.Errorf("failed to get adjacency map: %w", err)
	}

	for _, outEdges := range outEdges {
		size += len(outEdges)
	}

	return size, nil
}

func (d *directed[K, T]) edgesAreEqual(a, b Edge[T]) bool {
	aSourceHash := d.hash(a.Source)
	aTargetHash := d.hash(a.Target)
	bSourceHash := d.hash(b.Source)
	bTargetHash := d.hash(b.Target)

	return aSourceHash == bSourceHash && aTargetHash == bTargetHash
}

func (d *directed[K, T]) createsCycle(source, target K) (bool, error) {
	// If the underlying store implements CreatesCycle, use that fast path.
	if cc, ok := d.store.(interface {
		CreatesCycle(source, target K) (bool, error)
	}); ok {
		return cc.CreatesCycle(source, target)
	}

	// Slow path.
	return CreatesCycle(Graph[K, T](d), source, target)
}
