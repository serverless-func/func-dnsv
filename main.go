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
		w := func(text string) {
			log.Println(text)
		}
		md, err := Render(w)
		if err != nil {
			c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(err.Error()))
			return
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(md))
	})
	r.GET("/index.html", func(c *gin.Context) {
		w := func(text string) {
			c.Writer.Write([]byte(fmt.Sprintf("%s <br>", text)))
			c.Writer.(http.Flusher).Flush()
		}
		resultChan := make(chan string)
		go func() {
			time.Sleep(3 * time.Second) // 模拟耗时任务
			resultChan <- Convert(w)    // 发送处理结果
		}()
		// 立即返回 Loading 页面
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte("<html><body><div id='loading'>"))
		c.Writer.Write([]byte("loading ... <br>"))
		c.Writer.(http.Flusher).Flush() // 立即发送响应

		// 等待异步任务完成并获取结果
		result := <-resultChan

		// 清空loading
		c.Writer.Write([]byte(`
			</div>
			<script>
				document.getElementById("loading").remove();
			</script>
			</body></html>`))
		c.Writer.(http.Flusher).Flush()

		// 返回最终 HTML
		c.Writer.Write([]byte(result))
		c.Writer.(http.Flusher).Flush()
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
