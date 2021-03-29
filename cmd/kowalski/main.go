package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/csmith/kowalski/v3"
)

var (
	token       = flag.String("token", "", "Discord bot token")
	goodModel   = flag.String("good-model", "models/combined.wl", "Path of the 'good' model")
	backupModel = flag.String("backup-model", "models/urbandictionary.wl", "Path of the 'backup' model")
	prefix      = flag.String("prefix", "!", "Character(s) to require before commands")

	checkers []*kowalski.SpellChecker
)

func main() {
	flag.Parse()

	checkers = []*kowalski.SpellChecker{
		loadModel(*goodModel),
		loadModel(*backupModel),
	}

	dg, err := discordgo.New(fmt.Sprintf("Bot %s", *token))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(handleMessage)

	if err := dg.Open(); err != nil {
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

func loadModel(path string) (res *kowalski.SpellChecker) {
	f, err := os.Open(path)
	if err != nil {
		log.Panicf("Failed to open model: %v", err)
	}
	defer f.Close()

	res, err = kowalski.LoadSpellChecker(f)
	if err != nil {
		log.Panicf("Failed to load model: %v", err)
	}
	return res
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	line := m.Content
	command, arguments, ok := parseCommand(line)
	if !ok {
		return
	}

	if c, ok := textCommands[command]; ok {
		c(arguments, func(format string, a ...interface{}) {
			sendMessage(s, m, fmt.Sprintf(format, a...))
		})
	}

	if c, ok := fileCommands[command]; ok {
		var urls []string
		for i := range m.Attachments {
			urls = append(urls, m.Attachments[i].URL)
		}

		if len(urls) > 0 {
			c(arguments, urls, func(format string, a ...interface{}) {
				sendMessage(s, m, fmt.Sprintf(format, a...))
			})
		}
	}
}

func parseCommand(input string) (command, arguments string, ok bool) {
	parts := strings.SplitN(input, " ", 2)
	command = strings.ToLower(parts[0])
	if !strings.HasPrefix(command, *prefix) {
		return
	}

	command = strings.TrimPrefix(command, *prefix)
	if len(parts) > 1 {
		arguments = parts[1]
	}
	ok = true
	return
}

func subtract(input, exclusions []string) []string {
	var res []string
	for i := range input {
		excluded := false
		for j := range exclusions {
			if exclusions[j] == input[i] {
				excluded = true
				break
			}
		}
		if !excluded {
			res = append(res, input[i])
		}
	}
	return res
}

func merge(words [][]string) []string {
	var res []string
	for i := range words {
		for j := range words[i] {
			if i > 0 {
				res = append(res, fmt.Sprintf("_%s_", words[i][j]))
			} else {
				res = append(res, fmt.Sprintf("**%s**", words[i][j]))
			}
		}
	}
	sort.Strings(res)
	return res
}

func countReps(input []string) []string {
	sort.Strings(input)

	var res []string
	var last string
	var count int
	for i := range input {
		if input[i] == last {
			count++
			continue
		} else if count > 1 {
			res = append(res, fmt.Sprintf("**%s × %d**", last, count))
			count = 0
		} else if count == 1 {
			res = append(res, last)
			count = 0
		}
		last = input[i]
		count = 1
	}

	if count > 1 {
		res = append(res, fmt.Sprintf("%s × %d", last, count))
	} else if count == 1 {
		res = append(res, last)
	}
	return res
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

func isValidT9(word string) bool {
	if len(word) == 0 {
		return false
	}

	for _, r := range word {
		if r < '2' || r > '9' {
			return false
		}
	}
	return true
}
