package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	mut      sync.Mutex
	clients  map[string]*Client
	rate     rate.Limit
	burst    int
	lifetime time.Duration
}

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(
	requestsPerSecond float64, burst int,
) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*Client),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burst,
		lifetime: 3 * time.Minute,
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mut.Lock()

		for ip, client := range rl.clients {
			if time.Since(client.lastSeen) > rl.lifetime {
				delete(rl.clients, ip)
			}
		}

		rl.mut.Unlock()
	}
}

func (rl *RateLimiter) Middleware(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(
				w,
				"invalid client address",
				http.StatusInternalServerError,
			)
			return
		}

		rl.mut.Lock()

		client, exists := rl.clients[ip]

		if !exists {
			client = &Client{
				limiter: rate.NewLimiter(
					rl.rate,
					rl.burst,
				),
				lastSeen: time.Now(),
			}

			rl.clients[ip] = client
		}

		client.lastSeen = time.Now()

		rl.mut.Unlock()

		if !client.limiter.Allow() {
			http.Error(
				w,
				"rate limit exceeded",
				http.StatusTooManyRequests,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

