package output

import (
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type TotPSecretRenderer struct {
	ColCountCalculator CalcTerminalColumnsCount
	Writer             io.Writer
	Format             string
}

func (cr *TotPSecretRenderer) RenderTotPSecret(key *models.TotPSecretOutput) error {
	return RenderByFormat(
		cr.Format,
		cr.Writer,
		key,
		func() error {
			return cr.renderTotPSecretToHumanFormat(key)
		},
	)
}

func (cr *TotPSecretRenderer) renderTotPSecretToHumanFormat(key *models.TotPSecretOutput) error {
	if key == nil {
		return nil
	}
	err := RenderHeader(cr.Writer, "One time password secret\n")
	if err != nil {
		return err
	}

	err = RenderHeader(cr.Writer, key.Comment)
	if err != nil {
		return err
	}

	rows := []RowData{
		key,
	}

	err = RenderTable(cr.Writer, &models.TotPSecretOutput{}, rows, cr.ColCountCalculator)
	if err != nil {
		return err
	}

	return nil
}
