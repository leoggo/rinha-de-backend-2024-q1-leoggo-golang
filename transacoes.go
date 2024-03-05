package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type dadosTransacao struct {
	Tipo        string `json:"tipo"`
	Descricao   string `json:"descricao"`
	RealizadaEm string `json:"realizada_em"`
	Valor       int64  `json:"valor"`
}

type dadosdaConta struct {
	Limite int64 `json:"limite"`
	Saldo  int64 `json:"saldo"`
}

func doTransacao(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	clientID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	transacao := dadosTransacao{}

	err = json.NewDecoder(r.Body).Decode(&transacao)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	if len(transacao.Descricao) == 0 || len(transacao.Descricao) > 10 {
		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	resposta := dadosdaConta{}

	row := db.QueryRow(context.Background(), "SELECT limite, saldo from clientes WHERE id=$1;", clientID)

	if err := row.Scan(&resposta.Limite, &resposta.Saldo); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	if transacao.Tipo == "c" {
		resposta.Saldo += transacao.Valor

		_, err := db.Exec(context.Background(), "UPDATE clientes SET saldo=$1 WHERE id=$2;", resposta.Saldo, clientID)

		if err != nil {
			panic(err)
		}

	} else if transacao.Tipo == "d" {
		resposta.Saldo -= transacao.Valor

		if resposta.Saldo < (-1 * resposta.Limite) {
			w.WriteHeader(http.StatusUnprocessableEntity)

			return
		}

		_, err := db.Exec(context.Background(), "UPDATE clientes SET saldo=$1 WHERE id=$2;", resposta.Saldo, clientID)

		if err != nil {
			panic(err)
		}

	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	_, err = db.Exec(context.Background(), "INSERT INTO transacoes (valor, tipo, descricao, id) VALUES ($1, $2, $3, $4);", transacao.Valor, transacao.Tipo, transacao.Descricao, clientID)

	if err != nil {
		panic(err)
	}

	retval := strings.Builder{}

	retval.WriteString("{\"limite\":")

	retval.WriteString(strconv.FormatInt(resposta.Limite, 10))

	retval.WriteString(",\"saldo\":")

	retval.WriteString(strconv.FormatInt(resposta.Saldo, 10))

	retval.WriteString("}")

	w.Write([]byte(retval.String()))
}
