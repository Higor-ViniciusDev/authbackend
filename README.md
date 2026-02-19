# Auth Backend

API REST em **Go** responsável por toda a lógica de autenticação do sistema. Gerencia cadastro de usuários, login com JWT, verificação de e-mail e notificações em tempo real via WebSocket.

## O que faz

- **POST /validation** — Login: valida credenciais e retorna um JWT.
- **POST /register** — Cadastro: cria um novo usuário e dispara um e-mail de verificação via RabbitMQ.
- **GET /verify?token=JWT** — Verificação de e-mail: decodifica o token e marca o usuário como verificado no banco.
- **POST /resend-verification** — Reenvia o e-mail de verificação.
- **WS /ws/verify-status** — WebSocket que notifica o frontend em tempo real quando o e-mail é verificado.

Utiliza **PostgreSQL** como banco de dados, **RabbitMQ** para mensageria (filas de e-mail pendente e exchange fanout para notificação de verificação) e **Chi** como router HTTP.

## Como rodar localmente

### 1. Subir as dependências (Postgres + RabbitMQ)

```bash
docker-compose up -d
```

### 2. Criar a tabela no banco

Execute o conteúdo de `init.sql` no PostgreSQL (database `autenticacao`).

### 3. Configurar variáveis de ambiente

Crie o arquivo `cmd/api/.env`:

```env
WEB_SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=autenticacao
```

> **Obs:** a conexão com o RabbitMQ está hardcoded em `pkg/rabbitmq/rabbitmq.go` apontando para `rabbitmq:5672`. Para rodar fora do Docker, altere o host para `localhost`.

### 4. Build e execução

```bash
go build -o auth-backend ./cmd/api
./auth-backend
```

O servidor sobe na porta configurada (padrão `8080`).
