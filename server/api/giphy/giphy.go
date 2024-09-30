package giphy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const giphyAPIURL = "https://api.giphy.com/v1/gifs/search"

type GiphyResponse struct {
	Data []struct {
		ID     string `json:"id"`
		URL    string `json:"url"`
		Images struct {
			FixedHeight struct {
				URL string `json:"url"`
			} `json:"fixed_height"`
		} `json:"images"`
	} `json:"data"`
}

func SearchGiphyHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	searchTerm := queryParams.Get("q")
	limit := queryParams.Get("limit")

	if searchTerm == "" {
		http.Error(w, "Missing 'q' parameter", http.StatusBadRequest)
		return
	}

	apiKey := "eXTvJ6YPqsfWiNMYoftLpQrz8u9Za4ee"
	if apiKey == "" {
		http.Error(w, "GIPHY_API_KEY not set", http.StatusInternalServerError)
		return
	}

	params := url.Values{}
	params.Add("api_key", apiKey)
	params.Add("q", searchTerm)
	params.Add("limit", limit)

	fullURL := fmt.Sprintf("%s?%s", giphyAPIURL, params.Encode())

	fmt.Println("Request URL:", fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		http.Error(w, "Error fetching from Giphy", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	var giphyResp GiphyResponse
	if err := json.Unmarshal(body, &giphyResp); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(giphyResp.Data)
}
