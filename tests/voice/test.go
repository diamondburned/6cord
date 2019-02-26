package main

import (
	"bytes"
	"encoding/binary"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/hraban/opus.v2"
	"os"
	"path"
	"strconv"
	"time"
)

type recordedUser struct {
	user            *discordgo.User
	ssrc            uint32
	decoder         *opus.Decoder
	currentFilePath string
	fileOpened      time.Time
	file            *os.File
	lastTimestamp   uint32 // raw timestamp - sample no.
	baseTimestamp   timestamp
}

type timestamp struct {
	remoteTimestamp uint32 // first raw timestamp received
	localTimestamp  int64  // time when we received first frame in milliseconds
}

type mixData struct {
	pcm                   []int16
	remotePacketTimestamp uint32
	baseTimestamp         timestamp
}

type UserMap map[uint32]*recordedUser

const frameSize = 2 * 20 * 48000 / 1000
const timeoutFramesCount = (2 * 60 * 1000) / 20

var silence = []byte{0xF8, 0xFF, 0xFE}

func (bot *Bot) StartRecording() (chan bool, error) {
	channel, err := bot.Session.Channel(bot.Config.RecordedChannel)
	if err != nil {
		return nil, err
	}

	closeChan := make(chan bool)
	go recordChannel(channel, closeChan)

	return closeChan, err
}

/*func mixUsers(input chan mixData, close chan bool) {
	var startTime int64 = 0
	buffer := make([]int16, frameSize*100) // 2 seconds buffer

	mixFile, err := os.OpenFile("mix.pcm", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return
	}

	for {
		select {
		case data := <-input:
			localPacketTimestamp := data.baseTimestamp.localTimestamp + int64(((data.remotePacketTimestamp - data.baseTimestamp.remoteTimestamp) / 960) * 20)

			if startTime == 0 {
				startTime = data.baseTimestamp.localTimestamp
				buffer = append(buffer, data.pcm...)

				startTime += 20
			}

			if localPacketTimestamp >= startTime
		case <-close:
			mixFile.Close()
			return
		}
	}
}*/

func recordChannel(channel *discordgo.Channel, closeChan chan bool) {
	//mix := make(chan mixData)
	users := make(UserMap)

	/*log, err := os.OpenFile("packetlog2.csv", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return
	}

	log.Write([]byte("Timestamp;SSRC;Sequence\n"))*/

	pcm := make([]int16, frameSize)

	voice, err := bot.Session.ChannelVoiceJoin(channel.GuildID, channel.ID, true, false)
	if err != nil {
		return
	}

	voice.AddHandler(func(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
		user, err := bot.Session.User(vs.UserID)
		if err != nil {
			return
		}

		if rec, exists := users[uint32(vs.SSRC)]; !exists {
			rec, err := initUser(uint32(vs.SSRC), user)
			if err != nil {
				return
			}

			users[uint32(vs.SSRC)] = rec
		} else if rec.user == nil {
			err = rec.setUser(user)
			if err != nil {
				return
			}
		}
	})

	for i := 0; i < 10; i++ {
		voice.OpusSend <- silence
	}

	for {
		select {
		case packet := <-voice.OpusRecv:
			//log.Write([]byte(fmt.Sprintf("%d;%d;%d\n", packet.Timestamp, packet.SSRC, packet.Sequence)))

			// throw away silence
			if bytes.Equal(packet.Opus, silence) {
				continue
			}

			if _, exists := users[packet.SSRC]; !exists {
				rec, err := initUser(packet.SSRC, nil)
				if err != nil {
					return
				}

				users[packet.SSRC] = rec
			}

			user := users[packet.SSRC]

			if user.lastTimestamp != 0 {
				silentFrames := (packet.Timestamp - user.lastTimestamp) / 960

				if silentFrames < timeoutFramesCount {
					for i := uint32(0); i < silentFrames-1; i++ {
						binary.Write(user.file, binary.LittleEndian, new([frameSize]byte))
					}
				} else {
					user.openFile()
				}
			} else {
				user.baseTimestamp.remoteTimestamp = packet.Timestamp
				user.baseTimestamp.localTimestamp = time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
			}

			user.lastTimestamp = packet.Timestamp

			_, err := user.decoder.Decode(packet.Opus, pcm)
			if err != nil {
				return
			}

			//mix <- mixData{pcm: pcm, baseTimestamp: user.baseTimestamp, remotePacketTimestamp: packet.Timestamp}

			if time.Now().UTC().After(user.fileOpened.Add(1 * time.Hour)) {
				user.openFile()
			}

			err = binary.Write(user.file, binary.LittleEndian, pcm)
			if err != nil {
				return
			}

		case <-closeChan:
			voice.Close()

			//log.Close()

			for _, user := range users {
				err := user.file.Close()
				if err != nil {
					return
				}
			}

			return
		}
	}
}

func initUser(ssrc uint32, user *discordgo.User) (*recordedUser, error) {
	var newUser = recordedUser{ssrc: ssrc, user: user}

	err = newUser.openFile()
	if err != nil {
		return nil, err
	}

	newUser.decoder, err = opus.NewDecoder(48000, 2)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (ru *recordedUser) openFile() error {
	if ru.file != nil {
		err = ru.file.Close()
		if err != nil {
			return err
		}
	}

	var identifier string

	if ru.user == nil {
		identifier = strconv.Itoa(int(ru.ssrc))
	} else {
		identifier = ru.user.Username
	}

	ru.currentFilePath = path.Clean(path.Join(bot.Config.OutputPath, identifier+"_"+time.Now().UTC().Format("2006-01-02_15-04-05")) + ".pcm")
	ru.file, err = os.OpenFile(ru.currentFilePath, os.O_RDWR|os.O_CREATE, 0755)
	ru.fileOpened = time.Now().UTC()

	return err
}

func (ru *recordedUser) setUser(user *discordgo.User) error {
	ru.user = user

	err = ru.file.Close()
	if err != nil {
		return err
	}

	newPath := path.Clean(path.Join(bot.Config.OutputPath, user.Username+"_"+time.Now().UTC().Format("2006-01-02_15-04-05")) + ".pcm")

	err = os.Rename(ru.currentFilePath, newPath)
	if err != nil {
		return err
	}

	ru.currentFilePath = newPath
	ru.file, err = os.OpenFile(ru.currentFilePath, os.O_RDWR|os.O_CREATE, 0755)

	return err
}
