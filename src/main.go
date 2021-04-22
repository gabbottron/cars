package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Vehicle struct {
	Id    int    `json:"id"`
	Make  string `json:"make" binding:"required"`
	Model string `json:"model" binding:"required"`
	Year  int    `json:"year" binding:"required"`
}

var length = 0 //to start at least 0 length and increase overtime
var storeData = make([]Vehicle, length)
var nextId = 1 // the next ID in the database

func main() {

	router := gin.Default()

	//GET '/'  --> all cars
	router.GET("/cars", func(c *gin.Context) {
		c.JSON(200, storeData)
	})

	//POST '/cars'  --> create cars
	router.POST("/cars", func(c *gin.Context) {
		var car Vehicle

		err := c.ShouldBindJSON(&car) //binds the input data into 'motor' var
		if err != nil {
			c.JSON(422, gin.H{"message": "Unprocessable entity!"})
			return
		}

		car.Id = nextId
		nextId++

		storeData = append(storeData, car)

		c.JSON(200, car)
	})

	//GET '/cars/:carid'  --> get single car
	router.GET("/cars/:carid", func(c *gin.Context) {
		carid, err := strconv.Atoi(c.Param("carid"))
		if err != nil {
			c.JSON(404, gin.H{
				"message": "car id not valid",
			})
			return
		}

		var car Vehicle
		for i := 0; i < len(storeData); i++ {
			if storeData[i].Id == carid {
				car = storeData[i]
				break
			}
		}
		if car.Id == 0 {
			// the car was not found
			c.JSON(404, gin.H{
				"message": "car not found",
			})
			return
		}

		c.JSON(200, car)
	})

	//PUT '/cars/:carid'  --> modify that single car
	router.PUT("/cars/:carid", func(c *gin.Context) {
		carid, err := strconv.Atoi(c.Param("carid"))
		if err != nil {
			log.Println(err.Error())
			c.JSON(404, gin.H{
				"message": "car id not valid",
			})
			return
		}

		var car Vehicle

		err = c.ShouldBindJSON(&car) //binds the input data into 'motor' var
		if err != nil {
			c.JSON(422, gin.H{"message": "Unprocessable entity!"})
			return
		}

		found := false
		for i := 0; i < len(storeData); i++ {
			if storeData[i].Id == carid {
				car.Id = carid
				storeData[i] = car
				found = true
				break
			}
		}
		if !found {
			// the car was not found
			c.JSON(404, gin.H{
				"message": "car not found",
			})
			return
		}
		c.JSON(200, car)
	})

	//DELETE '/cars/:carid'  --> delete that single car
	router.DELETE("/cars/:carid", func(c *gin.Context) {
		carid, err := strconv.Atoi(c.Param("carid"))
		if err != nil {
			c.JSON(404, gin.H{
				"message": "car id not valid",
			})
			return
		}

		found := false
		for i := 0; i < len(storeData); i++ {
			if storeData[i].Id == carid {
				storeData = append(storeData[:i], storeData[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			c.JSON(404, gin.H{
				"message": "car not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "car successfully deleted",
		})
	})

	router.Run()
}
