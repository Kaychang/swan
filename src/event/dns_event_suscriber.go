package event

import (
	"encoding/json"
	"sync"

	"github.com/Dataman-Cloud/swan-resolver/nameserver"
	"github.com/Dataman-Cloud/swan/src/types"
	"github.com/Sirupsen/logrus"
)

type DNSSubscriber struct {
	Key          string
	acceptors    map[string]types.ResolverAcceptor
	acceptorLock sync.RWMutex
}

func NewDNSSubscriber() *DNSSubscriber {
	subscriber := &DNSSubscriber{
		Key:       "dns",
		acceptors: make(map[string]types.ResolverAcceptor),
	}

	return subscriber
}

func (subscriber *DNSSubscriber) GetKey() string {
	return subscriber.Key
}

func (subscriber *DNSSubscriber) AddAcceptor(acceptor types.ResolverAcceptor) {
	subscriber.acceptorLock.Lock()
	subscriber.acceptors[acceptor.ID] = acceptor
	subscriber.acceptorLock.Unlock()
}

func (subscriber *DNSSubscriber) RemoveAcceptor(ID string) {
	subscriber.acceptorLock.Lock()
	delete(subscriber.acceptors, ID)
	subscriber.acceptorLock.Unlock()
}

func (subscriber *DNSSubscriber) Write(e *Event) error {
	rgEvent, err := BuildResolverEvent(e)
	if err != nil {
		return err
	}

	go subscriber.pushResloverEvent(rgEvent)

	return nil
}

func (subscriber *DNSSubscriber) InterestIn(e *Event) bool {
	if e.Type == EventTypeTaskHealthy {
		return true
	}

	if e.Type == EventTypeTaskUnhealthy {
		return true
	}

	return false
}

func (subscriber *DNSSubscriber) pushResloverEvent(event *nameserver.RecordGeneratorChangeEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		logrus.Infof("marshal reslover event got error: %s", err.Error())
		return
	}

	subscriber.acceptorLock.RLock()
	for _, acceptor := range subscriber.acceptors {
		if err := SendEventByHttp(acceptor.RemoteAddr, "POST", data); err != nil {
			logrus.Infof("send reslover event by http to %s got error: %s", acceptor.RemoteAddr, err.Error())
		} else {
			logrus.Debugf("send reslover event by http to %s success", acceptor.RemoteAddr)
		}
	}
	subscriber.acceptorLock.RUnlock()
}
