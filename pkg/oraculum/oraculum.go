package oraculum

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type response struct {
	Error    interface{}
	Packages map[string]Package
	Status   int
}

func LookupPackage(name string) (*Package, error) {
	u := fmt.Sprintf("https://packager-dashboard.fedoraproject.org/api/v2/packager_dashboard?packages=%v", name)
	httpClient := &http.Client{Timeout: 60 * time.Second}
	r, err := httpClient.Get(u)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var rs response
	if err := json.Unmarshal(body, &rs); err != nil {
		return nil, err
	}
	for n, p := range rs.Packages {
		if n == name {
			return &p, nil
		}
	}
	return nil, errors.New("package not found")
}
