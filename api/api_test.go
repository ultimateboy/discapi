package api_test

import (
	"fmt"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/verdverm/frisby"

	"github.com/ultimateboy/discapi/api"
	"github.com/ultimateboy/discapi/config"
)

func TestApp(t *testing.T) {
	cfg := &config.Config{
		InMemoryStorage: true,
		Debug:           true,
		APIPort:         12345,
	}
	api, err := api.NewAPI(cfg)
	if err != nil {
		t.Error(err)
	}
	go api.Start()
	time.Sleep(1 * time.Second)

	basepath := fmt.Sprintf("http://localhost:%d/", cfg.APIPort)

	frisby.Create("Get healthz").
		Get(basepath+"healthz").
		Send().
		ExpectStatus(200).
		ExpectJson("status", "ok")

	frisby.Create("Get unknown resource").
		Get(basepath+"this-path-wont-exist").
		Send().
		ExpectStatus(404).
		ExpectJson("message", "Resource Not Found")

	frisby.Create("Post authentication with bad data").
		Post(basepath + "auth").
		SetJson(map[string]string{"email": "bad@example.com", "password": "user-does-not-exist"}).
		Send().
		ExpectStatus(401)

	testUserID := ""
	frisby.Create("Post first user").
		Post(basepath + "users").
		SetJson(map[string]string{"email": "test@example.com", "password": "test-pass"}).
		Send().
		ExpectStatus(201).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
			testUserID, err = json.Get("id").String()
			if err != nil {
				t.Fatalf("failed to find id in user post")
			}
		})

	frisby.Create("Post user that already exists").
		Post(basepath + "users").
		SetJson(map[string]string{"email": "test@example.com", "password": "some-other-pass"}).
		Send().
		ExpectStatus(409)

	frisby.Create("Post authentication with incorrect password").
		Post(basepath + "auth").
		SetJson(map[string]string{"email": "test@example.com", "password": "bad-pass"}).
		Send().
		ExpectStatus(401)

	testJWTValue := ""
	frisby.Create("Post authentication with correct password").
		Post(basepath + "auth").
		SetJson(map[string]string{"email": "test@example.com", "password": "test-pass"}).
		Send().
		ExpectStatus(201).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
			testJWTValue, err = json.Get("jwt").String()
			if err != nil {
				t.Fatalf("failed to find jwt in authenticaiton post")
			}

			// check the plain-text password doesnt get sent back
			_, err = json.Get("password").String()
			if err == nil {
				t.Fatalf("plain-text password found in response to authentication post")
			}
		})

	frisby.Global.PrintReport()
}
