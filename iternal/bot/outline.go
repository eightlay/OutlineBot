package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var oulineHost = os.Getenv("OUTLINEAPI")

type createResponse struct {
	OutlineID string `json:"id"`
}

type transferedResponse struct {
	Transfered map[string]interface{} `json:"-"`
}

type Key struct {
	OutlineID string `json:"id"`
	AccessURL string `json:"accessUrl"`
}

func getURL(name string) string {
	switch name {
	case "newAccessKey":
		return oulineHost + "access-keys/"
	case "getDataTransfered":
		return oulineHost + "metrics/transfer/"
	case "getAccessURL":
		return oulineHost + "access-keys/"
	}
	return oulineHost
}

func newAccessKey(username string) (string, string, error) {
	resp, err := http.Post(
		getURL("newAccessKey"), "application/json",
		bytes.NewReader([]byte("")),
	)
	if err != nil {
		err := fmt.Errorf("can't register user @%v: %v", username, err)
		return "", "Ошибка регистрации. Попробуйте использовать команду /start позже", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("can't parse Outline create response for user @%v: %v", username, err)
		return "", "Ошибка регистрации. Попробуйте использовать команду /start позже", err
	}

	var cr createResponse
	json.Unmarshal(bodyBytes, &cr)
	return cr.OutlineID, "Регистрация прошла успешно", nil
}

func getDataTransfered(outlineID string) (string, error) {
	logError := "can't parse data transfered: %v"
	errorMessage := "Данные временно недоступны. Попробуйте позже"

	resp, err := http.Get(getURL("getDataTransfered"))
	if err != nil {
		return errorMessage, fmt.Errorf("can't get data transfered: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errorMessage, fmt.Errorf(logError, err)
	}

	var objmap map[string]json.RawMessage
	if err := json.Unmarshal(bodyBytes, &objmap); err != nil {
		return errorMessage, fmt.Errorf(logError, err)
	}

	var tr transferedResponse
	if err := json.Unmarshal(objmap["bytesTransferredByUserId"], &tr); err != nil {
		return errorMessage, fmt.Errorf(logError, err)
	}

	if err := json.Unmarshal(objmap["bytesTransferredByUserId"], &tr.Transfered); err != nil {
		panic(err)
	}

	v, found := tr.Transfered[outlineID]

	bts, ok := v.(float64)
	if !ok && found {
		msg := fmt.Sprintf("can't convert bytes to float64 for outline id '%v'", outlineID)
		return errorMessage, fmt.Errorf(logError, msg)
	}

	return formatBytes(bts, 2), nil
}

func getAccessURL(outlineID string) (string, error) {
	logError := "can't parse access keys: %v"
	errorMessage := "Данные временно недоступны. Попробуйте позже"

	resp, err := http.Get(getURL("getAccessURL"))
	if err != nil {
		return errorMessage, fmt.Errorf("can't get access keys: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errorMessage, fmt.Errorf(logError, err)
	}

	var objmap map[string]json.RawMessage
	if err := json.Unmarshal(bodyBytes, &objmap); err != nil {
		return errorMessage, fmt.Errorf(logError, err)
	}

	var keys []Key
	if err := json.Unmarshal(objmap["accessKeys"], &keys); err != nil {
		return errorMessage, fmt.Errorf(logError, err)
	}

	var ind int
	for i, key := range keys {
		if key.OutlineID == outlineID {
			ind = i
			break
		}
	}

	url := fmt.Sprintf("https://s3.amazonaws.com/outline-vpn/invite.html#%v", keys[ind].AccessURL)
	return url, nil
}
