package controllers_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/suite"

	"github.com/egsam98/users-todos/pkg/env"
	"github.com/egsam98/users-todos/pkg/testutils"
	"github.com/egsam98/users-todos/users/controllers"
	"github.com/egsam98/users-todos/users/db"
	env2 "github.com/egsam98/users-todos/users/utils/env"
)

const (
	migrationsFolder = "../migrations"
	errorEmptyJSON   = `{"error": {"body": "EOF"}}`
)

type usersControllerSuite struct {
	suite.Suite
	db     *sql.DB
	router *gin.Engine
}

func (s *usersControllerSuite) SetupSuite() {
	var environment env2.Environment
	env.InitEnvironment(&environment)

	database := db.Init(environment.Database.Driver, environment.Database.ConnTest)
	s.NoError(goose.Up(database, migrationsFolder))
	s.db = database

	gin.SetMode(gin.TestMode)
	s.router = controllers.Init(environment, db.New(database))
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
	res := testutils.RunGinTest(s.router, req)

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
			res := testutils.RunGinTest(s.router, req)
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
			res := testutils.RunGinTest(s.router, req)
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
	res := testutils.RunGinTest(s.router, req)
	s.Equal(http.StatusBadRequest, res.Code)
	s.JSONEq(`{"error": {"password_confirmation":"must match password field"}}`, res.Body.String())
}

func (s *usersControllerSuite) TestSignup_When_RequestBodyIsCorrect() {
	req := testutils.NewRequestJSON("POST", "/signup", gin.H{
		"username":              "someuser",
		"password":              "somepassword",
		"password_confirmation": "somepassword",
	})
	res := testutils.RunGinTest(s.router, req)
	s.Equal(http.StatusCreated, res.Code)
}

func (s *usersControllerSuite) TestSignin_When_RequestBodyIsEmpty() {
	req := testutils.NewRequestJSON("POST", "/signin", nil)
	res := testutils.RunGinTest(s.router, req)
	s.Equal(http.StatusBadRequest, res.Code)
	s.JSONEq(errorEmptyJSON, res.Body.String())
}

func (s *usersControllerSuite) TestSignin_When_WrongCredentialsProvided() {
	s.Run("when username error", func() {
		req := testutils.NewRequestJSON("POST", "/signin", gin.H{
			"username": nil,
			"password": "pass",
		})
		res := testutils.RunGinTest(s.router, req)
		s.Equal(http.StatusBadRequest, res.Code)
		s.JSONEq(`{"error": {"username":"must be non empty"}}`, res.Body.String())
	})
	s.Run("when password error", func() {
		req := testutils.NewRequestJSON("POST", "/signin", gin.H{
			"username": "someuser",
		})
		res := testutils.RunGinTest(s.router, req)
		s.Equal(http.StatusBadRequest, res.Code)
		s.JSONEq(`{"error": {"password":"must be non empty"}}`, res.Body.String())
	})
}

func (s *usersControllerSuite) TestSignin_When_UserNotFoundByUsernameAndPassword() {
	req := testutils.NewRequestJSON("POST", "/signin", gin.H{
		"username": "someuser",
		"password": "pass",
	})
	res := testutils.RunGinTest(s.router, req)
	s.Equal(http.StatusUnauthorized, res.Code)
	s.JSONEq(`{"error":"username or/and password is incorrect"}`, res.Body.String())
}

func (s *usersControllerSuite) TestSignin_When_UserFoundByUsernameAndPassword() {
	// Регистрация пользователя в базе
	username := "someuser"
	password := "password"
	res := testutils.RunGinTest(s.router, testutils.NewRequestJSON("POST", "/signup", gin.H{
		"username":              username,
		"password":              password,
		"password_confirmation": password,
	}))
	s.Equal(http.StatusCreated, res.Code)

	req := testutils.NewRequestJSON("POST", "/signin", gin.H{
		"username": username,
		"password": password,
	})
	res = testutils.RunGinTest(s.router, req)
	s.Equal(http.StatusOK, res.Code)

	var body gin.H
	s.NoError(json.NewDecoder(res.Body).Decode(&body))
	token, ok := body["token"]
	if !ok {
		s.Fail(`"token" key does not exist in JSON-response`)
	}
	s.IsType("string", token)
}
