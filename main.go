package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"
	"time"

	gst "gstreamer-sink"

	ilog "github.com/pion/ion-log"
	sdk "github.com/pion/ion-sdk-go"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
)

var (
	log = ilog.NewLogger(ilog.DebugLevel, "")
)

func init() {
	// This example uses Gstreamer's autovideosink element to display the received video
	// This element, along with some others, sometimes require that the process' main thread is used
	runtime.LockOSThread()
}

func runClientLoop(addr, session string) {

	// new sdk engine
	connector := sdk.NewConnector(addr)

	rtc, err := sdk.NewRTC(connector)
	if err != nil {
		panic(err)
	}

	// subscribe rtp from sessoin
	// comment this if you don't need save to file
	rtc.OnTrack = func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		go func() {
			ticker := time.NewTicker(time.Second * 3)
			for range ticker.C {
				rtcpSendErr := rtc.GetSubTransport().GetPeerConnection().WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}})
				if rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		codecName := strings.Split(track.Codec().RTPCodecCapability.MimeType, "/")[1]
		fmt.Printf("Track has started, of type %d: %s \n", track.PayloadType(), codecName)
		sink := "fakesink"
		switch track.Kind() {
		case webrtc.RTPCodecTypeAudio:
			sink = "autoaudiosink"
		case webrtc.RTPCodecTypeVideo:
			sink = "udpsink host=127.0.0.1 port=12345"
		}
		pipeline := gst.CreatePipeline(strings.ToLower(codecName), sink)
		pipeline.Start()
		defer pipeline.Stop()
		buf := make([]byte, 1400)
		for {
			i, _, readErr := track.Read(buf)
			if readErr != nil {
				log.Errorf("%v", readErr)
				return
			}

			pipeline.Push(buf[:i])
		}
	}

	// client join a session
	err = rtc.Join(session, sdk.RandomKey(4))

	// publish file to session if needed
	if err != nil {
		log.Errorf("error: %v", err)
	}

	select {}
}

func main() {
	// parse flag
	var session, addr string
	flag.StringVar(&addr, "addr", "localhost:5551", "ion-sfu grpc addr")
	flag.StringVar(&session, "session", "ion", "join session name")
	flag.Parse()

	go runClientLoop(addr, session)
	gst.StartMainLoop()
}
