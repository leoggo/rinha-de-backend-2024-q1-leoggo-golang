package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type saldoDaConta struct {
	DataExtrato string `json:"data_extrato"`
	Limite      int64  `json:"limite"`
	Total       int64  `json:"total"`
}

func doExtrato(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	clientID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	saldoEmConta := saldoDaConta{
		DataExtrato: time.Now().Format(time.RFC3339),
	}

	row := db.QueryRow(context.Background(), "SELECT limite, saldo from clientes WHERE id=$1;", clientID)

	if err := row.Scan(&saldoEmConta.Limite, &saldoEmConta.Total); err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	transacoes := make([]dadosTransacao, 0, 10)

	rows, err := db.Query(context.Background(), "SELECT valor, tipo, descricao, TO_CHAR(momento::TIMESTAMP AT TIME ZONE 'UTC', 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"') FROM transacoes WHERE id=$1 ORDER BY momento DESC LIMIT 10;", clientID)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	defer rows.Close()

	for rows.Next() {
		dados := dadosTransacao{}

		err = rows.Scan(&dados.Valor, &dados.Tipo, &dados.Descricao, &dados.RealizadaEm)

		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)

			return
		}

		transacoes = append(transacoes, dados)
	}

	if rows.Err() != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	retval := strings.Builder{}

	retval.WriteString("{\"saldo\":")

	retval.WriteString("{\"limite\":")

	retval.WriteString(strconv.FormatInt(saldoEmConta.Limite, 10))

	retval.WriteString(",\"total\":")

	retval.WriteString(strconv.FormatInt(saldoEmConta.Total, 10))

	retval.WriteString(", \"data_extrato\":\"")

	retval.WriteString(saldoEmConta.DataExtrato)

	retval.WriteString("\"}")

	retval.WriteString(",\"ultimas_transacoes\":")

	retval.WriteString("[")

	for index, transacao := range transacoes {
		if index > 0 {
			retval.WriteString(",")
		}

		retval.WriteString("{\"tipo\":\"")

		retval.WriteString(transacao.Tipo)

		retval.WriteString("\",\"descricao\":\"")

		retval.WriteString(transacao.Descricao)

		retval.WriteString("\",\"realizada_em\":\"")

		retval.WriteString(transacao.RealizadaEm)

		retval.WriteString("\",\"valor\":")

		retval.WriteString(strconv.FormatInt(transacao.Valor, 10))

		retval.WriteString("}")
	}

	retval.WriteString("]")

	retval.WriteString("}")

	w.Write([]byte(retval.String()))
}
