package middleware

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	visitors   sync.Map
	rps        rate.Limit
	burst      int
	expiration time.Duration
	cleanupInt time.Duration
}

func NewIPRateLimiter(rps int, burst int, expiration, cleanupInt time.Duration) *IPRateLimiter {
	rl := &IPRateLimiter{
		rps:        rate.Limit(rps),
		burst:      burst,
		expiration: expiration,
		cleanupInt: cleanupInt,
	}

	go rl.cleanupVisitors()
	return rl
}

func (rl *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	v, exists := rl.visitors.Load(ip)
	if !exists {
		limiter := rate.NewLimiter(rl.rps, rl.burst)
		rl.visitors.Store(ip, &visitor{limiter: limiter, lastSeen: time.Now()})
		return limiter
	}

	vis := v.(*visitor)
	vis.lastSeen = time.Now()
	return vis.limiter
}

func (rl *IPRateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)
			limiter := rl.getLimiter(ip)

			if !limiter.Allow() {
				log.Printf("IP %s excedeu o rate limit", ip)
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintf(w, "429 - Too Many Requests")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (rl *IPRateLimiter) cleanupVisitors() {
	for {
		time.Sleep(rl.cleanupInt)
		now := time.Now()
		rl.visitors.Range(func(key, value any) bool {
			v := value.(*visitor)
			if now.Sub(v.lastSeen) > rl.expiration {
				rl.visitors.Delete(key)
			}
			return true
		})
	}
}

func getIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
