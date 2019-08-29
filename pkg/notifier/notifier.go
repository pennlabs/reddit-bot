package notifier

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/turnage/graw/reddit"
)

type Notifier struct {
	client       *slack.Client
	triggerRegex string
	channel      string
}

func New(client *slack.Client, triggerSubs []string, channel string) *Notifier {
	return &Notifier{
		client:       client,
		triggerRegex: strings.Join(triggerSubs, "|"),
		channel:      channel,
	}
}

func (a *Notifier) Post(post *reddit.Post) error {
	matched, err := regexp.Match(a.triggerRegex, []byte(strings.ToLower(post.Title)))
	if err != nil {
		return err
	}
	if !matched {
		matched, err = regexp.Match(a.triggerRegex, []byte(strings.ToLower(post.SelfText)))
		if err != nil {
			return err
		}
	}

	if matched {
		return a.alertPost(post)
	}

	return nil
}

func (a *Notifier) Comment(comment *reddit.Comment) error {
	matched, err := regexp.Match(a.triggerRegex, []byte(strings.ToLower(comment.Body)))
	if err != nil {
		return err
	}
	if matched {
		return a.alertComment(comment)
	}

	return nil
}

func (a *Notifier) alertPost(post *reddit.Post) error {
	response := fmt.Sprintf(`%s mentioned Penn Labs in "%s". Link to post:
https://reddit.com%s`, post.Author, post.Title, post.Permalink)
	text := slack.MsgOptionText(response, false)
	_, _, err := a.client.PostMessage(a.channel, text)
	if err != nil {
		log.Println(err)
	}
	log.Println(response)
	return nil
}

func (a *Notifier) alertComment(comment *reddit.Comment) error {
	response := fmt.Sprintf(`%s mentioned Penn Labs in a comment. Link to comment:
https://reddit.com%s`, comment.Author, comment.Permalink)
	text := slack.MsgOptionText(response, false)
	_, _, err := a.client.PostMessage(a.channel, text)
	if err != nil {
		log.Println(err)
	}
	log.Println(response)
	return nil
}
