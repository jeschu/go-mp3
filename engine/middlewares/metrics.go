package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

func Metrics() func(c *gin.Context) {
	requestMetrics := &requestMetrics{
		uriHistorgams:         make(map[string]prometheus.Histogram),
		uriHistogramsMutex:    new(sync.Mutex),
		statusHistograms:      make(map[int]prometheus.Histogram),
		statusHistogramsMutex: new(sync.Mutex),
	}
	return func(c *gin.Context) {
		uriTimer := prometheus.NewTimer(requestMetrics.uriHistogram(c))
		statusTimer := prometheus.NewTimer(requestMetrics.statusHistogram(c))
		defer uriTimer.ObserveDuration()
		defer statusTimer.ObserveDuration()
	}
}

type requestMetrics struct {
	uriHistorgams         map[string]prometheus.Histogram
	uriHistogramsMutex    *sync.Mutex
	statusHistograms      map[int]prometheus.Histogram
	statusHistogramsMutex *sync.Mutex
}

func (requestMetrics *requestMetrics) uriHistogram(c *gin.Context) (histogram prometheus.Histogram) {
	requestMetrics.uriHistogramsMutex.Lock()
	defer requestMetrics.uriHistogramsMutex.Unlock()
	var (
		name = keyFromUri(c)
		ok   bool
	)
	if histogram, ok = requestMetrics.uriHistorgams[name]; !ok {
		histogram = promauto.NewHistogram(prometheus.HistogramOpts{
			Name: name,
			Help: "request duration diagram by uri",
		})
		requestMetrics.uriHistorgams[name] = histogram
	}
	return
}

func (requestMetrics *requestMetrics) statusHistogram(c *gin.Context) (histogram prometheus.Histogram) {
	requestMetrics.statusHistogramsMutex.Lock()
	defer requestMetrics.statusHistogramsMutex.Unlock()
	var (
		status = c.Writer.Status()
		ok     bool
	)
	if histogram, ok = requestMetrics.statusHistograms[status]; !ok {
		histogram = promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: "http",
			Subsystem: "status",
			Name:      strconv.Itoa(status),
			Help:      "request duration diagram by status code",
		})
		requestMetrics.statusHistograms[status] = histogram
	}
	return
}

func keyFromUri(c *gin.Context) string {
	parts := make([]string, 0)
	parts = append(parts, "http", c.Request.Method)
	for _, part := range strings.Split(c.Request.RequestURI, "/") {
		if un, err := url.PathUnescape(part); err != nil {
			if len(part) > 0 {
				parts = append(parts, part)
			}
		} else {
			if len(un) > 0 {
				parts = append(parts, un)
			}
		}
	}
	for i, part := range parts {
		parts[i] = strings.ReplaceAll(
			strings.ReplaceAll(part, "_", "-"),
			"+", "",
		)
	}
	return strings.Join(parts, "_")
}
