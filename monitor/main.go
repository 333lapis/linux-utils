package main

import (
	"flag"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"encoding/json"
)

type playerMetadata struct {
	artUrl		string
	title		string
	artist		string
	length		uint64

	album		string
	albumArtist	string
}

type playerState struct {
	status		bool
	position	float64
	volume		float64
	loop		string
	shuffle		bool
}

type batteryState struct {
	capacity int64
}

type pulseAudioState struct {
	mute bool
	volume string
}

func queryPulseAudioState() pulseAudioState {
	muteCmd, err := exec.Command("pactl", "-f", "json", "get-sink-mute", "@DEFAULT_SINK@").Output()
	volCmd, err := exec.Command("pactl", "-f", "json", "get-sink-volume", "@DEFAULT_SINK@").Output()

	var muteCmdRes map[string]interface{}
	var volCmdRes map[string]interface{}

	err = json.Unmarshal([]byte(muteCmd), &muteCmdRes)
	err = json.Unmarshal([]byte(volCmd), &volCmdRes)

	if err != nil {
		return pulseAudioState{}
	}

	volCmdRes = volCmdRes["volume"].(map[string]interface{})
	frontLeft := volCmdRes["front-left"].(map[string]interface{})

	return pulseAudioState{
		mute: muteCmdRes["mute"].(bool),
		volume: frontLeft["value_percent"].(string),
	}
}

func queryBatteryState(battery string) batteryState {
	batteryPath := fmt.Sprintf("/sys/class/power_supply/%s", battery)
	capacityString, err := exec.Command("cat", fmt.Sprintf("%s/capacity", batteryPath)).Output()
	capacity, err := strconv.ParseInt(strings.TrimRight(string(capacityString), "\n"), 10, 64)
	
	if err != nil {
		return batteryState{}
	}

	return batteryState{
		capacity:	capacity,
	}
}

func playerctlCmd(args ...string) string {
	result, err := exec.Command("playerctl", args...).Output()

	if err != nil {
		return ""
	}

	return strings.TrimRight(string(result), "\n")
}

func queryMetadata() playerMetadata {
	length, err := strconv.ParseUint(playerctlCmd("metadata", "mpris:length"), 10, 32)

	if err != nil {
		return playerMetadata{}
	}

	return playerMetadata{
		artUrl:		playerctlCmd("metadata", "mpris:artUrl"),
		title:		playerctlCmd("metadata", "xesam:title"),
		artist:		playerctlCmd("metadata", "xesam:artist"),
		length:		length,
		album:		playerctlCmd("metadata", "xesam:album"),
		albumArtist: playerctlCmd("metadata", "xesam:albumArtist"),
	}
}

func queryPlayerState() playerState {
	position, err := strconv.ParseFloat(playerctlCmd("position"), 64)
	volume, err := strconv.ParseFloat(playerctlCmd("volume"), 64)

	if err != nil {
		return playerState{}
	}

	return playerState{
		status:		playerctlCmd("status") == "Playing",
		position:	position,
		volume:		volume,
		loop:		playerctlCmd("loop"),
		shuffle:	playerctlCmd("shuffle") == "On",
	}
}

func formatTime(seconds int, leadingZero bool) string {

	var output string

	minutes := math.Floor(float64(seconds / 60))
	seconds = seconds % 60

	if leadingZero {
		output += strconv.FormatFloat(minutes, 'f', -1, 64) + ":"
	}
	if seconds < 10 {
		output += "0"
	}

	output += strconv.FormatInt(int64(seconds), 10)

	return output
}

func genOutput(format string) string {
	metadata := queryMetadata()
	playerState := queryPlayerState()
	batteryState := queryBatteryState("BAT1")
	pulseAudioState := queryPulseAudioState()
	t := time.Now()

	output := format
	
	// todo refactor lol

	// playerctl
	output = strings.ReplaceAll(output, "@p:t@", metadata.title)
	output = strings.ReplaceAll(output, "@p:a@", metadata.artist)
	output = strings.ReplaceAll(output, "@p:A@", metadata.albumArtist)
	output = strings.ReplaceAll(output, "@p:al@", metadata.album)
	output = strings.ReplaceAll(output, "@p:au@", metadata.artUrl)
	output = strings.ReplaceAll(output, "@p:l@", strconv.FormatUint(metadata.length, 10))
	output = strings.ReplaceAll(output, "@p:lF@", formatTime(int(metadata.length / 1_000_000), true))
	output = strings.ReplaceAll(output, "@p:s@", strconv.FormatBool(playerState.status))
	output = strings.ReplaceAll(output, "@p:p@", strconv.FormatFloat(playerState.position, 'f', -1, 64))
	output = strings.ReplaceAll(output, "@p:pF@", formatTime(int(playerState.position), true))
	output = strings.ReplaceAll(output, "@p:v@", strconv.FormatFloat(playerState.volume, 'f', -1, 64))
	output = strings.ReplaceAll(output, "@p:L@", playerState.loop)
	output = strings.ReplaceAll(output, "@p:S@", strconv.FormatBool(playerState.shuffle))

	// battery
	output = strings.ReplaceAll(output, "@b:c@", strconv.FormatInt(batteryState.capacity, 10))

	// time
	output = strings.ReplaceAll(output, "@t:h@", strconv.Itoa(t.Hour()))
	output = strings.ReplaceAll(output, "@t:m@", formatTime(t.Minute(), false))
	output = strings.ReplaceAll(output, "@t:s@", formatTime(t.Second(), false))

	// PulseAudio
	output = strings.ReplaceAll(output, "@pa:m@", strconv.FormatBool(pulseAudioState.mute))
	output = strings.ReplaceAll(output, "@pa:v@", pulseAudioState.volume)

	return output
}

func main() {

	var format string
	var pollInterval int

	flag.StringVar(&format, "f", "@t:h@:@t:m@:@t:s@ @p:a@ - @p:t@", "Format for the information to be outputted in")
	flag.IntVar(&pollInterval, "p", 1000, "Time in ms between each poll of playerctl")

	flag.Parse()

	for {
		fmt.Println(genOutput(format))
		time.Sleep(time.Duration(pollInterval) * time.Millisecond)
	}
}