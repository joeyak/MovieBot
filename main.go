package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot token")
	flag.StringVar(&dir, "d", "output", "Output dir for emotes")
	flag.StringVar(&port, "p", "", "Port that server sharing emotes is on")
	flag.DurationVar(&sleepTime, "s", time.Minute, "Time bot sleeps before checking emotes again")
	flag.Parse()
}

var (
	token        string
	dir          string
	port         string
	sleepTime    time.Duration
	closeActions []func()
)

const emojiURL = "https://cdn.discordapp.com/emojis/"

func main() {
	if token == "" {
		fmt.Println("Token required.")
		return
	}

	if dir == "" {
		fmt.Println("Dir required")
		return
	}

	if port != "" && port[0] != ':' {
		port = ":" + port
	}

	fmt.Printf("Bot is set to sleep for %s between pulls\n", sleepTime)

	go runBot()

	// If no port provided, then do not start http server
	if port == "" {
		go serveEmojiFolder()
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	go forceClose()
	fmt.Println("Shutting down...")

	if closeBot != nil {
		err := closeBot()
		if err != nil {
			fmt.Printf("Error closing bot: %v\n", err)
		}
	}
}

func forceClose() {
	time.Sleep(time.Second * 20)
	fmt.Println("20s has passed, forcing closure")
	os.Exit(0)
}

func serveEmojiFolder() {
	fmt.Printf("Starting file server at port %s\n", port)
	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.ListenAndServe(port, nil)
}
