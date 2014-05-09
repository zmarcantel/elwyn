package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/zmarcantel/elwyn/logging"
	"github.com/zmarcantel/elwyn/routes"
)

const (
	EXIT_NORMAL      int = 0
	EXIT_CONFIG_ERR  int = 100
	EXIT_STARTUP_ERR int = 110
	EXIT_GENERAL_ERR int = 200
)

//
// Global variables
//
var logger *logging.Router

//
// Main entry point
//    Handles CLI loading
//
func main() {
	// initialize the CLI options
	// returns any errors to be caugh by the checkError function
	checkError(cliInit(), "There was a command-line error", EXIT_CONFIG_ERR)

	// setup watchers for termination signals
	var _, lock = initSignalWatchers()

	// start logging how that the system is initialized
	// currently have no use for the actual logger object
	logger, err := logging.Initialize(Opts.LogPath, Opts.Quiet)
	checkError(err, "Could not load log file", EXIT_STARTUP_ERR)

	//
	logger.Banner("Starting Server")
	checkError(routes.Initialize(lock, logger, Opts.Port), "Failed to start HTTP Server", EXIT_STARTUP_ERR)
}

//
// Check provided error, and if exists print output and exit with code
//
func checkError(actual error, text string, code int) {
	if actual != nil {
		fmt.Printf("%s:\n%s\n", text, actual)
		Cleanup()
		os.Exit(code)
	}
}

//
// Create a channel to listen to OS signals on
//    Waits for SIGKILL or SIGTERM and handles accordingly
//
func initSignalWatchers() (sigLock chan os.Signal, deadLock chan error) {
	// make the signal channel
	sigLock = make(chan os.Signal, 10)
	signal.Notify(sigLock, os.Interrupt)

	// make the global error channel
	deadLock = make(chan error, 10)

	// this functions spins up in a separate context
	// blocks at the first line until a signal is passed on channel
	go func() {
		// we don't care about value here, just that a signal was caught
		<-sigLock
		fmt.Println("Exhibit received termination signal.")
		Cleanup()
		os.Exit(EXIT_NORMAL)
	}()

	// this functions spins up in a separate context
	// blocks at first line of every iteration of for loop until an error is passed
	// the for loop is included in the case on non-fatal errors -- keeps the system going
	go func() {
		for {
			err := <-deadLock
			checkError(err, "Received error from server", EXIT_GENERAL_ERR)
		}
	}()

	return sigLock, deadLock
}

//
// This function will clean up any open files, sockets, or other
// resources still open
//
func Cleanup() {
	logging.Close()
}
