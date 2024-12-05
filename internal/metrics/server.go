package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// StartMetricsServer запускает сервер метрик
func StartMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":"+addr, nil); err != nil {
			panic(err)
		}
	}()
}
