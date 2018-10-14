package notification

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// PlayNotification ...
func PlayNotification() {
	f, _ := os.Open("notification/notification.mp3")
	s, format, err := mp3.Decode(f)
	log.Println(err)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan struct{})
	speaker.Play(beep.Seq(s, beep.Callback(func() { close(done) })))
	<-done

}
