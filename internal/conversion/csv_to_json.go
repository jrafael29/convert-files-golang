package conversion

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
)

type CsvToJsonConverter struct{}

func (c *CsvToJsonConverter) Convert(input io.Reader, output io.Writer) error {
	reader := csv.NewReader(input)

	// Lê o cabeçalho (primeira linha do CSV)
	headers, err := reader.Read()
	if err != nil {
		return errors.New("erro ao ler o cabeçalho do CSV: " + err.Error())
	}

	// Lista para armazenar as linhas como mapas
	var records []map[string]string

	// Lê cada linha do CSV
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("erro ao ler uma linha do CSV: " + err.Error())
		}

		// Cria um mapa para a linha atual
		row := make(map[string]string)
		for i, value := range record {
			row[headers[i]] = value
		}

		records = append(records, row)
	}

	// Converte os registros para JSON
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return errors.New("erro ao converter para JSON: " + err.Error())
	}

	// Escreve o JSON no output
	_, err = output.Write(jsonData)
	if err != nil {
		return errors.New("erro ao escrever o JSON no output: " + err.Error())
	}

	return nil
}
