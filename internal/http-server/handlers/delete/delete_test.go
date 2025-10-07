package delete_test

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	delete2 "url-shortener/internal/http-server/handlers/delete"
	"url-shortener/internal/http-server/handlers/delete/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type args struct {
		log        *slog.Logger
		urlDeleter delete2.URLDeleter
	}
	tests := []struct {
		name         string
		alias        string
		respError    string
		mockError    error
		expectedCode int
		shouldMock   bool
	}{
		{
			name:         "Success",
			alias:        "test_alias",
			expectedCode: http.StatusOK,
			shouldMock:   true,
		},
		{
			name:         "Empty alias",
			alias:        "",
			expectedCode: http.StatusNotFound,
			shouldMock:   false,
		},
		{
			name:         "DeleteURL error",
			alias:        "test_alias",
			respError:    "internal error",
			mockError:    errors.New("internal error"),
			expectedCode: http.StatusOK,
			shouldMock:   true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			URLDeleterMock := mocks.NewURLDeleter(t)

			if tc.shouldMock {
				URLDeleterMock.On("DeleteURL", mock.AnythingOfType("string")).
					Return(tc.mockError).
					Once()
			}

			handler := delete2.New(slogdiscard.NewDiscardLogger(), URLDeleterMock)

			router := chi.NewRouter()
			router.Delete("/delete/{alias}", handler)

			req, err := http.NewRequest(http.MethodDelete, "/delete/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code)

			if tc.shouldMock {
				body := rr.Body.String()

				var resp delete2.Response

				require.NoError(t, json.Unmarshal([]byte(body), &resp))

				require.Equal(t, tc.respError, resp.Error)
			}
		})
	}
}
