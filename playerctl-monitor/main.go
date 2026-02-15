package main

import (
	"flag"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

func queryPlayerState() (playerState, error) {
	position, err := strconv.ParseFloat(playerctlCmd("position"), 64)
	volume, err := strconv.ParseFloat(playerctlCmd("volume"), 64)

	return playerState{
		status:		playerctlCmd("status") == "Playing",
		position:	position,
		volume:		volume,
		loop:		playerctlCmd("loop"),
		shuffle:	playerctlCmd("shuffle") == "On",
	}, err
}

func formatTime(seconds int) string {

	var output string

	minutes := math.Floor(float64(seconds / 60))
	seconds = seconds % 60

	output += strconv.FormatFloat(minutes, 'f', -1, 64) + ":"
	if seconds < 10 {
		output += "0"
	}
	output += strconv.FormatInt(int64(seconds), 10)

	return output
}

func genOutput(format string) string {
	metadata := queryMetadata()
	playerState, err := queryPlayerState()

	if err != nil {
		return ""
	}

	output := format
	
	// todo refactor lol
	output = strings.ReplaceAll(output, "@t@", metadata.title)
	output = strings.ReplaceAll(output, "@a@", metadata.artist)
	output = strings.ReplaceAll(output, "@A@", metadata.albumArtist)
	output = strings.ReplaceAll(output, "@al@", metadata.album)
	output = strings.ReplaceAll(output, "@au@", metadata.artUrl)
	output = strings.ReplaceAll(output, "@l@", strconv.FormatUint(metadata.length, 10))
	output = strings.ReplaceAll(output, "@lF@", formatTime(int(metadata.length / 1_000_000)))
	output = strings.ReplaceAll(output, "@s@", strconv.FormatBool(playerState.status))
	output = strings.ReplaceAll(output, "@p@", strconv.FormatFloat(playerState.position, 'f', -1, 64))
	output = strings.ReplaceAll(output, "@pF@", formatTime(int(playerState.position)))
	output = strings.ReplaceAll(output, "@v@", strconv.FormatFloat(playerState.volume, 'f', -1, 64))
	output = strings.ReplaceAll(output, "@L@", playerState.loop)
	output = strings.ReplaceAll(output, "@S@", strconv.FormatBool(playerState.shuffle))

	return output
}

func main() {

	var format string
	flag.StringVar(&format, "f", "@a@ - @t@", "Format for the information to be outputted in")
	pollInterval := flag.Int("p", 1000, "Time in ms between each poll of playerctl")

	flag.Parse()

	for {
		fmt.Println(genOutput(format))
		time.Sleep(time.Duration(*pollInterval) * time.Millisecond)
	}
}