package flockerbot

import (
	"fmt"
	"testing"
	"time"

	"github.com/sorcix/irc"
)

func TestBot(t *testing.T) {
	b := &Bot{
		ConnectAddress: "127.0.0.1:6667",
		Nick:           "flocker",
		User:           "flocker",
		Timeout:        90,
		TLS:            true,
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
	b.ConnectedHandler = func() {
		b.SendStruct(&irc.Message{
			Command: "JOIN",
			Params:  []string{"#test"},
		})
	}
	b.Setup()
	go b.StayConnected()
	for {
		time.Sleep(time.Second / 4)
		fmt.Printf(".")
		if b.Connected() {
			fmt.Printf("C")
			break
		}
		if err := b.Error(); err != nil {
			t.Fatalf("Connect error: %s", err)
		}
	}
	time.Sleep(time.Hour)
	if b.Error() != nil {
		t.Errorf("Connect: %s", b.Error())
	}
}
