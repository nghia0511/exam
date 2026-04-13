package handlers

import (
	"net/http"
	"shortener/service"

	"github.com/gin-gonic/gin"
)

type ShortLinkHandler struct {
	svc *service.ShortLinkService
}

func NewShortLinkHandler(svc *service.ShortLinkService) *ShortLinkHandler {
	return &ShortLinkHandler{svc: svc}
}

type createRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
}

func (h *ShortLinkHandler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "original_url is required"})
		return
	}

	link, err := h.svc.Create(req.OriginalURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        link.ID,
		"short_url": link.ShortURL,
	})
}

func (h *ShortLinkHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	link, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shortlink not found"})
		return
	}
	c.JSON(http.StatusOK, link)
}

func (h *ShortLinkHandler) Redirect(c *gin.Context) {
	id := c.Param("id")
	link, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shortlink not found"})
		return
	}
	c.Redirect(http.StatusFound, link.OriginalURL)
}
