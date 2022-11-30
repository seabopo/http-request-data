package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/seabopo/http-request-data/version"
)

var (
	listenFlag  = flag.String("listen", ":80", "address and port to listen")
	versionFlag = flag.Bool("version", false, "display version information")

	// stdoutW and stderrW are for overriding in test.
	stdoutW = os.Stdout
	stderrW = os.Stderr
)

func main() {
	flag.Parse()

	// Asking for the version?
	if *versionFlag {
		fmt.Fprintln(stderrW, version.HumanVersion)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) > 0 {
		fmt.Fprintln(stderrW, "Too many arguments!")
		os.Exit(127)
	}

	// Flag gets printed as a page
	mux := http.NewServeMux()
	mux.HandleFunc("/", httpLog(stdoutW, withAppHeaders(httpEcho())))

	// Health endpoint
	mux.HandleFunc("/health", withAppHeaders(httpHealth()))

	server := &http.Server{
		Addr:    *listenFlag,
		Handler: mux,
	}
	serverCh := make(chan struct{})
	go func() {
		log.Printf("[INFO] server is listening on %s\n", *listenFlag)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[ERR] server exited with: %s", err)
		}
		close(serverCh)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// Wait for interrupt
	<-signalCh

	log.Printf("[INFO] received interrupt, shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERR] failed to shutdown server: %s", err)
	}

	// If we got this far, it was an interrupt, so don't exit cleanly
	os.Exit(2)
}

func httpEcho() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "HTTP REQUEST DATA")
		fmt.Fprintln(w, "-----------------")
		fmt.Fprintln(w, "Request Time :", time.Now().String())
		fmt.Fprintln(w, "Remote Address :", r.RemoteAddr)
		fmt.Fprintln(w, "Protocol :", r.Proto)
		fmt.Fprintln(w, "Host :", r.Host)
		fmt.Fprintln(w, "Proxied Host :", r.URL.Host)
		fmt.Fprintln(w, "Path :", r.URL.Path)

		for name, values := range r.Header {
			for _, value := range values {
				fmt.Fprintln(w, name, ":", value)
			}
		}
	}
}

func httpHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":"ok"}`)
	}
}
