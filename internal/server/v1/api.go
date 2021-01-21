package v1

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/lucabecci/golang-blog-psql.git/internal/data"
)

func New() http.Handler {
	router := chi.NewRouter()

	ur := &UserRouter{Repository: &data.UserRespository{Data: data.New()}}
	pr := &PostRouter{Repository: &data.PostRepository{Data: data.New()}}

	router.Mount("/users", ur.Routes())
	router.Mount("/posts", pr.Routes())

	return router
}
