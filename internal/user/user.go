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

type OnlineUser struct {
	Username string `json:"username"`
	Color    int    `json:"color"`
}

var (
	Users chan OnlineUser
)

func GetInitialUsers(serverAddr string) {
	u := url.URL{Scheme: "http", Host: serverAddr, Path: "/online-users"}
	res, err := http.Get(u.String())
	if err != nil {
		log.Println("error: GetInitialUsers couldn't fetch data")
		panic(0)
	}
	var data []OnlineUser
	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("error: GetInitialUsers couldn't parse data")
		panic(0)
	}
	Users = make(chan OnlineUser, 100)
	for _, user := range data {
		if user.Username != "" {
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

func AddUser(newUser OnlineUser) {
	Users <- newUser
}

func RemoveUser(staleUser OnlineUser) {
	staleUser.Username = " " + staleUser.Username
	Users <- staleUser
}
