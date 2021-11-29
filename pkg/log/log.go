package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Logger provides the default logger for the whole maiao project
var Logger = logrus.New()

const logContextKey = "maiao.log"

func ForContext(ctx context.Context) *logrus.Entry {
	return Logger.WithFields(contextFields(ctx))
}

func WithContextFields(ctx context.Context, fields logrus.Fields) context.Context {
	f := contextFields(ctx)
	for k, v := range fields {
		f[k] = v
	}
	return context.WithValue(ctx, logContextKey, f)
}

func contextFields(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}
	v := ctx.Value(logContextKey)
	if v != nil {
		f, ok := v.(logrus.Fields)
		if ok {
			return f
		}
	}
	return fields
}
