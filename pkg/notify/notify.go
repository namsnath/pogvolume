package notify

import (
	"fmt"
	"log"
	"os/exec"
)

type notificationUrgency string

const (
	UrgencyLow      notificationUrgency = "low"
	UrgencyNormal   notificationUrgency = "normal"
	UrgencyCritical notificationUrgency = "critical"
)

type NotifyData struct {
	Group         string
	Title         string
	Text          string
	Urgency       notificationUrgency
	ProgressValue int
}

func Notify(data NotifyData) {
	args := []string{}

	if data.Title != "" {
		args = append(args, data.Title)
	} else {
		log.Fatal("Cannot notify with no Title (summary)")
		return
	}

	if data.Text != "" {
		args = append(args, data.Text)
	}

	if data.ProgressValue >= 0 {
		args = append(args, "-h", fmt.Sprintf("int:value:%d", data.ProgressValue))
	} else {
		log.Fatal("Cannot notify with negative ProgressValue")
		return
	}

	if data.Group != "" {
		args = append(args, "-h", fmt.Sprintf("string:x-dunst-stack-tag:%s", data.Group))
	}

	if data.Urgency == "" {
		data.Urgency = UrgencyNormal
	}

	args = append(args, "-u", string(data.Urgency))

	cmd := exec.Command("notify-send", args...)
	cmd.Run()
}
