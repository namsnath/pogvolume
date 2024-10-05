package main

import (
	// "flag"

	"github.com/namsnath/pogvolume/pkg/notify"
	"github.com/namsnath/pogvolume/pkg/pulseaudio"
)

func main() {
	// notify.Notify("a", "b")
	// wordPtr := flag.String("word", "foo", "a string")
	// notify.ProgressNotify("volume", "Title", "a", 50)

	notificationGroup := "volume"
	notificationTitle := "Volume"

	pa := pulseaudio.New()

	volume, err := pa.GetVolume("")

	n := notify.NotifyData{
		Group:   notificationGroup,
		Title:   notificationTitle,
		Text:    "",
		Urgency: notify.UrgencyNormal,
	}

	if err != nil {
		n.Urgency = notify.UrgencyCritical
		n.Text = "Could not fetch volume"
	} else {
		n.ProgressValue = volume
	}

	notify.Notify(n)
}
