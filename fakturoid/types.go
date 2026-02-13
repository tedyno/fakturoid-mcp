package fakturoid

import "encoding/json"

// --- Invoice ---

type Invoice struct {
	ID              int          `json:"id"`
	Number          string       `json:"number"`
	SubjectID       int          `json:"subject_id"`
	Status          string       `json:"status"`
	DocumentType    string       `json:"document_type"`
	IssuedOn        string       `json:"issued_on"`
	TaxableFulfillmentDue string `json:"taxable_fulfillment_due"`
	DueOn           string       `json:"due_on"`
	PaidOn          string       `json:"paid_on,omitempty"`
	Note            string       `json:"note,omitempty"`
	FootNote        string       `json:"footer_note,omitempty"`
	Currency        string       `json:"currency"`
	NativeTotal     string       `json:"native_total"`
	Total           string       `json:"total"`
	RemainingAmount string       `json:"remaining_amount"`
	Lines           []InvoiceLine `json:"lines,omitempty"`
	SubjectName     string       `json:"subject_name,omitempty"`
}

type InvoiceLine struct {
	Name      string      `json:"name"`
	Quantity  json.Number `json:"quantity"`
	UnitName  string      `json:"unit_name,omitempty"`
	UnitPrice json.Number `json:"unit_price"`
	VATRate   json.Number `json:"vat_rate"`
}

type CreateInvoiceRequest struct {
	SubjectID int           `json:"subject_id"`
	Lines     []InvoiceLine `json:"lines"`
	Currency  string        `json:"currency,omitempty"`
	Note      string        `json:"note,omitempty"`
	DueOn     string        `json:"due_on,omitempty"`
	IssuedOn  string        `json:"issued_on,omitempty"`
	TaxableFulfillmentDue string `json:"taxable_fulfillment_due,omitempty"`
}

type UpdateInvoiceRequest struct {
	SubjectID *int          `json:"subject_id,omitempty"`
	Lines     []InvoiceLine `json:"lines,omitempty"`
	Currency  string        `json:"currency,omitempty"`
	Note      string        `json:"note,omitempty"`
	DueOn     string        `json:"due_on,omitempty"`
}

type SendInvoiceRequest struct {
	Email     string `json:"email"`
	EmailCopy string `json:"email_copy,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Message   string `json:"message,omitempty"`
}

// --- Subject (Contact) ---

type Subject struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Street           string `json:"street,omitempty"`
	City             string `json:"city,omitempty"`
	Zip              string `json:"zip,omitempty"`
	Country          string `json:"country,omitempty"`
	RegistrationNo   string `json:"registration_no,omitempty"`
	VATNo            string `json:"vat_no,omitempty"`
	Email            string `json:"email,omitempty"`
	Phone            string `json:"phone,omitempty"`
	FullName         string `json:"full_name,omitempty"`
	Type             string `json:"type,omitempty"`
}

type CreateSubjectRequest struct {
	Name           string `json:"name"`
	Street         string `json:"street,omitempty"`
	City           string `json:"city,omitempty"`
	Zip            string `json:"zip,omitempty"`
	Country        string `json:"country,omitempty"`
	RegistrationNo string `json:"registration_no,omitempty"`
	VATNo          string `json:"vat_no,omitempty"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
}

type UpdateSubjectRequest struct {
	Name           string `json:"name,omitempty"`
	Street         string `json:"street,omitempty"`
	City           string `json:"city,omitempty"`
	Zip            string `json:"zip,omitempty"`
	Country        string `json:"country,omitempty"`
	RegistrationNo string `json:"registration_no,omitempty"`
	VATNo          string `json:"vat_no,omitempty"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
}

// --- Expense ---

type Expense struct {
	ID              int          `json:"id"`
	Number          string       `json:"number"`
	OriginalNumber  string       `json:"original_number,omitempty"`
	SubjectID       int          `json:"subject_id"`
	Status          string       `json:"status"`
	IssuedOn        string       `json:"issued_on"`
	DueOn           string       `json:"due_on"`
	PaidOn          string       `json:"paid_on,omitempty"`
	Currency        string       `json:"currency"`
	NativeTotal     string       `json:"native_total"`
	Total           string       `json:"total"`
	Lines           []ExpenseLine `json:"lines,omitempty"`
	SubjectName     string       `json:"subject_name,omitempty"`
}

type ExpenseLine struct {
	Name      string      `json:"name"`
	Quantity  json.Number `json:"quantity"`
	UnitName  string      `json:"unit_name,omitempty"`
	UnitPrice json.Number `json:"unit_price"`
	VATRate   json.Number `json:"vat_rate"`
}

// --- Account ---

type Account struct {
	Subdomain        string `json:"subdomain"`
	Plan             string `json:"plan"`
	PlanPrice        int    `json:"plan_price"`
	Email            string `json:"email"`
	InvoiceEmail     string `json:"invoice_email,omitempty"`
	Phone            string `json:"phone,omitempty"`
	Name             string `json:"name"`
	RegistrationNo   string `json:"registration_no,omitempty"`
	VATNo            string `json:"vat_no,omitempty"`
	Street           string `json:"street,omitempty"`
	City             string `json:"city,omitempty"`
	Zip              string `json:"zip,omitempty"`
	Country          string `json:"country,omitempty"`
	Currency         string `json:"currency"`
}

// --- Event ---

type Event struct {
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

// --- InvoicePayment ---

type InvoicePayment struct {
	ID       int    `json:"id"`
	PaidOn   string `json:"paid_on"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}
