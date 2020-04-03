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

	"github.com/rylio/ytdl"
)

const maxConcurrentDownloads = 5

type app struct {
	videos    chan string
	toConvert chan string
	started   chan bool
	ffmpeg    *exec.Cmd
	wg        sync.WaitGroup
}

func new() *app {
	v := make(chan string, maxConcurrentDownloads)
	t := make(chan string, maxConcurrentDownloads)
	s := make(chan bool, 1)
	f := exec.Command("ffmpeg", "-vn", "-i")
	return &app{
		videos:    v,
		toConvert: t,
		started:   s,
		ffmpeg:    f,
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

	filename := vid.Title + "." + best.Extension
	file, _ := os.Create(filename)
	defer file.Close()

	vid.Download(best, file)
	a.toConvert <- filename
	go a.convertVideo(vid.Title)
	a.wg.Add(1)
}

func (a *app) convertVideo(vidName string) {
	defer a.wg.Done()

	filename, open := <-a.toConvert
	if !open {
		fmt.Println("Channel closed")
		return
	}

	newName := vidName + ".mp3"
	a.ffmpeg.Args = append(a.ffmpeg.Args, filename, newName)

	out, err := a.ffmpeg.Output()
	fmt.Println("[mp3] Finished:", string(out), err)
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
