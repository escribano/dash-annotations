// +build darwin
package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"launch"
)

func activityHandler(next http.Handler) http.HandlerFunc {
	var mu sync.Mutex
	var lastActivity time.Time
	var done = make(chan struct{})

	go func() {
		for range done {
			mu.Lock()
			if time.Now().Sub(lastActivity) >= 5*time.Minute {
				os.Exit(0)
			}
			mu.Unlock()
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		lastActivity = time.Now()
		mu.Unlock()

		next.ServeHTTP(w, r)

		go func() {
			select {
			case <-time.After(10 * time.Second):
				done <- struct{}{}
			}
		}()
	}
}

func main() {
	if listen != "" {
		log.Printf("Listening on %q", listen)
		http.ListenAndServe(listen, logHandler(jsonHandler(mux)))
	}

	log.Printf("Launchd handling")
	listeners, err := launch.SocketListeners("DashSocket")
	if err != nil || len(listeners) == 0 {
		log.Fatalf("dash: error activating launch socket: %s", err)
	}
	if len(listeners) != 1 {
		log.Fatalf("dash: too many sockets in plist: %d", len(listeners))
	}

	var listener net.Listener = listeners[0]
	http.Serve(listener, activityHandler(logHandler(jsonHandler(mux))))
}
