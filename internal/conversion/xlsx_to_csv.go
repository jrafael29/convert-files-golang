package conversion

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/tealeg/xlsx/v3"
)

type XlsxToCsvConverter struct{}

func (c *XlsxToCsvConverter) Convert(input io.Reader, output io.Writer) error {
	// Lê todo o conteúdo do input (io.Reader) em memória
	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("erro ao ler o input: %v", err)
	}

	// Abre o arquivo XLSX a partir do buffer em memória
	xlFile, err := xlsx.OpenBinary(data)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo XLSX: %v", err)
	}

	// Verifica se o arquivo XLSX contém planilhas
	sheetLen := len(xlFile.Sheets)
	if sheetLen == 0 {
		return errors.New("este arquivo XLSX não contém nenhuma planilha")
	}

	// Seleciona a primeira planilha para conversão
	sheet := xlFile.Sheets[0]

	// Cria um writer CSV
	cw := csv.NewWriter(output)

	// Itera sobre cada linha da planilha e escreve no CSV
	var vals []string
	err = sheet.ForEachRow(func(row *xlsx.Row) error {
		if row != nil {
			vals = vals[:0] // Reseta o slice para a nova linha
			err := row.ForEachCell(func(cell *xlsx.Cell) error {
				str, err := cell.FormattedValue()
				if err != nil {
					return err
				}
				vals = append(vals, str)
				return nil
			})
			if err != nil {
				return err
			}
		}
		cw.Write(vals)
		return nil
	})

	if err != nil {
		return fmt.Errorf("erro ao processar o arquivo XLSX: %v", err)
	}

	// Finaliza a escrita no CSV
	cw.Flush()

	return cw.Error()
}
