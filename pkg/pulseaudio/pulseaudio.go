package pulseaudio

import (
	"errors"
	"log"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const maxVolume = 65536

var volumeRegex = regexp.MustCompile("^set-sink-volume (?P<sink>[^ ]+) (?P<volume>.*)")
var volumeRegexSinkIndex = volumeRegex.SubexpIndex("sink")
var volumeRegexVolumeIndex = volumeRegex.SubexpIndex("volume")

var muteRegex = regexp.MustCompile("^set-sink-mute (?P<sink>[^ ]+) (?P<mute>.*)")
var muteRegexSinkIndex = muteRegex.SubexpIndex("sink")
var muteRegexMuteIndex = muteRegex.SubexpIndex("mute")

var defaultSinkRegexp = regexp.MustCompile("^set-default-sink (?P<sink>[^ ]+)")

func volumeHexToPercent(hexVolume string) (int, error) {
	hexRegex := regexp.MustCompile("0[xX]")

	cleanedVolumeString := hexRegex.ReplaceAllString(hexVolume, "")
	volumeInt, intConvertErr := strconv.ParseInt(cleanedVolumeString, 16, 16)

	if intConvertErr != nil {
		return -1, intConvertErr
	}

	volumePercent := math.Round(float64(volumeInt) / maxVolume * 100)

	return int(volumePercent), nil
}

type pulseAudio struct {
	volumeMap   map[string]int
	muteMap     map[string]bool
	defaultSink string
}

func New() pulseAudio {
	pa := pulseAudio{
		volumeMap: make(map[string]int),
		muteMap:   make(map[string]bool),
	}

	pa.Update()

	return pa
}

func (pa *pulseAudio) Update() {
	out, err := exec.Command("pacmd", "dump").Output()

	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if volumeRegex.MatchString(line) {
			matches := volumeRegex.FindStringSubmatch(line)
			sinkName := matches[volumeRegexSinkIndex]

			volumePercent, volumeConvertErr := volumeHexToPercent(matches[volumeRegexVolumeIndex])
			if volumeConvertErr != nil {
				log.Printf("Could not convert volume to int, got err: %s", volumeConvertErr)
				continue
			}

			pa.volumeMap[sinkName] = volumePercent
		} else if defaultSinkRegexp.MatchString(line) {
			matches := defaultSinkRegexp.FindStringSubmatch(line)
			sinkName := matches[defaultSinkRegexp.SubexpIndex("sink")]

			pa.defaultSink = sinkName
		} else if muteRegex.MatchString(line) {
			matches := muteRegex.FindStringSubmatch(line)
			sinkName := matches[muteRegexSinkIndex]
			mutedString := matches[muteRegexMuteIndex]

			pa.muteMap[sinkName] = mutedString == "yes"
		}
	}
}

func (pa pulseAudio) GetVolume(sink string) (int, error) {
	if sink == "" {
		if pa.defaultSink == "" {
			return -1, errors.New("Could not find a requested sink, default sink is not defined")
		}

		sink = pa.defaultSink
	}

	if sinkVolume, ok := pa.volumeMap[sink]; ok {
		return sinkVolume, nil
	} else {
		return -1, errors.New("Could not fetch volume for requested or default sinks")
	}
}

func (pa pulseAudio) GetMute(sink string) (bool, error) {
	if sink == "" {
		if pa.defaultSink == "" {
			return false, errors.New("Could not find a requested sink, default sink is not defined")
		}

		sink = pa.defaultSink
	}

	if sinkMute, ok := pa.muteMap[sink]; ok {
		return sinkMute, nil
	} else {
		return false, errors.New("Could not fetch volume for requested or default sinks")
	}
}
