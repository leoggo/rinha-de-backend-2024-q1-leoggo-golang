package main

import "net/http"

func buildRoutes(serverMuxer *http.ServeMux) {
	serverMuxer.HandleFunc("POST /clientes/{id}/transacoes", doTransacao)
	serverMuxer.HandleFunc("GET /clientes/{id}/extrato", doExtrato)
}
