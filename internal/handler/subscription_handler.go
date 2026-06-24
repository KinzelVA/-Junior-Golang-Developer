package handler

import (
"log/slog"
"net/http"

"github.com/gin-gonic/gin"

"github.com/KinzelVA/-Junior-Golang-Developer/internal/model"
"github.com/KinzelVA/-Junior-Golang-Developer/internal/service"
)

type SubscriptionHandler struct {
service *service.SubscriptionService
log     *slog.Logger
}

func NewSubscriptionHandler(service *service.SubscriptionService, log *slog.Logger) *SubscriptionHandler {
return &SubscriptionHandler{
service: service,
log:     log,
}
}

func (h *SubscriptionHandler) RegisterRoutes(router *gin.RouterGroup) {
router.POST("/subscriptions", h.Create)
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
var request model.CreateSubscriptionRequest

if err := c.ShouldBindJSON(&request); err != nil {
h.log.Warn("invalid create subscription request", slog.String("error", err.Error()))

c.JSON(http.StatusBadRequest, gin.H{
"error": "invalid request body",
})
return
}

subscription, err := h.service.Create(c.Request.Context(), request)
if err != nil {
h.log.Warn("failed to create subscription", slog.String("error", err.Error()))

c.JSON(http.StatusBadRequest, gin.H{
"error": err.Error(),
})
return
}

h.log.Info(
"subscription created",
slog.String("subscription_id", subscription.ID),
slog.String("user_id", subscription.UserID),
slog.String("service_name", subscription.ServiceName),
)

c.JSON(http.StatusCreated, model.NewSubscriptionResponse(subscription))
}
