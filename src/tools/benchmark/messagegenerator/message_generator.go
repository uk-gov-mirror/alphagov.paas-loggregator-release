package messagegenerator

import (
	"github.com/cloudfoundry/loggregatorlib/logmessage"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"
	"time"
)

type ValueMetricGenerator struct{}

func NewValueMetricGenerator() *ValueMetricGenerator {
	return &ValueMetricGenerator{}
}

func (*ValueMetricGenerator) Generate() []byte {
	return BasicValueMetric()
}

func BasicValueMetric() []byte {
	message, _ := proto.Marshal(BasicValueMetricEnvelope())
	return message
}

func BasicValueMetricEnvelope() *events.Envelope {
	return &events.Envelope{
		Origin:    proto.String("fake-origin-2"),
		EventType: events.Envelope_ValueMetric.Enum(),
		ValueMetric: &events.ValueMetric{
			Name:  proto.String("fake-metric-name"),
			Value: proto.Float64(42),
			Unit:  proto.String("fake-unit"),
		},
	}
}

type LogMessageGenerator struct{}

func NewLogMessageGenerator() *LogMessageGenerator {
	return &LogMessageGenerator{}
}

func (*LogMessageGenerator) Generate() []byte {
	return BasicLogMessage()
}

func BasicLogMessage() []byte {
	message, _ := proto.Marshal(BasicLogMessageEnvelope())
	return message
}

func BasicLogMessageEnvelope() *events.Envelope {
	return &events.Envelope{
		Origin:    proto.String("fake-origin-2"),
		EventType: events.Envelope_LogMessage.Enum(),
		LogMessage: &events.LogMessage{
			Message:     []byte("test message"),
			MessageType: events.LogMessage_OUT.Enum(),
			Timestamp:   proto.Int64(time.Now().UnixNano()),
		},
	}
}

type LegacyLogGenerator struct{}

func NewLegacyLogGenerator() *LegacyLogGenerator {
	return &LegacyLogGenerator{}
}

func (*LegacyLogGenerator) Generate() []byte {
	return BasicLegacyLogMessage()
}

func BasicLegacyLogMessage() []byte {
	message, _ := proto.Marshal(BasicLegacyLogMessageEnvelope())
	return message
}

func BasicLegacyLogMessageEnvelope() *logmessage.LogEnvelope {

	return &logmessage.LogEnvelope{
		RoutingKey: proto.String("routing-key"),
		Signature:  []byte(""),
		LogMessage: &logmessage.LogMessage{
			Message:     []byte("test message"),
			MessageType: logmessage.LogMessage_OUT.Enum(),
			Timestamp:   proto.Int64(time.Now().UnixNano()),
			AppId:       proto.String("app-id"),
		},
	}
}
