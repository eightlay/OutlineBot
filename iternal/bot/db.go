package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var usersPath string = "users.json"
var users map[string]user

func readDB() {
	usersSource, err := ioutil.ReadFile(usersPath)
	if err != nil {
		writeDB([]byte("{}"))
	}

	json.Unmarshal(usersSource, &users)

	fmt.Println(users)
}

func writeDB(data ...[]byte) {
	var err error
	if len(data) == 0 {
		toWrite, _ := json.Marshal(users)
		err = ioutil.WriteFile(usersPath, toWrite, 0644)
	} else {
		err = ioutil.WriteFile(usersPath, data[0], 0644)
	}
	if err != nil {
		panic("can't create users database")
	}
}

func addUser(telegramID int64, outlineID string) {
	users[hash(telegramID)] = user{outlineID, false}
	writeDB()
}

func userExists(telegramID int64) bool {
	_, ok := users[hash(telegramID)]
	return ok
}

func SetAdmin(telegramID int64, admin bool) error {
	readDB()
	crypted := hash(telegramID)
	if usr, ok := users[crypted]; ok {
		usr.Admin = admin
		users[crypted] = usr
		writeDB()
		return nil
	}
	return fmt.Errorf("no user with id: %v", telegramID)
}

func getOutlineID(telegramID int64, username string) (string, error) {
	if !userExists(telegramID) {
		return "", fmt.Errorf("user @%v is not registered", username)
	}
	return users[hash(telegramID)].OutlineID, nil
}
