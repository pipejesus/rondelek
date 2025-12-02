package providers

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Grid interface {
	Rectangle(startColumn, endColumn, startRow, endRow int) rl.Rectangle
	GetColumns() int
	GetRows() int
	DrawDebug()
}

type Sampler interface {
	Quit()
	Record()
	StopRecording() int
	SaveCurrentSample()
	PlaySample(idx int) error
}

type Container struct {
	Grid    Grid
	Sampler Sampler
}

var (
	mu        sync.Mutex
	container *Container
)

func SetContainer(c *Container) {
	mu.Lock()
	defer mu.Unlock()
	container = c
}

func GetContainer() *Container {
	mu.Lock()
	c := container
	mu.Unlock()
	if c == nil {
		panic("providers: container is not initialized (SetContainer not called)")
	}
	return c
}

func ResetContainer() {
	mu.Lock()
	container = nil
	mu.Unlock()
}
