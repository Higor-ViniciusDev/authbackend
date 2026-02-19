package web

import (
	"fmt"
	"net/http"

	"github.com/Higor-ViniciusDev/auth/configuration/logger"
	"github.com/go-chi/chi/v5"
)

type HandlerInfo struct {
	Metodo      string
	Handler     http.HandlerFunc
	Middlewares []func(http.Handler) http.Handler
}

type WebServer struct {
	Porta    string
	Handlers map[string][]HandlerInfo // path -> HandlerInfo (method + handler)
	Rotas    chi.Router
}

func NovoWebServer(porta string) *WebServer {
	return &WebServer{
		Porta:    porta,
		Handlers: make(map[string][]HandlerInfo),
		Rotas:    chi.NewRouter(),
	}
}

func (w WebServer) RegistrarRota(caminho string, handlerFunc http.HandlerFunc, metodo string, middlewares ...func(http.Handler) http.Handler,
) {
	w.Handlers[caminho] = append(w.Handlers[caminho], HandlerInfo{
		Metodo:      metodo,
		Handler:     handlerFunc,
		Middlewares: middlewares,
	})
}

func (w WebServer) IniciarWebServer() {
	for rota, handlers := range w.Handlers {
		for _, infoHandle := range handlers {
			var handler http.Handler = infoHandle.Handler

			for i := len(infoHandle.Middlewares) - 1; i >= 0; i-- {
				handler = infoHandle.Middlewares[i](handler)
			}

			w.Rotas.Method(infoHandle.Metodo, rota, handler)

			//adapter for handler
			w.Rotas.Options(rota, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.ServeHTTP(w, r)
			}))
			logger.Info(fmt.Sprintf("registering route %v with method %v", rota, infoHandle.Metodo))
		}
	}

	logger.Info(fmt.Sprintf("starting server on port %v", w.Porta))
	err := http.ListenAndServe(w.Porta, w.Rotas)

	if err != nil {
		logger.Error("error starting webserver", err)
	}
}
