package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/csmith/kowalski/v3"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strings"
	"syscall"
)

var (
	token       = flag.String("token", "", "Discord bot token")
	goodModel   = flag.String("good-model", "models/combined.wl", "Path of the 'good' model")
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

	if strings.HasPrefix(line, "analysis") {
		go func() {
			res := kowalski.Analyse(checkers[0], strings.TrimSpace(strings.TrimPrefix(line, "analysis")))
			if len(res) == 0 {
				sendMessage(s, m, "Analysis: nothing interesting found")
			}
			sendMessage(s, m, fmt.Sprintf("Analysis:\n\t%s", strings.Join(res, "\n\t")))
		}()
	}

	if strings.HasPrefix(line, "letters") {
		go func() {
			res := kowalski.LetterDistribution(strings.TrimSpace(strings.TrimPrefix(line, "letters")))
			max := 0
			for i := range res {
				if res[i] > max {
					max = res[i]
				}
			}
			const targetWidth = 20
			message := strings.Builder{}
			message.WriteString("Letter distribution:\n```")
			for i := range res {
				message.WriteByte(byte(i + 'A'))
				message.WriteString(": ")
				for j := 0; j < int(targetWidth * (float64(res[i]) / float64(max))); j++ {
					message.WriteRune('█')
				}
				message.WriteString(fmt.Sprintf(" %d\n", res[i]))
			}
			message.WriteString("```")
			sendMessage(s, m, message.String())
		}()
	}

	if strings.HasPrefix(line, "wordsearch") {
		go func() {
			input := strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "wordsearch")), "\n")
			res := kowalski.MultiplexWordSearch(checkers, input)
			sendMessage(s, m, fmt.Sprintf(
				"Words found:\n\nNormal: %s\n\nUD: %s",
				strings.Join(countReps(res[0]), ", "),
				strings.Join(countReps(subtract(res[1], res[0])), ", ")))
		}()
	}

	if strings.HasPrefix(line, "shift") {
		go func() {
			input := strings.TrimSpace(strings.TrimPrefix(line, "shift"))
			res := kowalski.CaesarShifts(input)
			out := strings.Builder{}
			out.WriteString( "Caesar shifts:\n")
			for i, s := range res {
				score := kowalski.Score(checkers[0], s)
				if score > 0.5 {
					s = fmt.Sprintf("**%s**", s)
				}
				out.WriteString(fmt.Sprintf("\t%2d: %s\n", i + 1, s))
			}
			sendMessage(s, m, out.String())
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
