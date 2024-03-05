# Submissão para Rinha de Backend, Segunda Edição: 2024/Q1 - Controle de Concorrência

Uma api altamente feia, nojenta, desnecessária e horrorosa para a rinha do backend. Do jeito que todo gopher ama... sqn...

# Disclaimer

Esse repositório é feito com fins esportivos e publicado apenas para sua consulta.

Você NÃO DEVE UTILIZAR ELE COMO REFERÊNCIA PARA UMA API REAL.

Ele é bastante bugado e vai estragar seus trabalhos. Se ainda assim, você insistir, É SOB SUA CONTA E RISCO.

# Como usar

## Direto na maquina

Caso já tenha o compilador do go instalado, abra o terminal e rode:

```
go mod download # para baixar as dependênciaas

e

go run ./*.go # para rodar o programa
```

Caso queira rodar um container do postgres para testes, tente:

```
sudo docker run --rm -ti --name=postgres -p 5432:5432 -e POSTGRES_USER=user -e POSTGRES_PASSWORD=senha -e POSTGRES_DB=clientes -v $(pwd)/init.sql:/docker-entrypoint-initdb.d/init.sql:ro postgres:16.1-alpine3.19
```

## Docker composer

Quando aparecer a mensagem "webserver ready to go!" ele estará rodando

Verifique as configurações do postgres no main.go caso de algum problema. Verifique so o postgres não está rodando no localhost.

Ou, alternativamente builde a imagem "rinha" com:

```
sudo docker build -t rinha .
```

E rode o docker-compose. Preste atenção nas configurações nele também.