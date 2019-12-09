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

	for _, path := range []string{"/etc/cridaslack.conf", homeDir + "/.cridaslack.conf", "./slackcat.conf"} {
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

type Message struct {
	Channel   string `json:"channel"`
	Username  string `json:"username,omitempty"`
	Text      string `json:"text"`
	Parse     string `json:"parse"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

func (m Message) Encode() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

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

	msg := Message{
		Channel:   cmd.Channel,
		Username:  cmd.User,
		Parse:     "full",
		Text:      cmd.Message,
		IconEmoji: cmd.Icon,
	}

	err = msg.Post(cfg.Webhook)
	if err != nil {
		log.Fatalf("Post failed: %v", err)
	}
}
