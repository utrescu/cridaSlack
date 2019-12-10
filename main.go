package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"

	"github.com/utrescu/cridaSlack/cmd"
)

type Config struct {
	Webhook string `json:"webhook"`
}

func ReadConfig() (*Config, error) {
	homeDir := ""
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	for _, path := range []string{"/etc/cridaslack.conf", homeDir + "/.cridaslack.conf", "./cridaslack.conf"} {
		file, err := os.Open(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		json.NewDecoder(file)
		conf := Config{}
		err = json.NewDecoder(file).Decode(&conf)
		if err != nil {
			return nil, err
		}
		return &conf, nil
	}

	return nil, errors.New("Config file not found, provide one")
}

// TextMessage content
type TextMessage struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

// BlocksMessage is for formatting messages
type BlocksMessage struct {
	Type string       `json:"type,omitempty"`
	Text *TextMessage `json:"text,omitempty"`
}

// NewBlockMessage retorna un missatge en format mrkdwn
func NewBlockMessage(text string, format string) BlocksMessage {
	bloc := BlocksMessage{}
	bloc.Type = "section"
	bloc.Text = &TextMessage{text, format}
	return bloc
}

// Message represents the message to send to Slack
type Message struct {
	Channel   string          `json:"channel"`
	Username  string          `json:"username,omitempty"`
	Text      string          `json:"text"`
	Blocks    []BlocksMessage `json:"blocks,omitempty"`
	Parse     string          `json:"parse"`
	IconEmoji string          `json:"icon_emoji,omitempty"`
}

// Encode encodes the message to be sent
func (m Message) Encode() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Post sends the message to the provided WebHook
func (m Message) Post(Webhook string) error {
	encoded, err := m.Encode()
	if err != nil {
		return err
	}

	resp, err := http.PostForm(Webhook, url.Values{"payload": {encoded}})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Send failed")
	}
	return nil
}

func main() {

	cfg, err := ReadConfig()
	if err != nil {
		log.Fatalf("Could not read config: %v", err)
	}

	cmd.Execute()
	var msg Message

	if !cmd.Markdown {
		msg = Message{
			Channel:   cmd.Channel,
			Username:  cmd.User,
			Parse:     "full",
			Text:      cmd.Message,
			IconEmoji: cmd.Icon,
		}
	} else {
		var blocs []BlocksMessage
		blocs = append(blocs, NewBlockMessage("mrkdwn", cmd.Message))

		msg = Message{
			Channel:   cmd.Channel,
			Username:  cmd.User,
			Blocks:    blocs,
			Text:      cmd.Message,
			IconEmoji: cmd.Icon,
		}
	}

	err = msg.Post(cfg.Webhook)
	if err != nil {
		log.Fatalf("Post failed: %v", err)
	}
}
