package handler

import (
	"context"
	"net/http"
	"strconv"

	catalogpb "github.com/OsipyanG/market/protos/warehouse"
	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	client catalogpb.CatalogClient
}

func NewCatalogHandler(cs catalogpb.CatalogClient) *CatalogHandler {
	return &CatalogHandler{client: cs}
}

func (h *CatalogHandler) GetCatalog(c *gin.Context) {
	strOffset := c.Query("offset")
	strLimit := c.Query("limit")

	if strLimit == "" || strOffset == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})

		return
	}

	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})

		return
	}

	offset, err := strconv.Atoi(strLimit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})

		return
	}

	getCatalogRequest := &catalogpb.GetCatalogRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	catalog, err := h.client.GetCatalog(context.TODO(), getCatalogRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, catalog)
}
