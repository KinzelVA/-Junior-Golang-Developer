package service

import (
"context"
"fmt"
"time"

"github.com/google/uuid"

"github.com/KinzelVA/-Junior-Golang-Developer/internal/model"
"github.com/KinzelVA/-Junior-Golang-Developer/internal/repository"
)

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

var endDatePointer *time.Time

if request.EndDate != nil && *request.EndDate != "" {
endDate, err := model.ParseMonthYear(*request.EndDate)
if err != nil {
return nil, fmt.Errorf("invalid end_date: %w", err)
}

if endDate.Before(startDate) {
return nil, fmt.Errorf("end_date must be greater than or equal to start_date")
}

endDatePointer = &endDate
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
