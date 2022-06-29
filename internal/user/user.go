package user

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	Users chan string
)

func GetInitialUsers(serverAddr string) {
	u := url.URL{Scheme: "http", Host: serverAddr, Path: "/online-users"}
	res, err := http.Get(u.String())
	if err != nil {
		log.Println("error: GetInitialUsers couldn't fetch data")
		panic(0)
	}
	var data []string
	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("error: GetInitialUsers couldn't parse data")
		panic(0)
	}
	Users = make(chan string, 100)
	for _, user := range data {
		if user != "" {
			Users <- user
		}
	}
}

func IsValidUsername(username string) error {
	if username == "" {
		return errors.New("username must not be empty")
	}
	arr := strings.Split(username, " ")
	if len(arr) == 1 {
		return nil
	} else {
		return errors.New("username must not contain any space")
	}
}

func AddUser(username string) {
	Users <- username
}

func RemoveUser(username string) {
	Users <- " " + username
}
