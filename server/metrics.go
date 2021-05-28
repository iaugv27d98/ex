package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/circleci/ex/o11y"
)

type MetricProducer interface {
	// MetricName The name for this group of metrics
	//(Name might be cleaner, but is much more likely to conflict in implementations)
	MetricName() string
	// Gauges are instantaneous name value pairs
	Gauges(context.Context) map[string]float64
}

func traceMetrics(ctx context.Context, producers []MetricProducer) {
	// acquire a span from the context that called traceMetrics, this saves on
	// unnecessary spans, (we don't care about the time it takes to generate the metrics.)
	parentSpan := o11y.FromContext(ctx).GetSpan(ctx)
	for _, producer := range producers {
		traceMetric(ctx, parentSpan, producer)
	}
}

func traceMetric(ctx context.Context, span o11y.Span, producer MetricProducer) {
	producerName := strings.Replace(producer.MetricName(), "-", "_", -1)
	for f, v := range producer.Gauges(ctx) {
		scopedField := fmt.Sprintf("gauge.%s.%s", producerName, f)
		span.AddRawField(scopedField, v)
		span.RecordMetric(o11y.Gauge(scopedField, scopedField))
	}
}