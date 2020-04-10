package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/csmith/kowalski"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	token = flag.String("token", "", "Discord bot token")
	words *kowalski.Node
)

func main()  {
	flag.Parse()

	var err error
	words, err = kowalski.LoadWords()
	if err != nil {
		log.Panicf("Failed to load words: %v\n", err)
	}

	dg, err := discordgo.New("Bot " + *token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	line := strings.ToLower(m.Content)

	if strings.HasPrefix(line, "anagram") {
		go func() {
			word := strings.TrimSpace(strings.TrimPrefix(line, "anagram"))
			if isValidWord(word) {
				res := words.Anagrams(word)
				sendMessage(s, m, fmt.Sprintf("Anagrams for %s: %v", word, res))
			} else {
				sendMessage(s, m, fmt.Sprintf("Invalid word: %s", word))
			}
		}()
	}

	if strings.HasPrefix(line, "match") {
		go func() {
			word := strings.TrimSpace(strings.TrimPrefix(line, "match"))
			if isValidWord(word) {
				res := words.Match(word)
				sendMessage(s, m, fmt.Sprintf("Matches for %s: %v", word, res))
			} else {
				sendMessage(s, m, fmt.Sprintf("Invalid word: %s", word))
			}
		}()
	}
}

func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate, text string) {
	if len(text) > 2000 {
		text = fmt.Sprintf("%s <truncated>", text[0:1988])
	}
	if _, err := s.ChannelMessageSend(m.ChannelID, text); err != nil {
		log.Printf("Failed to send message: %v\n", err)
	}
}

func isValidWord(word string) bool {
	if len(word) == 0 {
		return false
	}

	for _, r := range word {
		if (r < 'a' || r > 'z') && r != '?' {
			return false
		}
	}
	return true
}
