package main

import (
	"encoding/json"
	"fmt"
)

type EventBody struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Event interface {
	execute(m *Main) error
}

func encodeEvent(ev Event) []byte {
	b := &EventBody{}
	switch ev.(type) {
	case *ProblemInsertEvent:
		b.Type = "problem.insert"
	case *ProblemUpdateEvent:
		b.Type = "problem.update"
	case *ProblemDeleteEvent:
		b.Type = "problem.delete"
	default:
		panic("Unknown event type")
	}
	p, err := json.Marshal(ev)
	if err != nil {
		panic(err)
	}
	b.Payload = p
	p, err = json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return p
}

func decodeEvent(payload []byte) (Event, error) {
	ep := &EventBody{}
	if err := json.Unmarshal(payload, ep); err != nil {
		return nil, err
	}
	var ev Event
	switch ep.Type {
	case "problem.insert":
		ev = &ProblemInsertEvent{}
	case "problem.update":
		ev = &ProblemUpdateEvent{}
	case "problem.delete":
		ev = &ProblemDeleteEvent{}
	default:
		return nil, fmt.Errorf("Unrecognized event type: %s", ep.Type)
	}
	if err := json.Unmarshal(ep.Payload, ev); err != nil {
		return nil, err
	}
	return ev, nil
}

func (m *Main) OnEvent(pos []byte, data []byte) error {
	ev, err := decodeEvent(data)
	if err != nil {
		return err
	}
	return ev.execute(m)
}
