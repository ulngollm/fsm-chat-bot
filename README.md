# Telegram Bot
Middleware для обработки flow с fsm для `tucnak/telebot`.

Работа со стейтами в примере сделана через `looplab/fsm`. Но можно использовать другое решение, т.к. декларация и специфическая работа с fsm вынесена из `Flow` в `handler`-ы 

## Установка
```bash
go get github.com/ulngollm/middleware
```

## Features
1. Сохранение стейта между сообщениями
2. Просто и понятно добавлять обработчики для стейтов

## Как использовать

```go
package main

import (
	"fmt"
	lflow "github.com/ulngollm/msg-constructor/internal/flow"
	"github.com/ulngollm/msg-constructor/internal/middleware"
	"github.com/ulngollm/msg-constructor/internal/state"
)

func run() {
	pool := lflow.NewPool()
	flowManager := lflow.New(pool)
	stateManager := state.NewStateManager()
	flowFinder := middleware.NewFlowFinder(flowManager)

	defaultFlowHandler := middleware.NewFlowHandler("default", stateManager)
	flowFinder.RegisterFlowHandler("default", defaultFlowHandler)

	defaultFlowHandler.AddStateHandler(stateFirst, handleFirst)
	defaultFlowHandler.AddStateHandler(stateSecond, handleSecond)
	defaultFlowHandler.AddStateHandler(stateThird, handlerThird)
	defaultFlowHandler.AddStateHandler(stateLast, handleLast)
	
//    handler toggles state
}


```


## Модификации
Как добавить другой storage для flow...

## todo
1. format as library
2. use interfaces for flow handling
3. simplify initialization
