package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	fallbackResponse = "fallback"
	path             = "/test"
	url             = "https://test.com"
)

func TestMapHandler(t *testing.T) {
	// Arrange
	pathToUrls := map[string]string{path: url}

	// Act
	t.Run("it redirects to correct url", func(t *testing.T) {
		result := runMapHandler(pathToUrls, path)

		// Assert
		assertStatus(t, result, http.StatusFound)
		assertURL(t, result, url)
	})
}

func fallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fallbackResponse)
}

func runMapHandler(pathToUrls map[string]string, path string) *http.Response {
	request, _ := http.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	mapHandler := createMapHandler(pathToUrls)
	mapHandler(response, request)

	return response.Result()
}

func createMapHandler(pathToUrls map[string]string) http.HandlerFunc {
	fallbackHandler := http.HandlerFunc(fallback)
	return MapHandler(pathToUrls, fallbackHandler)
}

func assertStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		t.Errorf("Expected status to be %d, got %d",
			want, resp.StatusCode)
	}
}

func assertURL(t *testing.T, resp *http.Response, want string) {
	t.Helper()
	url, err := resp.Location()

	if err != nil {
		t.Fatal("Could not read location", err)
	}

	if url.String() != want {
		t.Errorf("Expected url to be %s, got %s", url, want)
	}
}
