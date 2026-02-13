package fakturoid

import (
	"fmt"
	"net/url"
)

func (c *Client) GetInvoices(page int, params url.Values) ([]Invoice, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("page", fmt.Sprintf("%d", page))
	var result []Invoice
	err := c.do("GET", fmt.Sprintf("/invoices.json?%s", params.Encode()), nil, &result)
	return result, err
}

func (c *Client) GetInvoice(id int) (*Invoice, error) {
	var result Invoice
	err := c.do("GET", fmt.Sprintf("/invoices/%d.json", id), nil, &result)
	return &result, err
}

func (c *Client) SearchInvoices(query string, page int) ([]Invoice, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("page", fmt.Sprintf("%d", page))
	var result []Invoice
	err := c.do("GET", fmt.Sprintf("/invoices/search.json?%s", params.Encode()), nil, &result)
	return result, err
}

func (c *Client) CreateInvoice(req CreateInvoiceRequest) (*Invoice, error) {
	var result Invoice
	err := c.do("POST", "/invoices.json", req, &result)
	return &result, err
}

func (c *Client) UpdateInvoice(id int, req UpdateInvoiceRequest) (*Invoice, error) {
	var result Invoice
	err := c.do("PATCH", fmt.Sprintf("/invoices/%d.json", id), req, &result)
	return &result, err
}

func (c *Client) DeleteInvoice(id int) error {
	return c.do("DELETE", fmt.Sprintf("/invoices/%d.json", id), nil, nil)
}

func (c *Client) GetInvoicePayments(invoiceID int) ([]InvoicePayment, error) {
	var result []InvoicePayment
	err := c.do("GET", fmt.Sprintf("/invoices/%d/payments.json", invoiceID), nil, &result)
	return result, err
}
