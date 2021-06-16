package main

import (
	"encoding/json"
	"io"
	"log"
	"path/filepath"
	"runtime"

	"os"
)

type Configuration struct {
	DefDir      string
	TemplateDir string
	DefFilename string
	LogFinename string
}

var config Configuration

var logger *log.Logger

func init() {
	loadConfig()
	file, err := os.OpenFile(config.LogFinename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func loadConfig() {
	// ホームディレクトリを取得
	var configDir string
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		configDir = os.Getenv("APPDATA")
	} else {
		configDir = filepath.Join(home, ".config")
	}

	filename := filepath.Join(configDir, "goinit", "config.json")
	file, err := os.Open(filename)
	if err != nil {
		error("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		error("Cannot get configuration from file", err)
	}
}

// for logging
func info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func error(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}

// version
func version() string {
	return "0.1"
}

// ファイルをコピー
func CopyFile(fromFile string, toFile string) {
	f, err := os.Open(fromFile)
	if err != nil {
		log.Fatal(err)
	}

	t, err := os.Create(toFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(t, f)
	if err != nil {
		log.Fatal(err)
	}
}

// ホームディレクトリを取得
func GetHomeDir() string {
	var homeDir string
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		homeDir = os.Getenv("HOMEPATH")
	} else {
		homeDir = home
	}
	return homeDir
}
