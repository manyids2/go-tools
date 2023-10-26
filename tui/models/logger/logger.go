package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Logger data
type Logger struct {
	Datadir  string
	FileExt  string
	Loaded   chan bool
	LogFiles []string
}

// From args
func New(datadir, fileext string) *Logger {
	m := Logger{
		Datadir: datadir,
		FileExt: fileext,
		Loaded:  make(chan bool),
	}
	return &m
}

// Defaults
func Default() *Logger {
	return New("./.log", ".log")
}

// Print
func (m Logger) String() string {
	return fmt.Sprintf(
		`Logger:
	Datadir: %s
	FileExt: %s
	  LogFiles: %v`, m.Datadir, m.FileExt, m.LogFiles)
}

func (m *Logger) SetLogFiles() {
	// Check if logdir exists, else return without error
	entries, err := os.ReadDir(m.Datadir)
	if err != nil {
		log.Println("Could not read datadir: ", m.Datadir, err)
		return
	}

	// Iterate over log directory and append to LogFiles
	for _, e := range entries {
		fmt.Println(e.Name())
		if filepath.Ext(e.Name()) == m.FileExt {
			m.LogFiles = append(m.LogFiles, e.Name())
		}
	}

	// Inform that load is finished
	m.Loaded <- true
}
