package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func init() {
	guildNames = make(map[string]string)
}

var (
	guildNames map[string]string
	closeBot   func() error
)

func runBot() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating discord bot: ", err)
		return
	}
	closeBot = dg.Close

	dg.AddHandler(guildCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}
	defer dg.Close()

	// Give enough time for guilds to be added
	time.Sleep(time.Second)
	for {
		fmt.Println("Getting Emojis")

		for ID, name := range guildNames {
			fmt.Println("[" + name + "]")

			guild, err := dg.Guild(ID)
			if err != nil {
				fmt.Printf("Error getting guild: %v\n", err)
				delete(guildNames, ID)
				break
			}

			if len(guild.Emojis) > 0 {
				err = getEmojis(guild.Name, guild.Emojis)
				if err != nil {
					fmt.Printf("Could not get emojis for guild \"%s\": %v\n", guild.Name, err)
				}
				fmt.Printf("Downloaded %d emojis\n", len(guild.Emojis))
			} else {
				fmt.Println("No emojis for guild")
			}
		}

		time.Sleep(sleepTime)
	}
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	fmt.Printf("Adding guild \"%s\" to current names\n", event.Guild.Name)
	guildNames[event.Guild.ID] = event.Guild.Name
}

func getEmojis(server string, emojis []*discordgo.Emoji) error {
	serverDir := path.Join(dir, server)
	if _, err := os.Stat(serverDir); os.IsNotExist(err) {
		err := os.MkdirAll(serverDir, os.ModePerm)
		if err != nil {
			return errors.Errorf("could not make server directory: %v", err)
		}
	}

	for _, e := range emojis {
		path := path.Join(serverDir, e.Name)
		if e.Animated {
			path += ".gif"
		} else {
			path += ".png"
		}

		err := downloadEmoji(path, e.ID)
		if err != nil {
			return errors.Errorf("could not get emoji \"%s\": %v\n", e.Name, err)
		}
	}

	return nil
}

func downloadEmoji(path, id string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Errorf("could not create error: %v", err)
	}
	defer file.Close()

	response, err := http.Get(emojiURL + id)
	if err != nil {
		return errors.Errorf("could not get emoji: %v", err)
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return errors.Errorf("could not save emoji: %v", err)
	}

	return nil
}
