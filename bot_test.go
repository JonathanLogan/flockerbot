package flockerbot

import (
	"testing"
	"time"

	"github.com/sorcix/irc"
)

func TestBot(t *testing.T) {
	b := &Bot{
		ConnectAddress: "127.0.0.1:6667",
		Nick:           "flocker",
		User:           "flocker",
		Password:       "mypass",
		Timeout:        190,
	}
	b.Handler = func(msg *irc.Message) {
		if msg.Command == "PRIVMSG" {
			b.SendStruct(&irc.Message{
				Command:  "NOTICE",
				Params:   []string{b.ReplyTo(msg)},
				Trailing: msg.Trailing,
			})
		}
	}
	go b.Connect()
	for {
		if b.Connected() {
			break
		}
	}
	b.SendStruct(&irc.Message{
		Command: "JOIN",
		Params:  []string{"#test"},
	})
	time.Sleep(time.Hour)
	if b.Error() != nil {
		t.Errorf("Connect: %s", b.Error())
	}
}
