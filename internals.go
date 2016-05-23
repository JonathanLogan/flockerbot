package flockerbot

import (
	"strconv"
	"time"

	"github.com/JonathanLogan/flockerbot/fixbuffer"
)

// socketReader reads a string from a blocking io.Reader socket and sends it to channel c.
func (b *Bot) socketReader() {
	defer func() {
		recover()
	}()
	var err error
	var line []byte
	r := fixbuffer.New(b.socket, 2048, []byte("\n"))
ReadLoop:
	for {
		line, err = r.ReadBytes()
		if err == fixbuffer.ErrNotFound {
			continue ReadLoop
		}
		msg := &channelString{
			Data: string(line),
			Dir:  socketRead,
			Err:  err,
		}
		b.socketChan <- msg
		if err != nil {
			break ReadLoop
		}
	}
	b.socket.Close()
	return
}

// setNick sets the nick of the connection
func (b *Bot) setNick() {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	if !b.IsServer {
		nick := b.Nick
		if b.nickCount >= 0 {
			nick += strconv.Itoa(b.nickCount)
		}
		b.activeNick = nick
		b.SendString("NICK " + nick)
		b.nickCount++
	}
}

// setUser sets the user
func (b *Bot) setUser() {
	if b.IsServer {
		b.SendString("SERVER " + b.User + " 1")
	} else {
		b.SendString("USER " + b.User + " 0 * " + b.User)
	}
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	b.userSet = true
}

// sendPass sends the PASS command on connect
func (b *Bot) sendPass() {
	if b.Password != "" {
		b.SendString("PASS " + b.Password)
	}
}

// Timeout ticker.
func (b *Bot) ticker() {
	defer func() {
		recover()
	}()
SendLoop:
	for {
		time.Sleep(time.Second * 10)
		select {
		case b.socketChan <- nil:
			continue SendLoop
		default:
			break SendLoop
		}
	}
}

func now() int64 {
	return time.Now().UTC().Unix()
}
