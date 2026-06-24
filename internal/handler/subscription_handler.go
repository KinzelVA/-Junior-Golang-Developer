package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

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
	router.GET("/subscriptions", h.List)
	router.GET("/subscriptions/:id", h.GetByID)
	router.PUT("/subscriptions/:id", h.Update)
	router.DELETE("/subscriptions/:id", h.Delete)
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

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	subscription, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "subscription not found",
			})
			return
		}

		h.log.Warn("failed to get subscription", slog.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.NewSubscriptionResponse(subscription))
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	limit, err := parseIntQuery(c, "limit", 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limit must be integer",
		})
		return
	}

	offset, err := parseIntQuery(c, "offset", 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "offset must be integer",
		})
		return
	}

	request := model.ListSubscriptionsRequest{
		UserID:      stringPointerFromQuery(c, "user_id"),
		ServiceName: stringPointerFromQuery(c, "service_name"),
		Limit:       limit,
		Offset:      offset,
	}

	subscriptions, err := h.service.List(c.Request.Context(), request)
	if err != nil {
		h.log.Warn("failed to list subscriptions", slog.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":  model.NewSubscriptionResponses(subscriptions),
		"limit":  request.Limit,
		"offset": request.Offset,
		"count":  len(subscriptions),
	})
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var request model.UpdateSubscriptionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Warn("invalid update subscription request", slog.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	subscription, err := h.service.Update(c.Request.Context(), id, request)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "subscription not found",
			})
			return
		}

		h.log.Warn("failed to update subscription", slog.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.log.Info(
		"subscription updated",
		slog.String("subscription_id", subscription.ID),
		slog.String("user_id", subscription.UserID),
		slog.String("service_name", subscription.ServiceName),
	)

	c.JSON(http.StatusOK, model.NewSubscriptionResponse(subscription))
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "subscription not found",
			})
			return
		}

		h.log.Warn("failed to delete subscription", slog.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.log.Info("subscription deleted", slog.String("subscription_id", id))

	c.Status(http.StatusNoContent)
}

func parseIntQuery(c *gin.Context, key string, defaultValue int) (int, error) {
	rawValue := c.Query(key)
	if rawValue == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(rawValue)
}

func stringPointerFromQuery(c *gin.Context, key string) *string {
	value := c.Query(key)
	if value == "" {
		return nil
	}

	return &value
}
