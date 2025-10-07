package delete

import (
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Error("invalid request"))

			return
		}
		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Info("failed to delete url", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("deleted alias", slog.String("alias", alias))
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
