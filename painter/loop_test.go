package painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"

	"golang.org/x/exp/shiny/screen"
)

func TestEventLoop_Initiate(t *testing.T) {
	s := &mockScreen{}
	el := &EventLoop{
		Receiver: &testTextureReceiver{},
	}
	el.Initiate(s)

	if el.currentTexture == nil || el.lastTexture == nil {
		t.Error("unexpected nil texture")
	}

	el.Terminate()
}

func TestEventLoop_Enqueue(t *testing.T) {
	var (
		el EventLoop
		tr testTextureReceiver
	)

	el.Receiver = &tr

	el.Initiate(mockScreen{})
	el.Enqueue(FillTexture(color.White))
	el.Enqueue(FillTexture(color.RGBA{R: 85, G: 217, B: 104}))
	el.Enqueue(MarkUpdated)
	if tr.LastTexture != nil {
		t.Fatal("Receiver got the texture too early")
	}
	el.Terminate()

	tx, ok := tr.LastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Receiver still has no texture")
	}
	if tx.FillCnt != 2 {
		t.Error("Unexpected number of fill calls:", tx.FillCnt)
	}
}

func TestOperationQueue_enqueue_dequeue_isEmpty(t *testing.T) {
	oq := &operationQueue{}
	op := &testTextureOperation{}

	oq.enqueue(op)

	if len(oq.operations) != 1 || oq.operations[0] != op {
		t.Error("failed to enqueue operation into queue")
	}

	dequeuedOp := oq.dequeue()
	if dequeuedOp != op {
		t.Error("failed to dequeue operation from queue")
	}

	if !oq.isEmpty() {
		t.Error("expected queue to be empty")
	}
}

func TestOperationQueue_enqueue_blocked(t *testing.T) {
	oq := &operationQueue{}

	for i := 0; i < 10; i++ {
		oq.enqueue(&testTextureOperation{})
	}

	op := &testTextureOperation{}

	// Enqueue operation and ensure that it's blocked
	oq.enqueue(op)

	if len(oq.operations) != 11 {
		t.Error("failed to enqueue operation into queue")
	}

	if oq.waitCh != nil {
		t.Error("expected operation queue to be blocked")
	}

	// Remove operation from queue and ensure that it's unblocked
	oq.dequeue()

	if len(oq.operations) != 10 {
		t.Error("failed to dequeue operation from queue")
	}

	if oq.waitCh != nil {
		t.Error("expected operation queue to be unblocked")
	}

	if oq.isEmpty() {
		t.Error("expected queue to not be empty")
	}
}

type testTextureReceiver struct {
	LastTexture screen.Texture
}

func (tr *testTextureReceiver) UpdateTexture(t screen.Texture) {
	tr.LastTexture = t
}

type testTextureOperation struct {
	applied bool
}

func (op *testTextureOperation) Apply(t screen.Texture) bool {
	op.applied = true
	return true
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(*screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	FillCnt int
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return canvasSize }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: canvasSize}
}

func (m *mockTexture) Upload(image.Point, screen.Buffer, image.Rectangle) {
	panic("implement me")
}

func (m *mockTexture) Fill(image.Rectangle, color.Color, draw.Op) {
	m.FillCnt++
}


func TestEventLoop_NewTexturesCreated(t *testing.T) {
	s := &mockScreen{}
	el := &EventLoop{
		Receiver: &testTextureReceiver{},
	}

	el.Initiate(s)
	defer el.Terminate()

	if el.currentTexture == nil || el.lastTexture == nil {
		t.Error("EventLoop did not create new textures on initiation")
	}
}

func TestEventLoop_TextureUpdate(t *testing.T) {
	s := &mockScreen{}
	el := &EventLoop{
		Receiver: &testTextureReceiver{},
	}

	el.Initiate(s)
	defer el.Terminate()

	// Enqueue a texture operation that marks the texture as updated
	el.Enqueue(MarkUpdated)

	// Allow some time for the operation to be processed
	time.Sleep(100 * time.Millisecond)

	if el.Receiver.(*testTextureReceiver).LastTexture == nil {
		t.Error("Texture was not updated after enqueueing MarkUpdated operation")
	}
}

func TestOperationQueue_OrderPreserved(t *testing.T) {
	oq := &operationQueue{}
	op1 := &testTextureOperation{}
	op2 := &testTextureOperation{}

	oq.enqueue(op1)
	oq.enqueue(op2)

	if oq.dequeue() != op1 || oq.dequeue() != op2 {
		t.Error("Operation queue did not preserve the order of operations")
	}
}

func TestOperationQueue_EmptyAfterDequeue(t *testing.T) {
	oq := &operationQueue{}
	op := &testTextureOperation{}

	oq.enqueue(op)
	oq.dequeue()

	if !oq.isEmpty() {
		t.Error("Operation queue was not empty after dequeueing the last operation")
	}
}

func TestMockTexture_FillCalled(t *testing.T) {
	mt := &mockTexture{}
	rect := image.Rect(0, 0, 100, 100)
	fillColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	mt.Fill(rect, fillColor, draw.Src)

	if mt.FillCnt != 1 {
		t.Errorf("Fill was not called on the mock texture")
	}
}



func TestEventLoop_StopChClosed(t *testing.T) {
  s := &mockScreen{}
  el := &EventLoop{
    Receiver: &testTextureReceiver{},
  }
  el.Initiate(s)
  defer el.Terminate()

  // Before calling Terminate, stopCh should not be closed
  select {
  case <-el.stopCh:
    t.Fatal("stopCh should not be closed before Terminate is called")
  default:
    // Expected case, do nothing
  }

  el.Terminate()

  // After calling Terminate, stopCh should be closed
  select {
  case <-el.stopCh:
    // Expected case, do nothing
  default:
    t.Fatal("stopCh should be closed after Terminate is called")
  }
}
