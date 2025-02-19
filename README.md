# D1: GO
Desafio 1

## Objetivos
Demonstrar domínio no item abaixo:
- webserver http
- contexto
- banco de dados
- manipulação de arquivos

# Dinâmica cliente servidor
Crie dois arquivos:
- `client.go`
- `server.go`

# CLIENT PARÂMETROS .env
- `SERVER_PORT` porta do servidor
- `QUOTE_REQUEST_DELAY` Delay da requisição de cota (api externa)
- `QUOTE_TIMEOUT_REQUEST` Timeout da requisição de cota  (api externa)
- `DB_QUOTE_DELAY` Delay da requisição DB de cota
- `DB_QUOTE_TIMEOUT` Timeout da requisição DB de cota
- `DB_SQLITE` nome do banco de dados SQLite

# SERVER PARÂMETROS
- `SERVER_PORT` = porta do servidor
- `QUOTE_REQUEST_DELAY` = Delay da requisição de cota (app server)
- `QUOTE_TIMEOUT_REQUEST` = Timeout da requisição de cota (app server)

* A opção **DEBUG** está disponível nos arquivos client *main.go*, server *handler.go*.

# Requisitos
- [X] O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente.
- [X] O endpoint necessário gerado pelo server.go é: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.
- [X] O server.go deverá registrar no banco de dados SQLite cada cotação recebida.
- [X] O server.go usando "context" tem o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
- [X] O server.go usando "context" tem o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms
- [X] O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
- [X] O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON).
- [X] O client.go usando "context" tem o timeout máximo de 300ms para receber o resultado do server.go.
- [X] O client.go terá que salvar a cotação atual em um arquivo cotacao.txt no formato: Dólar: {valor}
- [X] Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.

# Observação
- O banco de dados SQLite já exite e com a tabela quote, não é preciso executar migration.sql.