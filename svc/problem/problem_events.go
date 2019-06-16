package main

import (
	"bytes"
	"encoding/json"

	"github.com/syzoj/syzoj-ng-go/svc/problem/model"
)

type ProblemInsertEvent struct {
	Id        string   `json:"id"`
	Title     string   `json:"title"`
	Statement string   `json:"statement"`
	Tags      []string `json:"tags"`
}

func (ev *ProblemInsertEvent) execute(m *Main) error {
	doc := &model.ProblemDoc{
		Title:     ev.Title,
		Statement: ev.Statement,
		Tags:      ev.Tags,
	}
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	resp, err := m.es.Create("problem", ev.Id, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Info(resp)
	return nil
}

type ProblemUpdateEvent struct {
	Id        string   `json:"id"`
	Title     string   `json:"title"`
	Statement string   `json:"statement"`
	Tags      []string `json:"tags"`
}

func (ev *ProblemUpdateEvent) execute(m *Main) error {
	doc := &model.ProblemDoc{
		Title:     ev.Title,
		Statement: ev.Statement,
		Tags:      ev.Tags,
	}
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	resp, err := m.es.Index("problem", bytes.NewBuffer(data), m.es.Index.WithDocumentID(ev.Id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Info(resp)
	return nil
}

type ProblemDeleteEvent struct {
	Id string `json:"id"`
}

func (ev *ProblemDeleteEvent) execute(m *Main) error {
	resp, err := m.es.Delete("problem", ev.Id)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Info(resp)
	return nil
}
