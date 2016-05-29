package goldengine

import "sync/atomic"
import sf "github.com/manyminds/gosfml"

var inputSetCounter uint32 = 1

type inputSet struct {
	id       uint32
	Handlers map[uint32]*inputHandler
}

func GenInputSet() *inputSet {
	id := inputSetCounter
	atomic.AddUint32(&inputSetCounter, 1)
	return &inputSet{
		id:       id,
		Handlers: make(map[uint32]*inputHandler),
	}
}

var inputHandlerCounter uint32 = 1

type inputHandler struct {
	id       uint32
	Commands map[sf.KeyCode]KeyCommandSet
}

func GenInputHandler() *inputHandler {
	id := inputHandlerCounter
	atomic.AddUint32(&inputHandlerCounter, 1)
	return &inputHandler{
		id:       id,
		Commands: make(map[sf.KeyCode]KeyCommandSet),
	}
}

type KeyCommand func()
type KeyCommandSet struct {
	KeyPressed  KeyCommand
	KeyReleased KeyCommand
}
