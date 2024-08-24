package conversion

import (
	"encoding/csv"
	"errors"
	"io"

	"github.com/xuri/excelize/v2"
)

type CsvToXlsxConverter struct{}

func (c *CsvToXlsxConverter) Convert(input io.Reader, output io.Writer) error {
	// Cria um novo arquivo Excel
	f := excelize.NewFile()
	defer f.Close()

	// Lê o CSV do input
	reader := csv.NewReader(input)

	// Adiciona as linhas do CSV no sheet
	sheet := "Sheet1"
	f.NewSheet(sheet)

	rowIndex := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("erro ao ler o CSV: " + err.Error())
		}

		// Escreve cada célula no arquivo Excel
		for colIndex, cellValue := range record {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex)
			if err != nil {
				return errors.New("erro ao converter coordenadas: " + err.Error())
			}
			f.SetCellValue(sheet, cell, cellValue)
		}
		rowIndex++
	}

	// Salva o arquivo Excel no output
	if err := f.Write(output); err != nil {
		return errors.New("erro ao salvar o arquivo Excel: " + err.Error())
	}

	return nil
}
