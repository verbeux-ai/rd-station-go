package rd_station

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrRequestFailed    = errors.New("http request execution failed")
	ErrReadResponseBody = errors.New("failed to read response body")
	ErrApiReturnedError = errors.New("api returned an error status")
	ErrDecodeResponse   = errors.New("failed to decode api response")
)

type Contact struct {
	Birthday            BirthdayResponse     `json:"birthday"`
	ContactCustomFields []ContactCustomField `json:"contact_custom_fields"`
	CreatedAt           string               `json:"created_at"`
	Deals               []ContactDeal        `json:"deals"`
	Emails              []Email              `json:"emails"`
	Facebook            *string              `json:"facebook"`
	ID                  string               `json:"id"`
	LegalBases          []LegalBasis         `json:"legal_bases"`
	LinkedIn            *string              `json:"linkedin"`
	Name                string               `json:"name"`
	Notes               string               `json:"notes"`
	OrganizationID      *string              `json:"organization_id"`
	Phones              []Phone              `json:"phones"`
	Skype               *string              `json:"skype"`
	Title               *string              `json:"title"`
	UpdatedAt           string               `json:"updated_at"`
}

type ContactCustomField struct {
	ID            string `json:"_id"`
	CreatedAt     string `json:"created_at"`
	CustomFieldID string `json:"custom_field_id"`
	UpdatedAt     string `json:"updated_at"`
	Value         string `json:"value"`
}

type ContactDeal struct {
	DealID           string `json:"_id"`
	ClosedAt         string `json:"closed_at"`
	DealLostReasonID string `json:"deal_lost_reason_id"`
	ID               string `json:"id"`
	Name             string `json:"name"`
	PredictionDate   string `json:"prediction_date"`
	Win              bool   `json:"win"`
}

type Email struct {
	ID        string `json:"_id"`
	CreatedAt string `json:"created_at"`
	Email     string `json:"email"`
	UpdatedAt string `json:"updated_at"`
}

type LegalBasis struct {
	Category string `json:"category"`
	Status   string `json:"status"`
	Type     string `json:"type"`
}

type Phone struct {
	CreatedAt                 string `json:"created_at"`
	Phone                     string `json:"phone"`
	Type                      string `json:"type"`
	UpdatedAt                 string `json:"updated_at"`
	WhatsApp                  bool   `json:"whatsapp"`
	WhatsAppFullInternacional string `json:"whatsapp_full_internacional"`
	WhatsAppURLWeb            string `json:"whatsapp_url_web"`
}

type ListContactsFilterRequest struct {
	Token     string `form:"token" query:"token"`
	Page      string `form:"page,omitempty" query:"page"`
	Limit     string `form:"limit,omitempty" query:"limit"`
	Order     string `form:"order,omitempty" query:"order"`
	Direction string `form:"direction,omitempty" query:"direction"`
	Email     string `form:"email,omitempty" query:"email"`
	// Contact name
	Q     string `form:"q,omitempty" query:"q"`
	Phone string `form:"phone,omitempty" query:"phone"`
	Title string `form:"title,omitempty" query:"title"`
}

type ListContactsFilterResponse struct {
	Contacts []Contact `json:"contacts"`
	HasMore  bool      `json:"has_more"`
	Total    float64   `json:"total"`
}

