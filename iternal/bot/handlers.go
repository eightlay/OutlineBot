package bot

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func addHandlers(b *tele.Bot) {
	b.Handle("/start", registerUser)
	b.Handle(&btnRegisterAdmin, registerAdmin)
	b.Handle(&btnUnregisterAdmin, unregisterAdmin)
	b.Handle(&btnUsage, getUsage)
	b.Handle(&btnAccess, getAccess)
	b.Handle(&btnPayment, getPaymentChoices)
	b.Handle(&btnTon, getPaymentTon)
	b.Handle(&btnCard, getPaymentCard)
	b.Handle(&btnHelp, getHelp)
	b.Handle(&btnHelpAccept, helpAccepted)
	b.Handle(&btnHelpDecline, helpDeclined)
	b.Handle(tele.OnText, answer)
	b.Handle(tele.OnPhoto, answer)
	b.Handle(tele.OnAudio, answer)
}

func registerUser(c tele.Context) error {
	userID := c.Sender().ID

	if userExists(userID) {
		c.Send("Вы уже зарегистрированы", getKeyboard(c, "menu"))
		return nil
	}

	username := c.Sender().Username

	outlineID, message, err := newAccessKey(username)
	if err != nil {
		return err
	}

	addUser(userID, outlineID)

	return c.Send(message, getKeyboard(c, "menu"))
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

	return c.Send(fmt.Sprintf("Потребление за последние 30 дней:\n%v", usage), getKeyboard(c, "menu"))
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
	urlMenu.Inline(urlMenu.Row(btnURL))

	message := "<b>Рекомендуем удалить ссылку после подключения!</b>\n"
	return c.Send(message, urlMenu)
}

func getPaymentChoices(c tele.Context) error {
	setState(c.Sender().ID, statePAY)
	return c.Send("Выберите способ оплаты", getKeyboard(c, "payments"))
}

func getPaymentTon(c tele.Context) error {
	if !isState(c.Sender().ID, statePAY) {
		c.Send("Перейдите в меню 'Платежи'", getKeyboard(c, "menu"))
		return nil
	}

	price, err := getPrice("TON")
	if err != nil {
		return err
	}

	url, payUrl, com, err := getTonURLs(c.Sender().ID, c.Sender().Username, price)
	if err != nil {
		return err
	}

	btnURL := urlMenu.URL("Реквизиты", url)
	btnPayURL := urlMenu.URL("Оплатить", payUrl)
	urlMenu.Inline(urlMenu.Row(btnURL), urlMenu.Row(btnPayURL))

	if err := c.Send(fmt.Sprintf("К оплате %v💎", price), urlMenu); err != nil {
		return err
	}

	msg := fmt.Sprintf("<b>При оплате по реквизитам добавьте в комментарий код оплаты:</b>\n%v", com)
	err = c.Send(msg, getKeyboard(c, "menu"))
	setState(c.Sender().ID, stateIDLE)
	return err
}

func getPaymentCard(c tele.Context) error {
	if !isState(c.Sender().ID, statePAY) {
		c.Send("Перейдите в меню 'Платежи'", getKeyboard(c, "menu"))
		return nil
	}

	price, err := getPrice("CARD")
	if err != nil {
		return err
	}

	url, com, err := getCardURL(c.Sender().ID, c.Sender().Username, price)
	if err != nil {
		return err
	}

	if err := c.Send(url); err != nil {
		return err
	}

	err = c.Send(fmt.Sprintf("<b>Добавьте в комментарий код оплаты:</b>\n%v", com), getKeyboard(c, "menu"))
	setState(c.Sender().ID, stateIDLE)
	return err
}

func getHelp(c tele.Context) error {
	msg := "Запрашивая обращение, Вы соглашаетесь на обработку персональных данных:\n"
	msg += "Ваш ID будет сохранен в базе обращений до получения ответа"
	return c.Send(msg, getKeyboard(c, "helpMenu"))
}

func helpAccepted(c tele.Context) error {
	setState(c.Sender().ID, stateHELP)
	return c.Send("Сформулируйте своё обращение одним сообщением")
}

func helpSend(c tele.Context) error {
	setState(c.Sender().ID, stateIDLE)
	if err := registerHelpRequest(c); err != nil {
		c.Send("Ошибка при регистрации обращения", getKeyboard(c, "menu"))
		return err
	}
	return c.Send("Ваше обращение будет рассмотрено в ближайшее время", getKeyboard(c, "menu"))
}

func helpDeclined(c tele.Context) error {
	return c.Send("Вы также можете обратиться за помощью по контакту в описании", getKeyboard(c, "menu"))
}

func answer(c tele.Context) error {
	if isState(c.Sender().ID, stateHELP) {
		return helpSend(c)
	}
	if isAdmin(c.Sender().ID) {
		if answerRequest(c) {
			msg := "Пожалуйста, выберите опцию из меню или ответьте на открытое обращение"
			return c.Send(msg, getKeyboard(c, "menu"))
		}
		return c.Send("Ответ на обращение отправлен", getKeyboard(c, "menu"))
	}
	return c.Send("Пожалуйста, выберите опцию из меню", getKeyboard(c, "menu"))
}
