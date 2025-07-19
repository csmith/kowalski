package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/csmith/kowalski/v6"
)

//go:embed static/*
var staticFiles embed.FS

var (
	port        = flag.Int("port", 8080, "HTTP port to listen on")
	goodModel   = flag.String("good-model", "models/combined.wl", "Path of the 'good' model")
	backupModel = flag.String("backup-model", "models/urbandictionary.wl", "Path of the 'backup' model")
	fstModel    = flag.String("fst-model", "", "Path to FST for fast word operations")

	checkers []*kowalski.SpellChecker
)

type Request struct {
	Command string `json:"command"`
	Input   string `json:"input"`
}

type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func init() {
	flag.Parse()
}

func main() {
	checkers = []*kowalski.SpellChecker{
		loadModel(*goodModel),
		loadModel(*backupModel),
	}

	if *fstModel != "" {
		initFST(*fstModel)
	}

	// Create a sub-filesystem that strips the "static" prefix
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to create static filesystem: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	mux.HandleFunc("/api/command", handleCommand)
	mux.HandleFunc("/api/image", handleImageCommand)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting web server on port %d", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}

func loadModel(path string) *kowalski.SpellChecker {
	f, err := os.Open(path)
	if err != nil {
		log.Panicf("Failed to open model: %v", err)
	}
	defer f.Close()

	res, err := kowalski.LoadSpellChecker(f)
	if err != nil {
		log.Panicf("Failed to load model: %v", err)
	}
	return res
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, Response{Success: false, Error: "Invalid request"})
		return
	}

	result, err := processCommand(req.Command, req.Input)
	if err != nil {
		writeJSON(w, Response{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, Response{Success: true, Result: result})
}

func handleImageCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		writeJSON(w, Response{Success: false, Error: "Failed to parse form"})
		return
	}

	command := r.FormValue("command")
	file, _, err := r.FormFile("image")
	if err != nil {
		writeJSON(w, Response{Success: false, Error: "No image provided"})
		return
	}
	defer file.Close()

	result, err := processImageCommand(command, file)
	if err != nil {
		writeJSON(w, Response{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, Response{Success: true, Result: result})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func processCommand(command, input string) (interface{}, error) {
	switch command {
	case "anagram":
		return processAnagram(input)
	case "analysis":
		return processAnalysis(input)
	case "chunk":
		return processChunk(input)
	case "letters":
		return processLetters(input)
	case "match":
		return processMatch(input)
	case "morse":
		return processMorse(input)
	case "multianagram":
		return processMultiAnagram(input)
	case "multimatch":
		return processMultiMatch(input)
	case "offbyone":
		return processOffByOne(input)
	case "shift":
		return processShift(input)
	case "t9":
		return processT9(input)
	case "transpose":
		return processTranspose(input)
	case "wordsearch":
		return processWordSearch(input)
	case "firstletters":
		return processFirstLetters(input)
	case "reverse":
		return processReverse(input)
	case "checkwords":
		return processCheckWords(input)
	case "fstanagram":
		if fstTransducer != nil {
			return processFstAnagram(input)
		}
		return nil, fmt.Errorf("FST model not loaded")
	case "fstregex":
		if fstTransducer != nil {
			return processFstRegex(input)
		}
		return nil, fmt.Errorf("FST model not loaded")
	case "fstmorse":
		if fstTransducer != nil {
			return processFstMorse(input)
		}
		return nil, fmt.Errorf("FST model not loaded")
	case "wordlink":
		if fstTransducer != nil {
			return processWordLink(input)
		}
		return nil, fmt.Errorf("FST model not loaded")
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}

func processImageCommand(command string, file io.Reader) (interface{}, error) {
	switch command {
	case "colours", "colors":
		return processColours(file)
	case "hidden":
		return processHiddenPixels(file)
	case "rgb":
		return processRGB(file)
	default:
		return nil, fmt.Errorf("unknown image command: %s", command)
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

func merge(words [][]string) []string {
	var res []string
	for i := range words {
		for j := range words[i] {
			if i > 0 {
				res = append(res, fmt.Sprintf("_%s_", words[i][j]))
			} else {
				res = append(res, words[i][j])
			}
		}
	}
	return res
}
