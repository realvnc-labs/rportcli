package controllers

type PromptReaderMock struct {
	ReadCount           int
	PasswordReadCount   int
	ReadOutputs         []string
	PasswordReadOutputs []string
	ErrToGive           error
}

func (prm *PromptReaderMock) ReadString() (string, error) {
	prm.ReadCount++

	if len(prm.ReadOutputs) < prm.ReadCount {
		return "", prm.ErrToGive
	}

	return prm.ReadOutputs[prm.ReadCount-1], prm.ErrToGive
}

func (prm *PromptReaderMock) ReadPassword() (string, error) {
	prm.PasswordReadCount++

	if len(prm.PasswordReadOutputs) < prm.PasswordReadCount {
		return "", prm.ErrToGive
	}

	return prm.PasswordReadOutputs[prm.PasswordReadCount-1], prm.ErrToGive
}
