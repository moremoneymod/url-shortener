package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name               string
		alias              string
		url                string
		expectedStatusCode int
		respError          string
		mockError          error
	}{
		{
			name:               "success",
			alias:              "test_alias",
			url:                "http://example.com/",
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Empty alias",
			alias:              "",
			url:                "http://example.com/",
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Empty url",
			url:                "",
			alias:              "test_alias",
			respError:          "field URL is a required field",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid URL",
			url:                "invalid URL",
			alias:              "test_alias",
			respError:          "field URL is not a valid URL",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "SaveURl error",
			alias:              "test_alias",
			url:                "http://example.com/",
			respError:          "failed to add url",
			mockError:          errors.New("unexpected error"),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "URL already exists",
			alias:              "test_alias",
			url:                "http://example.com/",
			respError:          "url already exists",
			mockError:          storage.ErrURLExist,
			expectedStatusCode: http.StatusConflict,
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.expectedStatusCode)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

		})
	}
}
