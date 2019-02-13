package main

import (
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	d, _ := discordgo.New(os.Args[1])

	d.AddHandler(func(s *discordgo.Session, ts *discordgo.TypingStart) {
		spew.Dump(ts)
	})

	defer d.Close()

	e := d.Open()
	if e != nil {
		panic(e)
	}

	time.Sleep(time.Minute * 10)
}
