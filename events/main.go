package events

type Events string
type Service string
type Vendor string

const (
	WAREHOUSE_COST_TRACKER Events = "WAREHOUSE_COST_TRACKER"
)

const (
	ASR Service = "ASR"
	TTS Service = "TTS"
	SLU Service = "SLU"
	LID Service = "LID"
)

const (
	GOOGLE Vendor = "GOOGLE"
	SKIT   Vendor = "SKIT"
	AZURE  Vendor = "AZURE"
	SHORT_UTTERANCE Vendor = "SHORT_UTTERANCE"
)

type CostEvent struct {
	ServiceType      Service `json:"service"`
	Vendor           Vendor  `json:"vendor"`
	ClientUUID       string  `json:"client_uuid"`
	FlowUUID         string  `json:"flow_uuid"`
	CallUUID         string  `json:"call_uuid"`
	ConversationUUID string  `json:"conversation_uuid"`
	NumHits          int     `json:"num_hits"`
	Cost             int     `json:"cost"`
}

func NewCostEvent(service Service, vendor Vendor, clientUUID string, flowUUID string, callUUID string, conversationUUID string) CostEvent {
	return CostEvent{
		ServiceType:      service,
		Vendor:           vendor,
		ClientUUID:       clientUUID,
		FlowUUID:         flowUUID,
		CallUUID:         callUUID,
		ConversationUUID: conversationUUID,
		NumHits:          1,
		Cost:             0,
	}
}

func NewCostEventWithNumHits(service Service, vendor Vendor, clientUUID string, flowUUID string, callUUID string, conversationUUID string, numHits int) CostEvent {
	return CostEvent{
		ServiceType:      service,
		Vendor:           vendor,
		ClientUUID:       clientUUID,
		FlowUUID:         flowUUID,
		CallUUID:         callUUID,
		ConversationUUID: conversationUUID,
		NumHits:          numHits,
		Cost:             0,
	}
}

func NewCostEventWithCost(service Service, vendor Vendor, clientUUID string, flowUUID string, callUUID string, conversationUUID string, cost int) CostEvent {
	return CostEvent{
		ServiceType:      service,
		Vendor:           vendor,
		ClientUUID:       clientUUID,
		FlowUUID:         flowUUID,
		CallUUID:         callUUID,
		ConversationUUID: conversationUUID,
		Cost:             cost,
		NumHits:          1,
	}
}
