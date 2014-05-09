package web

import ()

var errorChan chan error

func Initialize(errChan chan error) (err error) {
	errorChan = errChan
	return
}

func Home() string {
	return "Hello world!"
}
