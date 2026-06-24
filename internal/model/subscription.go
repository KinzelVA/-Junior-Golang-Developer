package model

import (
	"fmt"
	"time"
)

type Subscription struct {
	ID          string
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
	EndDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int     `json:"price" binding:"required,min=1"`
	UserID      string  `json:"user_id" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int     `json:"price" binding:"required,min=1"`
	UserID      string  `json:"user_id" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

type ListSubscriptionsRequest struct {
	UserID      *string
	ServiceName *string
	Limit       int
	Offset      int
}

type SubscriptionResponse struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func ParseMonthYear(value string) (time.Time, error) {
	parsed, err := time.Parse("01-2006", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("date must be in MM-YYYY format")
	}

	return time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

func FormatMonthYear(value time.Time) string {
	return value.Format("01-2006")
}

func NewSubscriptionResponse(subscription *Subscription) SubscriptionResponse {
	var endDate *string

	if subscription.EndDate != nil {
		formatted := FormatMonthYear(*subscription.EndDate)
		endDate = &formatted
	}

	return SubscriptionResponse{
		ID:          subscription.ID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID,
		StartDate:   FormatMonthYear(subscription.StartDate),
		EndDate:     endDate,
		CreatedAt:   subscription.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   subscription.UpdatedAt.Format(time.RFC3339),
	}
}

func NewSubscriptionResponses(subscriptions []Subscription) []SubscriptionResponse {
	responses := make([]SubscriptionResponse, 0, len(subscriptions))

	for i := range subscriptions {
		responses = append(responses, NewSubscriptionResponse(&subscriptions[i]))
	}

	return responses
}
