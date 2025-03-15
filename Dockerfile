# Use a imagem base do Golang
FROM golang:1.23.3

# Configurar o diretório de trabalho no contêiner
WORKDIR /go/src/finanapp

# Copiar os arquivos de dependências
COPY go.mod go.sum ./

# Baixar as dependências do projeto
RUN go mod download

# Copiar o restante do código da aplicação
COPY . .

# Compilar a aplicação
RUN go build -o main ./cmd/app/main.go

# Expor a porta onde o aplicativo será executado
EXPOSE 8080

# Iniciar o aplicativo
CMD ["./main"]
