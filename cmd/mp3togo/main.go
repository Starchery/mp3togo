// Package main allows the user to enter video URLs to download
// from Youtube.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"unicode"

	"github.com/rylio/ytdl"
)

const maxConcurrentDownloads = 5

type app struct {
	videos    chan string
	toConvert chan string
	started   chan bool
	wg        sync.WaitGroup
}

func new() *app {
	v := make(chan string, maxConcurrentDownloads)
	t := make(chan string, maxConcurrentDownloads)
	s := make(chan bool, 1)
	return &app{
		videos:    v,
		toConvert: t,
		started:   s,
	}
}

func (a *app) readVideos() {
	buf := bufio.NewReader(os.Stdin)
	started := false
	for {
		fmt.Print("[mp3] Enter a URL: ")
		line, err := buf.ReadString('\n')
		if strings.TrimSpace(line) == "" {
			break
		} else if err != nil {
			fmt.Println("[mp3] Invalid URL")
		} else {
			a.videos <- strings.TrimSpace(line)
			go a.downloadVideo()
			a.wg.Add(1)
			if !started {
				a.started <- true
				started = true
			}
		}
	}
	close(a.videos)
}

func (a *app) downloadVideo() {
	defer a.wg.Done()

	url, open := <-a.videos
	if !open {
		fmt.Println("Channel closed")
		close(a.toConvert)
		return
	}

	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		fmt.Println("Failed to get video info")
		return
	}

	best := vid.Formats.Best(ytdl.FormatAudioBitrateKey)[0]

	cleanName := cleanTitle(&vid.Title)
	filename := cleanName + "." + best.Extension
	file, _ := os.Create(filename)
	defer file.Close()

	fmt.Println("Starting download")
	vid.Download(best, file)
	fmt.Println("Download finished")
	a.toConvert <- filename
	go a.convertVideo(cleanName)
	a.wg.Add(1)
}

func cleanTitle(s *string) string {
	var newTitle strings.Builder
	for _, r := range *s {
		switch {
		case unicode.IsSpace(r):
			newTitle.WriteRune('_')
		case r == '-' || r == '(' || r == ')':
			newTitle.WriteRune(r) // nice punctuation is ok
		case unicode.IsPunct(r):
			break // skip scary punctuation like quotes
		default:
			newTitle.WriteRune(r)
		}
	}
	return newTitle.String()
}

func (a *app) convertVideo(vidName string) {
	defer a.wg.Done()

	filename, open := <-a.toConvert
	if !open {
		fmt.Println("Channel closed")
		return
	}

	newName := vidName + ".mp3"
	fmt.Println(filename, newName)
	ffmpeg := exec.Command("ffmpeg", "-vn", "-i", filename, newName)
	fmt.Println(ffmpeg)
	fmt.Println("Starting conversion")
	out, err := ffmpeg.Output()
	fmt.Println("Conversion finished")

	fmt.Println("[mp3] Finished:", string(out), err, &a.wg)
}

func main() {
	fmt.Println("MP3 to Go v0.01")
	a := new()

	go a.readVideos()
	<-a.started // wait for first video to be entered
	// otherwise we cant wait on the waitgroup
	// because its initialized to zero, and
	// main will instantly exit

	a.wg.Wait()
}
