package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Limiter struct {
	ipCount map[string]int
	sync.Mutex
}

var limiter Limiter

func init() {
	limiter.ipCount = make(map[string]int)
}

// Limit the number of request from particular client
// rate limited with 5 request at a time, after that user has to wait for 60Sec and for full reset user has to wait
// for 1 hour to make another request, currently this values are hardcoded in the code which can be made configurable
func Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the IP address for the current user.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get the # of times the visitor has visited in the last 60 seconds
		limiter.Lock()
		count, ok := limiter.ipCount[ip]
		if !ok {
			limiter.ipCount[ip] = 0
		}
		if count > 5 {
			limiter.Unlock()
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		} else {
			limiter.ipCount[ip]++
		}
		time.AfterFunc(time.Second*60, func() {
			limiter.Lock()
			limiter.ipCount[ip]--
			limiter.Unlock()
		})
		if limiter.ipCount[ip] == 5 {
			// set it to 10 so the decrement timers will only decrease it to
			// 5, and they stay blocked until the next timer resets it to 0
			limiter.ipCount[ip] = 10
			time.AfterFunc(time.Hour, func() {
				limiter.Lock()
				limiter.ipCount[ip] = 0
				limiter.Unlock()
			})
		}
		limiter.Unlock()
		next.ServeHTTP(w, r)
	})
}
