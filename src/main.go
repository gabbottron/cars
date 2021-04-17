package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Vehicle struct {
	Name  string `json:"title" binding:"required"`
	Model string `json:"type" binding:"required"`
	Year  string `json:"madeIn"`
}

func main() {
	router := gin.Default()

	log.Println(router)

	//GET '/'  --> all cars
	router.GET("/", func(c *gin.Context) {
		var motor Vehicle
		log.Println(motor)
		c.JSON(200, gin.H{
			"message":  "hit the get route",
			"vehicles": motor,
		})
	})

	//POST '/cars'  --> create cars
	router.POST("/cars", func(c *gin.Context) {
		var motor Vehicle

		err := c.ShouldBindJSON(&motor) //binds the input data into 'motor' var
		if err != nil {
			c.JSON(200, gin.H{"message": "Failed."})
			return
		}
		c.JSON(200, gin.H{
			"message": "hit the POST cars route",
			"title":   motor.Name,  //write in 'title' --> `json:"title" binding:"required"`
			"type":    motor.Model, //write in 'type' --> `json:"type" binding:"required"`
			"madeIn":  motor.Year,  //write in 'madeIn' --> `json:"madeIn"`
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
