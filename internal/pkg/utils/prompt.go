package utils

import (
	"fmt"
	"io"
	"os"
)

type Scanner interface {
	Scan() bool
	Text() string
	Err() error
}

type PasswordScanner func() ([]byte, error)

type PromptReader struct {
	Sc              Scanner
	SigChan         chan os.Signal
	PasswordScanner PasswordScanner
}

func (pr *PromptReader) ReadString() (string, error) {
	return pr.promptForStrValue(func() (string, error) {
		for pr.Sc.Scan() {
			return pr.Sc.Text(), nil
		}
		err := pr.Sc.Err()
		if err != nil {
			return "", err
		}

		return "", io.EOF
	})
}

func (pr *PromptReader) ReadPassword() (string, error) {
	return pr.promptForStrValue(func() (string, error) {
		inputBytes, err := pr.PasswordScanner()
		return string(inputBytes), err
	})
}

func (pr *PromptReader) Output(text string) {
	fmt.Print(text)
}

func (pr *PromptReader) promptForStrValue(reader func() (string, error)) (string, error) {
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
	case <-pr.SigChan:
		return "", io.EOF
	case msg := <-msgChan:
		return msg, nil
	case err := <-errChan:
		return "", err
	}
}
