package http_handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sewaustav/Payment/internal/payment/dto"
	"github.com/sewaustav/Payment/internal/payment/models"
	service "github.com/sewaustav/Payment/internal/payment/service/api"
)

type PaymentHttpHandler struct {
	service service.PaymentApiService
}

func NewHttpHandler(service service.PaymentApiService) *PaymentHttpHandler {
	return &PaymentHttpHandler{
		service: service,
	}
}

func (h *PaymentHttpHandler) GetMySubscriptionInfoHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no sach user"})
		return
	}

	user := models.UserIdentity{
		UserID: userID,
		Role: &role,
	}

	sub, err := h.service.GetSubscriptionInfoService(ctx, user) 
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *PaymentHttpHandler) GetMyPaymentsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no sach user"})
		return
	}

	user := models.UserIdentity{
		UserID: userID,
		Role: &role,
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	history, err := h.service.GetMyPaymentsService(ctx, user, limit, page)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка по причине разраб долбаеб"})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *PaymentHttpHandler) UpdateSubscriptionHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var body dto.PatchSubcriptionInfoDto
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	
	req := dto.UpadateSubcriptionInfoDto{
		Subscription: body.Subscription,
		IsAutoRenew: body.IsAutoRenew,
		IsRenew: false,
	}

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no sach user"})
		return
	}

	user := models.UserIdentity{
		UserID: userID,
		Role: &role,
	}

	if err := h.service.UpdateSubscriptionInfoService(ctx, user, req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	
	c.Status(http.StatusOK)
}

// admin only 
func (h *PaymentHttpHandler) DeleteUserHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no sach user"})
		return
	}

	user := models.UserIdentity{
		UserID: userID,
		Role: &role,
	}

	usrID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid case ID"})
		return
	}

	if err = h.service.DeleteUserService(ctx, user, usrID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusNoContent)
	
}



func (h *PaymentHttpHandler) GetUserProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no such user"})
		return
	}

	user := models.UserIdentity{
		UserID: userID,
		Role:   &role,
	}

	targetUserID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	profile, err := h.service.GetUserProfileService(ctx, user, targetUserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *PaymentHttpHandler) GetUsersPaymentsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no such user"})
		return
	}

	user := models.UserIdentity{
		UserID: userID,
		Role:   &role,
	}

	targetUserID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	payments, err := h.service.GetUsersPaymentsService(ctx, user, targetUserID, limit, page)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHttpHandler) GetPaymentByTransactionIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no such user"})
		return
	}

	user := models.UserIdentity{
		UserID: userID,
		Role:   &role,
	}

	transactionID := c.Param("transaction_id")
	if transactionID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "transaction_id is required"})
		return
	}

	payment, err := h.service.GetPaymentByTransactionIDService(ctx, user, transactionID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHttpHandler) GetPaymentByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no such user"})
		return
	}

	user := models.UserIdentity{
		UserID: userID,
		Role:   &role,
	}

	paymentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		return
	}

	payment, err := h.service.GetPaymentByIDService(ctx, user, paymentID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}