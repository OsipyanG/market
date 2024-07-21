package handler

import "github.com/gin-gonic/gin"

func SetupRouter(router *gin.Engine, shopcart *ShopcartHandler,
	order *OrderHandler, catalog *CatalogHandler, auth *AuthHandler,
) {
	v1 := router.Group("")
	{
		v1.POST("/auth/register", auth.NewUser)
		v1.POST("/auth/login", auth.Login)
		v1.POST("/auth/refresh", auth.UpdateTokens)
		v1.POST("/auth/logout", auth.Logout)

		v1.GET("/catalog", catalog.GetCatalog)

		authorized := v1.Group("/")
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
