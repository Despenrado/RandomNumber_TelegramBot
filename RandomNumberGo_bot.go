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
	modesMap  map[int64]*Solution
	help      = "help"
	solutions Solutions
)

type Solutions struct {
	Solution []Solution `xml:"Solution"`
}

type Solution struct {
	Name     string   `xml:"Name"`
	Min      int      `xml:"Min"`
	Max      int      `xml:"Max"`
	Quantity int      `xml:"Quantity"`
	Values   []string `xml:"Values"`
}

func main() {
	modesMap = make(map[int64]*Solution)
	solutions, _ = parseSolutions()

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
			callBackText := update.CallbackQuery.Data
			switch callBackText {
			case "roll":
				if val, ok := modesMap[int64(update.CallbackQuery.From.ID)]; ok {
					sum := 0
					msgText := "\n"
					for i := 0; i < int(val.Quantity); i++ {
						if len(val.Values) == 0 {
							tmp := rand.Intn(val.Max-val.Min) + val.Min
							sum += tmp
							msgText += strconv.Itoa(tmp) + "\n"
						} else {
							tmp := rand.Intn(val.Max)
							msgText += val.Values[tmp] + "\n"
						}
					}
					msgText += "sum= " + strconv.Itoa(sum)
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msgText)
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "please use /setmode before using "+callBackText)
					bot.Send(msg)
				}
			default:
				for _, mode := range solutions.Solution {
					if update.CallbackQuery.Data == mode.Name {
						tmp := mode
						modesMap[int64(update.CallbackQuery.From.ID)] = &tmp
						break
					}
				}
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ok, I remember"))
			}
		}
		if update.Message != nil {
			userName := update.Message.From.UserName
			chatID := update.Message.Chat.ID
			query := update.Message.Text
			log.Println(strconv.Itoa(int(chatID)) + ":" + userName + ":" + query)
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
					for _, mode := range solutions.Solution {
						var row []tgbotapi.InlineKeyboardButton
						btn := tgbotapi.NewInlineKeyboardButtonData(mode.Name, mode.Name)
						row = append(row, btn)
						buttons.InlineKeyboard = append(buttons.InlineKeyboard, row)
					}
					msg.ReplyMarkup = buttons
					bot.Send(msg)
				case "setquantity":
					if _, ok := modesMap[chatID]; ok {
						commands := strings.Split(query, " ")
						modesMap[chatID].Quantity, err = strconv.Atoi(commands[1])
						if err != nil {
							msg := tgbotapi.NewMessage(chatID, "wrong format")
							bot.Send(msg)
						}
					} else {
						msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
						bot.Send(msg)
					}
				case "setmin":
					if _, ok := modesMap[chatID]; ok {
						commands := strings.Split(query, " ")
						modesMap[chatID].Min, err = strconv.Atoi(commands[1])
						if err != nil {
							msg := tgbotapi.NewMessage(chatID, "wrong format")
							bot.Send(msg)
						}
					} else {
						msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
						bot.Send(msg)
					}
				case "setmax":
					if _, ok := modesMap[chatID]; ok {
						commands := strings.Split(query, " ")
						modesMap[chatID].Max, err = strconv.Atoi(commands[1])
						if err != nil {
							msg := tgbotapi.NewMessage(chatID, "wrong format")
							bot.Send(msg)
						}
					} else {
						msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
						bot.Send(msg)
					}
				case "setwords":
					if _, ok := modesMap[chatID]; ok {
						commands := strings.Split(query, "/setwords ")
						words := strings.Split(commands[0], ";")
						for _, word := range words {
							modesMap[chatID].Values = append(modesMap[chatID].Values, word)
						}
					} else {
						msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
						bot.Send(msg)
					}
				case "setminmaxqua":
					if _, ok := modesMap[chatID]; ok {
						commands := strings.Split(query, " ")
						modesMap[chatID].Min, err = strconv.Atoi(commands[1])
						modesMap[chatID].Max, err = strconv.Atoi(commands[2])
						modesMap[chatID].Quantity, err = strconv.Atoi(commands[3])
						if err != nil {
							msg := tgbotapi.NewMessage(chatID, "wrong format")
							bot.Send(msg)
						}
					} else {
						msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
						bot.Send(msg)
					}
				case "roll":
					if val, ok := modesMap[chatID]; ok {
						sum := 0
						msgText := "\n"
						for i := 0; i < val.Quantity; i++ {
							if len(val.Values) == 0 {
								tmp := rand.Intn(val.Max-val.Min) + val.Min
								sum += tmp
								msgText += strconv.Itoa(tmp) + "\n"
							} else {
								tmp := rand.Intn(val.Max)
								msgText += val.Values[tmp] + "\n"
							}
						}
						msgText += "sum= " + strconv.Itoa(sum)
						msg := tgbotapi.NewMessage(chatID, msgText)
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(chatID, "please use /setmode before using "+query)
						bot.Send(msg)
					}
				default:
					msg := tgbotapi.NewMessage(chatID, help)
					bot.Send(msg)
				}
			}

		}
		if val, ok := modesMap[chatID]; ok {
			msg := tgbotapi.NewMessage(chatID, "Your \"random\":\nMode: "+val.Name+"\nQuantity: "+strconv.Itoa(int(val.Quantity))+"\nMin: "+strconv.Itoa(val.Min)+"\nMax: "+strconv.Itoa(val.Max)+"\nWords: "+strings.Join(val.Values, ";")+"\n\n/roll for roll")
			bot.Send(msg)
		}
	}

}

func randomNumber(num1 string, num2 string) (string, error) {
	max, err := strconv.Atoi(num2)
	min, err := strconv.Atoi(num1)
	msg := strconv.Itoa(rand.Intn(max-min) + min)
	return msg, err
}

type Config struct {
	TelegramBotToken string
}

func parseSolutions() (Solutions, error) {
	file, err := os.Open("resources.xml")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	var tmpSolutions Solutions
	xml.Unmarshal(byteValue, &tmpSolutions)
	return tmpSolutions, err
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

// args := strings.Split(query, " ")
// 			if len(args) > 1 {
// 				args = append(args[:0], args[1:]...)
// 			}
// 			switch command {
// 			case "help":
// 				msg = "*there will be help message*"
// 			case "start":
// 				msg = "_there will be start message_"
// 			case "rand":
// 				msg, err = randomNumber(args[0], args[1])
// 				if err != nil {
// 					log.Println(err)
// 				}
// 			default:
// 				msg = "Error"
// 			}
