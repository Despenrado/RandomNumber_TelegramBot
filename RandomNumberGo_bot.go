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
	templatesMap map[int]*Template
	help         = "help message:\n" +
		"structure of queries: /command [parametr] [parametr] ....\n" +
		"example: /setmin 10 - this command sets minimum border to 10\n\n" +
		"commands:\n" +
		"help - shows this message\n" +
		"settemplate - shows list of templates\n" +
		"setmin [number]- sets minimum border\n" +
		"setmax [number] - sets maximum border\n" +
		"setquantity [number] - sets number of random numbers\n" +
		"setminmaxqua [min] [max] [quantity] - sets minimum, maximum and quantity\n" +
		"setwords [word1;word2;word3] - sets words for random choice\n" +
		"status - shows your current tamplate of rundom\n" +
		"random\n" +
		"roll - genertes random nuber/numbers via current template"
	templates Templates
)

type Config struct {
	TelegramBotToken string
}

type Templates struct {
	Template []Template `xml:"Template"`
}

type Template struct {
	Name      string   `xml:"Name"`
	Min       int      `xml:"Min"`
	Max       int      `xml:"Max"`
	Quantity  int      `xml:"Quantity"`
	Values    []string `xml:"Values"`
	ImagePath []string `xml:"ImagePath"`
}

func main() {
	templatesMap = make(map[int]*Template)
	templates, _ = parseTemplates()

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
				for _, template := range templates.Template {
					if query == template.Name {
						tmp := template
						templatesMap[userID] = &tmp
						break
					}
				}
				sendMessage(bot, chatID, "Ok, I remember")
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
				sendMessage(bot, chatID, "/help shows help message")
			} else {
				//for commands
				switch command {
				case "start":
					sendMessage(bot, chatID, "Welcome to RandomNumber_bot\n"+
						"/help shows help message")
				case "settemplate":
					msg := tgbotapi.NewMessage(chatID, "Select template")
					buttons := tgbotapi.InlineKeyboardMarkup{}
					for _, template := range templates.Template {
						var row []tgbotapi.InlineKeyboardButton
						btn := tgbotapi.NewInlineKeyboardButtonData(template.Name, template.Name)
						row = append(row, btn)
						buttons.InlineKeyboard = append(buttons.InlineKeyboard, row)
					}
					msg.ReplyMarkup = buttons
					bot.Send(msg)
				case "setquantity":
					if _, ok := templatesMap[userID]; ok {
						commands := strings.Split(query, " ")
						if len(commands) > 1 {
							templatesMap[userID].Quantity, err = strconv.Atoi(commands[1])
							if err != nil {
								sendMessage(bot, chatID, "wrong format")
							}
							status(userID, chatID, query, bot)
						} else {
							sendMessage(bot, chatID, "wrong format")
						}
					} else {
						sendMessage(bot, chatID, "please use /settemplate before using "+query)
					}
				case "setmin":
					if _, ok := templatesMap[userID]; ok {
						commands := strings.Split(query, " ")
						if len(commands) > 1 {
							templatesMap[userID].Min, err = strconv.Atoi(commands[1])
							if err != nil {
								sendMessage(bot, chatID, "wrong format")
							}
							status(userID, chatID, query, bot)
						} else {
							sendMessage(bot, chatID, "wrong format")

						}
					} else {
						sendMessage(bot, chatID, "please use /settemplate before using "+query)
					}
				case "setmax":
					if _, ok := templatesMap[userID]; ok {
						commands := strings.Split(query, " ")
						if len(commands) > 1 {
							templatesMap[userID].Max, err = strconv.Atoi(commands[1])
							if err != nil {
								sendMessage(bot, chatID, "wrong format")
							}
							status(userID, chatID, query, bot)
						} else {
							sendMessage(bot, chatID, "wrong format")
						}
					} else {
						sendMessage(bot, chatID, "please use /settemplate before using "+query)
					}
				case "setwords":
					if _, ok := templatesMap[userID]; ok {
						commands := strings.Split(query, "/setwords ")
						words := strings.Split(commands[1], ";")
						for _, word := range words {
							templatesMap[userID].Values = append(templatesMap[userID].Values, word)
						}
						status(userID, chatID, query, bot)
					} else {
						sendMessage(bot, chatID, "please use /settemplate before using "+query)
					}
				case "setminmaxqua":
					if _, ok := templatesMap[userID]; ok {
						commands := strings.Split(query, " ")
						if len(commands) > 3 {
							templatesMap[userID].Min, err = strconv.Atoi(commands[1])
							templatesMap[userID].Max, err = strconv.Atoi(commands[2])
							templatesMap[userID].Quantity, err = strconv.Atoi(commands[3])
							if err != nil {
								sendMessage(bot, chatID, "wrong format")
							}
							status(userID, chatID, query, bot)
						} else {
							sendMessage(bot, chatID, "wrong format")
						}
					} else {
						sendMessage(bot, chatID, "please use /settemplate before using "+query)
					}
				case "random":
					commands := strings.Split(query, " ")
					if len(commands) > 3 {
						min, err := strconv.Atoi(commands[1])
						max, err := strconv.Atoi(commands[2])
						if err != nil {
							sendMessage(bot, chatID, "wrong format")
						}
						sendMessage(bot, chatID, strconv.Itoa(rand.Intn(max+1-min)+min))
					} else {
						sendMessage(bot, chatID, "wrong format")
					}
				case "status":
					status(userID, chatID, query, bot)
				case "roll":
					roll(userID, chatID, query, bot)
				default:
					sendMessage(bot, chatID, help)
				}
			}
		}
	}

}

