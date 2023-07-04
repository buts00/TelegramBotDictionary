package main

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type Commands struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type RestResponseCommand struct {
	Result []Commands `json:"result"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type Word struct {
	Id          int
	word        string
	translation string
}
