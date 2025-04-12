# Teleflow
```bash
go get github.com/ulngollm/teleflow
```
## Overview
Простой и легковесный middleware для Telegram ботов на основе `tucnak/telebot`.

Позволяет сохранять состояние (контекст) между сообщениями. 

Middleware облегчает добавление обработчиков для разных состояний и поддерживает несколько параллельных флоу в приложении. 
Интеграция с любым решением FSM (например, `looplab/fsm`) позволяет реализовывать произвольные сценарии переходов, оставляя логику управления состояниями на усмотрение разработчика.

## Features
1. Сохранение стейта между сообщениями
2. Простое и понятное добавление обработчиков для стейтов
3. Поддержка нескольких параллельных флоу в приложении

## Getting Started
Добавление обработчиков реализовано как добавление роутов в популярных API-фреймворках:

```go
package main

import (
	"github.com/ulngollm/teleflow"
	tele "gopkg.in/telebot.v4"
)

func run() {
	pool := teleflow.NewMemoryPool()
	flowManager := teleflow.NewFlowManager(pool)
	router := teleflow.NewFlowRouter(flowManager)

	g := router.Group("default") // flow name
	g.AddHandler("opened", handle)
	g.AddHandler("waiting", handle)
	g.AddHandler("closed", handle)
}

func handle(c tele.Context) error {
    return c.Send(teleflow.GetCurrentFlow(c))
}

```
Рабочий пример - в [main.go](example/main.go)

