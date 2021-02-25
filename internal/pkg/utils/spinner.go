package utils

import (
	bspin "github.com/leaanthony/spinner"
)

type Spinner struct {
	baseSpinner *bspin.Spinner
}

func NewSpinner() *Spinner {
	bspinner := bspin.New()
	bspinner.SetSpinFrames([]string{`\`, `-`, `/`, `-`})
	bspinner.SetAbortMessage(InterruptMessage)
	return &Spinner{
		baseSpinner: bspinner,
	}
}

func (s *Spinner) Start(msg string) {
	s.baseSpinner.Start(msg)
}

func (s *Spinner) Update(msg string) {
	s.baseSpinner.UpdateMessage(msg)
}

func (s *Spinner) Stop(msg string) {
	s.baseSpinner.Success(msg)
}
