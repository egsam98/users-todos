package controllers_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/suite"

	"github.com/egsam98/users-todos/pkg/dbutils"
	"github.com/egsam98/users-todos/pkg/env"
	"github.com/egsam98/users-todos/pkg/testutils"
	"github.com/egsam98/users-todos/users/controllers"
	"github.com/egsam98/users-todos/users/db"
	env2 "github.com/egsam98/users-todos/users/utils/env"
	"github.com/egsam98/users-todos/users/utils/testmocks"
)

const (
	migrationsFolder = "../migrations"
	errorEmptyJSON   = `{"error": {"body": "EOF"}}`
)

type usersControllerSuite struct {
	suite.Suite
	db     *sql.DB
	q      *db.Queries
	router *gin.Engine
}

func (s *usersControllerSuite) SetupSuite() {
	var environment env2.Environment
	env.InitEnvironment(&environment)

	database := dbutils.Init(environment.Database.Driver, environment.Database.ConnTest)
	s.NoError(goose.Up(database, migrationsFolder))
	s.db = database
	s.q = db.New(database)

	gin.SetMode(gin.TestMode)
	s.router = controllers.Init(environment, s.q)
	s.router.GET("/users-test/:id", testutils.GinHandler(s.router, "GET", "/users/:id"))
}

func (s *usersControllerSuite) TearDownTest() {
	_, err := s.db.Exec(`delete from users`)
	s.NoError(err)
}

func TestUsersController(t *testing.T) {
	suite.Run(t, new(usersControllerSuite))
}

func (s *usersControllerSuite) TestSignup_When_RequestBodyIsEmpty() {
	req := testutils.NewRequestJSON("POST", "/signup", nil)
	res := testutils.RunHTTPTest(s.router, req)

	s.Equal(http.StatusBadRequest, res.Code)
	s.JSONEq(errorEmptyJSON, res.Body.String())
}

func (s *usersControllerSuite) TestSignup_When_WrongUsernameProvided() {
	word13length := "TTTTTTTTTTTTT"
	bodyVariants := []gin.H{
		{"username": nil, "password": "testing", "password_confirmation": "testing"},
		{"username": "", "password": "testing", "password_confirmation": "testing"},
		{"username": word13length, "password": "testing", "password_confirmation": "testing"},
	}

	for _, variant := range bodyVariants {
		s.Run(fmt.Sprintf("when username = %v", variant["username"]), func() {
			req := testutils.NewRequestJSON("POST", "/signup", variant)
			res := testutils.RunHTTPTest(s.router, req)
			s.Equal(http.StatusBadRequest, res.Code)
			s.JSONEq(`{"error": {"username": "must be non empty and have max length 12"}}`, res.Body.String())
		})
	}
}

func (s *usersControllerSuite) TestSignup_When_WrongPasswordProvided() {
	bodyVariants := []gin.H{
		{"username": "testing", "password": nil, "password_confirmation": nil},
		{"username": "testing", "password": "", "password_confirmation": ""},
		{"username": "testing", "password": "passw", "password_confirmation": "passw"},                 // length 5
		{"username": "testing", "password": "passwordpassw", "password_confirmation": "passwordpassw"}, // length 13
		{"username": "testing", "password": "password", "password_confirmation": "Password"},           // mismatch with confirmation
	}

	for i, variant := range bodyVariants[:len(bodyVariants)-1] {
		testName := fmt.Sprintf("when password = %v, password_confirmation = %v", variant["password"], variant["password_confirmation"])
		s.Run(testName, func() {
			req := testutils.NewRequestJSON("POST", "/signup", variant)
			res := testutils.RunHTTPTest(s.router, req)
			s.Equal(http.StatusBadRequest, res.Code)

			jsonRes := ""
			switch i {
			case 0, 1:
				jsonRes = `{"error": {"password": "must have length 6..12 symbols", "password_confirmation":"must match password field"}}`
			case 2, 3:
				jsonRes = `{"error": {"password": "must have length 6..12 symbols"}}`
			default:
				jsonRes = `{"password_confirmation":"must match password field"}`
			}
			s.JSONEq(jsonRes, res.Body.String())
		})
	}

	req := testutils.NewRequestJSON("POST", "/signup", bodyVariants[len(bodyVariants)-1])
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusBadRequest, res.Code)
	s.JSONEq(`{"error": {"password_confirmation":"must match password field"}}`, res.Body.String())
}

func (s *usersControllerSuite) TestSignup_When_RequestBodyIsCorrect() {
	req := testutils.NewRequestJSON("POST", "/signup", gin.H{
		"username":              "someuser",
		"password":              "somepassword",
		"password_confirmation": "somepassword",
	})
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusCreated, res.Code)
}

func (s *usersControllerSuite) TestSignin_When_RequestBodyIsEmpty() {
	req := testutils.NewRequestJSON("POST", "/signin", nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusBadRequest, res.Code)
	s.JSONEq(errorEmptyJSON, res.Body.String())
}

