package bot

func getURL(name string) string {
	switch name {
	case "newAccessKey":
		return oulineHost + "access-keys/"
	case "getDataTransfered":
		return oulineHost + "metrics/transfer/"
	case "getAccessURL":
		return oulineHost + "access-keys/"
	case "getPriceRUBUSD":
		return "https://www.cbr-xml-daily.ru/latest.js"
	case "getPriceTON":
		return "https://api.tonnames.org/ton-price"
	}
	return oulineHost
}
