package processor

import (
	"fmt"

	"github.com/MGYOSBEL/pathfinder/pkg/message"
)

func (p *Processor) Forward() {
	p.input.Subscribe(func(msg message.Message) {
		p.output.Publish(fmt.Sprintf("%s/%s", p.options.OutputTopic, msg.Topic), msg)
	})
}
