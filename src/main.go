package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type Vehicle struct {
	Id    int    `json:"id"`
	Make  string `json:"make" binding:"required"`
	Model string `json:"model" binding:"required"`
	Year  int    `json:"year" binding:"required"`
}

const (
	SelectAllCarsQuery string = `SELECT id, make, model, year FROM cars ORDER BY id ASC`

	SelectCarQuery string = `SELECT id, make, model, year FROM cars WHERE id = $1`

	DeleteCarQuery string = `DELETE FROM cars WHERE id = $1`

	InsertCarQuery string = `INSERT INTO cars 
		(make, model, year) 
		VALUES($1, $2, $3)
    RETURNING id, make, model, year`

	UpdateCarQuery string = `UPDATE cars
		SET
			make  = COALESCE($2, make),
			model = COALESCE($3, model),
      		year  = COALESCE($4, year)
		WHERE
			id = $1 
		RETURNING
			id, make, model, year`
)

var length = 0 //to start at least 0 length and increase overtime
var storeData = make([]Vehicle, length)
var nextId = 1 // the next ID in the database

func main() {

	router := gin.Default()

	// set up database connection
	dbConnStr := fmt.Sprintf("user=postgres password=mysecretpassword dbname=postgres host=localhost port=5439 sslmode=disable")
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Panic(err)
	}

	// test the connection
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	//GET '/'  --> all cars
	router.GET("/cars", func(c *gin.Context) {
		results := make([]Vehicle, 0)

		var rows *sql.Rows
		var err error

		rows, err = db.Query(SelectAllCarsQuery)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()

		// count of records returned
		count := 0

		for rows.Next() {
			var obj Vehicle

			err = rows.Scan(&obj.Id, &obj.Make, &obj.Model, &obj.Year)

			// if the scan was successful, load the row
			if err == nil {
				results = append(results, obj)
				count++
			}
		}

		// show count of successfully processed rows
		log.Println("Rows returned: " + strconv.Itoa(count))

		if err = rows.Err(); err != nil {
			// Abnormal termination of the rows loop
			// close should be called automatically in this case
			log.Println(err)
		}

		c.JSON(200, results)
		//c.JSON(200, storeData)
	})

	//POST '/cars'  --> create cars
	router.POST("/cars", func(c *gin.Context) {
		var car Vehicle

		err := c.ShouldBindJSON(&car) //binds the input data into 'motor' var
		if err != nil {
			c.JSON(422, gin.H{"message": "Unprocessable entity!"})
			return
		}

		/*
					InsertCarQuery string = `INSERT INTO cars
					(make, model, year)
					VALUES($1, $2, $3)
			    RETURNING id, make, model, year`
		*/

		err = db.QueryRow(InsertCarQuery, car.Make, car.Model, car.Year).Scan(
			&car.Id, &car.Make, &car.Model, &car.Year)

		if err != nil {
			c.JSON(500, gin.H{"message": "Internal server error!"})
		}

		/*
			car.Id = nextId
			nextId++

			storeData = append(storeData, car)
		*/
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

		err = db.QueryRow(SelectCarQuery, carid).Scan(&car.Id,
			&car.Make, &car.Model, &car.Year)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"message": "Car not found!"})
				return
			}
			c.JSON(500, gin.H{"message": "Internal server error!"})
			return
		}

		/*

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
		*/

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

		err = db.QueryRow(UpdateCarQuery, carid, car.Make, car.Model, car.Year).Scan(
			&car.Id, &car.Make, &car.Model, &car.Year)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"message": "Car not found!"})
				return
			}
			c.JSON(500, gin.H{"message": "Internal server error!"})
			return
		}

		/*
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
		*/
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

		res, err := db.Exec(DeleteCarQuery, carid)
		if err != nil {
			c.JSON(500, gin.H{"message": "Internal server error!"})
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			c.JSON(500, gin.H{"message": "Internal server error!"})
			return
		}
		if count == 0 {
			c.JSON(404, gin.H{
				"message": "car not found",
			})
			return
		}

		/*
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
		*/
		c.JSON(200, gin.H{
			"message": "car successfully deleted",
		})

	})

	router.Run()
}
