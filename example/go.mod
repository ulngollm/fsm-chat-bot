module github.com/ulngollm/teleflow/example

go 1.23.5

replace github.com/ulngollm/teleflow => ../

require (
	github.com/jessevdk/go-flags v1.6.1
	github.com/looplab/fsm v1.0.2
	github.com/ulngollm/teleflow v0.0.0-20250412170803-9115d81a382e
	gopkg.in/telebot.v4 v4.0.0-beta.4
)

require golang.org/x/sys v0.32.0 // indirect
