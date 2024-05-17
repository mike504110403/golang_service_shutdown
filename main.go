package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var wg sync.WaitGroup

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	wg.Add(1)
	go func() {
		// service connections
		log.Println("========== Server start ==========")
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 接收關服務訊號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	srv.SetKeepAlivesEnabled(false)

	log.Println("========== Shutdown Server ... ==========")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("========== Server Shutdown: ", err, " ==========")
	}

	log.Println("========== Watting unfinished requests ==========")
	wg.Wait()
	log.Println("========== Server exiting ==========")
}
