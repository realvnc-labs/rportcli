package utils

import (
	"bufio"
	"io"
	"os"
	"os/signal"
	"syscall"
)

type PromptReader struct {
}

func (pr *PromptReader) ReadString() (string, error) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	msgChan := make(chan string, 1)
	errChan := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msgChan <- scanner.Text()
			return
		}
		if err := scanner.Err(); err != nil {
			errChan <- err
		}
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
