package plex

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	APIBaseURL   = "https://plex.tv/api"
	APIBaseURLv2 = APIBaseURL + "/v2"
)

var (
	// ClientIdentifier is used as the value for the X-Plex-Client-Identifier header
	ClientIdentifier string
	// Product is used as the value for the X-Plex-Product header
	Product string
)

func GetServers(authToken string) ([]Server, error) {
	var (
		client http.Client
		resp   *http.Response
	)
	req, err := http.NewRequest(http.MethodGet, APIBaseURL+"/servers", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Plex-Client-Identifier", ClientIdentifier)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Token", authToken)
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()
	var mcXML MediaContainerXML
	err = xml.Unmarshal(body, &mcXML)
	if err != nil {
		return nil, err
	}
	return mcXML.Servers, nil
}

func GetTVLibraries(baseURL, authToken string) (libraries []Library, err error) {
	var (
		l      = log.WithField("function", "GetTVLibraries")
		client http.Client
		resp   *http.Response
		req    *http.Request
	)
	req, err = http.NewRequest(http.MethodGet, baseURL+"/library/sections", nil)
	if err != nil {
		l.WithError(err).Error("error creating libraries request")
		return
	}
	req.Header.Set("X-Plex-Client-Identifier", ClientIdentifier)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Token", authToken)
	resp, err = client.Do(req)
	if err != nil {
		l.WithError(err).Error("error requesting libraries")
		return
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		l.WithError(err).Error("could not read response body")
		return
	}
	_ = resp.Body.Close()
	var libraryResponse LibraryResponse
	err = json.Unmarshal(body, &libraryResponse)
	if err != nil {
		l.WithError(err).Error("error decoding json response")
		return
	}
	for _, library := range libraryResponse.Data.Sections {
		if library.Type == "show" {
			libraries = append(libraries, library)
		}
	}
	log.Debugf("found %d tv show libraries", len(libraries))
	return
}

func GetTVSeries(baseURL, authToken string, libraryKey int, filter bool) (series []Series, err error) {
	var (
		l       = log.WithField("function", "GetTVSeries")
		client  http.Client
		resp    *http.Response
		req     *http.Request
		fullURL = fmt.Sprintf("%s/library/sections/%d/all", baseURL, libraryKey)
	)
	if filter {
		fullURL = fmt.Sprintf("%s?type=2", fullURL)
	}
	req, err = http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		l.WithError(err).Error("error creating series request")
		return
	}
	req.Header.Set("X-Plex-Client-Identifier", ClientIdentifier)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Token", authToken)
	resp, err = client.Do(req)
	if err != nil {
		l.WithError(err).Error("error requesting series")
		return
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		l.WithError(err).Error("could not read response body")
		return
	}
	_ = resp.Body.Close()
	var seriesResponse SeriesResponse
	err = json.Unmarshal(body, &seriesResponse)
	if err != nil {
		l.WithError(err).Error("error decoding series response")
		return
	}
	series = seriesResponse.Data.Series
	return
}

func Scrobble(url, accessToken string, ratingKey int) {
	var (
		l       = log.WithField("function", "MarkWatched")
		client  http.Client
		fullURL = fmt.Sprintf("%s/:/scrobble?identifier=com.plexapp.plugins.library&key=%d", url, ratingKey)
	)
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		l.WithError(err).Error("error creating series request")
		return
	}
	req.Header.Set("X-Plex-Client-Identifier", ClientIdentifier)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Token", accessToken)
	_, err = client.Do(req)
	if err != nil {
		l.WithError(err).Error("error marking series as watched")
		return
	}
}
