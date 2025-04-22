package rd_station

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Deal struct {
	ID                   string        `json:"id"`
	AmountMonthly        float64       `json:"amount_montly"`
	AmountTotal          float64       `json:"amount_total"`
	AmountUnique         float64       `json:"amount_unique"`
	ClosedAt             string        `json:"closed_at"`
	Deals                []Deal        `json:"deals"`
	CreatedAt            string        `json:"created_at"`
	DealCustomFields     []interface{} `json:"deal_custom_fields"`
	DealProducts         []DealProduct `json:"deal_products"`
	DealStage            DealStage     `json:"deal_stage"`
	Hold                 string        `json:"hold"`
	Interactions         int           `json:"interactions"`
	LastActivityAt       string        `json:"last_activity_at"`
	LastActivityContent  string        `json:"last_activity_content"`
	Markup               string        `json:"markup"`
	MarkupCreated        string        `json:"markup_created"`
	MarkupLastActivities string        `json:"markup_last_activities"`
	Name                 string        `json:"name"`
	PredictionDate       string        `json:"prediction_date"`
	Rating               int           `json:"rating"`
	StopTimeLimit        interface{}   `json:"stop_time_limit"`
	UpdatedAt            string        `json:"updated_at"`
	User                 User          `json:"user"`
	UserChanged          bool          `json:"user_changed"`
	Win                  string        `json:"win"`
}

type DealProduct struct {
	ID           string  `json:"id"`
	Amount       int     `json:"amount"`
	BasePrice    float64 `json:"base_price"`
	CreatedAt    string  `json:"created_at"`
	Description  string  `json:"description"`
	Discount     float64 `json:"discount"`
	DiscountType string  `json:"discount_type"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	ProductID    string  `json:"product_id"`
	Recurrence   string  `json:"recurrence"`
	Total        float64 `json:"total"`
	UpdatedAt    string  `json:"updated_at"`
}

type DealStage struct {
	InternalID string `json:"_id"`
	CreatedAt  string `json:"created_at"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Nickname   string `json:"nickname"`
	UpdatedAt  string `json:"updated_at"`
}