func status(userID int, chatID int64, query string, bot *tgbotapi.BotAPI) {
	if val, ok := templatesMap[userID]; ok {
		sendMessage(bot, chatID, "Your \"random\":\nTemplate: "+val.Name+"\nQuantity: "+strconv.Itoa(int(val.Quantity))+"\nMin: "+
			strconv.Itoa(val.Min)+"\nMax: "+strconv.Itoa(val.Max)+"\nWords: "+strings.Join(val.Values, ";")+
			"\n\n/help for help\n/roll for roll")
	} else {
		sendMessage(bot, chatID, "Template not found. Please use /settemplate before using "+query)
	}
}

func roll(userID int, chatID int64, query string, bot *tgbotapi.BotAPI) {
	if val, ok := templatesMap[userID]; ok {
		if val.Max > val.Min || len(val.Values) > 0 {
			sum := 0
			msgText := "\n"
			for i := 0; i < val.Quantity; i++ {
				if len(val.Values) == 0 {
					tmp := rand.Intn(val.Max+1-val.Min) + val.Min
					sum += tmp
					if val.Min == 1 && val.Quantity < 3 {
						switch val.Max {
						case 4:
							if len(val.ImagePath) > tmp && val.ImagePath[tmp] != "" {
								sendImage(bot, chatID, val.ImagePath[tmp])
							}
						case 6:
							if len(val.ImagePath) > tmp && val.ImagePath[tmp] != "" {
								sendImage(bot, chatID, val.ImagePath[tmp])
							}
						case 8:
							if len(val.ImagePath) > tmp && val.ImagePath[tmp] != "" {
								sendImage(bot, chatID, val.ImagePath[tmp])
							}
						case 10:
							if len(val.ImagePath) > tmp && val.ImagePath[tmp] != "" {
								sendImage(bot, chatID, val.ImagePath[tmp])
							}
						case 12:
							if len(val.ImagePath) > tmp && val.ImagePath[tmp] != "" {
								sendImage(bot, chatID, val.ImagePath[tmp])
							}
						case 20:
							if len(val.ImagePath) > tmp && val.ImagePath[tmp] != "" {
								sendImage(bot, chatID, val.ImagePath[tmp])
							}
						case 100:

						}
					} else {
						msgText += strconv.Itoa(tmp) + "\n"
					}
				} else {
					tmp := rand.Intn(len(val.Values))
					msgText += val.Values[tmp] + "\n"
				}
			}
			msgText += "sum= " + strconv.Itoa(sum) + "\navg= " + strconv.FormatFloat(float64(sum)/float64(val.Quantity), 'f', -4, 32) + "\n/roll again"
			sendMessage(bot, chatID, msgText)

		} else {
			sendMessage(bot, chatID, "please use /setmin or /setmax to change numbers, because your max number less than min, before using "+query)
		}
	} else {
		sendMessage(bot, chatID, "please use /settemplate before using "+query)
	}
}

func parseTemplates() (Templates, error) {
	file, err := os.Open("resources.xml")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	var tmpTemplates Templates
	err = xml.Unmarshal(byteValue, &tmpTemplates)
	if err != nil {
		log.Println(err)
	}
	log.Println("length of templates slice=" + strconv.Itoa(len(tmpTemplates.Template)))
	return tmpTemplates, err
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

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

func sendImage(bot *tgbotapi.BotAPI, chatID int64, imgPath string) {
	msg := tgbotapi.NewPhotoUpload(chatID, imgPath) //NewInputMediaPhoto(media string)
	bot.Send(msg)
}
