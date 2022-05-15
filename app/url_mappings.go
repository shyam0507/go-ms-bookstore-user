package app

import (
	"go-ms-bookstore-user/controllers/ping"

	"go-ms-bookstore-user/controllers/users"
)

func MapUrls() {
	router.GET("/ping", ping.Ping)
	router.GET("users/:user_id", users.GetUser)
	router.POST("/users", users.CreateUser)
	router.PUT("/users/:user_id", users.UpdateUser)
	router.PATCH("/users/:user_id", users.UpdateUser)
	router.DELETE("/users/:user_id", users.DeleteUser)
	router.GET("/internal/users/search", users.SearchUser)

}
