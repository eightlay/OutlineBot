package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akamensky/argparse"
	"github.com/eightlay/outline-telegram-bot/iternal/bot"
)

func main() {
	// Command line arguments
	parser := argparse.NewParser("OutlineBot", "Outline Telegram bot")
	startCmd := parser.NewCommand("run", "start bot")
	setCmd := parser.NewCommand("admin", "give admin rights")
	adminID := setCmd.Int("u", "user", &argparse.Options{Help: "user's Telegram ID"})
	depriveAdmin := setCmd.Flag("d", "deprive", &argparse.Options{Help: "deprive admin rights"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	// Logging
	f, err := os.OpenFile("bot.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	defer f.Close()

	// Logic
	if startCmd.Happened() {
		// Start bot
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		bot.StartBot()
	} else if setCmd.Happened() {
		// Set admin
		err := bot.SetAdmin(int64(*adminID), !*depriveAdmin)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("success")
		}
	} else {
		// Bad arguments
		err := fmt.Errorf("bad arguments, please check usage")
		fmt.Print(parser.Usage(err))
	}
}
