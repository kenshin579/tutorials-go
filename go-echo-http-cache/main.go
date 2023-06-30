package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/SporkHubr/echo-http-cache"
	"github.com/SporkHubr/echo-http-cache/adapter/redis"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	{ID: 3, Name: "Bob Johnson", Email: "bob@example.com"},
}

func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

func getUser(c echo.Context) error {
	id := c.Param("id")

	for _, user := range users {
		if id == strconv.Itoa(user.ID) {
			return c.JSON(http.StatusOK, user)
		}
	}

	return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found"})
}

func main() {
	ringOpt := &redis.RingOptions{
		Addrs: map[string]string{
			"stock-api.advenoh.pe.kr": ":16379",
		},
	}

	// todo: for some reason, it does not work.
	// redis: 2023/07/01 00:56:55 ring.go:325: ring shard state changed: Redis<:16379 db:0> is down
	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(redis.NewAdapter(ringOpt)),
		cache.ClientWithTTL(10*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err == nil {
		router := echo.New()
		router.Use(cacheClient.Middleware())

		// Routes
		router.GET("/users", getUsers)
		router.GET("/users/:id", getUser)

		// Start the server
		router.Start(":8080")
	}

}
