package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

var ctxMarker struct{}

func ToContextFields(ctx context.Context, fields Fields) context.Context {
	return context.WithValue(ctx, ctxMarker, Ctx(ctx).WithFields(fields))
}

// Ctx creates an entry from the standard logger and adds a context to it.
// Add a single field(trace_id) to the Entry
func Ctx(ctx context.Context) *Entry {
	if entry, ok := ctx.Value(ctxMarker).(*Entry); ok && entry != nil {
		return entry
	}
	entry := logrus.WithContext(ctx)
	return entry
}
