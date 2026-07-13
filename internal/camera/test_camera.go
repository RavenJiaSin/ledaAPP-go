package camera

import (
	"sync"

	"gocv.io/x/gocv"
)


type TestCamera struct {
	mu sync.RWMutex

	frames map[string]gocv.Mat
}


func NewTestCamera() *TestCamera {

	return &TestCamera{
		frames: make(map[string]gocv.Mat),
	}
}


func (c *TestCamera) Open(
	name string,
	source interface{},
) error {
	return nil
}


func (c *TestCamera) Close(name string) {

	c.mu.Lock()
	defer c.mu.Unlock()

	if mat, ok := c.frames[name]; ok {
		mat.Close()
		delete(c.frames, name)
	}
}


func (c *TestCamera) CloseAll() {

	c.mu.Lock()
	defer c.mu.Unlock()

	for name, mat := range c.frames {
		mat.Close()
		delete(c.frames, name)
	}
}


func (c *TestCamera) LastFrame(
	name string,
) (gocv.Mat, bool) {

	c.mu.RLock()
	defer c.mu.RUnlock()

	mat, ok := c.frames[name]

	if !ok {
		return gocv.Mat{}, false
	}

	return mat.Clone(), true
}


func (c *TestCamera) InjectFrame(
	name string,
	mat gocv.Mat,
) {

	c.mu.Lock()
	defer c.mu.Unlock()

	if old, ok := c.frames[name]; ok {
		old.Close()
	}

	c.frames[name] = mat.Clone()
}
