package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("spotify")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Please install shpotify", err)
		return
	}

	fmt.Println(GetCurrent())

}

func GetCurrent() (artist, track string) {
	out, err := exec.Command("spotify", "status").Output()
	if err != nil {
		fmt.Println("could not get song", err)
		return "", ""
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Artist") {
			artist = strings.Split(line, ": ")[1]
		} else if strings.Contains(line, "Track") {
			track = strings.Split(line, ": ")[1]
			track = strings.Split(track, "-")[0]
		}
	}
	return artist, track
}
