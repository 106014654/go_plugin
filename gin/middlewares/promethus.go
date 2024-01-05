package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewardBuilder struct {
	Namespace  string
	Subsystem  string
	Name       string
	Help       string
	InstanceId string
}

func NewMiddlewardBuilder(namespace string, subsystem string, name string, help string, instanceId string) *MiddlewardBuilder {
	return &MiddlewardBuilder{Namespace: namespace, Subsystem: subsystem, Name: name, Help: help, InstanceId: instanceId}
}

func (m *MiddlewardBuilder) Build() gin.HandlerFunc {
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name + "_resp_time",
		Help:      m.Help,
		ConstLabels: map[string]string{
			"instance_id": m.InstanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
	}, []string{"method", "pattern", "status"})

	prometheus.MustRegister(summary)

	return func(ctx *gin.Context) {
		startTime := time.Now()
		defer func() {
			dur := time.Since(startTime)

			pattern := ctx.FullPath()
			if pattern == "" {
				pattern = "unknow"
			}

			summary.WithLabelValues(
				ctx.Request.Method,
				ctx.FullPath(),
				strconv.Itoa(ctx.Writer.Status()),
			).Observe(float64(dur.Milliseconds()))
		}()
		ctx.Next()
	}
}
