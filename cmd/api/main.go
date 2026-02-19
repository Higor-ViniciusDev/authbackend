package main

import (
	"fmt"
	"os"

	"github.com/Higor-ViniciusDev/auth/configuration/database"
	"github.com/Higor-ViniciusDev/auth/configuration/logger"
	"github.com/Higor-ViniciusDev/auth/internal/consumer"
	"github.com/Higor-ViniciusDev/auth/internal/entity"
	"github.com/Higor-ViniciusDev/auth/internal/infra/api/controller"
	"github.com/Higor-ViniciusDev/auth/internal/infra/api/web"
	handlers "github.com/Higor-ViniciusDev/auth/internal/infra/events/handler"
	"github.com/Higor-ViniciusDev/auth/internal/infra/repository"
	"github.com/Higor-ViniciusDev/auth/internal/infra/ws"
	"github.com/Higor-ViniciusDev/auth/internal/middleware"
	user_usecase "github.com/Higor-ViniciusDev/auth/internal/usecase/user"
	"github.com/Higor-ViniciusDev/auth/pkg/events"
	"github.com/Higor-ViniciusDev/auth/pkg/rabbitmq"
	"github.com/joho/godotenv"
)

func main() {
	defer logger.GetLogger().Sync()

	_ = godotenv.Load("./cmd/api/.env")

	webServerPort := os.Getenv("WEB_SERVER_PORT")
	if webServerPort == "" {
		webServerPort = "8080"
	}

	// ── Token ────────────────────────────────────────────────────
	token := entity.NewToken()

	// ── RabbitMQ ─────────────────────────────────────────────────
	producerChannel, err := rabbitmq.OpenChannel()
	if err != nil {
		logger.Error("error opening RabbitMQ producer channel", err)
		panic(err)
	}

	// email.pending consumer channel: consumes and sends emails
	emailPendingChannel, err := rabbitmq.OpenChannel()
	if err != nil {
		logger.Error("error opening email.pending channel", err)
		panic(err)
	}

	// email.verified consumer channel: consumes fanout and notifies WebSocket
	emailVerifiedChannel, err := rabbitmq.OpenChannel()
	if err != nil {
		logger.Error("error opening email.verified channel", err)
		panic(err)
	}

	// ── WebSocket Manager ────────────────────────────────────────
	wsManager := ws.NewManager()

	// ── Event Dispatcher ─────────────────────────────────────────
	createdHandler := handlers.NewCreateuserHandler(producerChannel)
	verifiedHandler := handlers.NewEmailVerifiedHandler(producerChannel)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.RegistrarHandler("UserCreated", createdHandler)
	eventDispatcher.RegistrarHandler("UserEmailVerified", verifiedHandler)

	// ── Consumers in goroutines ───────────────────────────────────
	go consumer.NewEmailPendingConsumer(emailPendingChannel).Start()
	go consumer.NewEmailVerifiedConsumer(emailVerifiedChannel, wsManager).Start()

	// ── Dependency injection ───────────────────────────────────
	db := database.NewConnect()
	userRepo := repository.NewUserRepository(db)
	usecase := user_usecase.NewUserUseCase(userRepo, eventDispatcher, token)
	authHandler := controller.NewAuthLoginHandler(usecase)

	// ── Routes ─────────────────────────────────────────────────────
	webServer := web.NovoWebServer(fmt.Sprintf(":%v", webServerPort))

	webServer.RegistrarRota("/validation", authHandler.Autenticacao, "POST", middleware.CorsPolicy)
	webServer.RegistrarRota("/register", authHandler.Register, "POST", middleware.CorsPolicy)
	webServer.RegistrarRota("/verify", authHandler.VerifyEmail, "GET", middleware.CorsPolicy)
	webServer.RegistrarRota("/resend-verification", authHandler.ResendVerification, "POST", middleware.CorsPolicy)

	webServer.RegistrarRota("/ws/verify-status", wsManager.HandleConnection, "GET")

	webServer.IniciarWebServer()
}
