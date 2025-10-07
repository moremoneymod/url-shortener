package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=URLGetter
type URLGetter interface {
	GetURL(string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)
		// привязались к роутеру
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Error("invalid request"))

			return
		}
		resUrl, err := urlGetter.GetURL(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Error("url not found"))

			return

		}

		if err != nil {
			log.Info("failed to get url", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resUrl))

		http.Redirect(w, r, resUrl, http.StatusFound)
	}
}
