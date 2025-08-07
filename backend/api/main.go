package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

type server struct {
	mux *mux.Router
}

func (s *server) routes() {
	r := s.mux
	r.HandleFunc("/healthz", s.handleHealth()).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/traffic/vehicle", s.handleVehicle()).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/traffic/incident", s.handleIncident()).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/routes/optimal", s.handleOptimalRoute()).Methods(http.MethodPost)
    r.HandleFunc("/ws", s.handleWS())
}

func (s *server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func (s *server) handleVehicle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"message": "vehicle event accepted"})
	}
}

func (s *server) handleIncident() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"message": "incident accepted"})
	}
}

func (s *server) handleOptimalRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		json.NewEncoder(w).Encode(map[string]string{"error": "route optimization not implemented"})
	}
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *server) handleWS() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Printf("ws upgrade error: %v", err)
            return
        }
        defer conn.Close()
        // Simple heartbeat loop
        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                break
            }
            // Echo
            _ = conn.WriteMessage(websocket.TextMessage, append([]byte("ack:"), msg...))
        }
    }
}

func initConfig() {
	viper.SetEnvPrefix("routeiq")
	viper.AutomaticEnv()
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("READ_TIMEOUT_SEC", 15)
	viper.SetDefault("WRITE_TIMEOUT_SEC", 15)
	viper.SetDefault("IDLE_TIMEOUT_SEC", 60)
}

func main() {
	initConfig()

	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":8080"
	}

	s := &server{mux: mux.NewRouter()}
	s.routes()

	c := cors.New(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Addr:         addr,
		Handler:      c.Handler(s.mux),
		ReadTimeout:  time.Duration(viper.GetInt("READ_TIMEOUT_SEC")) * time.Second,
		WriteTimeout: time.Duration(viper.GetInt("WRITE_TIMEOUT_SEC")) * time.Second,
		IdleTimeout:  time.Duration(viper.GetInt("IDLE_TIMEOUT_SEC")) * time.Second,
	}

	log.Printf("API listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
	_ = context.Background()
}
