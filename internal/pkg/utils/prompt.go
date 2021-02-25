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
	return pr.promptForStrValue(func() (string, error) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			return scanner.Text(), nil
		}
		err := scanner.Err()
		if err != nil {
			return "", err
		}

		return "", io.EOF
	})
}

func (pr *PromptReader) ReadPassword() (string, error) {
	return pr.promptForStrValue(func() (string, error) {
		inputBytes, err := ReadPassword()
		return string(inputBytes), err
	})
}

func (pr *PromptReader) promptForStrValue(reader func() (string, error)) (string, error) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	msgChan := make(chan string, 1)
	errChan := make(chan error, 1)
	go func() {
		msg, err := reader()
		if err != nil {
			errChan <- err
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
