package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/KinzelVA/-Junior-Golang-Developer/internal/model"
	"github.com/KinzelVA/-Junior-Golang-Developer/internal/repository"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type SubscriptionService struct {
	repository *repository.SubscriptionRepository
}

func NewSubscriptionService(repository *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repository: repository,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, request model.CreateSubscriptionRequest) (*model.Subscription, error) {
	if _, err := uuid.Parse(request.UserID); err != nil {
		return nil, fmt.Errorf("user_id must be valid UUID")
	}

	startDate, err := model.ParseMonthYear(request.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	endDatePointer, err := parseOptionalEndDate(request.EndDate, startDate)
	if err != nil {
		return nil, err
	}

	subscription := &model.Subscription{
		ServiceName: request.ServiceName,
		Price:       request.Price,
		UserID:      request.UserID,
		StartDate:   startDate,
		EndDate:     endDatePointer,
	}

	if err := s.repository.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id string) (*model.Subscription, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("id must be valid UUID")
	}

	subscription, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			return nil, ErrSubscriptionNotFound
		}

		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) List(ctx context.Context, request model.ListSubscriptionsRequest) ([]model.Subscription, error) {
	if request.UserID != nil && *request.UserID != "" {
		if _, err := uuid.Parse(*request.UserID); err != nil {
			return nil, fmt.Errorf("user_id must be valid UUID")
		}
	}

	if request.Limit <= 0 {
		request.Limit = 20
	}

	if request.Limit > 100 {
		request.Limit = 100
	}

	if request.Offset < 0 {
		request.Offset = 0
	}

	return s.repository.List(ctx, request)
}

func (s *SubscriptionService) Update(ctx context.Context, id string, request model.UpdateSubscriptionRequest) (*model.Subscription, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("id must be valid UUID")
	}

	if _, err := uuid.Parse(request.UserID); err != nil {
		return nil, fmt.Errorf("user_id must be valid UUID")
	}

	startDate, err := model.ParseMonthYear(request.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	endDatePointer, err := parseOptionalEndDate(request.EndDate, startDate)
	if err != nil {
		return nil, err
	}

	subscription := &model.Subscription{
		ID:          id,
		ServiceName: request.ServiceName,
		Price:       request.Price,
		UserID:      request.UserID,
		StartDate:   startDate,
		EndDate:     endDatePointer,
	}

	if err := s.repository.Update(ctx, subscription); err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			return nil, ErrSubscriptionNotFound
		}

		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("id must be valid UUID")
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			return ErrSubscriptionNotFound
		}

		return err
	}

	return nil
}

func (s *SubscriptionService) TotalCost(ctx context.Context, request model.TotalCostRequest) (int, error) {
	if request.PeriodStart == "" {
		return 0, fmt.Errorf("period_start is required")
	}

	if request.PeriodEnd == "" {
		return 0, fmt.Errorf("period_end is required")
	}

	periodStart, err := model.ParseMonthYear(request.PeriodStart)
	if err != nil {
		return 0, fmt.Errorf("invalid period_start: %w", err)
	}

	periodEnd, err := model.ParseMonthYear(request.PeriodEnd)
	if err != nil {
		return 0, fmt.Errorf("invalid period_end: %w", err)
	}

	if periodEnd.Before(periodStart) {
		return 0, fmt.Errorf("period_end must be greater than or equal to period_start")
	}

	if request.UserID != nil && *request.UserID != "" {
		if _, err := uuid.Parse(*request.UserID); err != nil {
			return 0, fmt.Errorf("user_id must be valid UUID")
		}
	}

	filter := model.TotalCostFilter{
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		UserID:      request.UserID,
		ServiceName: request.ServiceName,
	}

	return s.repository.TotalCost(ctx, filter)
}

func parseOptionalEndDate(rawEndDate *string, startDate time.Time) (*time.Time, error) {
	if rawEndDate == nil || *rawEndDate == "" {
		return nil, nil
	}

	endDate, err := model.ParseMonthYear(*rawEndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date: %w", err)
	}

	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end_date must be greater than or equal to start_date")
	}

	return &endDate, nil
}
