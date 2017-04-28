package metricemitter

import (
	"fmt"
	v2 "plumbing/v2"
	"sync/atomic"
	"time"
)

type CounterMetric struct {
	client   *client
	name     string
	sourceID string
	tags     map[string]*v2.Value
	delta    uint64
}

type MetricOption func(*CounterMetric)

func NewCounterMetric(name, sourceID string, opts ...MetricOption) *CounterMetric {
	m := &CounterMetric{
		name:     name,
		sourceID: sourceID,
		tags:     make(map[string]*v2.Value),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *CounterMetric) Increment(c uint64) {
	atomic.AddUint64(&m.delta, c)
}

func (m *CounterMetric) GetDelta() uint64 {
	return atomic.LoadUint64(&m.delta)
}

func (m *CounterMetric) WithEnvelope(fn func(*v2.Envelope) error) error {
	d := m.GetDelta()

	if err := fn(m.toEnvelope(d)); err != nil {
		return err
	}

	atomic.AddUint64(&m.delta, -d)
	return nil
}

func (m *CounterMetric) toEnvelope(delta uint64) *v2.Envelope {
	return &v2.Envelope{
		SourceId:  m.sourceID,
		Timestamp: time.Now().UnixNano(),
		Message: &v2.Envelope_Counter{
			Counter: &v2.Counter{
				Name: m.name,
				Value: &v2.Counter_Delta{
					Delta: delta,
				},
			},
		},
		Tags: m.tags,
	}
}

func WithVersion(major, minor uint) MetricOption {
	return WithTags(map[string]string{
		"metric_version": fmt.Sprintf("%d.%d", major, minor),
	})
}

func WithTags(tags map[string]string) MetricOption {
	return func(m *CounterMetric) {
		for k, v := range tags {
			m.tags[k] = &v2.Value{
				Data: &v2.Value_Text{
					Text: v,
				},
			}
		}
	}
}
