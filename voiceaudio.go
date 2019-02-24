package main

//import (
//"log"

//"github.com/gordonklaus/portaudio"
//"github.com/rumblefrog/discordgo"
//"gitlab.com/diamondburned/dgvoice"
//)

//var inVoice int64

//func toggleVoiceJoin(chID int64) {
//if inVoice != 0 {
//d.GatewayManager.ChannelVoiceLeave(chID)
//inVoice = 0
//if chID != 0 {
//Record(chID)
//}
//} else {
//Record(chID)
//}
//}

//// Record ..
//func Record(chID int64) {
//c, err := d.State.Channel(chID)
//if err != nil {
//Warn("Voice error: " + err.Error())
//}

//dgv, err := d.GatewayManager.ChannelVoiceJoin(
//c.GuildID, c.ID, true, false,
//)

//if err != nil {
//Warn(err.Error())
//return
//}

//inVoice = chID

//portaudio.Initialize()
//defer portaudio.Terminate()

//out := make([]int16, 960)

//stream, err := portaudio.OpenDefaultStream(0, 2, 48000, len(out), &out)
//if err != nil {
//Warn(err.Error())
//return
//}

//defer stream.Close()

//if err := stream.Start(); err != nil {
//Warn(err.Error())
//return
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
//Warn(err.Error())
//break
//}
//}

////portaudio.Initialize()
////defer portaudio.Terminate()

////in := make([]int32, 64)
////stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
////if (err) != nil {
////Warn(err.Error())
////return
////}

////defer stream.Close()

////if err := stream.Start(); err != nil {
////Warn(err.Error())
////return
////}

////for {
////if err := stream.Read(); err != nil {
////Warn(err.Error())
////return
////}

////if err := binary.Write(f, binary.BigEndian, in); err != nil {
////Warn(err.Error())
////return
////}

////nSamples += len(in)
////select {
////case <-sig:
////return
////default:
////}
////}
////chk(stream.Stop())
//}
