package handler

import (
	"context"
	"net/http"

	jwt "github.com/OsipyanG/market/protos/jwt"
	shopcart "github.com/OsipyanG/market/protos/shopcart"
	"github.com/gin-gonic/gin"
)

type ShopcartHandler struct {
	client shopcart.UserShopcartClient
}

func NewShopcartHandler(sc shopcart.UserShopcartClient) *ShopcartHandler {
	return &ShopcartHandler{client: sc}
}

func (h *ShopcartHandler) AddProduct(c *gin.Context) {
	requestByIDWithProduct := &shopcart.RequestByIDWithProduct{}

	if err := c.ShouldBindJSON(requestByIDWithProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json body unmarhaling"})

		return
	}

	claims := c.MustGet("claims")

	parsedJwtClaims, ok := claims.(*jwt.JWTClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid claims"})

		return
	}

	requestByIDWithProduct.JwtClaims = parsedJwtClaims

	_, err := h.client.AddProduct(context.TODO(), requestByIDWithProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, "product added")
}

func (h *ShopcartHandler) DeleteProduct(c *gin.Context) {
	requestByIDWithProduct := &shopcart.RequestByIDWithProduct{}

	if err := c.ShouldBindJSON(requestByIDWithProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json body unmarhaling"})

		return
	}

	productID, ok := c.Params.Get("product_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})

		return
	}

	requestByIDWithProduct.Product.Id = productID

	claims := c.MustGet("claims")

	parsedJwtClaims, ok := claims.(*jwt.JWTClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid claims"})

		return
	}

	requestByIDWithProduct.JwtClaims = parsedJwtClaims

	_, err := h.client.DeleteProduct(context.TODO(), requestByIDWithProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, "product deleted")
}

func (h *ShopcartHandler) Clear(c *gin.Context) {
	requestByID := &shopcart.RequestByID{}

	claims := c.MustGet("claims")

	parsedJwtClaims, ok := claims.(*jwt.JWTClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid claims"})

		return
	}

	requestByID.JwtClaims = parsedJwtClaims

	_, err := h.client.Clear(context.TODO(), requestByID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, "shopcart cleared")
}

func (h *ShopcartHandler) GetProducts(c *gin.Context) {
	requestByID := &shopcart.RequestByID{}

	claims := c.MustGet("claims")

	parsedJwtClaims, ok := claims.(*jwt.JWTClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid claims"})

		return
	}

	requestByID.JwtClaims = parsedJwtClaims

	products, err := h.client.GetProducts(context.TODO(), requestByID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, products)
}
