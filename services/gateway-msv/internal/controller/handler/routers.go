package handler

import "github.com/gin-gonic/gin"

func SetupRouter(router *gin.Engine, shopcart *ShopcartHandler, order *OrderHandler, catalog *CatalogHandler, auth *AuthHandler) {
	api := router.Group("api")
	{
		api.POST("/auth/register", auth.NewUser)
		api.POST("/auth/login", auth.Login)
		api.POST("/auth/refresh", auth.UpdateTokens)
		api.POST("/auth/logout", auth.Logout)

		api.GET("/catalog", catalog.GetCatalog)

		authorized := api.Group("/")
		authorized.Use(AuthMiddleware(auth.authClient))
		{
			authorized.PATCH("/auth/password", auth.UpdatePassword)
			authorized.GET("/auth/users", auth.GetAllUsersWithLevel)
			authorized.PATCH("/auth/users/:user_id", auth.SetAccessLevel)
			authorized.DELETE("/auth/users/:user_id", auth.DeleteUser)

			authorized.PATCH("/shopcart", shopcart.AddProduct)
			authorized.DELETE("/shopcart/:product_id", shopcart.DeleteProduct)
			authorized.GET("/shopcart", shopcart.GetProducts)
			authorized.DELETE("/shopcart", shopcart.Clear)

			authorized.POST("/order", order.CreateOrder)
			authorized.GET("/order/:order_id", order.GetOrder)
			authorized.GET("/order", order.GetOrders)
			authorized.PATCH("/order/:order_id", order.UpdateOrderStatus)

			authorized.GET("/delivery/pending", order.GetAllDeliveries)
			authorized.GET("/delivery", order.GetAllPendingOrders)
		}
	}
}
