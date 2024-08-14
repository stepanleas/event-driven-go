package middleware

import (
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	messagesProcessingDuration = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "messages",
			Name:       "processing_duration_seconds",
			Help:       "The total time spent processing messages",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"topic", "handler"},
	)

	messagesProcessingFailedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "messages",
			Name:      "processing_failed_total",
			Help:      "The total number of messages processing failures",
		},
		[]string{"topic", "handler"},
	)

	messagesProcessedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "messages",
			Name:      "processed_total",
			Help:      "The total number of processed messages",
		},
		[]string{"topic", "handler"},
	)
)

func TracingMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		topic := message.SubscribeTopicFromCtx(msg.Context())
		handler := message.HandlerNameFromCtx(msg.Context())

		ctx := msg.Context()

		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Metadata))

		ctx, span := otel.Tracer("").Start(
			ctx,
			fmt.Sprintf("topic: %s, handler: %s", topic, handler),
			trace.WithAttributes(
				attribute.String("topic", topic),
				attribute.String("handler", handler),
			),
		)
		defer span.End()

		msg.SetContext(ctx)

		msgs, err := next(msg)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		return msgs, err
	}
}

func PrometheusMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) (events []*message.Message, err error) {
		start := time.Now()

		topic := message.SubscribeTopicFromCtx(msg.Context())
		handler := message.HandlerNameFromCtx(msg.Context())

		labels := prometheus.Labels{"topic": topic, "handler": handler}

		messagesProcessedCounter.With(labels).Inc()

		msgs, err := next(msg)
		if err != nil {
			messagesProcessingFailedCounter.With(labels).Inc()
		}

		messagesProcessingDuration.With(labels).Observe(time.Since(start).Seconds())

		return msgs, err
	}
}
