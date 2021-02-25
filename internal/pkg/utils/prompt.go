package utils

import (
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/chzyer/readline"
)

type PromptReader struct {
}

func (pr *PromptReader) ReadString() (string, error) {
	rl, err := readline.New("> ")
	if err != nil {
		return "", err
	}
	defer rl.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	msgChan := make(chan string, 1)
	errChan := make(chan error, 1)
	go func() {
		msg, err := rl.Readline()
		if err != nil {
			errChan <- err
			return
		}

		msgChan <- msg
	}()

	select {
	case <-sigs:
		return "", io.EOF
	case msg := <-msgChan:
		return msg, nil
	case err := <-errChan:
		return "", err
	}
}

func (pr *PromptReader) ReadPassword() (string, error) {
	inputBytes, err := ReadPassword()
	return string(inputBytes), err
}
