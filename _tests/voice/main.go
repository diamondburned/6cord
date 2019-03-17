package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

func main() {
	d, _ := discordgo.New(os.Args[1])
	if err := d.Open(); err != nil {
		panic(err)
	}

	defer d.Close()

	c, err := d.State.Channel(os.Args[2])
	if err != nil {
		panic(err)
	}

	dgv, err := d.ChannelVoiceJoin(
		c.GuildID, c.ID, false, false,
	)

	defer dgv.Close()

	if err != nil {
		panic(err)
	}

	//dgv.AddHandler(func(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
	//spew.Dump(vc, vs)
	//})

	recv := make(chan *discordgo.Packet, 2)
	go ReceivePCM(dgv, recv)

	send := make(chan []int16, 2)
	go dgvoice.SendPCM(dgv, send)

	dgv.Speaking(true)
	defer dgv.Speaking(false)

	for {
		p, ok := <-recv
		if !ok {
			return
		}

		send <- p.PCM
	}

	//portaudio.Initialize()
	//defer portaudio.Terminate()

	//out := make([]int16, 960)

	//stream, err := portaudio.OpenDefaultStream(0, 2, 48000, len(out), &out)
	//if err != nil {
	//panic(err)
	//}

	//defer stream.Close()

	//if err := stream.Start(); err != nil {
	//panic(err)
	//}

	//defer stream.Stop()

	//recv := make(chan *discordgo.Packet, 2)
	//go dgvoice.ReceivePCM(dgv, recv)

	//for {
	//p, ok := <-recv
	//if !ok {
	//return
	//}

	//log.Println(len(p.PCM))

	//copy(out, p.PCM)

	//if err := stream.Write(); err != nil {
	//panic(err)
	//}
	//}
}

var (
	speakers    map[uint32]*gopus.Decoder
	opusEncoder *gopus.Encoder
	mu          sync.Mutex
)

func ReceivePCM(v *discordgo.VoiceConnection, c chan *discordgo.Packet) {
	if c == nil {
		return
	}

	var err error

	for {
		if v.Ready == false || v.OpusRecv == nil {
			OnError(fmt.Sprintf("Discordgo not to receive opus packets. %+v : %+v", v.Ready, v.OpusSend), nil)
			return
		}

		log.Println("I'm desperate")

		p, ok := <-v.OpusRecv
		if !ok {
			log.Println("Closed ch")
			return
		}

		log.Println("Received")

		if speakers == nil {
			speakers = make(map[uint32]*gopus.Decoder)
		}

		_, ok = speakers[p.SSRC]
		if !ok {
			speakers[p.SSRC], err = gopus.NewDecoder(48000, 2)
			if err != nil {
				OnError("error creating opus decoder", err)
				continue
			}
		}

		p.PCM, err = speakers[p.SSRC].Decode(p.Opus, 960, false)
		if err != nil {
			OnError("Error decoding opus data", err)
			continue
		}

		c <- p
	}
}

func OnError(str string, err error) {
	println(str)
	panic(err)
}
