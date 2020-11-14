package controllers_test

import (
	context2 "context"
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/suite"

	"github.com/egsam98/users-todos/pkg/dbutils"
	"github.com/egsam98/users-todos/pkg/env"
	"github.com/egsam98/users-todos/pkg/testutils"
	"github.com/egsam98/users-todos/todos/controllers"
	"github.com/egsam98/users-todos/todos/controllers/requests"
	"github.com/egsam98/users-todos/todos/controllers/responses"
	"github.com/egsam98/users-todos/todos/db"
	env2 "github.com/egsam98/users-todos/todos/utils/env"
)

const (
	migrationsFolder       = "../migrations"
	userID           int32 = 1
)

type todosControllerSuite struct {
	suite.Suite
	db     *sql.DB
	q      *db.Queries
	router *gin.Engine
}

func (s *todosControllerSuite) SetupSuite() {
	var environment env2.Environment
	env.InitEnvironment(&environment)

	database := dbutils.Init(environment.Database.Driver, environment.Database.ConnTest)
	s.NoError(goose.Up(database, migrationsFolder))
	s.db = database
	s.q = db.New(database)

	gin.SetMode(gin.TestMode)
	s.router = controllers.Init(environment, s.q)

	addUserID := func(ctx *gin.Context) { ctx.Set("userID", userID) }
	r := s.router.Group("/", addUserID)
	r.POST("/todos-test", testutils.GinHandler(s.router, "POST", "/todos"))
	r.PUT("/todos-test/:id", testutils.GinHandler(s.router, "PUT", "/todos/:id"))
	r.DELETE("/todos-test/:id", testutils.GinHandler(s.router, "DELETE", "/todos/:id"))
	r.GET("/todos-test", testutils.GinHandler(s.router, "GET", "/todos"))
	r.POST("/todos-test/before", testutils.GinHandler(s.router, "POST", "/todos/before"))
}

func TestTodosController(t *testing.T) {
	suite.Run(t, new(todosControllerSuite))
}

func (s *todosControllerSuite) TearDownTest() {
	_, err := s.db.Exec(`delete from todos`)
	s.NoError(err)
}

func (s *todosControllerSuite) TestCreateTodo_When_JWTIsRequired() {
	testutils.JwtAuthRequired(s.T(), s.router, "POST", "/todos")
}

func (s *todosControllerSuite) TestCreateTodo_When_BodyIsEmpty() {
	req := testutils.NewRequestJSON("POST", "/todos-test", nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusBadRequest, res.Code)
	s.JSONEq(`{"error": {"body":"EOF"}}`, res.Body.String())
}

func (s *todosControllerSuite) TestCreateTodo_When_BodyInvalid() {
	s.Run("required params", func() {
		req := testutils.NewRequestJSON("POST", "/todos-test", gin.H{})
		res := testutils.RunHTTPTest(s.router, req)
		s.Equal(http.StatusBadRequest, res.Code)
		s.JSONEq(`{"error": {"title":"must be non empty"}}`, res.Body.String())
	})

	s.Run("invalid deadline", func() {
		req := testutils.NewRequestJSON("POST", "/todos-test", gin.H{
			"title":    "A",
			"deadline": "10-10-2020 15:00:13",
		})
		res := testutils.RunHTTPTest(s.router, req)
		s.Equal(http.StatusBadRequest, res.Code)
		s.JSONEq(`{"error": {"deadline":"must be int64"}}`, res.Body.String())
	})

	s.Run("invalid description", func() {
		req := testutils.NewRequestJSON("POST", "/todos-test", gin.H{
			"title":       "A",
			"description": 1,
		})
		res := testutils.RunHTTPTest(s.router, req)
		s.Equal(http.StatusBadRequest, res.Code)
		s.JSONEq(`{"error": {"description":"must be string"}}`, res.Body.String())
	})
}