type User struct {
	InternalID string `json:"_id"`
	Email      string `json:"email"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Nickname   string `json:"nickname"`
}

type ListDealsFilterRequest struct {
	Token     string `form:"token" query:"token"`
	Page      string `form:"page,omitempty" query:"page"`
	Limit     string `form:"limit,omitempty" query:"limit"`         // Default value: 20. Maximum value: 200
	Order     string `form:"order,omitempty" query:"order"`         // Default value: "created_at"
	Direction string `form:"direction,omitempty" query:"direction"` // "asc" or "desc", default is "desc"

	// Name is the deal name for searching
	Name string `form:"name,omitempty" query:"name"`

	// ExactName when "true", searches for the exact deal name defined in the Name field
	ExactName string `form:"exact_name,omitempty" query:"exact_name"`

	// Win filters deals by status: "true" (won), "false" (lost) or "null" (open)
	Win string `form:"win,omitempty" query:"win"`

	UserID string `form:"user_id,omitempty" query:"user_id"`

	// ClosedAt when "true" returns "won" or "lost" deals, when "false" returns "open" or "paused" deals
	ClosedAt string `form:"closed_at,omitempty" query:"closed_at"`

	// ClosedAtPeriod when "true", must be used together with StartDate and EndDate to filter by closing period
	ClosedAtPeriod string `form:"closed_at_period,omitempty" query:"closed_at_period"`

	// CreatedAtPeriod when "true", must be used together with StartDate and EndDate to filter by creation period
	CreatedAtPeriod string `form:"created_at_period,omitempty" query:"created_at_period"`

	// PredictionDatePeriod when "true", must be used together with StartDate and EndDate to filter by prediction period
	PredictionDatePeriod string `form:"prediction_date_period,omitempty" query:"prediction_date_period"`

	// StartDate defines the beginning of the period for date filters in ISO 8601 format, e.g.: "2020-12-14T15:00:00"
	StartDate string `form:"start_date,omitempty" query:"start_date"`

	// EndDate defines the end of the period for date filters in ISO 8601 format, e.g.: "2020-12-14T15:00:00"
	EndDate string `form:"end_date,omitempty" query:"end_date"`

	CampaignID       string `form:"campaign_id,omitempty" query:"campaign_id"`
	DealStageID      string `form:"deal_stage_id,omitempty" query:"deal_stage_id"`
	DealLostReasonID string `form:"deal_lost_reason_id,omitempty" query:"deal_lost_reason_id"`
	DealPipelineID   string `form:"deal_pipeline_id,omitempty" query:"deal_pipeline_id"`
	Organization     string `form:"organization,omitempty" query:"organization"`

	// Hold when "true" returns only "paused" deals
	Hold string `form:"hold,omitempty" query:"hold"`

	// ProductPresence filters by related products:
	// "false" (no related products)
	// "true" (one or more related products)
	// Or a list of product IDs separated by commas, e.g.: "5esdsds,d767dsdssd"
	ProductPresence string `form:"product_presence,omitempty" query:"product_presence"`

	// NextPage is a token for the next page of results, obtained from a previous API call
	// Use this value to navigate to the next page of results for the same search
	NextPage string `form:"next_page,omitempty" query:"next_page"`
}

type ListDealsFilterResponse struct {
	Deals    []Deal `json:"deals"`
	HasMore  bool   `json:"has_more"`
	NextPage string `json:"next_page"`
	Total    int    `json:"total"`
}

func (s *Client) ListDealsFilter(ctx context.Context, filter ListDealsFilterRequest) (*ListDealsFilterResponse, error) {
	queryString, err := StructToQueryString(filter)
	if err != nil {
		return nil, fmt.Errorf("error creating query string from filter: %w", err)
	}

	fullPath := listDealsEndpoint
	if queryString != "" {
		fullPath += "?" + queryString
	}

	resp, err := s.request(ctx, nil, http.MethodGet, fullPath)
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to list deals: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to list deals (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to list deals (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload ListDealsFilterResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding list deals response: %w", ErrDecodeResponse, err)
	}

	return &responsePayload, nil
}

type DealProductData struct {
	Amount       *int     `json:"amount,omitempty"`
	BasePrice    *float64 `json:"base_price,omitempty"`
	Description  *string  `json:"description,omitempty"`
	DiscountType *string  `json:"discount_type,omitempty"`
	Name         *string  `json:"name,omitempty"`
	Price        *float64 `json:"price,omitempty"`
	Recurrence   *string  `json:"recurrence,omitempty"`
	Total        *float64 `json:"total,omitempty"`
}

type DealSourceData struct {
	ID                   *string                   `json:"_id,omitempty"`
	DistributionSettings *DistributionSettingsData `json:"distribution_settings,omitempty"`
}

type DistributionSettingsData struct {
	Owner        *OwnerData        `json:"owner,omitempty"`
	Organization *OrganizationData `json:"organization,omitempty"`
}

type OwnerData struct {
	Email *string `json:"email,omitempty"`
	ID    *string `json:"id,omitempty"`
	Type  *string `json:"type,omitempty"`
}

type OrganizationData struct {
	ID *string `json:"_id,omitempty"`
}

type CampaignData struct {
	ID    *string `json:"_id,omitempty"`
	Deals []Deal  `json:"deals,omitempty"`
}

type CreateDealRequest struct {
	Deal     CreateDealData `json:"deal"`
	Campaign *CampaignData  `json:"campaign,omitempty"`
}

type CreateDealData struct {
	Name             string             `json:"name"`
	Contacts         *[]Contact         `json:"contacts,omitempty"`
	DealCustomFields *[]interface{}     `json:"deal_custom_fields,omitempty"`
	DealStageID      *string            `json:"deal_stage_id,omitempty"`
	PredictionDate   *string            `json:"prediction_date,omitempty"`
	Rating           *int               `json:"rating,omitempty"`
	UserID           *string            `json:"user_id,omitempty"`
	DealProducts     *[]DealProductData `json:"deal_products,omitempty"`
	DealSource       *DealSourceData    `json:"deal_source,omitempty"`
}

type CreateDealResponse struct {
	ID                  string                     `json:"id"`
	AmountMontly        float64                    `json:"amount_montly"`
	AmountTotal         float64                    `json:"amount_total"`
	AmountUnique        float64                    `json:"amount_unique"`
	BestMomentToTouch   *bool                      `json:"best_moment_to_touch,omitempty"`
	CCfErrors           map[string]interface{}     `json:"c_cf_errors,omitempty"`
	Campaign            *CampaignResponse          `json:"campaign,omitempty"`
	CampaignID          *string                    `json:"campaign_id,omitempty"`
	ClosedAt            *string                    `json:"closed_at,omitempty"`
	DealErrors          map[string]interface{}     `json:"deal_errors,omitempty"`
	CreatedAt           string                     `json:"created_at"`
	DealCustomFields    []DealCustomFieldResponse  `json:"deal_custom_fields"`
	DealLostReasonID    *string                    `json:"deal_lost_reason_id,omitempty"`
	DealProducts        []DealProductResponse      `json:"deal_products"`
	DealSource          *DealSourceResponse        `json:"deal_source,omitempty"`
	DealStage           *DealStageResponse         `json:"deal_stage,omitempty"`
	DealStageHistories  []DealStageHistoryResponse `json:"deal_stage_histories"`
	Errors              map[string]interface{}     `json:"errors,omitempty"`
	FromRdsmIntegration *bool                      `json:"from_rdsm_integration,omitempty"`
	Hold                *string                    `json:"hold,omitempty"`
	Interactions        int                        `json:"interactions"`
	LastNoteContent     *string                    `json:"last_note_content,omitempty"`
	Name                string                     `json:"name"`
	Organization        *OrganizationResponse      `json:"organization,omitempty"`
	PredictionDate      *string                    `json:"prediction_date,omitempty"`
	Rating              *float64                   `json:"rating,omitempty"`
	Resume              *string                    `json:"resume,omitempty"`
	StopTimeLimit       *StopTimeLimitResponse     `json:"stop_time_limit,omitempty"`
	UpdatedAt           string                     `json:"updated_at"`
	URL                 *string                    `json:"url,omitempty"`
	User                *UserResponse              `json:"user,omitempty"`
	Visible             *bool                      `json:"visible,omitempty"`
	Win                 *string                    `json:"win,omitempty"`
}

func (s *Client) CreateDeal(ctx context.Context, deal CreateDealRequest) (*CreateDealResponse, error) {
	resp, err := s.request(ctx, deal, http.MethodPost, createDealEndpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to create deal: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to create deal (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to create deal (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload CreateDealResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding create deal response: %w", ErrDecodeResponse, err)
	}

	return &responsePayload, nil
}

type UpdateDealResponse struct {
	ID                  string                     `json:"id"`
	AmountMontly        float64                    `json:"amount_montly"`
	AmountTotal         float64                    `json:"amount_total"`
	AmountUnique        float64                    `json:"amount_unique"`
	BestMomentToTouch   *bool                      `json:"best_moment_to_touch,omitempty"`
	CCfErrors           map[string]interface{}     `json:"c_cf_errors,omitempty"`
	Campaign            *CampaignResponse          `json:"campaign,omitempty"`
	CampaignID          *string                    `json:"campaign_id,omitempty"`
	ClosedAt            *string                    `json:"closed_at,omitempty"`
	DealErrors          map[string]interface{}     `json:"deal_errors,omitempty"`
	CreatedAt           string                     `json:"created_at"`
	DealCustomFields    []DealCustomFieldResponse  `json:"deal_custom_fields"`
	DealLostNote        *string                    `json:"deal_lost_note,omitempty"`
	DealLostReasonID    *string                    `json:"deal_lost_reason_id,omitempty"`
	DealProducts        []DealProductResponse      `json:"deal_products"`
	DealSource          *DealSourceResponse        `json:"deal_source,omitempty"`
	DealStage           *DealStageResponse         `json:"deal_stage,omitempty"`
	DealStageHistories  []DealStageHistoryResponse `json:"deal_stage_histories"`
	Errors              map[string]interface{}     `json:"errors,omitempty"`
	FromRdsmIntegration *bool                      `json:"from_rdsm_integration,omitempty"`
	Hold                *string                    `json:"hold,omitempty"`
	Interactions        int                        `json:"interactions"`
	LastNoteContent     *string                    `json:"last_note_content,omitempty"`
	Name                string                     `json:"name"`
	Organization        *OrganizationResponse      `json:"organization,omitempty"`
	PredictionDate      *string                    `json:"prediction_date,omitempty"`
	Rating              *float64                   `json:"rating,omitempty"`
	Resume              *string                    `json:"resume,omitempty"`
	StopTimeLimit       *StopTimeLimitResponse     `json:"stop_time_limit,omitempty"`
	UpdatedAt           string                     `json:"updated_at"`
	URL                 *string                    `json:"url,omitempty"`
	User                *UserResponse              `json:"user,omitempty"`
	Visible             *bool                      `json:"visible,omitempty"`
	Win                 *string                    `json:"win,omitempty"`
}

type CampaignResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DealCustomFieldResponse struct {
	CreatedAt   string              `json:"created_at"`
	CustomField CustomFieldResponse `json:"custom_field"`
	UpdatedAt   string              `json:"updated_at"`
	Value       interface{}         `json:"value"`
}

type CustomFieldResponse struct {
	CustomFieldID string `json:"custom_field_id"`
}

type DealProductResponse struct {
	ID           string  `json:"id"`
	Amount       int     `json:"amount"`
	BasePrice    float64 `json:"base_price"`
	CreatedAt    string  `json:"created_at"`
	Description  string  `json:"description"`
	Discount     float64 `json:"discount"`
	DiscountType string  `json:"discount_type"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	ProductID    string  `json:"product_id"`
	Recurrence   string  `json:"recurrence"`
	Total        float64 `json:"total"`
	UpdatedAt    string  `json:"updated_at"`
}

type DealSourceResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DealSourceID string `json:"deal_source_id"`
}

type DealStageResponse struct {
	DealPipelineID string `json:"deal_pipeline_id"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	Nickname       string `json:"nickname"`
}

type DealStageHistoryResponse struct {
	DealStageID string  `json:"deal_stage_id"`
	EndDate     *string `json:"end_date"`
	ID          string  `json:"id"`
	StartDate   string  `json:"start_date"`
}

type OrganizationCustomFieldResponse struct {
	ID            string      `json:"id"`
	CreatedAt     string      `json:"created_at"`
	CustomFieldID string      `json:"custom_field_id"`
	UpdatedAt     string      `json:"updated_at"`
	Value         interface{} `json:"value"`
}

type StopTimeLimitResponse struct {
	ExpirationDateTime *string `json:"expiration_date_time,omitempty"`
	Expired            *bool   `json:"expired,omitempty"`
	ExpiredDays        *int    `json:"expired_days,omitempty"`
}

type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type UpdateDealRequest struct {
	Campaign *UpdateCampaignRequestData `json:"campaign,omitempty"`
	Deal     UpdateDealRequestData      `json:"deal"`
}

type UpdateCampaignRequestData struct {
	ID *string `json:"_id,omitempty"`
}

type UpdateDealRequestData struct {
	DealCustomFields *[]UpdateDealCustomFieldRequestData `json:"deal_custom_fields,omitempty"`
	DealLostNote     *string                             `json:"deal_lost_note,omitempty"`
	DealLostReasonID *string                             `json:"deal_lost_reason_id,omitempty"`
	Hold             *string                             `json:"hold,omitempty"`
	Name             *string                             `json:"name,omitempty"`
	OrganizationID   *string                             `json:"organization_id,omitempty"`
	PredictionDate   *string                             `json:"prediction_date,omitempty"`
	Rating           *float64                            `json:"rating,omitempty"`
	UserID           *string                             `json:"user_id,omitempty"`
	Win              *string                             `json:"win,omitempty"`
	DealSource       *UpdateDealSourceRequestData        `json:"deal_source,omitempty"`
	DealStageID      *string                             `json:"deal_stage_id,omitempty"`
}

type UpdateDealCustomFieldRequestData struct {
	CustomFieldID string      `json:"custom_field_id"`
	Value         interface{} `json:"value"`
}

type UpdateDealSourceRequestData struct {
	ID          *string `json:"_id,omitempty"`
	DealStageID *string `json:"deal_stage_id,omitempty"`
}

func (s *Client) UpdateDeal(ctx context.Context, dealID string, deal UpdateDealRequest) (*UpdateDealResponse, error) {
	resp, err := s.request(ctx, deal, http.MethodPut, fmt.Sprintf(updateDealByIDEndpoint, dealID))
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to update deal: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to update deal (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to update deal (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload UpdateDealResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding update deal response: %w", ErrDecodeResponse, err)
	}

	return &responsePayload, nil
}
