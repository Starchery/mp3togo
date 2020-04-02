# mp3togo
<p>Bulk downloader and converter for Youtube videos. Written in Go.</p>

---
This program is intended for beatmakers who often need to download many youtube videos at a time
and convert them all to a suitable audio format for easy sampling.

This is a rewrite of my sample downloading program [flvto](https://github.com/Starchery/flvto). 
I've been working on it in a local repo, and was increasingly unsatisfied with the speed (or lack thereof), 
general inefficiency, and all the dependencies. So I'm starting over, in Go this time, because why not.

This is a work in progress. It's currently based on the [ytdl](https://github.com/rylio/ytdl) library and calls
on FFmpeg to do the conversion. Until I can figure out how to do it myself, that means you'll need to
have FFmpeg installed to use this. Why wouldn't you have FFmpeg, anyway? It's a great program.

## Installation
TBA. For now, see [Building.](https://github.com/Starchery/mp3togo#building)

## Building
Requirements: Go 1.14. 1.13 probably works too.

```
$ git clone https://github.com/Starchery/mp3togo
$ cd mp3togo
$ go run cmd/mp3togo/main.go
```

## What's with the name?
I don't know. It's in Go, so it's *to go?* Or something? Debate your auntie.

## To-Do
- [x] Download videos from Youtube
- [x] Concurrent downloads
- [ ] Convert downloaded videos to .mp3
- [ ] Restrict filenames to ASCII characters and no whitespace
- [ ] Better name
- [ ] Probably have to write my own youtube downloader library
- [ ] GUI :)
