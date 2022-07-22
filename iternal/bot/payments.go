package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
)

type priceResponse struct {
	Price float64 `json:"price"`
}

type currencyResponse struct {
	Rates map[string]interface{} `json:"rates"`
}

func getTonURLs(telegramID int64, username string, price float64) (string, string, string, error) {
	pc, err := getUserPaymentCode(telegramID, username)
	if err != nil {
		return "", "", "", err
	}
	url := os.Getenv("TON")
	payUrl := os.Getenv("TON") + "💎" + strconv.FormatFloat(price, 'f', -1, 64) + "/" + pc
	return url, payUrl, pc, nil
}

func getCardURL(telegramID int64, username string, price float64) (string, string, error) {
	pc, err := getUserPaymentCode(telegramID, username)
	if err != nil {
		return "", "", err
	}
	url := fmt.Sprintf("Номер карты: %v\nК оплате %v₽", os.Getenv("CARD"), math.Ceil(price))
	return url, pc, nil
}

func getPrice(platform string) (float64, error) {
	switch platform {
	case "TON":
		tonusd, err := getRateTONUSD()
		if err != nil {
			return 0, err
		}
		return 1. / tonusd, nil
	case "CARD":
		rubusd, err := getRateRUBUSD()
		if err != nil {
			return 0, err
		}
		return 1/rubusd + 1, nil
	}
	return 0, fmt.Errorf("no such platform '%v'", platform)
}

func getRateTONUSD() (float64, error) {
	resp, err := http.Get(getURL("getPriceTON"))
	if err != nil {
		return 0, fmt.Errorf("can't get data from tonnames: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("can't parse response from tonnames %v", err)
	}

	var price priceResponse
	if err := json.Unmarshal(bodyBytes, &price); err != nil {
		return 0, fmt.Errorf("can't unmarshal ton price response %v", err)
	}

	return price.Price, nil
}

func getRateRUBUSD() (float64, error) {
	respRub, err := http.Get(getURL("getPriceRUBUSD"))
	if err != nil {
		return 0, fmt.Errorf("can't get data from cbr: %v", err)
	}
	defer respRub.Body.Close()

	bodyBytes, err := ioutil.ReadAll(respRub.Body)
	if err != nil {
		return 0, fmt.Errorf("can't parse data from cbr: %v", err)
	}

	var objmap map[string]json.RawMessage
	if err := json.Unmarshal(bodyBytes, &objmap); err != nil {
		return 0, fmt.Errorf("can't unmarshal rub price: %v", err)
	}

	var pr currencyResponse
	if err := json.Unmarshal(objmap["rates"], &pr); err != nil {
		return 0, fmt.Errorf("can't unmarshal rub price: %v", err)
	}
	if err := json.Unmarshal(objmap["rates"], &pr.Rates); err != nil {
		panic(err)
	}

	v, found := pr.Rates["USD"]

	rubusd, ok := v.(float64)
	if !ok && found {
		return 0, fmt.Errorf("can't convert bytes to float32 for RUBUSD")
	}

	return rubusd, nil
}
