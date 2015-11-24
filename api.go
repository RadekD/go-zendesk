package zendesk

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

type Auth struct {
	Username  string
	Password  string
	Subdomain string

	UseApiToken bool
}

func api(auth Auth, method string, path string, params string) ([]byte, error) {
	r := &http.Transport{}
	client := &http.Client{
		Transport: r,
	}

	var URL string

	if strings.HasPrefix(path, "http") {
		URL = path
	} else {
		if strings.HasPrefix(auth.Subdomain, "http") {
			URL = auth.Subdomain + "/api/v2" + path
		} else {
			URL = "https://" + auth.Subdomain + "/api/v2" + path
		}
	}

	req, err := http.NewRequest(method, URL, bytes.NewBufferString(params))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	if auth.UseApiToken {
		req.SetBasicAuth(auth.Username+"/token", auth.Password)
	} else {
		req.SetBasicAuth(auth.Username, auth.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	return data, nil
}
