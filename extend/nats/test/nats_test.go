package test

import "testing"

func TestNatsConnect(t *testing.T) {
	nats.New()

	nats.NewJetStream()

}
