package handler

import (
	"convert-files/internal/conversion"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type ConversionHandler struct {
	converters map[string]conversion.Converter
}

func NewConversionHandler() *ConversionHandler {
	return &ConversionHandler{
		converters: map[string]conversion.Converter{
			"csv-to-xlsx": &conversion.CsvToXlsxConverter{},
			"csv-to-json": &conversion.CsvToJsonConverter{},
			"xlsx-to-csv": &conversion.XlsxToCsvConverter{},
		},
	}
}

func (h *ConversionHandler) Convert(w http.ResponseWriter, r *http.Request) {
	// Verifica o método HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtém o tipo de conversão da query string
	conversionType := r.URL.Query().Get("type")
	if conversionType == "" {
		http.Error(w, "Tipo de conversão não especificado", http.StatusBadRequest)
		return
	}

	// Verifica se o tipo de conversão é suportado
	converter, exists := h.converters[conversionType]
	if !exists {
		http.Error(w, "Tipo de conversão não suportado", http.StatusBadRequest)
		return
	}

	// Lê o arquivo de entrada
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao ler o arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Cria um arquivo temporário para armazenar o resultado da conversão
	tempFile, err := os.CreateTemp("", "converted-*")
	if err != nil {
		http.Error(w, "Erro ao criar arquivo temporário", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name()) // Remove o arquivo temporário ao final
	defer tempFile.Close()

	// Realiza a conversão
	err = converter.Convert(file, tempFile)
	if err != nil {
		http.Error(w, "Erro na conversão: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Define o nome do arquivo de saída
	outputFileName := filepath.Base(tempFile.Name()) + ".converted"

	// Define o cabeçalho para download do arquivo convertido
	w.Header().Set("Content-Disposition", "attachment; filename="+outputFileName)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Envia o arquivo convertido para o cliente
	tempFile.Seek(0, io.SeekStart)
	if _, err := io.Copy(w, tempFile); err != nil {
		http.Error(w, "Erro ao enviar arquivo convertido", http.StatusInternalServerError)
		return
	}

	log.Printf("Conversão %s realizada com sucesso", conversionType)
}
