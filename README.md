# Telegram Bot
Middleware для обработки flow с fsm для `tucnak/telebot`.

Пакет решает проблему сохранения состояния в сценариях, где нужно не терять контекст между сообщениями.  
Переключение состояний и инициализацию переходов нужно реализовывать самостоятельно в обработчиках (см AddHandler). 

В [примере](example/main.go) работа со переходами реализована через `looplab/fsm`. 
Можно использовать любое другое решение, т.к. декларация и специфическая работа с fsm вынесена в обработчики. 

## Установка
```bash
go get github.com/ulngollm/teleflow
```

## Features
1. Сохранение стейта между сообщениями
2. Простое и понятное добавление обработчиков для стейтов
3. Поддержка нескольких параллельных флоу в приложении

## Как использовать
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

