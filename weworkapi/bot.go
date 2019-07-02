package weworkapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	UnknowMessageType      = errors.New("unknown message types")
)

func NewBot(key string) *Bot {
	return &Bot{
		key: key,
		weworkBotAPI: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send",
	}
}

func (bot *Bot) SetBotApi(weworkBotAPI string) {
	if weworkBotAPI == "" {
		weworkBotAPI = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"
	}
	bot.weworkBotAPI = weworkBotAPI
}

func (bot *Bot) SetMessage(msg interface{}) error {
	switch msg.(type) {
	case MessageText:
		bot.messageType = MsgtypeText
		bot.message = msg.(MessageText)
	case MessageImage:
		bot.messageType = MsgtypeImage
		bot.message = msg.(MessageImage)
	case MessageMarkdown:
		bot.messageType = MsgtypeMarkdown
		bot.message = msg.(MessageMarkdown)
	case MessageNews:
		bot.messageType = MsgtypeNews
		bot.message = msg.(MessageNews)
	default:
		return UnknowMessageType
	}
	return nil
}

func (bot *Bot) Send() ([]byte, error) {
	tr := &http.Transport{    //解决x509: certificate signed by unknown authority
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	query := url.Values{}
	query.Add("key", bot.key)
	message, _ := json.Marshal(bot.message)
	url := bot.weworkBotAPI + "?" + query.Encode()
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: tr,
	}
	fmt.Println(url,string(message))
	response, err := client.Post(url, "application/json", bytes.NewReader(message))
	errorReceiver(err)
	defer response.Body.Close()
	if err != nil {
		errorReceiver(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	errorReceiver(err)
	return body, nil
}



func errorReceiver( err error) {
	if err != nil {
		log.Fatal(err)
	}
}