package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"gopkg.in/telebot.v3/middleware"

	tele "gopkg.in/telebot.v3"
)

type BotConfig struct {
	Token        string
	AdminId      int64
	GroupId      int64
	YummyId      int64
	Domain       string
	OrderHourEnd int
}

var (
	bot      *tele.Bot
	cfg      *BotConfig
	menuDate = time.Now()
)

func onRects(c tele.Context) error {
	payload := c.Message().Payload
	if payload == "" {
		text := cfg.Domain + "/rects-tool"
		currentRects, err := GetRects()
		if err != nil {
			return c.Send(text)
		}
		jsonRects, err := json.Marshal(currentRects)
		if err != nil {
			return c.Send(text)
		}
		text = text + "\n`" + string(jsonRects) + "`"
		return c.Send(text, &tele.SendOptions{
			ParseMode: tele.ModeMarkdown,
		})
	}
	var rects []FloatRect
	if err := json.Unmarshal([]byte(payload), &rects); err != nil {
		return c.Send("parse error")
	}
	if len(rects) == 0 {
		return c.Send("no rects")
	}
	SaveRects(payload)
	return c.Send("ok")
}

func onPhoto(c tele.Context) error {
	senderId := c.Message().Sender.ID
	if senderId != cfg.YummyId && senderId != cfg.AdminId {
		return nil
	}

	minskHour := GetMinskHour()
	fmt.Println(minskHour)
	if senderId != cfg.AdminId && (minskHour < 9 || minskHour > 12) {
		log.Println("photo not in time")
		return nil
	}

	photo := c.Message().Photo

	reader, err := bot.File(&photo.File)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		return c.Send("error")
	}

	rects, err := GetRects()
	if err != nil {
		return nil
	}

	items := GetTextFromImage(buf.Bytes(), rects)

	today := Today()
	tomorrowDate := today.AddDate(0, 0, 1)
	if menuDate.Before(tomorrowDate) {
		menuDate = tomorrowDate
	}

	itemsJson, _ := json.Marshal(items)
	menu := Menu{
		Items:        string(itemsJson),
		PublishDate:  today,
		DeliveryDate: menuDate,
	}

	SaveMenu(menu)
	authButton := tele.InlineButton{
		Text: "Создать заказ",
		Login: &tele.Login{
			URL:         cfg.Domain,
			WriteAccess: false,
		},
	}

	return c.Send(
		"Нажмите на кнопку ниже, чтобы создать заказ",
		&tele.SendOptions{
			ReplyMarkup: &tele.ReplyMarkup{
				InlineKeyboard: [][]tele.InlineButton{
					{authButton},
				},
			},
		})
}

func onText(c tele.Context) error {
	senderId := c.Message().Sender.ID
	if senderId != cfg.YummyId && senderId != cfg.AdminId {
		return nil
	}

	r := regexp.MustCompile(`.*?[Мм]еню на.*?(\d{1,2}\.\d{1,2}\.\d{2,4}).*`)
	m := r.FindStringSubmatch(c.Message().Text)
	if len(m) != 2 {
		return nil
	}
	dateStr := m[1]
	date, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return nil
	}

	if date.Before(time.Now()) {
		return nil
	}
	menuDate = date
	return nil
}

func PostOrderInChat(order OrderRequest) error {
	group := tele.Chat{
		ID: cfg.GroupId,
	}
	memberOf, err := bot.ChatMemberOf(&group, &tele.User{ID: order.UserId})
	if err != nil {
		return err
	}

	allowedRoles := map[tele.MemberStatus]bool{
		tele.Creator:       true,
		tele.Administrator: true,
		tele.Member:        true,
	}
	if !allowedRoles[memberOf.Role] {
		return errors.New("not allowed")
	}

	minskHour := GetMinskHour()
	if minskHour >= cfg.OrderHourEnd {
		return errors.New("too late")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s](tg://user?id=%d):\n", order.Name, order.UserId))
	for _, item := range order.Items {
		sb.WriteString(fmt.Sprintf("- %s x%d\n", item.Name, item.Quantity))
	}
	_, err = bot.Send(
		&tele.Chat{
			ID: cfg.GroupId,
		},
		sb.String(),
		&tele.SendOptions{
			ParseMode: tele.ModeMarkdown,
		},
	)
	return err
}

func onMenu(c tele.Context) error {
	today := Today()
	if c.Message().Payload == "" {
		menu, err := GetMenu(today)
		if err != nil {
			return c.Send("Menu not found")
		}
		items := make([]string, 0)
		if err := json.Unmarshal([]byte(menu.Items), &items); err != nil {
			return c.Send("Menu parse error")
		}
		text := "`/menu\n" + strings.Join(items, "\n") + "`"
		return c.Send(text, &tele.SendOptions{
			ParseMode: tele.ModeMarkdown,
		})
	}
	lines := strings.Split(c.Message().Text, "\n")[1:]
	linesJson, err := json.Marshal(lines)
	if err != nil {
		return c.Send("Menu marshal error")
	}
	tomorrow := today.AddDate(0, 0, 1)
	if menuDate.Before(tomorrow) {
		menuDate = tomorrow
	}
	SaveMenu(Menu{
		Items:        string(linesJson),
		PublishDate:  today,
		DeliveryDate: menuDate,
	})
	return c.Send(fmt.Sprintf("Save %d items", len(lines)))
}

func RunBot(config BotConfig) {
	pref := tele.Settings{
		Token:  config.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	cfg = &config
	bot = b

	adminOnly := b.Group()
	adminOnly.Use(middleware.Whitelist(config.AdminId))
	adminOnly.Handle("/rects", onRects)
	adminOnly.Handle("/menu", onMenu)

	b.Handle(tele.OnText, onText)
	b.Handle(tele.OnPhoto, onPhoto)

	b.Start()
}
