package main

import (
	"convert-files/internal/handler"
	"log"
	"net/http"
)

func main() {
	// Cria uma nova instância do handler de conversão
	conversionHandler := handler.NewConversionHandler()

	// Registra a rota de conversão
	http.HandleFunc("/convert", conversionHandler.Convert)

	// Define a porta e inicia o servidor
	port := ":8080"
	log.Printf("Servidor rodando na porta %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
