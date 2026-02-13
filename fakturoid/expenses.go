package fakturoid

import (
	"fmt"
	"net/url"
)

func (c *Client) GetExpenses(page int, params url.Values) ([]Expense, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("page", fmt.Sprintf("%d", page))
	var result []Expense
	err := c.do("GET", fmt.Sprintf("/expenses.json?%s", params.Encode()), nil, &result)
	return result, err
}

func (c *Client) GetExpense(id int) (*Expense, error) {
	var result Expense
	err := c.do("GET", fmt.Sprintf("/expenses/%d.json", id), nil, &result)
	return &result, err
}
