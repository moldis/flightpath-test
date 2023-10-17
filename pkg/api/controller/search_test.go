package controller

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestSearchController(t *testing.T) {
	var err error

	tests := []struct {
		name      string
		route     string
		wantErr   bool
		wantRoute []string
	}{
		{
			name:      "Single route",
			route:     "[[\"SFO\", \"EWR\"]]",
			wantErr:   false,
			wantRoute: []string{"SFO", "EWR"},
		},
		{
			name:      "Few routes",
			route:     "[[\"ATL\", \"EWR\"], [\"SFO\", \"ATL\"]]",
			wantErr:   false,
			wantRoute: []string{"SFO", "ATL", "EWR"},
		},
		{
			name:      "Multiple routes",
			route:     "[[\"IND\", \"EWR\"], [\"SFO\", \"ATL\"], [\"GSO\", \"IND\"], [\"ATL\", \"GSO\"]]",
			wantErr:   false,
			wantRoute: []string{"SFO", "ATL", "GSO", "IND", "EWR"},
		},
		{
			name:      "Duplicated routes",
			route:     "[[\"IND\", \"EWR\"], [\"SFO\", \"ATL\"], [\"SFO\", \"ATL\"], [\"GSO\", \"IND\"], [\"ATL\", \"GSO\"]]",
			wantErr:   false,
			wantRoute: []string{"SFO", "ATL", "GSO", "IND", "EWR"},
		},
		{
			name:    "Cycling routes",
			route:   "[[\"IND\", \"EWR\"], [\"SFO\", \"ATL\"], [\"SFO\", \"ATL\"], [\"SFO\", \"SFO\"], [\"GSO\", \"IND\"], [\"ATL\", \"GSO\"]]",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var segments [][]string
			err = json.Unmarshal([]byte(test.route), &segments)
			assert.NoError(t, err)

			controller := SearchController{}
			res, err := controller.calculate(context.Background(), segments)
			if !test.wantErr {
				assert.NoError(t, err)
				assert.True(t, reflect.DeepEqual(res, test.wantRoute))
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestSearchHttpResponses(t *testing.T) {
	tests := []struct {
		name         string
		route        string
		wantResponse string
		wantCode     int
	}{
		{
			name:         "Single route",
			route:        `[["SFO", "EWR"]]`,
			wantResponse: `{"short_path":["SFO","EWR"],"full_path":["SFO","EWR"]}`,
			wantCode:     200,
		},
		{
			name:         "Few routes",
			route:        `[["ATL", "EWR"], ["SFO", "ATL"]]`,
			wantResponse: `{"short_path":["SFO","EWR"],"full_path":["SFO","ATL","EWR"]}`,
			wantCode:     200,
		},
		{
			name:         "Multiple routes",
			route:        `[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]`,
			wantResponse: `{"short_path":["SFO","EWR"],"full_path":["SFO","ATL","GSO","IND","EWR"]}`,
			wantCode:     200,
		},
		{
			name:         "Duplicated routes",
			route:        `[["IND", "EWR"], ["SFO", "ATL"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]`,
			wantResponse: `{"short_path":["SFO","EWR"],"full_path":["SFO","ATL","GSO","IND","EWR"]}`,
			wantCode:     200,
		},
		{
			name:         "Cycling routes",
			route:        `[["IND", "EWR"], ["SFO", "ATL"], ["SFO", "ATL"], ["SFO", "SFO"], ["GSO", "IND"], ["ATL", "GSO"]]`,
			wantResponse: `{"error":"edge would create a cycle"}`,
			wantCode:     400,
		},
		{
			name:         "wrong payload routes",
			route:        `["IND", "EWR"]`,
			wantResponse: `{"error":"wrong payload"}`,
			wantCode:     400,
		},
		{
			name:         "empty payload",
			route:        ``,
			wantResponse: `{"error":"empty payload"}`,
			wantCode:     400,
		},
		{
			name:         "wrong payload 2",
			route:        `["IND", "EWR", "FDF"]`,
			wantResponse: `{"error":"wrong payload"}`,
			wantCode:     400,
		},
		{
			name:         "disconnected routes will take random",
			route:        `[["IND", "FDF"], ["DAD", "EED"]]`,
			wantResponse: `{"short_path":["DAD","EED"],"full_path":["DAD","EED"]}`,
			wantCode:     200,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := SearchController{}
			bodyReader := strings.NewReader(test.route)
			req := httptest.NewRequest("GET", "http://example.com/test", bodyReader)
			w := httptest.NewRecorder()
			controller.Search(w, req)

			assert.Equal(t, w.Code, test.wantCode)
			assert.Equal(t, w.Body.String(), test.wantResponse+"\n")
		})
	}
}
