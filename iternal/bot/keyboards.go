package bot

import (
	tele "gopkg.in/telebot.v3"
)

var (
	menu       = &tele.ReplyMarkup{ResizeKeyboard: true}
	btnUsage   = menu.Text("Потребление")
	btnAccess  = menu.Text("Подключение")
	btnPayment = menu.Text("Платежи")
	btnHelp    = menu.Text("Поддержка")

	btnRegisterAdmin   = menu.Text("Получать обращения")
	btnUnregisterAdmin = menu.Text("Перестать получать обращения")

	urlMenu = &tele.ReplyMarkup{}

	paymentsMenu = &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	btnTon       = paymentsMenu.Text("TON")
	btnCard      = paymentsMenu.Text("Карта")

	helpMenu       = &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	btnHelpAccept  = menu.Text("Согласиться")
	btnHelpDecline = menu.Text("Отказаться")
)

func initKeyboards() {
	paymentsMenu.Reply(
		paymentsMenu.Row(btnTon),
		paymentsMenu.Row(btnCard),
	)

	helpMenu.Reply(
		helpMenu.Row(btnHelpAccept),
		helpMenu.Row(btnHelpDecline),
	)
}

func getKeyboard(c tele.Context, keyboard string) *tele.ReplyMarkup {
	switch keyboard {
	case "menu":
		if isAdmin(c.Sender().ID) {
			menu.Reply(
				menu.Row(btnUsage, btnAccess),
				menu.Row(btnPayment, btnHelp),
				menu.Row(btnRegisterAdmin),
				menu.Row(btnUnregisterAdmin),
			)
		} else {
			menu.Reply(
				menu.Row(btnUsage, btnAccess),
				menu.Row(btnPayment, btnHelp),
			)
		}
	case "payments":
		return paymentsMenu
	case "helpMenu":
		return helpMenu
	}

	return menu
}
