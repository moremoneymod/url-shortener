package redirect_test

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/redirect/mocks"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type args struct {
		log       *slog.Logger
		urlGetter redirect.URLGetter
	}
	tests := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test",
			url:   "http://example.com/",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).
					Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToUrl, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			require.Equal(t, tc.url, redirectedToUrl)

		})
	}
}

func TestNew_EmptyAlias(t *testing.T) {
	urlGetterMock := mocks.NewURLGetter(t)

	handler := redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock)

	r := chi.NewRouter()
	r.Get("/{alias}", handler)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	urlGetterMock.AssertExpectations(t)

}

func TestNew_URLNotFound(t *testing.T) {
	urlGetterMock := mocks.NewURLGetter(t)

	alias := "non_existing_alias"
	urlGetterMock.On("GetURL", alias).
		Return("", errors.New("url not found")).
		Once()

	handler := redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock)

	r := chi.NewRouter()
	r.Get("/{alias}", handler)

	req, err := http.NewRequest(http.MethodGet, "/"+alias, nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	urlGetterMock.AssertExpectations(t)

}