func (s *todosControllerSuite) TestCreateTodo() {
	title := "A"
	desc := "desc"
	var deadline int64 = 1234
	req := testutils.NewRequestJSON("POST", "/todos-test", requests.NewTodo{
		Title:       title,
		Description: &desc,
		Deadline:    &deadline,
	})

	expectedRes := responses.Todo{
		Title:       title,
		Description: &desc,
		Deadline:    &deadline,
		UserID:      userID,
	}

	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusCreated, res.Code)

	var actualRes responses.Todo
	testutils.DecodeBody(res, &actualRes)

	expectedRes.ID = actualRes.ID
	s.Equal(expectedRes, actualRes, "expected Todo doesn't match actual one")
}

func (s *todosControllerSuite) TestUpdateTodo_When_JWTIsRequired() {
	testutils.JwtAuthRequired(s.T(), s.router, "PUT", "/todos/0")
}

func (s *todosControllerSuite) TestUpdateTodo_When_IDNotAnInteger() {
	req := testutils.NewRequestJSON("PUT", "/todos-test/a", nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusBadRequest, res.Code)
	s.JSONEq(`{"error":"id must be integer"}`, res.Body.String())
}

func (s *todosControllerSuite) TestUpdateTodo_When_BodyInvalid() {
	s.Run("when required keys're not provided", func() {
		req := testutils.NewRequestJSON("PUT", "/todos-test/1", gin.H{"title": "title"})
		res := testutils.RunHTTPTest(s.router, req)
		s.Equal(http.StatusBadRequest, res.Code)
		s.JSONEq(`{"error":"'description', 'deadline' must present"}`, res.Body.String())
	})

	s.Run("when invalid data", func() {
		req := testutils.NewRequestJSON("PUT", "/todos-test/1", gin.H{
			"title":       1,
			"description": "a",
			"deadline":    1234,
		})
		res := testutils.RunHTTPTest(s.router, req)
		s.Equal(http.StatusBadRequest, res.Code)
		s.JSONEq(`{"error": {"title":"must be string"}}`, res.Body.String())
	})
}

func (s *todosControllerSuite) TestUpdateTodo_When_TodoNotFound() {
	req := testutils.NewRequestJSON("PUT", "/todos-test/0", gin.H{
		"title":       "title",
		"description": "a",
		"deadline":    1234,
	})
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusNotFound, res.Code)
	s.JSONEq(`{"error":"todo is not found"}`, res.Body.String())
}

func (s *todosControllerSuite) TestUpdateTodo_When_NoAccess() {
	todo, err := s.q.CreateTodo(context2.TODO(), db.CreateTodoParams{UserID: userID + 1})
	s.NoError(err)

	req := testutils.NewRequestJSON("PUT", fmt.Sprintf("/todos-test/%d", todo.ID), gin.H{
		"title":       "title",
		"description": "a",
		"deadline":    1234,
	})
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusForbidden, res.Code)
	s.JSONEq(`{"error":"you have no access to this todo"}`, res.Body.String())
}

func (s *todosControllerSuite) TestUpdateTodo() {
	todo, err := s.q.CreateTodo(context2.TODO(), db.CreateTodoParams{UserID: userID})
	s.NoError(err)

	title := "title"
	desc := "a"
	req := testutils.NewRequestJSON("PUT", fmt.Sprintf("/todos-test/%d", todo.ID), gin.H{
		"title":       title,
		"description": desc,
		"deadline":    nil,
	})

	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusOK, res.Code)

	expected := responses.NewTodo(todo)
	expected.Title = title
	expected.Description = &desc

	var actual responses.Todo
	testutils.DecodeBody(res, &actual)

	s.Equal(expected, actual)
}

func (s *todosControllerSuite) TestDeleteTodo_When_IDNotAnInteger() {
	req := testutils.NewRequestJSON("DELETE", "/todos-test/a", nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusBadRequest, res.Code)
	s.JSONEq(`{"error":"id must be integer"}`, res.Body.String())
}

