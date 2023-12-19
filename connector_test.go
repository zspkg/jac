package jac

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

// -- Test API --

type (
	// GET models
	getTestResponse struct {
		Foo string `json:"foo"`
	}

	// POST models
	postTestRequestBody struct {
		X int64 `json:"x"`
		Y int64 `json:"y"`
	}
	postTestResponse struct {
		Z int64 `json:"z"`
	}

	// PATCH models
	patchTestRequestBody struct {
		Foo string `json:"foo"`
	}
	patchTestResponse struct {
		Bar string `json:"bar"`
	}
)

func newTestRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/integrations", func(r chi.Router) {
		r.Route("/some-logic", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ape.Render(w, getTestResponse{Foo: "bar"})
			})
			r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
				var request postTestRequestBody
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					ape.RenderErr(w, problems.BadRequest(err)...)
				}

				ape.Render(w, postTestResponse{Z: request.X + request.Y})
			})
			r.Post("/multiply", func(w http.ResponseWriter, r *http.Request) {
				var request postTestRequestBody
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					ape.RenderErr(w, problems.BadRequest(err)...)
				}

				ape.Render(w, postTestResponse{Z: request.X * request.Y})
			})
			r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
				var request patchTestRequestBody
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					ape.RenderErr(w, problems.BadRequest(err)...)
				}

				ape.Render(w, patchTestResponse{Bar: request.Foo})
			})
		})
	})

	return r
}

// -- Actual tests --

const (
	testUrlBase = "/integrations/some-logic"
)

var testRouter = newTestRouter()

func getTestJac(server *httptest.Server) Jac {
	var (
		testJac = NewJac(server.URL)
	)

	return testJac
}

func TestJacer_Get(t *testing.T) {
	testServer := httptest.NewServer(testRouter)
	defer testServer.Close()

	testJac := getTestJac(testServer)

	t.Run("get request", func(t *testing.T) {
		var testResponse getTestResponse
		_, err := testJac.Get(RequestParams{Endpoint: testUrlBase}, &testResponse)
		assert.Nil(t, err, "expected nil error when unmarshalling response")
		assert.Equal(t, getTestResponse{Foo: "bar"}, testResponse)
	})
}

func TestJacer_Post(t *testing.T) {
	testServer := httptest.NewServer(testRouter)
	defer testServer.Close()

	var (
		testJac = getTestJac(testServer)

		testRequest                  = postTestRequestBody{X: 10, Y: 50}
		testAddExpectedResponse      = postTestResponse{Z: 60}
		testMultiplyExpectedResponse = postTestResponse{Z: 500}
	)

	requestAsBytes, err := json.Marshal(testRequest)
	assert.Nil(t, err, "expected nil error when marshalling a request body")

	t.Run("post add request", func(t *testing.T) {
		var testAddResponse postTestResponse
		_, err = testJac.Post(
			RequestParams{
				Endpoint: fmt.Sprintf("%s/%s", testUrlBase, "add"),
				Body:     requestAsBytes,
			},
			&testAddResponse,
		)
		assert.Nil(t, err, "expected nil error when unmarshalling add response")
		assert.Equal(t, testAddResponse, testAddExpectedResponse)
	})
	t.Run("post multiply request", func(t *testing.T) {
		var testMultiplyResponse postTestResponse
		_, err = testJac.Post(
			RequestParams{
				Endpoint: fmt.Sprintf("%s/%s", testUrlBase, "multiply"),
				Body:     requestAsBytes,
			},
			&testMultiplyResponse,
		)
		assert.Nil(t, err, "expected nil error when unmarshalling multiply response")
		assert.Equal(t, testMultiplyResponse, testMultiplyExpectedResponse)
	})
}
