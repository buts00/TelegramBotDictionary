package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + BotToken
	offset := 0

	for {
		allUpdates, err := updates(botUrl, offset)
		if err != nil {
			log.Fatal("some problem", err)
		}
		for _, update := range allUpdates {
			err = respond(botUrl, update)
			offset = update.UpdateID + 1
		}
		fmt.Println(allUpdates)

	}
}

func updates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil

}

func respond(botUrl string, update Update) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = update.Message.Text
	if botMessage.Text == "/words" {
		allWords := AllWords(Database(), botMessage.ChatId)
		if len(allWords) != 0 {
			joinWords := []string{}
			for _, word := range allWords {
				s := word.word + " / " + word.translation
				joinWords = append(joinWords, s)
			}
			botMessage.Text = strings.Join(joinWords, "\n")
		} else {
			botMessage.Text = "You have not words yet("
		}

	} else if botMessage.Text == "/random" {
		botMessage.Text = RandomWord(botMessage.ChatId)
	} else if len(strings.Split(botMessage.Text, " ")) == 3 {
		messageArr := strings.Split(botMessage.Text, " ")
		if messageArr[0] == "/add" {
			InsertWord(botMessage.Text, botMessage.ChatId)
			botMessage.Text = "The new word has been saved"
		} else {
			botMessage.Text = "fuck you"
		}
	} else if botMessage.Text == "/help" {
		allCommandsArr, err := listCommand(botUrl)
		if err != nil {
			return err
		}
		if len(allCommandsArr) == 0 {
			botMessage.Text = "There are no commands yet("
		} else {
			allCommands := []string{}
			for _, command := range allCommandsArr {
				newCommand := "/" + command.Command + " - " + command.Description
				allCommands = append(allCommands, newCommand)
			}
			botMessage.Text = strings.Join(allCommands, "\n")
		}

	} else {
		botMessage.Text = "fuck you"
	}
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

func listCommand(botUrl string) ([]Commands, error) {

	resp, err := http.Get(botUrl + "/getMyCommands")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponseCommand RestResponseCommand
	err = json.Unmarshal(body, &restResponseCommand)
	if err != nil {
		return nil, err
	}
	return restResponseCommand.Result, nil

}
