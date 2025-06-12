package middlewares

// // func CORSMiddleware() gin.HandlerFunc {
// // 	return func(c *gin.Context) {
// // 		c.Header("Access-Control-Allow-Origin", "*") // Allow all origins (*), or set specific domain
// // 		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// // 		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

// // 		// Handle preflight requests for OPTIONS method
// // 		if c.Request.Method == "OPTIONS" {
// // 			c.AbortWithStatus(204)
// // 			return
// // 		}

// // 		c.Next()
// // 	}
// }
