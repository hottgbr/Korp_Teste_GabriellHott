package invoice

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
