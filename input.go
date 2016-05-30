package goldengine

import (
	"sync/atomic"

	sf "github.com/manyminds/gosfml"
)

var keyboardSetCounter uint32 = 1

//InputCollection : Collection of Sets handling Mouse and Keyboard Input.
//Used to emulate keyboard input.
type InputCollection struct {
	keyboardSets map[uint32]*KeyboardSet
}

//GenInputCollection : Convienece function for making an InputCollection
func GenInputCollection() *InputCollection {
	collection := &InputCollection{
		keyboardSets: make(map[uint32]*KeyboardSet),
	}
	return collection
}

//KeyPressed : Puts KeyPressed
func (collection *InputCollection) KeyPressed(code sf.KeyCode) {
	for _, set := range collection.keyboardSets {
		if set.active {
			set.KeyPressed(code)
		}
	}
}

//KeyReleased : Puts KeyReleased
func (collection *InputCollection) KeyReleased(code sf.KeyCode) {
	for _, set := range collection.keyboardSets {
		if set.active {
			set.KeyReleased(code)
		}
	}
}

//ActivateSet : KeyboardSet will recive input
func (collection *InputCollection) ActivateSet(set *KeyboardSet) {
	set.active = true
}

//DeactivateSet : KeyboardSet will nolong recieve input
func (collection *InputCollection) DeactivateSet(set *KeyboardSet) {
	set.active = false
}

//DeactivateSets : All sets will no longer recive input
func (collection *InputCollection) DeactivateSets() {
	for _, set := range collection.keyboardSets {
		set.active = false
	}
}

//InstallKeyboardSet : Add set to collection
func (collection *InputCollection) InstallKeyboardSet(set *KeyboardSet) {
	collection.keyboardSets[set.id] = set
	set.active = true
}

//UninstallKeyboardSet : Remove set from collection
func (collection *InputCollection) UninstallKeyboardSet(set *KeyboardSet) {
	delete(collection.keyboardSets, set.id)
}

//KeyboardSet : Set of KeyboardHandlers. Technically only need 1 but it's usefull
//for grouping
type KeyboardSet struct {
	id       uint32
	active   bool
	handlers map[uint32]*KeyboardHandler
}

//KeyPressed : Says a key should be pressed
func (set *KeyboardSet) KeyPressed(code sf.KeyCode) {
	for _, handler := range set.handlers {
		if cmd, ok := handler.KeyPressedCommands[code]; ok {
			go cmd()
		}
	}
}

//KeyReleased : Says a key should be released
func (set *KeyboardSet) KeyReleased(code sf.KeyCode) {
	for _, handler := range set.handlers {
		if cmd, ok := handler.KeyReleasedCommands[code]; ok {
			go cmd()
		}
	}
}

//AddHandler : adds a Handler to a keyboardSet
func (set *KeyboardSet) AddHandler(handler *KeyboardHandler) {
	set.handlers[handler.id] = handler
}

//RemoveHandler : remove a KeybaordHandler from a keyboardSet
func (set *KeyboardSet) RemoveHandler(handler *KeyboardHandler) {
	delete(set.handlers, handler.id)
}

//GenKeyboardSet : Creates a New KeyboardSet
func GenKeyboardSet() *KeyboardSet {
	id := keyboardSetCounter
	atomic.AddUint32(&keyboardSetCounter, 1)
	return &KeyboardSet{
		id:       id,
		handlers: make(map[uint32]*KeyboardHandler),
	}
}

var keyboardHandlerCounter uint32 = 1

//KeyboardHandler : functions to call whe na key is pressed and released
type KeyboardHandler struct {
	id                  uint32
	KeyPressedCommands  map[sf.KeyCode]KeyCommand
	KeyReleasedCommands map[sf.KeyCode]KeyCommand
}

//RegisterKeyPressedCommand : call htis function when a key si pressed
func (handler *KeyboardHandler) RegisterKeyPressedCommand(code sf.KeyCode, cmd KeyCommand) {
	handler.KeyPressedCommands[code] = cmd
}

//RegisterKeyReleasedCommand : call this function when a key is released
func (handler *KeyboardHandler) RegisterKeyReleasedCommand(code sf.KeyCode, cmd KeyCommand) {
	handler.KeyReleasedCommands[code] = cmd
}

//ClearKeyPressedCommand : Clear function ofr this key
func (handler *KeyboardHandler) ClearKeyPressedCommand(code sf.KeyCode) {
	delete(handler.KeyPressedCommands, code)
}

//ClearKeyReleasedCommand : Clear function for this key
func (handler *KeyboardHandler) ClearKeyReleasedCommand(code sf.KeyCode) {
	delete(handler.KeyReleasedCommands, code)
}

//GenInputHandler : Creates a Handler for commands
func GenInputHandler() *KeyboardHandler {
	id := keyboardHandlerCounter
	atomic.AddUint32(&keyboardHandlerCounter, 1)
	return &KeyboardHandler{
		id:                  id,
		KeyPressedCommands:  make(map[sf.KeyCode]KeyCommand),
		KeyReleasedCommands: make(map[sf.KeyCode]KeyCommand),
	}
}

//KeyCommand : Function called when a button is pressed or released
type KeyCommand func()
