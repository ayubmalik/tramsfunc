package tramsfunc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(TFGM_API_KEY_HEADER) != "some key" {
			w.WriteHeader(401)
			return
		}

		if r.URL.Path == "/Metrolinks" {
			fmt.Fprintf(w, "{ \"Value\": [{ \"ID\": 99, \"Line\": \"line 99\"}] }")
			return
		}

		re := regexp.MustCompile(`/Metrolinks/(\d+)`)
		matches := re.FindStringSubmatch(r.URL.Path)
		if len(matches) == 2 {
			fmt.Fprintf(w, `{"ID": %s, "Line": "line %s"}`, matches[1], matches[1])
			return
		}
		t.Errorf("unexpected method call %v", r.URL)
	}))

	t.Cleanup(func() {
		ts.Close()
	})

	client := newClient("some key", ts.URL)

	t.Run("allMetrolinks", func(t *testing.T) {
		metrolinks, err := client.allMetrolinks()
		if err != nil {
			t.Errorf("got error want no err: %v", err)
		}

		if len(metrolinks) != 1 {
			t.Errorf("got empty result want at least 1")
		}
	})

	t.Run("metrolinksById", func(t *testing.T) {
		metrolinks, err := client.metrolinksById("3", "5", "7")
		if err != nil {
			t.Errorf("got error want no err: %v", err)
		}

		if len(metrolinks) != 3 {
			t.Errorf("got %d results want %d", len(metrolinks), 3)
		}

	})
}

func TestAPI(t *testing.T) {

	t.Run("when API key not set", func(t *testing.T) {
		os.Setenv(TFGM_API_KEY, "")
		os.Setenv(TFGM_API_URL, "some url")

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		API(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	})

	t.Run("when API URL not set", func(t *testing.T) {
		os.Setenv(TFGM_API_KEY, "some key")
		os.Setenv(TFGM_API_URL, "")

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		API(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	})

	t.Run("calls all metrolinks when no ID param", func(t *testing.T) {
		os.Setenv(TFGM_API_KEY, "some key")
		os.Setenv(TFGM_API_URL, "http://localhost")

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		API(res, req)

		// TODO: use http server but for now check
		// expected URL contains path even though
		// we probably got connections refuse error :)
		assert.Contains(t, res.Body.String(), "/Metrolinks")
	})

	t.Run("calls metrolinks by ID", func(t *testing.T) {
		os.Setenv(TFGM_API_KEY, "some key")
		os.Setenv(TFGM_API_URL, "http://localhost")

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/?id=3&id=5&id=7", nil)

		API(res, req)

		// TODO: use http server
		assert.Contains(t, res.Body.String(), "/Metrolinks/3")
	})
}
