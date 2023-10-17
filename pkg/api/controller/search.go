package controller

import (
	"artemb/flights-path/pkg/api/response"
	"artemb/flights-path/pkg/graph"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sort"
	"time"
)

type SearchController struct {
	Logger *zap.Logger
}

type SearchResponse struct {
	ShortPath []string `json:"short_path"`
	FullPath  []string `json:"full_path"`
}

func (c *SearchController) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO not using validator here, since it's simple structure
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteJSONResponse(w, r, http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}
	if len(body) == 0 {
		response.WriteJSONResponse(w, r, http.StatusBadRequest, response.ErrorResponse{Error: "empty payload"})
		return
	}

	var segments [][]string
	err = json.Unmarshal(body, &segments)
	if err != nil {
		response.WriteJSONResponse(w, r, http.StatusBadRequest, response.ErrorResponse{Error: "wrong payload"})
		return
	}

	if len(segments) == 0 {
		response.WriteJSONResponse(w, r, http.StatusBadRequest, response.ErrorResponse{Error: "wrong segments in payload"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	result, err := c.calculate(ctx, segments)
	if err != nil {
		response.WriteJSONResponse(w, r, http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if len(result) == 0 {
		response.WriteJSONResponse(w, r, http.StatusBadRequest, response.ErrorResponse{Error: "can't find route"})
		return
	}

	response.WriteJSONResponse(w, r, http.StatusOK, SearchResponse{FullPath: result, ShortPath: []string{result[0], result[len(result)-1]}})
}

func (c *SearchController) calculate(ctx context.Context, segments [][]string) ([]string, error) {
	sort.Slice(segments, func(i, j int) bool {
		if segments[i][0] < segments[j][0] {
			return true
		}
		if segments[i][0] > segments[j][0] {
			return false
		}
		return segments[i][1] < segments[j][1]
	})

	var err error
	g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())
	for _, el := range segments {
		source := el[0]
		target := el[1]

		_, err = g.Vertex(source)
		if err != nil {
			err = g.AddVertex(source)
			if err != nil {
				return nil, err
			}
		}

		_, err = g.Vertex(target)
		if err != nil {
			err = g.AddVertex(target)
			if err != nil {
				return nil, err
			}
		}

		_, err = g.Edge(source, target)
		if err != nil {
			err = g.AddEdge(source, target)
			if err != nil {
				return nil, err
			}
		}
	}

	var dfs []string
	for _, el := range segments {
		start := el[0]
		var innerDfs []string
		err = graph.DFS(g, start, func(value string) bool {
			innerDfs = append(innerDfs, value)
			return false
		})
		if err != nil {
			println(err)
		}
		if len(dfs) < len(innerDfs) {
			dfs = innerDfs
		}
	}

	return dfs, nil
}
