package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/labstack/echo"
)

type Server struct {
	Hub    *AppHub
	router *http.ServeMux
	ctx    context.Context
}

func NewHttpServer(cluster *cluster.Cluster, ctx context.Context) *Server {
	router := http.NewServeMux()
	hub := serveHub(router, cluster.ActorSystem, ctx)

	echo := echo.New()
	router.Handle("/", echo)

	serveApi(echo, cluster)
	serveStaticFiles(echo)

	return &Server{
		Hub:    hub,
		router: router,
		ctx:    ctx,
	}
}

func (s *Server) ListenAndServe() <-chan bool {
	done := make(chan bool)

	go listenAndServe(s.router, done, s.ctx)

	return done
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
