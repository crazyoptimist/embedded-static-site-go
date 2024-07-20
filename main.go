package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

//go:embed public
var filesystem embed.FS

var port string
var help bool

func main() {
	flag.StringVar(&port, "port", "9999", "webserver port")
	flag.BoolVar(&help, "help", false, "show usage guide")
	flag.Parse()

	if help {
		fmt.Println(`
--port <custom port> Runs the webserver using a custom port.
--help Shows the usage guide`)
		return
	}

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	router := http.NewServeMux()

	publicAsRoot, _ := fs.Sub(filesystem, "public")

	router.Handle("/", http.FileServer(http.FS(publicAsRoot)))

	server := &http.Server{
		Addr:     fmt.Sprintf("0.0.0.0:%s", port),
		Handler:  router,
		ErrorLog: logger,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not listen on port %s: %v\n", port, err)
		}
	}()

	go func() {
		// Give some time to the webserver
		time.Sleep(time.Second)
		openBrowser(fmt.Sprintf("http://localhost:%s", port))
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down the HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("Failed to shutdown:", err)
	}
	fmt.Println("Server successfully shut down.")

}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		panic(err)
	}

}
