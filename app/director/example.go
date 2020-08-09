package director

import (
	"context"
	"log"
	"time"
)

// ExampleDirector является абстрактным приложением.
type ExampleDirector struct {
	ctx context.Context
}

// NewServerApp - конструктор приложения.
func NewServerApp(ctx context.Context) *ExampleDirector {
	return &ExampleDirector{ctx: ctx}
}

// Run - основной цикл программы.
func (a *ExampleDirector) Run() error {
	go a.exampleWorker()

	<-a.ctx.Done()
	log.Println("[DEBUG] application terminated")
	a.Shutdown()
	return nil
}

// Shutdown - выключает сервер, закрывая канал.
func (a *ExampleDirector) Shutdown() {
	log.Println("[DEBUG] ExampleDirector shutdown successfully")
}

// exampleWorker выводит сообщение о своей работе на экран.
func (a *ExampleDirector) exampleWorker() {
	for {
		select {
		case <-time.After(time.Second):
			log.Println("[INFO] work in progress...")
		case <-a.ctx.Done():
			break
		}
	}
}
