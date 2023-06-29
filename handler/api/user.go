package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserTaskCategory(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Register(c *gin.Context) {
	var user model.UserRegister

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	var recordUser = model.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	recordUser, err := u.userService.Register(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse("register success"))
}

func (u *userAPI) Login(c *gin.Context) {
	user := model.UserLogin{}
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid decode json"})
		return
	}

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid decode json"})
		return
	}

	userLogin := model.User{
		Email: user.Email,
		Password: user.Password,
	}

	tkn, err := u.userService.Login(&userLogin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "error internal server"})
		return
	}
	expiration := time.Now().Add(5 * time.Minute)

	http.SetCookie(c.Writer, &http.Cookie{
		Name: "session_token",
		Value: *tkn,
		Expires: expiration,
	})

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "login success"})
}

func (u *userAPI) GetUserTaskCategory(c *gin.Context) {
	users, err := u.userService.GetUserTaskCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "error internal server"})
		return
	}
	
	c.JSON(http.StatusOK, users)
}
