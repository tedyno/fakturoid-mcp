package fakturoid

func (c *Client) GetAccount() (*Account, error) {
	var result Account
	err := c.do("GET", "/account.json", nil, &result)
	return &result, err
}

func (c *Client) GetEvents(page int) ([]Event, error) {
	var result []Event
	err := c.do("GET", "/events.json", nil, &result)
	return result, err
}
