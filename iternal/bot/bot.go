package bot

import (
	"fmt"
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

var (
	menu       = &tele.ReplyMarkup{ResizeKeyboard: true}
	btnUsage   = menu.Text("Потребление")
	btnAccess  = menu.Text("Подключение")
	btnPayment = menu.Text("Платежи")
	btnHelp    = menu.Text("Поддержка")

	urlMenu = &tele.ReplyMarkup{}
)

func registerUser(c tele.Context) error {
	userID := c.Sender().ID

	if userExists(userID) {
		c.Send("Вы уже зарегистрированы", menu)
		return nil
	}

	username := c.Sender().Username

	outlineID, message, err := newAccessKey(username)
	if err != nil {
		return err
	}

	addUser(userID, outlineID)

	return c.Send(message, menu)
}

func getUsage(c tele.Context) error {
	outlineID, err := getOutlineID(c.Sender().ID, c.Sender().Username)
	if err != nil {
		return err
	}

	usage, err := getDataTransfered(outlineID)
	if err != nil {
		return err
	}

	return c.Send(fmt.Sprintf("Потребление за последние 30 дней:\n%v", usage), menu)
}

func getAccess(c tele.Context) error {
	outlineID, err := getOutlineID(c.Sender().ID, c.Sender().Username)
	if err != nil {
		return err
	}

	url, err := getAccessURL(outlineID)
	if err != nil {
		return err
	}

	btnURL := urlMenu.URL("Подключиться", url)

	urlMenu.Inline(
		urlMenu.Row(btnURL),
	)

	message := "<b>Рекомендуем удалить ссылку после подключения!</b>\n"
	return c.Send(message, urlMenu)
}

func StartBot() {
	readDB()

	tg, ok := os.LookupEnv("TGTOKEN")
	if !ok {
		panic(fmt.Errorf("can't find telegram token in env vars"))
	}

	pref := tele.Settings{
		Token:     tg,
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeHTML,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	menu.Reply(
		menu.Row(btnUsage, btnAccess),
		menu.Row(btnPayment, btnHelp),
	)

	b.Handle("/start", registerUser)
	b.Handle(&btnUsage, getUsage)
	b.Handle(&btnAccess, getAccess)

	b.Start()
}
