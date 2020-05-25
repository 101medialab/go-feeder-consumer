package feederConsumer

import (
	"sync"
)

type FeederConsumer struct {
	WaitGroup        sync.WaitGroup
	DataChannel      chan interface{}
	NFeeders         int
	FeederCallback   func() (interface{}, bool)
	NConsumers       int
	ConsumerCallback func(interface{})
}

func (fc *FeederConsumer) Run() {
	defer close(fc.DataChannel)
	
	for i := 0; i < fc.NConsumers; i++ {
		go func() {
			for newDataSet := range fc.DataChannel {
				fc.ConsumerCallback(newDataSet)

				fc.WaitGroup.Done()
			}
		}()
	}

	feederProcessWaitGroup := sync.WaitGroup{}
	feederProcessWaitGroup.Add(fc.NFeeders)

	for i := 0; i < fc.NFeeders; i++ {
		go func() {
			for {
				newDataSet, isFeedingFinished := fc.FeederCallback()

				fc.WaitGroup.Add(1)

				fc.DataChannel <- newDataSet

				if isFeedingFinished {
					break
				}
			}

			feederProcessWaitGroup.Done()
		}()
	}

	feederProcessWaitGroup.Wait()

	fc.WaitGroup.Wait()
}
