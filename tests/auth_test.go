//go:build auth
// +build auth

package tests

import (
	"avito-test-task/internal/delivery/handlers"
	"avito-test-task/internal/domain"
	"avito-test-task/internal/repo"
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-bdd/gobdd"
	"github.com/go-telegram/bot"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var uid string

func register(t gobdd.StepTest, ctx gobdd.Context, tg, password, role string) {
	request := domain.RegisterUserRequest{tg, os.Getenv("PASSWORD"), role}
	body, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("%v", err)
	}
	reader := bytes.NewReader(body)

	resp, err := http.Post("http://0.0.0.0:8081/register", "application/json", reader)
	if err != nil {
		t.Fatalf("%v", err)
	}

	var r domain.RegisterUserResponse
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		t.Fatalf("%v", err)
	}

	uid = r.UserID.String()
	ctx.Set("userId", r.UserID)
}

func login(t gobdd.StepTest, ctx gobdd.Context, userId, password string) {
	id, err := ctx.Get("userId")
	if err != nil {
		t.Errorf("%v", err)
	}

	request := domain.LoginUserRequest{id.(uuid.UUID), os.Getenv("PASSWORD")}

	body, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("%v", err)
	}
	reader := bytes.NewReader(body)

	resp, err := http.Post("http://0.0.0.0:8081/login", "application/json", reader)
	var r domain.LoginUserResponse

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		t.Fatalf("%v", err)
	}

	ctx.Set("message", r.Message)
}

func checkLogin(t gobdd.StepTest, ctx gobdd.Context, expected string) {
	message, err := ctx.Get("message")
	if err != nil {
		t.Errorf("%v", err)
	}
	msg := message.(string)

	if !strings.Contains(msg, expected) {
		t.Errorf("doesn t send code")
	}
}

func verificate(t gobdd.StepTest, ctx gobdd.Context) {
	tg, _ := bot.New(os.Getenv("TOKEN"))
	params := bot.GetChatParams{
		ChatID: "1186604465",
	}
	ch, err := tg.GetChat(context.Background(), &params)

	code := ch.PinnedMessage.Text
	intCode, _ := strconv.Atoi(code)

	request := domain.FinalLoginUserRequest{uuid.MustParse(uid), intCode}

	body, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("%v", err)
	}
	reader := bytes.NewReader(body)

	resp, err := http.Post("http://0.0.0.0:8081/finallogin", "application/json", reader)
	if err != nil {
		t.Errorf("%v", err)
	}
	var r domain.FinalLoginUserResponse

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		t.Fatalf("%v", err)
	}

	ctx.Set("token", r.Token)
}

func badVerificate(t gobdd.StepTest, ctx gobdd.Context) {
	request := domain.FinalLoginUserRequest{uuid.MustParse(uid), 0}
	body, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("%v", err)
	}
	reader := bytes.NewReader(body)

	resp, err := http.Post("http://0.0.0.0:8081/finallogin", "application/json", reader)
	if err != nil {
		t.Errorf("%v", err)
	}
	var r handlers.ErrorResponse

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		t.Fatalf("%v", err)
	}

	ctx.Set("code", r.Code)
}

func checkVerificate(t gobdd.StepTest, ctx gobdd.Context) {
	token, err := ctx.Get("token")
	if err != nil {
		t.Errorf("%v", err)
	}
	tkn := token.(string)
	if tkn == "" {
		t.Errorf("bad token")
	}
}

func badLogin(t gobdd.StepTest, ctx gobdd.Context, userId, password string) {
	request := domain.LoginUserRequest{uuid.New(), os.Getenv("PASSWORD")}

	body, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("%v", err)
	}
	reader := bytes.NewReader(body)

	resp, err := http.Post("http://0.0.0.0:8081/login", "application/json", reader)
	var r handlers.ErrorResponse

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		t.Fatalf("%v", err)
	}

	ctx.Set("code", r.Code)
}

func checkBadLogin(t gobdd.StepTest, ctx gobdd.Context) {
	code, err := ctx.Get("code")
	if err != nil {
		t.Errorf("%v", err)
	}

	intCode := code.(int)
	if intCode != handlers.LoginUserError {
		t.Errorf("bad response %d", intCode)
	}
}

func loginVer(t gobdd.StepTest, ctx gobdd.Context, userId, password string) {
	request := domain.LoginUserRequest{uuid.MustParse(uid), os.Getenv("PASSWORD")}

	body, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("%v", err)
	}
	reader := bytes.NewReader(body)

	resp, err := http.Post("http://0.0.0.0:8081/login", "application/json", reader)
	var r domain.LoginUserResponse

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		t.Fatalf("%v", err)
	}

	ctx.Set("message", r.Message)
}

func TestScenarios(t *testing.T) {
	host := os.Getenv("POSTGRES_TEST_HOST")
	port := os.Getenv("POSTGRES_TEST_PORT")
	connString := "postgres://test-user:test-password@" + host + ":" + port + "/test-db?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), connString)
	defer pool.Close()
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	retryAdapter := repo.NewPostgresRetryAdapter(pool, 3, time.Second*3)
	userRepo := repo.NewPostrgesUserRepo(pool, retryAdapter)

	suite := gobdd.NewSuite(t)
	suite.AddStep(`I register with ([\da-zA-Z0-9@\_]+), ([\da-zA-Z0-9]+), ([\da-zA-Z0-9]+)`, register)

	suite.AddStep(`I log with ([\da-zA-Z0-9@\-\_]+) and ([\da-zA-Z0-9@\-]+)`, login)
	suite.AddStep(`I get message '([a-zA-Z ]+)'`, checkLogin)
	suite.AddStep(`I provide user_id and code`, verificate)
	suite.AddStep(`I get token`, checkVerificate)
	suite.AddStep(`I log without registration with ([\da-zA-Z0-9@\-\_]+) and ([\da-zA-Z0-9@\-]+)`, badLogin)
	suite.AddStep(`I get error`, checkBadLogin)
	suite.AddStep(`I login with ([\da-zA-Z0-9@\-\_]+) and ([\da-zA-Z0-9@\-]+)`, loginVer)
	suite.AddStep(`I get message '([a-zA-Z ]+)'`, checkLogin)
	suite.AddStep(`I incorrectly provide user_id and code`, badVerificate)
	suite.AddStep(`I get error`, checkBadLogin)
	suite.Run()

	users, err := userRepo.GetAll(context.Background(), 0, 100, zap.NewNop())
	if err != nil {
		t.Errorf("%v", err)
	}

	for _, u := range users {
		if u.Mail == "@N_AR24" {
			uid = u.UserID.String()
		}
	}

	err = userRepo.DeleteByID(context.Background(), uuid.MustParse(uid), zap.NewNop())
	if err != nil {
		t.Errorf("%v", err)
	}
}
