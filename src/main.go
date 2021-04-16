package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	log.Println(router)

	//GET '/'  --> all cars
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hit the get route",
		})
	})

	//POST '/cars'  --> create cars
	router.POST("/cars", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hit the POST cars route",
		})
	})

	//GET '/cars/:carid'  --> get single car
	router.GET("/cars/:carid", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hit the  GET single car route",
		})
	})

	//PUT '/cars/:carid'  --> modify that single car
	router.PUT("/cars/:carid", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hit the PUT single car route",
		})
	})

	//DELETE '/cars/:carid'  --> delete that single car
	router.DELETE("/cars/:carid", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hit the DELETE single car route",
		})
	})

	router.Run()
}
