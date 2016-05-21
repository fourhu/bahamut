// Author: Antoine Mercadal
// See LICENSE file for full LICENSE
// Copyright 2016 Aporeto.

package bahamut

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/aporeto-inc/elemental"
	"github.com/go-zoo/bone"
)

var defaultBahamut *Bahamut

// DefaultBahamut returns the defaut Bahamut.
// Needless to say I don't like this. but that will be ok for now.
func DefaultBahamut() *Bahamut {
	return defaultBahamut
}

// Bahamut is crazy
type Bahamut struct {
	address         string
	apiServer       *apiServer
	pushServer      *pushServer
	multiplexer     *bone.Mux
	enablePush      bool
	enableProfiling bool
	stop            chan bool
	processors      map[string]Processor
	authenticator   Authenticator
}

// NewBahamut creates a new Bahamut.
func NewBahamut(address string, routes []*Route, enabledAPI, enablePush, enableProfiling bool) *Bahamut {

	mux := bone.New()

	var apiServer *apiServer
	if enabledAPI {
		apiServer = newAPIServer(address, mux, routes)
	}

	var pushServer *pushServer
	if enablePush {
		pushServer = newPushServer(address, mux)
	}

	if enableProfiling {
		mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	}

	srv := &Bahamut{
		address:         address,
		apiServer:       apiServer,
		pushServer:      pushServer,
		multiplexer:     mux,
		enablePush:      enablePush,
		enableProfiling: enableProfiling,
		stop:            make(chan bool),
		processors:      make(map[string]Processor),
	}

	defaultBahamut = srv

	return srv
}

// RegisterProcessor registers a new Processor for a particular Identity.
func (b *Bahamut) RegisterProcessor(processor Processor, identity elemental.Identity) error {

	if _, ok := b.processors[identity.Name]; ok {
		return fmt.Errorf("identity %s already has a registered processor", identity)
	}

	b.processors[identity.Name] = processor

	return nil
}

// UnregisterProcessor unregisters a registered Processor for a particular identity.
func (b *Bahamut) UnregisterProcessor(identity elemental.Identity) error {

	if _, ok := b.processors[identity.Name]; !ok {
		return fmt.Errorf("no registered processor for identity %s", identity)
	}

	delete(b.processors, identity.Name)

	return nil
}

// ProcessorForIdentity returns the registered Processor for a particular identity.
func (b *Bahamut) ProcessorForIdentity(identity elemental.Identity) (Processor, error) {

	if _, ok := b.processors[identity.Name]; !ok {
		return nil, fmt.Errorf("no registered processor for identity %s", identity)
	}

	return b.processors[identity.Name], nil
}

// Push pushes the given events to all active sessions.
func (b *Bahamut) Push(events ...*elemental.Event) {

	if !b.enablePush {
		panic("you cannot push events as it is not enabled.")
	}

	b.pushServer.pushEvents(events...)
}

// SetAuthenticator sets the Authenticator to use for the Bahamut server.
func (b *Bahamut) SetAuthenticator(authenticator Authenticator) {
	b.authenticator = authenticator
}

// Authenticator returns the current authenticator
func (b *Bahamut) Authenticator() (Authenticator, error) {

	if b.authenticator == nil {
		return nil, fmt.Errorf("no authenticator set")
	}

	return b.authenticator, nil
}

// Start starts the Bahamut server.
func (b *Bahamut) Start() {

	if b.enableProfiling {
		log.WithFields(log.Fields{
			"endpoint": b.address + "/debug/pprof/",
		}).Info("starting profiling server")
	}

	if b.pushServer != nil {
		go b.pushServer.start()
	}

	log.WithFields(log.Fields{
		"endpoint": b.address,
	}).Info("starting bahamut")

	go func() {

		if err := http.ListenAndServe(b.address, b.multiplexer); err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("unable to start the bahamut")
		}
	}()

	<-b.stop
}

// Stop stops the Bahamut server.
func (b *Bahamut) Stop() {

	if b.pushServer != nil {
		b.pushServer.Stop()
	}

	b.stop <- true
}
