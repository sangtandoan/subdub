package chrono

import (
	"log/slog"
	"sync"
)

type Background struct {
	Wg     *sync.WaitGroup
	logger *slog.Logger
}

func NewBackground() *Background {
	return &Background{
		Wg: &sync.WaitGroup{},
	}
}

func (b *Background) Run(fn func()) {
	b.Wg.Add(1)

	go func() {
		defer b.Wg.Done()

		defer func() {
			if err := recover(); err != nil {
				b.logger.Error("Background task panicked", slog.Any("error", err))
			}
		}()

		fn()
	}()
}
