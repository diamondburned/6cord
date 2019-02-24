package main

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	d, _ := discordgo.New(os.Args[1])
	if err := d.Open(); err != nil {
		panic(err)
	}

	defer d.Close()

	d.AddHandler(func(s *discordgo.Session, t *discordgo.TypingStart) {
		spew.Dump(t)
	})
}
