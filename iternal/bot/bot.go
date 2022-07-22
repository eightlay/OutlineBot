package bot

import (
	"fmt"
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

func StartBot() {
	readDB()
	initKeyboards()

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

	addHandlers(b)

	b.Start()
}
