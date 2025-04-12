package main

import "github.com/ulngollm/teleflow"

type FlowController interface {
	checkoutState(flow teleflow.Flow, event string) error
}
