package rest

import (
	"fmt"
	"net/http"

	. "fs/service/config"
	"fs/service/metrics"		
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log.Info("Handling request " + fmt.Sprintf("[%s] %s %s %s",
			r.Method,
			r.Host,
			r.URL.Path,
			r.URL.RawQuery))

		w.Header().Set("X-Request-Id", metrics.RequestId())

		metrics.RequestsCounter.Inc()

        next.ServeHTTP(w, r)
    })
}
