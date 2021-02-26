package owners

import (
	"encoding/json"
	"fmt"
	"net/http"
	urlParse "net/url"
	"path"
	"strings"
	"time"

	sdk "github.com/DeviaVir/servente-sdk"
)

var supportedAPIs = []string{
	"v1",
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// endpoint helpers
type endpoint struct {
	url         string
	enabledAPIs []string
}

func (endpoint *endpoint) v1Request(api *API, p string) (*sdk.JSONData, error) {
	u, err := urlParse.Parse(endpoint.url)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "api", "v1", p)

	sdkData, err := api.request(endpoint.url, u.Path)
	if err != nil {
		return nil, err
	}
	return sdkData, nil
}

// API helpers
func (api *API) discoverPassword(url string) (string, string) {
	p := strings.Split(url, "?")
	return p[0], p[1]
}

func (api *API) request(e, p string) (*sdk.JSONData, error) {
	url, password := api.discoverPassword(e)

	// build URL string from endpoint and path
	u, err := urlParse.Parse(url)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, p)

	// do the HTTP request
	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Servente-Access-Key", password)
	req.Header.Add("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sdkResponse sdk.JSONResponse

	// decode into our target
	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&sdkResponse)
	if err != nil {
		return nil, err
	}

	if sdkResponse.Error {
		return nil, fmt.Errorf(sdkResponse.Data.Message)
	}

	return sdkResponse.Data, nil
}

func (api *API) discoverAPIs(url string) (*sdk.JSONData, error) {
	sdkData, err := api.request(url, "api")
	if err != nil {
		return nil, err
	}
	return sdkData, nil
}

func (api *API) handle(p string) ([]*sdk.JSONData, error) {
	// @TODO: check if our endpoints return a supported API version
	var endpoints []*endpoint
	for _, url := range api.Endpoints {
		data, err := api.discoverAPIs(url)
		if err != nil {
			return nil, err
		}

		var enabledAPIs []string
		for _, version := range data.APIVersions {
			for _, supportedVersion := range supportedAPIs {
				if version.ID == supportedVersion {
					enabledAPIs = append(enabledAPIs, version.ID)
				}
			}
		}

		// skip this one if we don't have any enabled API's
		if len(enabledAPIs) < 1 {
			continue
		}

		endpoint := &endpoint{
			url:         url,
			enabledAPIs: enabledAPIs,
		}
		endpoints = append(endpoints, endpoint)
	}

	var datas []*sdk.JSONData
	for _, endpoint := range endpoints {
		for _, enabledAPI := range endpoint.enabledAPIs {
			switch enabledAPI {
			case "v1":
				data, err := endpoint.v1Request(api, p)
				if err != nil {
					return nil, err
				}
				datas = append(datas, data)
			}
		}
	}

	return datas, nil
}
