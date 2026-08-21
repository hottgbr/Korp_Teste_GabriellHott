package invoice

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var input CreateInvoiceInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	createdInvoice, err := h.service.Create(
		c.Request.Context(),
		input,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrInvoiceItemsRequired),
			errors.Is(err, ErrInvalidProductID),
			errors.Is(err, ErrInvalidQuantity):

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "failed to create invoice",
				},
			)
		}

		return
	}

	c.JSON(
		http.StatusCreated,
		createdInvoice,
	)
}
func (h *Handler) List(c *gin.Context) {
	invoices, err := h.service.List(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to list invoices",
			},
		)
		return
	}

	c.JSON(http.StatusOK, invoices)
}

func (h *Handler) FindByID(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid invoice id",
			},
		)
		return
	}

	invoice, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to retrieve invoice",
			},
		)
		return
	}

	c.JSON(http.StatusOK, invoice)
}

func (h *Handler) Close(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid invoice id",
			},
		)
		return
	}

	closedInvoice, err := h.service.Close(
		c.Request.Context(),
		id,
	)

	if err != nil {
		switch {
		case errors.Is(
			err,
			ErrInvoiceAlreadyClosed,
		):
			c.JSON(
				http.StatusConflict,
				gin.H{
					"error": err.Error(),
				},
			)

		case errors.Is(
			err,
			ErrStockUpdateFailed,
		):
			c.JSON(
				http.StatusServiceUnavailable,
				gin.H{
					"error": "stock service unavailable or stock update failed",
				},
			)

		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "failed to close invoice",
				},
			)
		}

		return
	}

	c.JSON(
		http.StatusOK,
		closedInvoice,
	)
}
