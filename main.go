package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-redis/db"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	db.RedisInit()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
  
	// Routes
	e.POST("/insert", InsertRedis)
	e.GET("/get", GetRedis)
  
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

type Response struct{
	Data interface{}
	Status string
}

type Req struct {
	Name string
	Age string
}

var key = "app_ardi"
func GetRedis(c echo.Context) error {
	id := c.QueryParam("id")

	rdb := db.RedisConnect()
	val, err := rdb.HGet(context.Background(), key, id).Result()
	if err == redis.Nil{
		return c.JSON(http.StatusNotFound, "Data ga ada")
	}else if err != nil{
		return c.JSON(http.StatusBadRequest, "Err Get Redis")

	}
	var requestRedis Req
	err = json.Unmarshal([]byte(val), &requestRedis)
	if err != nil{
		return c.JSON(http.StatusNotFound, fmt.Sprintf("err unmarshal %s", err.Error()))
	}
	respon := Response{
		Data: requestRedis,
		Status: "Success",
	}
	return c.JSON(http.StatusOK, respon)
}
func InsertRedis(c echo.Context) error {
	id := c.QueryParam("id")
	name := c.QueryParam("name")
	age := c.QueryParam("age")

	rdb := db.RedisConnect()
	reqRedis := Req{
		Name: name,
		Age: age,
	}

	request, _ := json.Marshal(reqRedis)
	err := rdb.HSet(context.Background(),key, id, request).Err()
	if err != nil{
		return fmt.Errorf("error set redis %s", err)
	}

	resp := Response{
		Data: id,
		Status: "Success",
	}

	return c.JSON(http.StatusOK, resp)
  }