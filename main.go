package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	cmd := exec.Command("spotify")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Please install shpotify: 'brew install shpotify'", err)
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/", SongHandler)
	r.HandleFunc("/public/{folder}/{file}", PublicHandler)
	r.HandleFunc("/live", live)

	fmt.Println("Server started on localhost:8000")
	err = http.ListenAndServe(":8000", r)

	if err != nil {
		fmt.Println(err)
		return
	}

}

func getPath() string {
	return os.Getenv("GOPATH") + "/src/github.com/fjukstad/now-playing"
}

func SongHandler(w http.ResponseWriter, r *http.Request) {
	artist, track := GetCurrent()

	path := getPath()
	files, err := ioutil.ReadDir(path + "/public/img")
	if err != nil {
		log.Fatal(err)
	}

	randomFile := files[rand.Intn(len(files))]

	np := NowPlaying{artist, track, randomFile.Name()}

	indexTemplate := template.Must(template.ParseFiles(path + "/index.html"))
	indexTemplate.Execute(w, np)

}

type NowPlaying struct {
	Artist string
	Track  string
	Image  string
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

func PublicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folder := vars["folder"]
	file := vars["file"]

	path := getPath()
	base := path + "/public/"
	filename := base + folder + "/" + file

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		w.Write([]byte("Could not find file " + filename))
	} else {
		w.Write(f)
	}
}
