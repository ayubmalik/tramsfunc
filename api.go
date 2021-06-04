package tramsfunc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	API_KEY_HEADER = "Ocp-Apim-Subscription-Key"
)

// API is the entry point for GCP Functions.
func API(w http.ResponseWriter, r *http.Request) {
	var apiKey string
	if apiKey = os.Getenv("TFGM_API_KEY"); apiKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "TFGM API Key is not set")
		return
	}
	fmt.Fprintf(w, "Hello, World! API KEY = %s\n", apiKey)
}

type client struct {
	key        string
	url        string
	httpClient http.Client
}

func newClient(apiKey, apiURL string) *client {
	httpClient := http.Client{Timeout: 3 * time.Second}

	return &client{
		key:        apiKey,
		url:        apiURL,
		httpClient: httpClient,
	}
}

func (c client) allMetrolinks() ([]Metrolink, error) {
	var metrolinks Metrolinks
	err := c.callAPI("/Metrolinks", &metrolinks)
	if err != nil {
		return nil, err
	}
	return metrolinks.Value, nil
}

func (c client) metrolinksById(ids ...string) ([]Metrolink, error) {
	metrolinks := make([]Metrolink, 0)
	for _, id := range ids {
		var m Metrolink
		err := c.callAPI("/Metrolinks/"+id, &m)
		if err != nil {
			return nil, err
		}
		metrolinks = append(metrolinks, m)
	}
	return metrolinks, nil
}

func (c client) callAPI(path string, value interface{}) error {
	req, err := http.NewRequest(http.MethodGet, c.url+path, nil)
	if err != nil {
		return err
	}
	req.Header.Add(API_KEY_HEADER, c.key)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(value)
	if err != nil {
		return err
	}
	return nil
}

type Metrolinks struct {
	Value []Metrolink `json:"value"`
}

type Metrolink struct {
	Id              int
	Line            string
	TLAREF          string
	PIDREF          string
	StationLocation string
	AtcoCode        string
	Direction       string
	Dest0           string
	Carriages0      string
	Status0         string
	Wait0           string
	Dest1           string
	Carriages1      string
	Status1         string
	Wait1           string
	Dest2           string
	Carriages2      string
	Status2         string
	Wait2           string
	Dest3           string
	Carriages3      string
	Status3         string
	MessageBoard    string
	Wait3           string
	LastUpdated     time.Time
}
