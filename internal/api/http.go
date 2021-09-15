package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type HttpApi struct {
	Hub    *AppHub
	router *http.ServeMux
	ctx    context.Context
}

func NewHttpApi(actorSystem *actor.ActorSystem, ctx context.Context) *HttpApi {
	router := http.NewServeMux()
	hub := serveHub(router, actorSystem, ctx)
	serveStaticFiles(router)

	return &HttpApi{
		Hub:    hub,
		router: router,
		ctx:    ctx,
	}
}

func (api *HttpApi) ListenAndServe() <-chan bool {
	done := make(chan bool)

	go listenAndServe(api.router, done, api.ctx)

	return done
}

func serveStaticFiles(router *http.ServeMux) {
	fs := http.FileServer(http.Dir("./public"))
	router.Handle("/", fs)
}

func listenAndServe(router *http.ServeMux, done chan<- bool, ctx context.Context) {
	address := "localhost:8080"
	server := &http.Server{Addr: address, Handler: router}

	go func() {
		fmt.Printf("Http server starting to listen at http://%s\n", address)
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}

		done <- true
	}()

	<-ctx.Done()
	fmt.Println("Shutting down http server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(shutdownCtx)
}
