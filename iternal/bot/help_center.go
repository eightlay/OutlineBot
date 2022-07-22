package bot

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

var admins []int64
var openRequests map[int64][]helpRequest = map[int64][]helpRequest{}

type helpRequest struct {
	storedUser  *tele.Message
	storedAdmin *tele.Message
}

func answerRequest(c tele.Context) bool {
	if !c.Message().IsReply() && isState(c.Sender().ID, stateIDLE) {
		return false
	}

	rplID := c.Message().ReplyTo.ID

	for i, r := range openRequests[c.Sender().ID] {
		if rplID == r.storedAdmin.ID {
			c.Bot().Reply(r.storedUser, c.Message().Text)
			c.Bot().Delete(r.storedAdmin)
			openRequests[c.Sender().ID] = append(
				openRequests[c.Sender().ID][:i], openRequests[c.Sender().ID][i+1:]...,
			)
		}
	}
	return true
}

func registerHelpRequest(c tele.Context) error {
	for _, a := range admins {
		chat, err := c.Bot().ChatByID(a)
		if err != nil {
			return fmt.Errorf("can't find chat with admin %v: %v", a, err)
		}

		msg, err := c.Bot().Forward(chat, c.Message())
		if err != nil {
			return fmt.Errorf("can't forward message to admin %v: %v", a, err)
		}

		openRequests[a] = append(
			openRequests[a], helpRequest{c.Message(), msg},
		)
	}
	return nil
}

func registerAdmin(c tele.Context) error {
	telegramID := c.Sender().ID
	if isAdmin(telegramID) {
		for _, a := range admins {
			if a == telegramID {
				return nil
			}
		}
		admins = append(admins, telegramID)
		openRequests[telegramID] = []helpRequest{}
		c.Send("Авторизация прошла успешно")
	} else {
		c.Send("Вы не являетесь администратором")
	}
	return nil
}

func unregisterAdmin(c tele.Context) error {
	for i, a := range admins {
		if a == c.Sender().ID {
			admins = append(admins[:i], admins[i+1:]...)
			c.Send("Вы больше не будете получать обращения")
			break
		}
	}
	return nil
}
