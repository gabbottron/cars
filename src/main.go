package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type Vehicle struct {
	Id    string `json:"id" binding:"required"`
	Name  string `json:"title" binding:"required"`
	Model string `json:"type" binding:"required"`
	Year  int    `json:"madeIn"`
}

var length = 0 //to start at least 0 length and increase overtime
var storeData = make([]Vehicle, length)

func main() {

	router := gin.Default()

	//GET '/'  --> all cars
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "hit the get route",
			"vehicles": storeData,
		})
	})

	//POST '/cars'  --> create cars
	router.POST("/cars", func(c *gin.Context) {
		var motor Vehicle

		err := c.ShouldBindJSON(&motor) //binds the input data into 'motor' var
		if err != nil {
			c.JSON(400, gin.H{"message": "Failed."})
			return
		}
		storeData = append(storeData, motor)

		c.JSON(200, gin.H{
			"message": "hit the POST cars route & success post",
			"id":      motor.Id,    //write in 'id' --> `json:"id"`
			"title":   motor.Name,  //write in 'title' --> `json:"title" binding:"required"`
			"type":    motor.Model, //write in 'type' --> `json:"type" binding:"required"`
			"madeIn":  motor.Year,  //write in 'madeIn' --> `json:"madeIn"`
		})
	})

	//GET '/cars/:carid'  --> get single car
	router.GET("/cars/:carid", func(c *gin.Context) {
		carid := c.Param("carid")
		var car Vehicle

		for i := 0; i < len(storeData); i++ {
			//at the end and id not found
			if i == len(storeData)-1 && storeData[i].Id != carid {
				c.JSON(404, gin.H{
					"message": "car id not found",
				})
				return
			}
			if storeData[i].Id == carid {
				car = storeData[i]
			}
		}
		c.JSON(200, gin.H{
			"message": "hit the  GET single car route",
			"car":     car,
		})
	})

	//PUT '/cars/:carid'  --> modify that single car
	router.PUT("/cars/:carid", func(c *gin.Context) {

		body, _ := ioutil.ReadAll(c.Request.Body)
		var car Vehicle
		err := json.Unmarshal(body, &car) //since body is byte[] --> unmarshalling to change to json byte data to struct

		for i := 0; i < len(storeData); i++ {
			if storeData[i].Id != car.Id && err == nil {
				c.JSON(404, gin.H{
					"message": "car id not found",
				})
				return
			} else {
				storeData[i] = car
				c.JSON(200, gin.H{
					"car": storeData[i],
				})
			}
		}

	})

	//DELETE '/cars/:carid'  --> delete that single car
	router.DELETE("/cars/:carid", func(c *gin.Context) {
		carid := c.Param("carid")

		for i := 0; i < len(storeData); i++ {

			//at the end and id not found
			if i == len(storeData)-1 && storeData[i].Id != carid {
				c.JSON(404, gin.H{
					"message": "car id not found",
				})
				return
			}
			if storeData[i].Id == carid {
				storeData = append(storeData[:i], storeData[i+1:]...)
			}
		}
		c.JSON(200, gin.H{
			"message":  "hit the get route",
			"vehicles": storeData,
		})
	})

	router.Run()
}
