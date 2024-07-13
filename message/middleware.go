package message

import (
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/sirupsen/logrus"
)

func CorrelationIDMiddleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		ctx := msg.Context()

		reqCorrelationId := msg.Metadata.Get("correlation_id")
		if reqCorrelationId == "" {
			reqCorrelationId = watermill.NewUUID()
		}

		ctx = log.ToContext(ctx, logrus.WithFields(logrus.Fields{"correlation_id": reqCorrelationId}))
		ctx = log.ContextWithCorrelationID(ctx, reqCorrelationId)

		msg.SetContext(ctx)

		return h(msg)
	}
}

func LoggingMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		logger := log.FromContext(msg.Context()).WithField("message_uuid", msg.UUID)

		logger.Info("Handling a message")

		msgs, err := next(msg)
		if err != nil {
			logger.WithError(err).Error("Message handling error")
		}

		return msgs, err
	}
}

func RetryMiddleware(logger watermill.LoggerAdapter) middleware.Retry {
	return middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Millisecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          logger,
	}
}
