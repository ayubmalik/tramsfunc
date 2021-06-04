package tramsfunc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(API_KEY_HEADER) != "some key" {
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
