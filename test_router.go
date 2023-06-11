package jac

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

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
