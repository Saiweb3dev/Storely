package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    // Create a Gin router
    router := gin.Default()

    // Define a route for GET requests to "/greet"
    router.GET("/greet", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, from Gin!",
        })
    })

    // Start the server on port 3000
    router.Run(":3000")
}
