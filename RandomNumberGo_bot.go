package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	modesMap map[int]*Mode
	help     = "help"
	mods     Modes
)

type Config struct {
	TelegramBotToken string
}

type Modes struct {
	Mode []Mode `xml:"Mode"`
}

type Mode struct {
	Name     string   `xml:"Name"`
	Min      int      `xml:"Min"`
	Max      int      `xml:"Max"`
	Quantity int      `xml:"Quantity"`
	Values   []string `xml:"Values"`
}

func main() {
	modesMap = make(map[int]*Mode)
	mods, _ = parseModes()

	config, err := parceConfig()
	if err != nil {
		log.Panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			userID := update.CallbackQuery.From.ID
			chatID := update.CallbackQuery.Message.Chat.ID
			query := update.CallbackQuery.Data
			log.Println(strconv.Itoa(int(chatID)) + ":" + strconv.Itoa(userID) + ":" + query)
			switch query {
			case "roll":
				roll(userID, chatID, query, bot)
			default:
				for _, mode := range mods.Mode {
					if query == mode.Name {
						tmp := mode
						modesMap[userID] = &tmp
						break
					}
				}
				bot.Send(tgbotapi.NewMessage(chatID, "Ok, I remember"))
				status(userID, chatID, query, bot)
			}
		}
		if update.Message != nil {
			userID := update.Message.From.ID
			chatID := update.Message.Chat.ID
			query := update.Message.Text
			log.Println(strconv.Itoa(int(chatID)) + ":" + strconv.Itoa(userID) + ":" + query)
			var command = ""
			command = update.Message.Command()
			if command == "" {
				//for text
				bot.Send(tgbotapi.NewMessage(chatID, help))
			} else {
				//for commands
				switch command {
				case "setmode":
					msg := tgbotapi.NewMessage(chatID, "Select mode")
					buttons := tgbotapi.InlineKeyboardMarkup{}
					for _, mode := range mods.Mode {
						var row []tgbotapi.InlineKeyboardButton
						btn := tgbotapi.NewInlineKeyboardButtonData(mode.Name, mode.Name)
						row = append(row, btn)
						buttons.InlineKeyboard = append(buttons.InlineKeyboard, row)
					}
					status(userID, chatID, query, bot)
					msg.ReplyMarkup = buttons
					bot.Send(msg)
				case "setquantity":
					if _, ok := modesMap[userID]; ok {
						commands := strings.Split(query, " ")
						if len(commands) > 1 {
							modesMap[userID].Quantity, err = strconv.Atoi(commands[1])
							if err != nil {
								msg := tgbotapi.NewMessage(chatID, "wrong format")
								bot.Send(msg)
							}
							status(userID, chatID, query, bot)
						} else {
							msg := tgbotapi.NewMessage(chatID, "wrong format")
							bot.Send(msg)
						}
					} else {
						msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
						bot.Send(msg)
					}
				case "setmin":
					if _, ok := modesMap[userID]; ok {
						commands := strings.Split(query, " ")
						if len(commands) > 1 {
							modesMap[userID].Min, err = strconv.Atoi(commands[1])
							if err != nil {
								msg := tgbotapi.NewMessage(chatID, "wrong format")
								bot.Send(msg)
							}
							status(userID, chatID, query, bot)
						} else {
							msg := tgbotapi.NewMessage(chatID, "wrong format")
							bot.Send(msg)
						}
					} else {
						msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
						bot.Send(msg)
					}
				case "setmax":
					if _, ok := modesMap[userID]; ok {
						commands := strings.Split(query, " ")
						if len(commands) > 1 {
							modesMap[userID].Max, err = strconv.Atoi(commands[1])
							if err != nil {
								msg := tgbotapi.NewMessage(chatID, "wrong format")
								bot.Send(msg)
							}
							status(userID, chatID, query, bot)
						} else {
							msg := tgbotapi.NewMessage(chatID, "wrong format")
							bot.Send(msg)
						}
					} else {
						msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
						bot.Send(msg)
					}
				case "setwords":
					if _, ok := modesMap[userID]; ok {
						commands := strings.Split(query, "/setwords ")
						log.Println("==================================================" + commands[1])
						words := strings.Split(commands[1], ";")
						for _, word := range words {
							modesMap[userID].Values = append(modesMap[userID].Values, word)
						}
						status(userID, chatID, query, bot)
					} else {
						msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
						bot.Send(msg)
					}
				case "setminmaxqua":
					if _, ok := modesMap[userID]; ok {
						commands := strings.Split(query, " ")
						if len(commands) > 3 {
							modesMap[userID].Min, err = strconv.Atoi(commands[1])
							modesMap[userID].Max, err = strconv.Atoi(commands[2])
							modesMap[userID].Quantity, err = strconv.Atoi(commands[3])
							if err != nil {
								msg := tgbotapi.NewMessage(chatID, "wrong format")
								bot.Send(msg)
							}
							status(userID, chatID, query, bot)
						} else {
							msg := tgbotapi.NewMessage(chatID, "wrong format")
							bot.Send(msg)
						}
					} else {
						msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
						bot.Send(msg)
					}
				case "status":
					status(userID, chatID, query, bot)
				case "roll":
					roll(userID, chatID, query, bot)
				default:
					msg := tgbotapi.NewMessage(chatID, help)
					bot.Send(msg)
				}
			}
		}
	}

}

func status(userID int, chatID int64, query string, bot *tgbotapi.BotAPI) {
	if val, ok := modesMap[userID]; ok {
		msg := tgbotapi.NewMessage(chatID, "Your \"random\":\nMode: "+val.Name+"\nQuantity: "+strconv.Itoa(int(val.Quantity))+"\nMin: "+
			strconv.Itoa(val.Min)+"\nMax: "+strconv.Itoa(val.Max)+"\nWords: "+strings.Join(val.Values, ";")+
			"\n\n/help for help\n/roll for roll")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Mode not found. Please use /setmode before using "+query)
		bot.Send(msg)
	}
}

func roll(userID int, chatID int64, query string, bot *tgbotapi.BotAPI) {
	if val, ok := modesMap[userID]; ok {
		if val.Max > val.Min || len(val.Values) > 0 {
			sum := 0
			msgText := "\n"
			for i := 0; i < val.Quantity; i++ {
				if len(val.Values) == 0 {
					tmp := rand.Intn(val.Max+1-val.Min) + val.Min
					sum += tmp
					msgText += strconv.Itoa(tmp) + "\n"
				} else {
					tmp := rand.Intn(len(val.Values))
					msgText += val.Values[tmp] + "\n"
				}
			}
			msgText += "sum= " + strconv.Itoa(sum)
			msg := tgbotapi.NewMessage(chatID, msgText)
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "please use /setmin or /setmax to change numbers, because your max number less than min, before using "+query)
			bot.Send(msg)
		}
	} else {
		msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
		bot.Send(msg)
	}
}

func randomNumber(num1 string, num2 string) (string, error) {
	max, err := strconv.Atoi(num2)
	min, err := strconv.Atoi(num1)
	msg := strconv.Itoa(rand.Intn(max-min) + min)
	return msg, err
}

func parseModes() (Modes, error) {
	file, err := os.Open("resources.xml")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	var tmpModes Modes
	xml.Unmarshal(byteValue, &tmpModes)
	return tmpModes, err
}

func parceConfig() (Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Panic(err)
	}
	return config, err
}
