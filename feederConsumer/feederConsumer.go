package feederConsumer

import (
	"sync"
)

type Feeders struct {
	Count    int
	Callback func() (interface{}, bool)
}

type Consumers struct {
	Count    int
	Callback func(interface{})
}

type FeederConsumer struct {
	WaitGroup   sync.WaitGroup
	DataChannel chan interface{}

	Feeders   []Feeders
	Consumers Consumers
}

func (fc *FeederConsumer) Run() {
	defer close(fc.DataChannel)

	for i := 0; i < fc.Consumers.Count; i++ {
		go func() {
			for newDataSet := range fc.DataChannel {
				fc.Consumers.Callback(newDataSet)

				fc.WaitGroup.Done()
			}
		}()
	}

	feederProcessWaitGroup := sync.WaitGroup{}

	for i := 0; i < len(fc.Feeders); i++ {
		feederProcessWaitGroup.Add(fc.Feeders[i].Count)

		for j := 0; j < fc.Feeders[i].Count; j++ {
			go func() {
				for {
					newDataSet, isFeedingFinished := fc.Feeders[i].Callback()

					fc.WaitGroup.Add(1)

					fc.DataChannel <- newDataSet

					if isFeedingFinished {
						break
					}
				}

				feederProcessWaitGroup.Done()
			}()
		}
	}

	feederProcessWaitGroup.Wait()

	fc.WaitGroup.Wait()
}
