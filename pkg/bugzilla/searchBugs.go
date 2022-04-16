package bugzilla

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Bug struct {
	AssignedTo       string                   `json:"assigned_to"`
	AssignedToDetail map[string]interface{}   `json:"assigned_to_detail"`
	CC               []string                 `json:"cc"`
	CCDetail         []map[string]interface{} `json:"cc_detail"`
	DependsOn        []int                    `json:"depends_on"`
	Summary          string                   `json:"summary"`
}

type bugs struct {
	Bugs []Bug
}

type SearchBugRequest struct {
	Products   []string
	Components []string
	Statuses   []string
}

func SearchBug(baseurl string, req SearchBugRequest) ([]Bug, error) {
	u, err := url.Parse(baseurl + "/rest/bug")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for _, i := range req.Products {
		q.Add("product", i)
	}
	for _, i := range req.Components {
		q.Add("component", i)
	}
	for _, i := range req.Statuses {
		q.Add("status", i)
	}
	u.RawQuery = q.Encode()
	httpClient := &http.Client{Timeout: 60 * time.Second}
	r, err := httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var bs bugs
	if err := json.Unmarshal(body, &bs); err != nil {
		return nil, err
	}
	return bs.Bugs, nil
}
