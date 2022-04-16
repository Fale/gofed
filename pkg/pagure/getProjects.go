package pagure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type GetProjectsParameters struct {
	BaseUrl   string
	NameSpace string
	Pattern   string
	Fork      string
}

func GetProjects(pars GetProjectsParameters) ([]Project, error) {
	projects := []Project{}
	values := url.Values{}
	values.Set("per_page", "100")
	if len(pars.NameSpace) > 0 {
		values.Set("namespace", pars.NameSpace)
	}
	if len(pars.Pattern) > 0 {
		values.Set("pattern", pars.Pattern)
	}
	if len(pars.Fork) > 0 {
		values.Set("fork", pars.Fork)
	}
	page := 1
	for {
		t := projectPage{}
		values.Set("page", strconv.Itoa(page))
		url := fmt.Sprintf("%v/api/0/projects?%v", pars.BaseUrl, values.Encode())
		if err := getPage(url, &t); err != nil {
			return nil, err
		}
		projects = append(projects, t.Projects...)
		// fmt.Printf("page %v, pages %v\n", t.Pagination.Page, t.Pagination.Pages)
		if page == t.Pagination.Pages {
			break
		}
		page++
	}
	return projects, nil
}

func getPage(url string, body interface{}) error {
	httpClient := &http.Client{Timeout: 60 * time.Second}
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(body)
}

type projectPage struct {
	Args       projectArgs
	Pagination projectPagination
	Projects   []Project
}

type projectArgs struct{}

type projectPagination struct {
	First    string
	Last     string
	Next     *string
	Page     int
	Pages    int
	Per_page int
	Prev     *string
}

type Project struct {
	AccessGroups map[string][]string `json:"access_groups"`
	AccessUsers  map[string][]string `json:"access_users"`
	FullName     string              `json:"fullname"`
	Name         string              `json:"name"`
	Namespace    string              `json:"namespace"`
	User         struct {
		FullUrl  string `json:"full_url"`
		FullName string `json:"fullname"`
		Name     string `json:"name"`
		UrlPath  string `json:"url_path"`
	}
}
