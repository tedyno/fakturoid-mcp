package fakturoid

import (
	"fmt"
	"net/url"
)

func (c *Client) GetSubjects(page int) ([]Subject, error) {
	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))
	var result []Subject
	err := c.do("GET", fmt.Sprintf("/subjects.json?%s", params.Encode()), nil, &result)
	return result, err
}

func (c *Client) GetSubject(id int) (*Subject, error) {
	var result Subject
	err := c.do("GET", fmt.Sprintf("/subjects/%d.json", id), nil, &result)
	return &result, err
}

func (c *Client) SearchSubjects(query string, page int) ([]Subject, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("page", fmt.Sprintf("%d", page))
	var result []Subject
	err := c.do("GET", fmt.Sprintf("/subjects/search.json?%s", params.Encode()), nil, &result)
	return result, err
}

func (c *Client) CreateSubject(req CreateSubjectRequest) (*Subject, error) {
	var result Subject
	err := c.do("POST", "/subjects.json", req, &result)
	return &result, err
}

func (c *Client) UpdateSubject(id int, req UpdateSubjectRequest) (*Subject, error) {
	var result Subject
	err := c.do("PATCH", fmt.Sprintf("/subjects/%d.json", id), req, &result)
	return &result, err
}

func (c *Client) DeleteSubject(id int) error {
	return c.do("DELETE", fmt.Sprintf("/subjects/%d.json", id), nil, nil)
}
