package painter

import (
	"image"
	"log"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// TextureReceiver receives a texture that has been prepared as a result of executing operations in the event loop.
type TextureReceiver interface {
	UpdateTexture(t screen.Texture)
}

// EventLoop manages an event loop for creating textures through executing operations from an internal queue.
type EventLoop struct {
	Receiver TextureReceiver

	currentTexture screen.Texture // Texture currently being formed
	lastTexture    screen.Texture // Texture last sent to the Receiver

	opQueue operationQueue

	stopCh  chan struct{}
	requestStop bool
}

var canvasSize = image.Pt(800, 800)

// Initiate starts the event loop. This method should be called before any other methods on it.
func (el *EventLoop) Initiate(screenProvider screen.Screen) {
	el.currentTexture, _ = screenProvider.NewTexture(canvasSize)
	el.lastTexture, _ = screenProvider.NewTexture(canvasSize)

	el.stopCh = make(chan struct{})

	go func() {
		for !el.requestStop || !el.opQueue.isEmpty() {
			op := el.opQueue.dequeue()
			ready := op.Apply(el.currentTexture)

			if ready {

				log.Println("Texture updated, calling UpdateTexture")

				el.Receiver.UpdateTexture(el.currentTexture)
				el.currentTexture, el.lastTexture = el.lastTexture, el.currentTexture

				log.Println("Texture swap complete")
			}
		}
		close(el.stopCh)
	}()
}

// Enqueue adds a new operation to the internal queue.
func (el *EventLoop) Enqueue(op TextureOperation) {
	el.opQueue.enqueue(op)
}

// Terminate signals the event loop to stop and waits for it to finish.
func (el *EventLoop) Terminate() {
	el.Enqueue(TextureFunc(func(t screen.Texture) {
		el.requestStop = true
	}))
	<-el.stopCh
}

// operationQueue is a custom message queue for texture operations.
type operationQueue struct {
	operations []TextureOperation
	mutex      sync.Mutex
	waitCh     chan struct{}
}

func (oq *operationQueue) enqueue(op TextureOperation) {
	oq.mutex.Lock()
	defer oq.mutex.Unlock()

	oq.operations = append(oq.operations, op)

	if oq.waitCh != nil {
		close(oq.waitCh)
		oq.waitCh = nil
	}
}

func (oq *operationQueue) dequeue() TextureOperation {
	oq.mutex.Lock()
	defer oq.mutex.Unlock()

	for len(oq.operations) == 0 {
		oq.waitCh = make(chan struct{})
		oq.mutex.Unlock()
		<-oq.waitCh
		oq.mutex.Lock()
	}

	op := oq.operations[0]
	oq.operations[0] = nil
	oq.operations = oq.operations[1:]
	return op
}

func (oq *operationQueue) isEmpty() bool {
	oq.mutex.Lock()
	defer oq.mutex.Unlock()

	return len(oq.operations) == 0
}
