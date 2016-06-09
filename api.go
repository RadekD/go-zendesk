package zendesk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Auth struct {
	Username  string
	Password  string
	Subdomain string

	UseAPIToken bool
}

//NewAuth creates Auth
func NewAuth(Username, Password, Subdomain string, UseAPIToken bool) *Auth {
	return &Auth{Username, Password, Subdomain, UseAPIToken}
}

/*
{
  "details": {
    "value": [
      {
        "type": "blank",
        "description": "can't be blank"
      },
      {
        "type": "invalid",
        "description": " is not properly formatted"
      }
    ]
  },
  "description": "RecordValidation errors",
  "error": "RecordInvalid"
}
*/

//Error represents zendesk api error
type Error struct {
	Description string `json:"description"`
	Details     struct {
		Value []struct {
			Type        string
			Description string
		} `json:"value"`
	}
}

func (e Error) Error() string {
	return e.Description
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
	if auth.UseAPIToken {
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
	if resp.StatusCode >= 400 {
		err := Error{}

		err2 := json.Unmarshal(data, &err)
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}
	return data, nil
}
