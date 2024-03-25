package console

import (
	"fmt"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/message"
)

type consolePublisher struct {
	logger interfaces.Logger
}

func NewConsolePublisher(logger interfaces.Logger) interfaces.StreamProvider {
	return &consolePublisher{
		logger: logger,
	}
}

func (p *consolePublisher) Connect() error {
	return nil
}

func (p *consolePublisher) Publish(msg message.PublishWrapper) error {
	fmt.Printf("action: %s, messageType: %s, %v\n", msg.Action, msg.MessageType, string(msg.Content))
	return nil
}

func (p *consolePublisher) Close() error {
	return nil
}
