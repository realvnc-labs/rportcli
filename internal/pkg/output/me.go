package output

import "io"

type MeRenderer struct {
	Writer io.Writer
	Format string
}

func (m *MeRenderer) RenderMe(os KvProvider) error {
	return RenderByFormat(
		m.Format,
		m.Writer,
		os,
		func() error {
			RenderKeyValues(m.Writer, os)
			return nil
		},
	)
}
