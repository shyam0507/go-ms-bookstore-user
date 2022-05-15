package users

import (
	"fmt"
	"go-ms-bookstore-user/domain/users"
	"go-ms-bookstore-user/services"
	"go-ms-bookstore-user/utils/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("User Id should be a number")
		c.JSON(err.Status, err)
		return
	}

	user, err := services.UserService.GetUser(userId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, user.Marshall(c.GetHeader("X-Public") == "true"))
}

func CreateUser(c *gin.Context) {
	var user users.User

	// bytes, err := ioutil.ReadAll(c.Request.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// // fmt.Println(string(bytes))

	// // convert the byte to JSON
	// err = json.Unmarshal(bytes, &user)

	err := c.ShouldBindJSON(&user)

	if err != nil {
		fmt.Println(err.Error())
		restErr := errors.NewBadRequestError("invalid JSON body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.UserService.CreateUser(user)
	if saveErr != nil {
		fmt.Println(saveErr)
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func UpdateUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("User Id should be a number")
		c.JSON(err.Status, err)
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid JSON body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId

	isPartial := c.Request.Method == http.MethodPatch
	fmt.Println(isPartial)

	result, err := services.UserService.UpdateUser(isPartial, user)
	if err != nil {
		fmt.Println(err)
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))

}

func DeleteUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("User Id should be a number")
		c.JSON(err.Status, err)
	}

	user := users.User{Id: userId}

	err := services.UserService.DeleteUser(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func SearchUser(c *gin.Context) {
	status, found := c.GetQuery("status")

	if !found {
		err := errors.NewBadRequestError("Please pass correct search query")
		c.JSON(err.Status, err)
	}

	users, err := services.UserService.SearchUser(status)

	if err != nil {
		fmt.Println(err)
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, users.Marshall(c.GetHeader("X-Public") == "true"))
}
