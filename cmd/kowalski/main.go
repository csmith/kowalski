package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/csmith/kowalski/v2"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strings"
	"syscall"
)

var (
	token    = flag.String("token", "", "Discord bot token")
	goodModel = flag.String("good-model", "models/combined.wl", "Path of the 'good' model")
	backupModel = flag.String("backup-model", "models/urbandictionary.wl", "Path of the 'backup' model")

	checkers []*kowalski.SpellChecker
)

func main() {
	flag.Parse()

	checkers = []*kowalski.SpellChecker{
		loadModel(*goodModel),
		loadModel(*backupModel),
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	line := strings.ToLower(m.Content)

	if strings.HasPrefix(line, "anagram") {
		go func() {
			word := strings.TrimSpace(strings.TrimPrefix(line, "anagram"))
			if isValidWord(word) {
				res := merge(kowalski.MultiplexAnagram(checkers, word, kowalski.Dedupe))
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
				res := merge(kowalski.MultiplexMatch(checkers, word, kowalski.Dedupe))
				sendMessage(s, m, fmt.Sprintf("Matches for %s: %v", word, res))
			} else {
				sendMessage(s, m, fmt.Sprintf("Invalid word: %s", word))
			}
		}()
	}

	if strings.HasPrefix(line, "morse") {
		go func() {
			word := strings.TrimSpace(strings.TrimPrefix(line, "morse"))
			res := merge(kowalski.MultiplexFromMorse(checkers, word, kowalski.Dedupe))
			sendMessage(s, m, fmt.Sprintf("Matches for %s: %v", word, res))
		}()
	}

	if strings.HasPrefix(line, "t9") {
		go func() {
			word := strings.TrimSpace(strings.TrimPrefix(line, "t9"))
			if isValidT9(word) {
				res := merge(kowalski.MultiplexFromT9(checkers, word, kowalski.Dedupe))
				sendMessage(s, m, fmt.Sprintf("Matches for %s: %v", word, res))
			} else {
				sendMessage(s, m, fmt.Sprintf("Invalid word: %s", word))
			}
		}()
	}


	if strings.HasPrefix(line, "memstats") {
		go func() {
			mem := &runtime.MemStats{}
			runtime.ReadMemStats(mem)
			sendMessage(s, m, fmt.Sprintf("Memory usage: TotalAlloc=%d Sys=%d HeapSys=%d StackSys=%d", mem.TotalAlloc, mem.Sys, mem.HeapSys, mem.StackSys))
		}()
	}
}

func merge(words [][]string) []string {
	var res []string
	for i := range words {
		for j := range words[i] {
			if i > 0 {
				res = append(res, fmt.Sprintf("%sᵁᴰ", words[i][j]))
			} else {
				res = append(res, words[i][j])
			}
		}
	}
	sort.Strings(res)
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