func (s *usersControllerSuite) TestSignin_When_WrongCredentialsProvided() {
	s.Run("when username error", func() {
		req := testutils.NewRequestJSON("POST", "/signin", gin.H{
			"username": nil,
			"password": "pass",
		})
		res := testutils.RunHTTPTest(s.router, req)
		s.Equal(http.StatusBadRequest, res.Code)
		s.JSONEq(`{"error": {"username":"must be non empty"}}`, res.Body.String())
	})
	s.Run("when password error", func() {
		req := testutils.NewRequestJSON("POST", "/signin", gin.H{
			"username": "someuser",
		})
		res := testutils.RunHTTPTest(s.router, req)
		s.Equal(http.StatusBadRequest, res.Code)
		s.JSONEq(`{"error": {"password":"must be non empty"}}`, res.Body.String())
	})
}

func (s *usersControllerSuite) TestSignin_When_UserNotFoundByUsernameAndPassword() {
	req := testutils.NewRequestJSON("POST", "/signin", gin.H{
		"username": "someuser",
		"password": "pass",
	})
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusUnauthorized, res.Code)
	s.JSONEq(`{"error":"username or/and password is incorrect"}`, res.Body.String())
}

func (s *usersControllerSuite) TestSignin_When_UserFoundByUsernameAndPassword() {
	// Регистрация пользователя в базе
	username := "someuser"
	password := "password"
	res := testutils.RunHTTPTest(s.router, testutils.NewRequestJSON("POST", "/signup", gin.H{
		"username":              username,
		"password":              password,
		"password_confirmation": password,
	}))
	s.Equal(http.StatusCreated, res.Code)

	req := testutils.NewRequestJSON("POST", "/signin", gin.H{
		"username": username,
		"password": password,
	})
	res = testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusOK, res.Code)

	var body gin.H
	testutils.DecodeBody(res, &body)
	token, ok := body["token"]
	if !ok {
		s.Fail(`"token" key does not exist in JSON-response`)
	}
	s.IsType("string", token)
}

func (s *usersControllerSuite) TestFetchUser_When_JWTIsRequired() {
	testutils.JwtAuthRequired(s.T(), s.router, "GET", "/users/1")
}

func (s *usersControllerSuite) TestFetchUser_When_IDNotAnInteger() {
	req := testutils.NewRequestJSON("GET", "/users-test/a", nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusBadRequest, res.Code)
	s.JSONEq(`{"error": "user ID must be integer"}`, res.Body.String())
}

func (s *usersControllerSuite) TestFetchUser_When_UserNotFound() {
	req := testutils.NewRequestJSON("GET", "/users-test/0", nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusNotFound, res.Code)
	s.JSONEq(`{"error": "user ID=0 is not found"}`, res.Body.String())
}

func (s *usersControllerSuite) TestFetchUser() {
	user, err := db.New(s.db).CreateUser(context.TODO(), db.CreateUserParams{
		Username: "someuser",
		Password: "%f$dzzrd$%d",
	})
	s.NoError(err)

	req := testutils.NewRequestJSON("GET", fmt.Sprintf("/users-test/%d", user.ID), nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusOK, res.Code)
	s.JSONEq(fmt.Sprintf(`{"id":%d, "username":"someuser"}`, user.ID), res.Body.String())
}

func (s *usersControllerSuite) TestAuth_When_NoAuthorizationHeader() {
	req := testutils.NewRequestJSON("POST", "/auth", nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusForbidden, res.Code)
	s.JSONEq(`{"error": "Authorization header is not provided"}`, res.Body.String())
}

func (s *usersControllerSuite) TestAuth_When_TokenInvalid() {
	req := testutils.NewRequestJSON("POST", "/auth", nil)
	req.Header.Set("Authorization", "t.t.t")
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusForbidden, res.Code)

	s.T().Logf("Response body: %v", res.Body)

	var jsonBody gin.H
	testutils.DecodeBody(res, &jsonBody)
	if err, ok := jsonBody["error"]; !ok {
		s.Fail("no 'error' key in response body")
	} else {
		if _, ok := err.(map[string]interface{})["jwt"]; !ok {
			s.Fail("no 'error.jwt' key in response body")
		}
	}
}

func (s *usersControllerSuite) TestAuth() {
	tokenServiceMock := testmocks.NewTokenServiceMock(s.q)
	token := "some_token"
	expectedUser := &db.User{
		ID:       0,
		Username: token,
	}
	tokenServiceMock.On("Parse", token).Return(expectedUser, nil)

	controller := &controllers.UsersController{TokenService: tokenServiceMock}
	s.router.POST("/users-test/auth", controller.Auth(false))
	req := testutils.NewRequestJSON("POST", "/users-test/auth", nil)
	req.Header.Set("Authorization", token)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusOK, res.Code)

	s.T().Logf("Response body: %v", res.Body)

	var user db.User
	testutils.DecodeBody(res, &user)
	s.Equal(*expectedUser, user)
}
