package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
	"github.com/fsnotify/fsnotify"
)

var runningCmd *exec.Cmd
var restartFlag int32
var isInitialized bool = false

func main() {
	execCmd        := flag.String("exec", "go run .", "Command to execute")
	watchDir       := flag.String("watch", ".", "Directory to watch for changes")
	excludeDirsRaw := flag.String("exclude", ".git,vendor,node_modules", "Comma-separated list of folders to exclude")
	debounceMs     := flag.Int("debounce", 300, "Debounce time in milliseconds")
	flag.Parse()

	excludeDirs    := strings.Split(*excludeDirsRaw, ",")
	debounce       := time.Duration(*debounceMs) * time.Millisecond
	watcher, err   := fsnotify.NewWatcher()
	
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		if runningCmd != nil {
			runningCmd.Process.Kill()
		}
		os.Exit(0)
	}()

	// Watch all subdirectories
	err = filepath.WalkDir(*watchDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		for _, exclude := range excludeDirs {
			if strings.Contains(path, exclude) {
				return nil
			}
		}
		if d.IsDir() {
			watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	restart := make(chan bool, 1)
	go debounceRestart(restart, debounce, func() {
		rebuildAndRun(*execCmd)
	})
	
	go rebuildAndRun(*execCmd)

	fmt.Println("Watching for changes in ", *watchDir)
	for {
		select {
		case event := <-watcher.Events:
			if event.Op & fsnotify.Write == fsnotify.Write && strings.HasSuffix(event.Name, ".go") {
				fmt.Println("Change detected in:", event.Name)
				restart <- true
			}
		case err := <-watcher.Errors:
			fmt.Println("Watcher error:", err)
		}
	}
}

func rebuildAndRun(command string) {
	if !atomic.CompareAndSwapInt32(&restartFlag, 0, 1) {
		return
	}
	defer atomic.StoreInt32(&restartFlag, 0)

	if runningCmd != nil && runningCmd.Process != nil {
		runningCmd.Process.Kill()
		runningCmd.Wait()
	}

	if isInitialized {
		fmt.Println("🔨 Rebuilding & restarting...")
	} else {
		isInitialized = true		
	}
	parts	  := strings.Split(command, " ")
	runningCmd = exec.Command(parts[0], parts[1:]...)
	
	stdoutPipe, _ := runningCmd.StdoutPipe()
	stderrPipe, _ := runningCmd.StderrPipe()
	
	go pipeWithPrefix(stdoutPipe, "->")
	go pipeWithPrefix(stderrPipe, "->")

	err := runningCmd.Start()
	if err != nil {
		fmt.Println("❌ Failed to start:", err)
	} else {
		fmt.Println("🚀 Running [", command ," | PID:", runningCmd.Process.Pid, "] \n ")
	}
}

func debounceRestart(trigger <-chan bool, delay time.Duration, fn func()) {
	var timer *time.Timer
	for range trigger {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(delay, fn)
	}
}

func pipeWithPrefix(pipe io.ReadCloser, prefix string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Printf("%s %s\n", prefix, scanner.Text())
	}
}