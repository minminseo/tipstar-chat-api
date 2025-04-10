package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/minminseo/tipstar-chat-api/adapter/rest"
	"github.com/minminseo/tipstar-chat-api/usecase"
)

func NewRouter(
	editUC usecase.EditMessageInput,
	deleteUC usecase.DeleteMessageInput,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(rest.JWTMiddleware)

	r.Put("/message/edit", rest.NewEditMessageHandler(editUC).ServeHTTP)
	r.Delete("/message/delete", rest.NewDeleteMessageHandler(deleteUC).ServeHTTP)

	return r
}
