package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	root = "https://pypi.org/pypi"
)

// Client is a type the describes a PyPi repository
type Client struct {
	client *http.Client
}

// NewClient is a function creates a new Repository
func NewClient(c *http.Client) *Client {
	return &Client{
		client: c,
	}
}

// Response is a type that represents the result of a PyPi JSON API Project|Release request
// See: https://warehouse.readthedocs.io/api-reference/json/
// TODO(dazwilkin) Incomplete
type Response struct {
	Info       Info     `json:"info"`
	LastSerial int      `json:"last_serial"`
	Releases   Releases `json:"releases"`
	URLs       URLs     `json:"urls"`
}

// Info is a sub-type of Response
// TODO(dazwilkin) Incomplete
type Info struct {
	Author      string   `json:"author"`
	AuthorEmail string   `json:"author_email"`
	Classifiers []string `json:"classifiers"`
	PackageURL  string   `json:"package_url"`
}

// Releases is a sub-type of Response
type Releases map[string]Packages

// Packages is a sub-type of Response
type Packages []Package

// Package is a sub-type of Response
type Package struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

// URLs is a sub-type of Response
type URLs []Package

// Get is a function that GETs a path from the PyPi API endpoint returning a Response
func (c *Client) Get(path string) (Response, error) {
	resp, err := c.client.Get(path)
	if err != nil {
		return Response{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	result := Response{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Response{}, err
	}
	return result, nil

}

// Project is a function that uses PyPi JSON API 'Project' endpoint returning a Response
func (c *Client) Project(name string) (Response, error) {
	return c.Get(fmt.Sprintf("%s/%s/json", root, name))
}

// Release is a function that uses PyPi JSON API 'Release' endpoint returning a Response
func (c Client) Release(name, version string) (Response, error) {
	return c.Get(fmt.Sprintf("%s/%s/%s/json", root, name, version))
}

// Packages is a function that returns the Response's version's Packages
func (r Response) Packages(version string) (Packages, error) {
	if packages, ok := r.Releases[version]; ok {
		return packages, nil
	}
	return nil, fmt.Errorf("Version [%s] not found", version)
}

// Package is a function that returns the unique Package matching 'c'
func (p Packages) Package(c Package) (Package, error) {
	if c.Filename == "" && c.URL == "" {
		return Package{}, fmt.Errorf("undefined criteria will never match")
	}
	var match func(Package) bool
	if c.Filename != "" && c.URL != "" {
		match = func(p Package) bool {
			return p.Filename == c.Filename && p.URL == c.URL
		}
	}
	if c.Filename == "" {
		match = func(p Package) bool {
			return p.URL == c.URL
		}
	}
	if c.URL == "" {
		match = func(p Package) bool {
			return p.Filename == c.Filename
		}
	}

	for _, e := range p {
		if match(e) {
			return e, nil
		}
	}
	return Package{}, fmt.Errorf("criteria defined in Package [%v] not found", c)
}