func (s *todosControllerSuite) TestDeleteTodo_When_TodoNotFound() {
	s.Run("todo doesn't exist in database", func() {
		req := testutils.NewRequestJSON("DELETE", "/todos-test/0", nil)
		res := testutils.RunHTTPTest(s.router, req)
		s.Equal(http.StatusNotFound, res.Code)
		s.JSONEq(`{"error":"todo is not found"}`, res.Body.String())
	})

	s.Run("todo exists but it belongs to another user", func() {
		todo, err := s.q.CreateTodo(context2.TODO(), db.CreateTodoParams{
			UserID: userID + 1,
		})
		s.NoError(err)

		req := testutils.NewRequestJSON("DELETE", fmt.Sprintf("/todos-test/%d", todo.ID), nil)
		res := testutils.RunHTTPTest(s.router, req)
		s.Equal(http.StatusNotFound, res.Code)
		s.JSONEq(`{"error":"todo is not found"}`, res.Body.String())
	})
}

func (s *todosControllerSuite) TestDeleteTodo() {
	todo, err := s.q.CreateTodo(context2.TODO(), db.CreateTodoParams{UserID: userID})
	s.NoError(err)

	req := testutils.NewRequestJSON("DELETE", fmt.Sprintf("/todos-test/%d", todo.ID), nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusOK, res.Code)
}

func (s *todosControllerSuite) TestFetchAll_When_NoTodos() {
	req := testutils.NewRequestJSON("GET", "/todos-test", nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusOK, res.Code)
	s.JSONEq(`[]`, res.Body.String())
}

func (s *todosControllerSuite) TestFetchAll() {
	ctx := context2.TODO()

	deadlines := []sql.NullTime{
		{Time: time.Now(), Valid: true},
		{Valid: false},
		{Time: time.Now().Add(time.Hour), Valid: true},
	}
	todos := make([]db.Todo, len(deadlines))
	for i, _ := range deadlines {
		var err error
		todos[i], err = s.q.CreateTodo(ctx, db.CreateTodoParams{
			Deadline: deadlines[i],
			UserID:   userID,
		})
		s.NoError(err)
	}

	req := testutils.NewRequestJSON("GET", "/todos-test", nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusOK, res.Code)

	var actual []responses.Todo
	testutils.DecodeBody(res, &actual)
	s.ElementsMatch(responses.NewTodos(todos), actual)

	// Todo с имеющимся значением deadline имеет приоритет на todo со значением nil
	isSortedByDeadline := sort.SliceIsSorted(actual, func(i, j int) bool {
		deadline1 := actual[i].Deadline
		deadline2 := actual[j].Deadline

		if deadline1 == nil {
			return false
		}
		if deadline2 == nil {
			return true
		}
		return *deadline1 < *deadline2
	})
	s.True(isSortedByDeadline, "todos are not sorted by deadline")
}

func (s *todosControllerSuite) TestFetchBeforeDeadline_When_RequestBodyIsEmpty() {
	req := testutils.NewRequestJSON("POST", "/todos-test/before", nil)
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusBadRequest, res.Code)
	s.JSONEq(`{"error": {"body":"EOF"}}`, res.Body.String())
}

func (s *todosControllerSuite) TestFetchBeforeDeadline_When_DeadlineInvalid() {
	req := testutils.NewRequestJSON("POST", "/todos-test/before", gin.H{
		"deadline": "10-10-1998 10:23:32",
	})
	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusBadRequest, res.Code)
	s.JSONEq(`{"error": {"deadline":"must be int64"}}`, res.Body.String())
}

func (s *todosControllerSuite) TestFetchBeforeDeadline() {
	withCurrentDeadline := func() db.Todo {
		ctx := context2.TODO()
		todo, err := s.q.CreateTodo(ctx, db.CreateTodoParams{
			UserID:   userID,
			Deadline: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		})
		s.NoError(err)
		return todo
	}

	todo1 := withCurrentDeadline()
	time.Sleep(time.Second)

	req := testutils.NewRequestJSON("POST", "/todos-test/before", requests.Deadline{
		Value: time.Now().Unix(),
	})
	_ = withCurrentDeadline()

	res := testutils.RunHTTPTest(s.router, req)
	s.Equal(http.StatusOK, res.Code)

	var actual []responses.Todo
	testutils.DecodeBody(res, &actual)
	s.Len(actual, 1)
	s.Equal(todo1.ID, actual[0].ID)
}
