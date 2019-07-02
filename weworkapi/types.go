package weworkapi

import "strings"

type Bot struct {
	weworkBotAPI string  "http://qyapi.weixin.qq.com/cgi-bin/webhook/send"
	key string
	message interface{}
	messageType string
}

// 消息类型
const (
	MsgtypeText string = "text"
	MsgtypeMarkdown string = "markdown"
	MsgtypeImage string = "image"
	MsgtypeNews string = "news"

)

// 文本类型
type MessageText struct {
	Msgtype string `json:"msgtype"`
	Text Text `json:"text"`
}

type Text struct {
	Content string `json:"content"`
	MentionedList []string `json:"mentioned_list"`
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}


// markdown类型
type MessageMarkdown struct {
	Msgtype string `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}


// 图片类型
type MessageImage struct {
	Msgtype string `json:"msgtype"`
	Image Image `json:"image"`
}

type Image struct {
	Base64 string `json:"base64"`
	Md5 string `json:"md5"`
}


// 图文类型
type MessageNews struct {
	Msgtype string `json:"msgtype"`
	News News `json:"news"`
}

type Articles struct {
	Title string `json:"title"`
	Description string `json:"description"`
	URL string `json:"url"`
	Picurl string `json:"picurl"`
}

type News struct {
	Articles []Articles `json:"articles"`
}

func MarkDownMessageColorInfo(msg string) string {
	return "<font color=\"info\">"+msg+"</font>"
}

func MarkDownMessageColorComment(msg string) string {
	return "<font color=\"comment\">"+msg+"</font>"
}

func MarkDownMessageColorWarning(msg string) string {
	return "<font color=\"warning\">"+msg+"</font>"
}

func MarkDownMessageLink(name string, url string) string {
	return "["+name+"]("+url+")"
}


func MarkDownMessageCode(code string) string {
	return "`"+code+"`"
}

func MarkDownMessageBold(msg string) string {
	return "**"+msg+"**"
}

func MarkDownMessageTitle(msg string, level int) string {
	return strings.Repeat("#", level) + " " + msg
}
