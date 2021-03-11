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

func (s *Spinner) StopSuccess(msg string) {
	s.baseSpinner.Success(msg)
}

func (s *Spinner) StopError(msg string) {
	s.baseSpinner.Error(msg)
}

type NullSpinner struct {
}

func (ns *NullSpinner) Start(msg string) {
}

func (ns *NullSpinner) Update(msg string) {
}

func (ns *NullSpinner) StopSuccess(msg string) {
}

func (ns *NullSpinner) StopError(msg string) {
}
