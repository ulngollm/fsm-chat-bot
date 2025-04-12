package main

import "github.com/ulngollm/teleflow"

type FlowController interface {
	CheckoutState(flow *teleflow.Flow, e string) error
}