func (s *Client) ListContactsFilter(ctx context.Context, filter ListContactsFilterRequest) (*ListContactsFilterResponse, error) {
	queryString, err := StructToQueryString(filter)
	if err != nil {
		return nil, fmt.Errorf("error creating query string from filter: %w", err)
	}

	fullPath := listContactsEndpoint
	if queryString != "" {
		fullPath += "?" + queryString
	}

	resp, err := s.request(ctx, nil, http.MethodGet, fullPath)
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to list contacts: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to list contacts (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to list contacts (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload ListContactsFilterResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding list contacts response: %w", ErrDecodeResponse, err)
	}

	return &responsePayload, nil
}

type BirthdayData struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type EmailData struct {
	Email string `json:"email"`
}

type PhoneData struct {
	Phone string `json:"phone"`
	Type  string `json:"type,omitempty"`
}

type CreateContactData struct {
	Birthday            *BirthdayData         `json:"birthday,omitempty"`
	ContactCustomFields *[]ContactCustomField `json:"contact_custom_fields,omitempty"`
	DealIDs             *[]string             `json:"deal_ids,omitempty"`
	Emails              *[]EmailData          `json:"emails,omitempty"`
	Facebook            *string               `json:"facebook,omitempty"`
	LegalBases          *[]LegalBasis         `json:"legal_bases,omitempty"`
	LinkedIn            *string               `json:"linkedin,omitempty"`
	Name                string                `json:"name"`
	Phones              *[]PhoneData          `json:"phones,omitempty"`
	OrganizationID      *string               `json:"organization_id,omitempty"`
	Skype               *string               `json:"skype,omitempty"`
}

type CreateContactRequest struct {
	Contact CreateContactData `json:"contact"`
}

type CreateContactResponse struct {
	ID                  string               `json:"id"`
	InternalID          string               `json:"_id"`
	Birthday            BirthdayResponse     `json:"birthday"`
	ContactCF           interface{}          `json:"contact_c_f"`
	ContactCustomFields []ContactCustomField `json:"contact_custom_fields"`
	CreatedAt           string               `json:"created_at"`
	DealIDs             []string             `json:"deal_ids"`
	Emails              []Email              `json:"emails"`
	Facebook            string               `json:"facebook"`
	LegalBases          []LegalBasis         `json:"legal_bases"`
	LinkedIn            string               `json:"linkedin"`
	Name                string               `json:"name"`
	Notes               string               `json:"notes"`
	Organization        OrganizationResponse `json:"organization"`
	OrganizationID      string               `json:"organization_id"`
	Phones              []Phone              `json:"phones"`
	Skype               string               `json:"skype"`
	Title               string               `json:"title"`
	UpdatedAt           string               `json:"updated_at"`
}

type BirthdayResponse struct {
	ID        string `json:"_id"`
	CreatedAt string `json:"created_at"`
	Day       int    `json:"day"`
	Month     int    `json:"month"`
	UpdatedAt string `json:"updated_at"`
	Year      int    `json:"year"`
}

type OrganizationResponse struct {
	ID     string `json:"_id"`
	Name   string `json:"name"`
	Resume string `json:"resume"`
	URL    string `json:"url"`
}

func (s *Client) CreateContact(ctx context.Context, contact CreateContactRequest) (*CreateContactResponse, error) {
	resp, err := s.request(ctx, contact, http.MethodPost, createContactEndpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to create contact: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to create contact (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to create contact (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload CreateContactResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding create contact response: %w", ErrDecodeResponse, err)
	}

	return &responsePayload, nil
}

type UpdateContactData struct {
	Birthday            *BirthdayData        `json:"birthday,omitempty"`
	ContactCustomFields []ContactCustomField `json:"contact_custom_fields,omitempty"`
	DealIDs             []string             `json:"deal_ids,omitempty"`
	Emails              []EmailData          `json:"emails,omitempty"`
	Facebook            *string              `json:"facebook,omitempty"`
	LegalBases          []LegalBasis         `json:"legal_bases,omitempty"`
	LinkedIn            *string              `json:"linkedin,omitempty"`
	Name                string               `json:"name,omitempty"`
	OrganizationID      *string              `json:"organization_id,omitempty"`
	Phones              []Phone              `json:"phones,omitempty"`
	Skype               *string              `json:"skype,omitempty"`
	Title               *string              `json:"title,omitempty"`
}

type UpdateContactRequest struct {
	Contact UpdateContactData `json:"contact"`
}

type UpdateContactResponse struct {
	ID                  string               `json:"id"`
	InternalID          string               `json:"_id"`
	Birthday            BirthdayResponse     `json:"birthday"`
	ContactCF           interface{}          `json:"contact_c_f"`
	ContactCustomFields []ContactCustomField `json:"contact_custom_fields"`
	CreatedAt           string               `json:"created_at"`
	DealIDs             []string             `json:"deal_ids"`
	Emails              []Email              `json:"emails"`
	Facebook            string               `json:"facebook"`
	LegalBases          []LegalBasis         `json:"legal_bases"`
	LinkedIn            string               `json:"linkedin"`
	Name                string               `json:"name"`
	Notes               string               `json:"notes"`
	Organization        OrganizationResponse `json:"organization"`
	OrganizationID      string               `json:"organization_id"`
	Phones              []Phone              `json:"phones"`
	Skype               string               `json:"skype"`
	Title               string               `json:"title"`
	UpdatedAt           string               `json:"updated_at"`
}

func (s *Client) UpdateContact(ctx context.Context, contactID string, contact UpdateContactRequest) (*UpdateContactResponse, error) {
	resp, err := s.request(ctx, contact, http.MethodPut, fmt.Sprintf(updateContactByIDEndpoint, contactID))
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to update contact: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to update contact (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to update contact (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload UpdateContactResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding update contact response: %w", ErrDecodeResponse, err)
	}

	return &responsePayload, nil
}
