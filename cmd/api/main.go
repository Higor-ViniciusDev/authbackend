package main

import (
	"fmt"
	"os"

	"github.com/Higor-ViniciusDev/auth/configuration/database"
	"github.com/Higor-ViniciusDev/auth/configuration/logger"
	"github.com/Higor-ViniciusDev/auth/internal/infra/api/controller"
	"github.com/Higor-ViniciusDev/auth/internal/infra/api/web"
	handlers "github.com/Higor-ViniciusDev/auth/internal/infra/events/handler"
	"github.com/Higor-ViniciusDev/auth/internal/infra/repository"
	"github.com/Higor-ViniciusDev/auth/internal/middleware"
	user_usecase "github.com/Higor-ViniciusDev/auth/internal/usecase/user"
	"github.com/Higor-ViniciusDev/auth/pkg/events"
	"github.com/Higor-ViniciusDev/auth/pkg/rabbitmq"
	"github.com/joho/godotenv"
)

func main() {
	defer logger.GetLogger().Sync()

	// .env é opcional — no Kubernetes, variáveis vêm do deployment.yaml
	_ = godotenv.Load("./cmd/api/.env")

	webServerPort := os.Getenv("WEB_SERVER_PORT")
	if webServerPort == "" {
		webServerPort = "8080"
	}
	//Injection rabbitmq
	rabbitMqCanal, _ := rabbitmq.OpenChannel()
	userHandler := handlers.NewCreateuserHandler(rabbitMqCanal)

	eventorDisparador := events.NewEventDispatcher()
	eventorDisparador.RegistrarHandler("UserCreated", userHandler)

	//Injection
	db := database.NewConnect()
	userRepo := repository.NewUserRepository(db)
	usecase := user_usecase.NewUserUseCase(userRepo, eventorDisparador)
	authHandler := controller.NewAuthLoginHandler(usecase)

	webServer := web.NovoWebServer(fmt.Sprintf(":%v", webServerPort))

	webServer.RegistrarRota("/validation", authHandler.Autenticacao, "POST", middleware.CorsPolicy)
	webServer.RegistrarRota("/register", authHandler.Register, "POST", middleware.CorsPolicy)
	webServer.RegistrarRota("/resend-verification", nil, "POST", middleware.CorsPolicy)
	webServer.IniciarWebServer()
}
