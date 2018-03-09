package recaptcha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type (
	recaptcha struct {
		secret string
	}
)

const validationURL = "https://www.google.com/recaptcha/api/siteverify"

func New(secret string) *recaptcha {
	return &recaptcha{secret}
}

func (r *recaptcha) IsHuman(f url.Values) error {
	postData := url.Values{
		"secret":   {r.secret},
		"response": {f.Get("g-recaptcha-response")},
	}

	resp, err := http.PostForm(validationURL, postData)
	if err != nil {
		return fmt.Errorf("failed recaptcha request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read recaptcha response body: %v", err)
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return fmt.Errorf("error reading JSON: %v", err)
	}
	success, ok := m["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("did not validate, response was %v", m)
	}
	return nil
}
