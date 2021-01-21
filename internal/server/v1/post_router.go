package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/lucabecci/golang-blog-psql.git/pkg/post"
	"github.com/lucabecci/golang-blog-psql.git/pkg/response"
)

type PostRouter struct {
	Repository post.Repository
}

func (pr *PostRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var p post.Post
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = pr.Repository.Create(ctx, &p)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), p.ID))
	response.JSON(w, r, http.StatusCreated, response.Map{"post": p})
}

func (pr *PostRouter) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, err := pr.Repository.GetAll(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"posts": posts})
}

func (pr *PostRouter) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(userID)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	p, err := pr.Repository.GetOne(ctx, uint(id))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"post": p})
}

func (pr *PostRouter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(userID)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var p post.Post
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = pr.Repository.Update(ctx, uint(id), p)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, nil)
}

func (pr *PostRouter) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(userID)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	err = pr.Repository.Delete(ctx, uint(id))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{})
}

func (pr *PostRouter) GetByUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDString := chi.URLParam(r, "userId")

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	posts, err := pr.Repository.GetByUser(ctx, uint(userID))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"posts": posts})
}

func (pr *PostRouter) Routes() http.Handler {
	router := chi.NewRouter()

	router.Get("/user/{userId}", pr.GetByUserHandler)

	router.Get("/", pr.GetAllHandler)

	router.Post("/", pr.CreateHandler)

	router.Get("/{id}", pr.GetOneHandler)

	router.Put("/{id}", pr.UpdateHandler)

	router.Delete("/{id}", pr.DeleteHandler)

	return router
}
