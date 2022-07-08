package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"yummy-bot/ocr"

	"gopkg.in/telebot.v3/middleware"

	tele "gopkg.in/telebot.v3"
)

type Configuration struct {
	TelegramToken string
	AdminId       int64
	GroupId       int64
	YummyId       int64
	Domain        string
	OrderHourEnd  int
	AbbyyUsername string
	AbbyyPassword string
}

var (
	bot      *tele.Bot
	cfg      *Configuration
	menuDate = time.Now()
)

func onRects(c tele.Context) error {
	payload := c.Message().Payload
	if payload == "" {
		currentRects, err := GetRects()
		if err != nil {
			currentRects = make([]ocr.FloatRect, 0)
		}
		return c.Send(rectsToUri(currentRects), &tele.SendOptions{
			ParseMode: tele.ModeMarkdown,
		})
	}
	rects, err := loadRectsFromUri(payload)
	if err != nil {
		return c.Send("parse error")
	}
	if len(rects) == 0 {
		return c.Send("no rects")
	}
	rectsJson, err := json.Marshal(rects)
	if err != nil {
		return c.Send("rects marshal error")
	}
	SaveRects(string(rectsJson))
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

	rects, err := GetRects()
	if err != nil {
		_, err = bot.Send(&tele.Chat{
			ID: cfg.AdminId,
		}, "ERROR: could not get rects")
		return err
	}

	items, err := ocr.GetTextFromImageAbbyy(reader, rects, cfg.AbbyyUsername, cfg.AbbyyPassword)
	if err != nil {
		log.Println("could not extract lines from image")
		return err
	}
	items = append(items, "хлеб")

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
	secretKey := sha256.Sum256([]byte(cfg.TelegramToken))
	secretKeyHmac := hmac.New(sha256.New, secretKey[:])
	secretKeyHmac.Write([]byte(order.DataCheckString))
	hash := secretKeyHmac.Sum(nil)
	if hex.EncodeToString(hash) != order.Hash {
		log.Printf("auth fail: %v", order)
	}
	userId, _ := strconv.ParseInt(order.UserId, 10, 64)

	minskHour := GetMinskHour()
	if minskHour >= cfg.OrderHourEnd {
		return errors.New("заказы принимаются до " + strconv.Itoa(cfg.OrderHourEnd) + ":00")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s](tg://user?id=%d):\n", order.Name, userId))
	for _, item := range order.Items {
		sb.WriteString(fmt.Sprintf("- %s x%d\n", item.Name, item.Quantity))
	}
	_, err := bot.Send(
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

func RunBot(config Configuration) {
	b, err := tele.NewBot(tele.Settings{
		Token:  config.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	cfg = &config
	bot = b

	adminOnly := b.Group()
	adminOnly.Use(middleware.Whitelist(config.AdminId, config.YummyId))
	adminOnly.Handle("/rects", onRects)
	adminOnly.Handle("/menu", onMenu)

	b.Handle(tele.OnText, onText)
	b.Handle(tele.OnPhoto, onPhoto)

	b.Start()
}