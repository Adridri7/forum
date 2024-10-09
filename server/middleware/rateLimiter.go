package middleware

import (
	"net/http"
	"sync"
	"time"
)

var Rl = NewRateLimiter(5, 10*time.Second)

func RateLimiterMiddleware(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ip := r.RemoteAddr

		if !Rl.Allow(ip) {
			http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}

}

type RateLimiter struct {
	Requests map[string]int
	Limit    int
	Window   time.Duration
	Mu       sync.Mutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	// Crée un nouveau RateLimiter avec une limite de requêtes et une fenêtre de temps donnée
	return &RateLimiter{
		Requests: make(map[string]int),
		Limit:    limit,
		Window:   window,
	}
}

// Fonction pour nettoyer les anciennes entrées
func (rl *RateLimiter) Cleanup() {
	rl.Mu.Lock()
	defer rl.Mu.Unlock()

	for ip := range rl.Requests {
		rl.Requests[ip] = 0 // Réinitialise le compteur après la fenêtre de temps
	}
}

// Vérifie si une requête dépasse la limite
func (rl *RateLimiter) Allow(ip string) bool {
	rl.Mu.Lock()
	defer rl.Mu.Unlock()

	// Augmente le nombre de requêtes pour cette adresse IP
	rl.Requests[ip]++

	// Si le nombre de requêtes est supérieur à la limite, refuse
	if rl.Requests[ip] > rl.Limit {
		return false
	}

	return true
}
