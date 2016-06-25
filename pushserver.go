// Author: Antoine Mercadal
// See LICENSE file for full LICENSE
// Copyright 2016 Aporeto.

package bahamut

import (
	"bytes"
	"encoding/json"

	"golang.org/x/net/websocket"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/aporeto-inc/elemental"
	"github.com/go-zoo/bone"
)

type pushServer struct {
	address       string
	sessions      map[string]*pushSession
	register      chan *pushSession
	unregister    chan *pushSession
	events        chan *elemental.Event
	close         chan bool
	multiplexer   *bone.Mux
	kafkaInfo     *KafkaInfo
	kafkaProducer sarama.SyncProducer
}

func newPushServer(address string, multiplexer *bone.Mux, kafkaInfo *KafkaInfo) *pushServer {

	srv := &pushServer{
		address:     address,
		sessions:    map[string]*pushSession{},
		register:    make(chan *pushSession),
		unregister:  make(chan *pushSession),
		close:       make(chan bool),
		events:      make(chan *elemental.Event),
		multiplexer: multiplexer,
		kafkaInfo:   kafkaInfo,
	}

	srv.multiplexer.Handle("/events", websocket.Handler(srv.handleConnection))

	return srv
}

// adds a new push session to register in the push server
func (n *pushServer) registerSession(session *pushSession) {

	n.register <- session
}

// adds a new push session to unregister from the push server
func (n *pushServer) unregisterSession(session *pushSession) {

	n.unregister <- session
}

// unpublish all local sessions from the global registry system
func (n *pushServer) handleConnection(ws *websocket.Conn) {

	session := newSession(ws, n)
	n.registerSession(session)
	session.listen()
}

// push a new event. If the global push system is available, it will be used.
// otherwise, only local sessions will receive the push
func (n *pushServer) pushEvents(events ...*elemental.Event) {

	// if we don't have a valid kafka producer, we simply manually push the events
	if n.kafkaProducer == nil {
		for _, e := range events {
			n.events <- e
		}
		return
	}

	for _, event := range events {

		buffer := &bytes.Buffer{}
		if err := json.NewEncoder(buffer).Encode(event); err != nil {
			log.WithFields(log.Fields{
				"producer": n.kafkaProducer,
				"event":    event,
			}).Error("unable to encode event")
		}

		message := &sarama.ProducerMessage{
			Topic: n.kafkaInfo.Topic,
			Key:   sarama.StringEncoder("namespace=default"),
			Value: sarama.ByteEncoder(buffer.Bytes()),
		}

		n.kafkaProducer.SendMessage(message)
	}
}

// starts the push server
func (n *pushServer) start() {

	if n.kafkaInfo != nil {
		n.kafkaProducer = n.kafkaInfo.makeProducer()
		// n.kafkaInfo.createTopicsIfNeeded([]string{"bahamut:push:events"})

		defer n.kafkaProducer.Close()

		log.WithFields(log.Fields{
			"info": n.kafkaInfo,
		}).Info("global push system is active")
	} else {
		log.Warn("global push system is inactive")
	}

	log.WithFields(log.Fields{
		"endpoint": n.address + "/events",
	}).Info("starting event server")

	for {
		select {

		case session := <-n.register:

			if _, ok := n.sessions[session.id]; ok {
				break
			}

			n.sessions[session.id] = session

			log.WithFields(log.Fields{
				"total":  len(n.sessions),
				"client": session.socket.RemoteAddr(),
			}).Info("started push session")

		case session := <-n.unregister:

			if _, ok := n.sessions[session.id]; !ok {
				break
			}

			delete(n.sessions, session.id)

			log.WithFields(log.Fields{
				"total":  len(n.sessions),
				"client": session.socket.RemoteAddr(),
			}).Info("closed session")

		case event := <-n.events:
			go func() {
				for _, session := range n.sessions {
					buffer := &bytes.Buffer{}
					if err := json.NewEncoder(buffer).Encode(event); err != nil {
						log.WithFields(log.Fields{
							"data": event,
						}).Error("unable to encode event data")
						return
					}

					session.events <- buffer.String()
				}
			}()

		case <-n.close:

			for _, session := range n.sessions {
				session.close()
			}
			n.sessions = map[string]*pushSession{}

			return
		}
	}
}

// stops the push server
func (n *pushServer) stop() {

	n.close <- true
}
