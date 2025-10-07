package tests

import (
	"net/http"
	"net/url"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/random"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	e := httpexpect.Default(t, "http://"+host)

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("user", "pass").
		Expect().
		Status(http.StatusCreated).
		JSON().
		Object().
		ContainsKey("alias")
}

func TestURLShortener_SaveRedirectRemove(t *testing.T) {
	testcases := []struct {
		name               string
		url                string
		alias              string
		error              string
		expectedStatusCode int
	}{
		{
			name:               "Valid url",
			url:                gofakeit.URL(),
			alias:              gofakeit.Word() + gofakeit.Word(),
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Invalid url",
			url:                "test",
			alias:              gofakeit.Word(),
			error:              "field URL is not a valid URL",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Empty alias",
			url:                gofakeit.URL(),
			alias:              "",
			expectedStatusCode: http.StatusCreated,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).WithBasicAuth("user", "pass").
				Expect().Status(tc.expectedStatusCode).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEqual(tc.alias)
				alias = resp.Value("alias").String().Raw()
			}

			testRedirect(t, alias, tc.url)
			if tc.alias != "" {
				resp = e.DELETE("/url/"+tc.alias).
					WithBasicAuth("user", "pass").
					Expect().JSON().Object()

				if tc.error != "" {
					resp.NotContainsKey("alias")

					resp.Value("error").String().IsEqual(tc.error)
				}

				require.Equal(t, alias, resp.Value("alias").String().Raw())

				e.GET(u.String() + "/" + tc.alias).Expect().Status(http.StatusNotFound)
			}

		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)
	require.Equal(t, redirectedURL, urlToRedirect)
}
