package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())

	setupRouter(r)

	start(&http.Server{
		Addr:    fmt.Sprintf(":%s", env("FC_SERVER_PORT", "9000")),
		Handler: r,
	})
}

func setupRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte("pong"))
	})
	r.GET("/index.md", func(c *gin.Context) {
		md, err := Render()
		if err != nil {
			c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(err.Error()))
			return
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(md))
	})
}

func start(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {

			log.Printf("listen: %s\n", err)
		}
	}()

	log.Printf("Start Server @ %s", srv.Addr)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown:%s", err)
	}
	<-ctx.Done()
	log.Print("Server exiting")
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func failed(msg string) gin.H {
	return gin.H{
		"msg":       msg,
		"timestamp": time.Now().Unix(),
	}
}

func data(data interface{}) gin.H {
	return gin.H{
		"msg":       "success",
		"data":      data,
		"timestamp": time.Now().Unix(),
	}
}
