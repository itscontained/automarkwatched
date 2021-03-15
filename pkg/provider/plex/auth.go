package plex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GeneratePin() (*PinData, error) {
	var (
		client http.Client
		resp   *http.Response
	)
	form := url.Values{
		"strong":                   {"true"},
		"X-Plex-Product":           {Product},
		"X-Plex-Client-Identifier": {ClientIdentifier},
	}
	req, err := http.NewRequest(http.MethodPost, APIBaseURLv2+"/pins", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

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

	var pinData PinData
	err = json.Unmarshal(body, &pinData)
	if err != nil {
		return nil, err
	}
	if len(pinData.Errors) > 0 {
		return nil, pinData.Errors[0]

	}
	return &pinData, nil
}

func ConstructAuthAppUrl(code string) (authAppURL string) {
	params := url.Values{
		"clientID":                 {ClientIdentifier},
		"code":                     {code},
		"context[device][product]": {Product},
	}
	return fmt.Sprintf("https://app.plex.tv/auth#?%s", params.Encode())
}

func CheckPin(data PinData) (authToken string, err error) {
	var (
		l      = log.WithField("function", "CheckPin")
		client http.Client
		resp   *http.Response
		req    *http.Request
	)
	form := url.Values{
		"code":                     {data.Code},
		"X-Plex-Client-Identifier": {ClientIdentifier},
	}
	req, err = http.NewRequest(http.MethodGet, APIBaseURLv2+"/pins/"+strconv.Itoa(data.ID), strings.NewReader(form.Encode()))
	if err != nil {
		l.WithError(err).Error("error creating pin request")
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		l.WithError(err).Error("error requesting pin")
		return
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		l.WithError(err).Error("could not read response body")
		return
	}
	_ = resp.Body.Close()

	var pin PinData
	err = json.Unmarshal(body, &pin)
	if err != nil {
		l.WithError(err).Error("error decoding pin response")
		return
	}
	if pin.AuthToken != "" {
		authToken = pin.AuthToken
	}
	return
}

func GetUser(authToken string) (user User, err error) {
	var (
		l      = log.WithField("function", "GetUser")
		client http.Client
		resp   *http.Response
		req    *http.Request
	)
	form := url.Values{
		"X-Plex-Token":             {authToken},
		"X-Plex-Product":           {Product},
		"X-Plex-Client-Identifier": {ClientIdentifier},
	}
	req, err = http.NewRequest(http.MethodGet, APIBaseURLv2+"/user", strings.NewReader(form.Encode()))
	if err != nil {
		l.WithError(err).Error("error creating user request")
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		l.WithError(err).Error("error requesting user")
		return
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		l.WithError(err).Error("could not read response body")
		return
	}
	_ = resp.Body.Close()

	err = json.Unmarshal(body, &user)
	if err != nil {
		l.WithError(err).Error("error decoding user response")
		return
	}
	return
}
