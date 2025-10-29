package processor

import "github.com/MGYOSBEL/pathfinder/pkg/pubsub"

type Processor struct {
	options Options
	output  pubsub.Publisher
	input   pubsub.Subscriber
}

type Options struct {
	InputTopic  string
	OutputTopic string
}

func New(p pubsub.Publisher, s pubsub.Subscriber, opts Options) *Processor {
	return &Processor{
		options: opts,
		output:  p,
		input:   s,
	}
}

func (p *Processor) Process() error {
	p.Forward()
	return nil
}
