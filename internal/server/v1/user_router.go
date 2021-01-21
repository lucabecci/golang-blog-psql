package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/lucabecci/golang-blog-psql.git/pkg/response"
	"github.com/lucabecci/golang-blog-psql.git/pkg/user"
)

type UserRouter struct {
	Repository user.Repository
}

func (ur UserRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = ur.Repository.Create(ctx, &u)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	u.Password = ""
	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), u.ID))
	response.JSON(w, r, http.StatusCreated, response.Map{"user": u})
}

func (ur *UserRouter) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := ur.Repository.GetAll(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"users": users})
}

func (ur *UserRouter) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(userID)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	u, err := ur.Repository.GetOne(ctx, uint(id))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"user": u})
}

func (ur *UserRouter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(userID)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var u user.User
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = ur.Repository.Update(ctx, uint(id), u)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, nil)
}

func (ur *UserRouter) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(userID)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	err = ur.Repository.Delete(ctx, uint(id))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{})
}

func (ur *UserRouter) Routes() http.Handler {
	router := chi.NewRouter()

	router.Get("/", ur.GetAllHandler)

	router.Post("/", ur.CreateHandler)

	router.Get("/{id}", ur.GetOneHandler)

	router.Put("/{id}", ur.UpdateHandler)

	router.Delete("/{id}", ur.DeleteHandler)
	return router
}
