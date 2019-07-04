package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"wbot/weworkapi"

	"gopkg.in/go-playground/webhooks.v5/gitlab"
)


var req *http.Request
var bot *weworkapi.Bot
func init() {
	bot = weworkapi.NewBot(getenv("KEY"))
}

func main() {
	hook, err := gitlab.New(gitlab.Options.Secret(getenv("X-Gitlab-Token")))
	errorReceiver(err)
	req, _ = http.NewRequest("POST", "", strings.NewReader(getenv("payload")))
	req.Header.Add("X-Gitlab-Event", getenv("X-Gitlab-Event"))
	req.Header.Add("X-Gitlab-Token", getenv("X-Gitlab-Token"))
	payload, err := hook.Parse(req,
		gitlab.PipelineEvents,
		gitlab.MergeRequestEvents,
		gitlab.CommentEvents,
		gitlab.TagEvents,
		gitlab.PushEvents,
		gitlab.SystemHookEvents,
		gitlab.BuildEvents)
	errorReceiver(err)

	switch payload.(type) {
	case gitlab.PipelineEventPayload:
		content := payload.(gitlab.PipelineEventPayload)
		RunPipeline(content)
	case gitlab.MergeRequestEventPayload:
		content := payload.(gitlab.MergeRequestEventPayload)
		RunMergeRequest(&content)
	case gitlab.TagEventPayload:
		content := payload.(gitlab.TagEventPayload)
		RunTag(&content)
		

	}

}

// pipeline 运行
func RunPipeline(payload gitlab.PipelineEventPayload){
	if payload.ObjectAttributes.Status == "running" {
		fmt.Println("continue running")
	}
	var message interface{}
	target := payload.ObjectAttributes.Status
	switch target {
	case "failed":
		message = weworkapi.MessageMarkdown{
			Msgtype: weworkapi.MsgtypeMarkdown,
			Markdown: weworkapi.Markdown{
				Content: fmt.Sprintf(
					"# 项目:%s (%s)Pipeline %s\n > 最新提交: %s \n > 时间:  %s \n > 作者: %s",
					payload.Project.Name,
					payload.Project.Description,
					weworkapi.MarkDownMessageColorWarning(" 运行失败"),
					payload.Commit.Message + "  " + weworkapi.MarkDownMessageLink(payload.Commit.ID[0:8], payload.Commit.URL),
					payload.Commit.Timestamp,
					payload.Commit.Author.Name + "(" + payload.Commit.Author.Email + ")",
				),
			},
		}
		break
	case "success":
		message = weworkapi.MessageMarkdown{
			Msgtype: weworkapi.MsgtypeMarkdown,
			Markdown: weworkapi.Markdown{
				Content: fmt.Sprintf(
					"# 项目:%s (%s)Pipeline %s\n > 最新提交: %s \n > 时间:  %s \n > 作者: %s",
					payload.Project.Name,
					payload.Project.Description,
					weworkapi.MarkDownMessageColorInfo(" 运行成功"),
					payload.Commit.Message + weworkapi.MarkDownMessageLink(payload.Commit.ID[0:8], payload.Commit.URL),
					payload.Commit.Timestamp,
					payload.Commit.Author.Name + "(" + payload.Commit.Author.Email + ")",
				),
			},
		}
		break
	}
	bot.SetMessage(message)
	res, err := bot.Send()
	errorReceiver(err)
	fmt.Println(string(res))
}

// merge request请求
func RunMergeRequest(payload *gitlab.MergeRequestEventPayload) {
	target := payload.ObjectAttributes.Action
	var message interface{}
	//var message1 interface{}
	switch target {
	case "open":
		//message = weworkapi.MessageText{
		//	Msgtype:weworkapi.MsgtypeText,
		//	Text:weworkapi.Text{
		//		Content: fmt.Sprintf(
		//			"你有一个新的Merge Rquest请求需要处理: [%s] ==> [%s], 作者: %s ",
		//			payload.ObjectAttributes.SourceBranch,
		//			payload.ObjectAttributes.TargetBranch,
		//			payload.User.Name + "(" + payload.User.UserName + ")",
		//			),
		//	},
		//}

		message = weworkapi.MessageMarkdown{
			Msgtype: weworkapi.MsgtypeMarkdown,
			Markdown: weworkapi.Markdown{
				Content: fmt.Sprintf(
					"# 项目:%s (%s)Merge Request %s\n > 最新提交: %s \n > 时间:  %s \n > 作者: %s \n 合并请求: %s ==> %s \n 处理地址: %s",
					payload.Project.Name,
					payload.Project.Description,
					weworkapi.MarkDownMessageColorComment(" 待处理"),
					payload.ObjectAttributes.LastCommit.Message + "  " + weworkapi.MarkDownMessageLink(payload.ObjectAttributes.LastCommit.ID[0:8], payload.ObjectAttributes.LastCommit.URL),
					payload.ObjectAttributes.LastCommit.Timestamp,
					payload.ObjectAttributes.LastCommit.Author.Name + "(" + payload.ObjectAttributes.LastCommit.Author.Email + ")",
					payload.ObjectAttributes.SourceBranch,
					payload.ObjectAttributes.TargetBranch,
					weworkapi.MarkDownMessageLink("请点击", payload.ObjectAttributes.URL),
				),
			},
		}
		break
	case "merge":
		message = weworkapi.MessageMarkdown{
			Msgtype: weworkapi.MsgtypeMarkdown,
			Markdown: weworkapi.Markdown{
				Content: fmt.Sprintf(
					"# 项目:%s (%s)Merge Request %s\n > 最新提交: %s \n > 处理时间:  %s \n > 作者: %s \n 合并请求: %s ==> %s \n 处理人: %s",
					payload.Project.Name,
					payload.Project.Description,
					weworkapi.MarkDownMessageColorInfo(" 已合并"),
					payload.ObjectAttributes.LastCommit.Message + "  " + weworkapi.MarkDownMessageLink(payload.ObjectAttributes.LastCommit.ID[0:8], payload.ObjectAttributes.LastCommit.URL),
					payload.ObjectAttributes.UpdatedAt,
					payload.ObjectAttributes.LastCommit.Author.Name + "(" + payload.ObjectAttributes.LastCommit.Author.Email + ")",
					payload.ObjectAttributes.SourceBranch,
					payload.ObjectAttributes.TargetBranch,
					payload.User.Name + "(" + payload.User.UserName + ")",
				),
			},
		}
		break
	}
	bot.SetMessage(message)
	res, err := bot.Send()
	errorReceiver(err)
	fmt.Println(string(res))
}

// set tag event
func RunTag(payload *gitlab.TagEventPayload)  {
	var message interface{}
	ref := strings.Split(payload.Ref, "/")
	message = weworkapi.MessageMarkdown{
		Msgtype: weworkapi.MsgtypeMarkdown,
		Markdown: weworkapi.Markdown{
			Content: fmt.Sprintf(
				"# 项目:%s (%s)Merge Request %s\n > 版本: %s \n > 处理人: %s",
				payload.Project.Name,
				payload.Project.Description,
				weworkapi.MarkDownMessageColorInfo("有新版本发布啦"),
				ref[len(ref)-1:],
				payload.UserName,

			),
		},
	}
	bot.SetMessage(message)
	res, err := bot.Send()
	errorReceiver(err)
	fmt.Println(string(res))
}



func getenv(s string) string {
	return os.Getenv(s)
}

func errorReceiver( err error) {
	if err != nil {
		log.Fatal(err)
	}
}