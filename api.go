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
	TFGM_API_KEY        = "TFGM_API_KEY"
	TFGM_API_URL        = "TFGM_API_URL"
	TFGM_API_KEY_HEADER = "Ocp-Apim-Subscription-Key"
)

// API is the entry point for GCP Functions.
func API(w http.ResponseWriter, r *http.Request) {

	var apiKey, apiURL string

	if apiKey = os.Getenv(TFGM_API_KEY); apiKey == "" {
		handleError(w, http.StatusInternalServerError, "TFGM API Key is not set")
		return
	}

	if apiURL = os.Getenv(TFGM_API_URL); apiURL == "" {
		handleError(w, http.StatusInternalServerError, "TFGM API URL is not set")
		return
	}

	var (
		client     *tfgmClient = newTFGMClient(apiKey, apiURL)
		ids        []string    = r.URL.Query()["id"]
		metrolinks []Metrolink
		err        error
	)

	if len(ids) == 0 {
		metrolinks, err = client.allMetrolinks()
	} else {
		metrolinks, err = client.metrolinksById(ids...)
	}

	if err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
	}

	if err = json.NewEncoder(w).Encode(metrolinks); err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
	}
}

func handleError(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	fmt.Fprint(w, err)
}

type tfgmClient struct {
	key        string
	url        string
	httpClient http.Client
}

func newTFGMClient(apiKey, apiURL string) *tfgmClient {
	httpClient := http.Client{Timeout: 3 * time.Second}

	return &tfgmClient{
		key:        apiKey,
		url:        apiURL,
		httpClient: httpClient,
	}
}

// allMetrolinks() returns all available metrolinks
func (c tfgmClient) allMetrolinks() ([]Metrolink, error) {
	var metrolinks Metrolinks
	err := c.callAPI("/Metrolinks", &metrolinks)
	if err != nil {
		return nil, err
	}
	return metrolinks.Value, nil
}

// metrolinksById returns metrolinks for the given IDs
func (c tfgmClient) metrolinksById(ids ...string) ([]Metrolink, error) {
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

func (c tfgmClient) callAPI(path string, value interface{}) error {
	req, err := http.NewRequest(http.MethodGet, c.url+path, nil)
	if err != nil {
		return err
	}
	req.Header.Add(TFGM_API_KEY_HEADER, c.key)

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

// Metrolinks is the JSON response object from the TTFGM API
type Metrolinks struct {
	Value []Metrolink `json:"value"`
}

// Metrolink represents a TFGM Metrolink tram location
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
