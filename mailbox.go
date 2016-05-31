package goldengine

import "sync"

type MessageType string
type Message struct {
	Message MessageType
	Content interface{}
}

type PostOffice struct {
	penpals map[*MailBox]map[MessageType]map[*MailBox]struct{}
	mu      *sync.RWMutex
}

func (office *PostOffice) RecieveMessage(sender *MailBox, msg Message) {

	boxes := office.GetListeners(sender, msg.Message)
	for box, _ := range boxes {
		box.PostMessage(msg)
	}
}

func (office *PostOffice) AssignMailBox(box *MailBox) {
	box.Remove()
	box.office = office
}

func (office *PostOffice) GetListeners(sender *MailBox, msg MessageType) map[*MailBox]struct{} {
	office.mu.RLock()
	defer office.mu.RUnlock()
	penpals, ok := office.penpals[sender]
	if !ok || penpals == nil {
		return nil
	}
	boxes, _ := penpals[msg]
	return boxes
}

func (office *PostOffice) Subscribe(sender *MailBox, reciever *MailBox, msg MessageType) {
	isListening := office.IsSubscribed(sender, reciever, msg)
	if isListening {
		return
	}
	office.mu.Lock()
	defer office.mu.Lock()
	penpals, ok := office.penpals[sender]
	if !ok || penpals == nil {
		return
	}
	boxes, ok := penpals[msg]
	if !ok || boxes == nil {
		penpals[msg] = make(map[*MailBox]struct{})
		return
	}
	penpals[msg][reciever] = struct{}{}
}

func (office *PostOffice) IsSubscribed(sender *MailBox, reciever *MailBox, msg MessageType) bool {
	office.mu.RLock()
	defer office.mu.RUnlock()
	penpals, ok := office.penpals[sender]
	if !ok || penpals == nil {
		return false
	}
	boxes, ok := penpals[msg]
	if !ok || boxes == nil {
		return false
	}
	_, ok = boxes[reciever]
	if ok {
		return true
	}
	return false
}

func (office *PostOffice) Remove(box *MailBox) {
	if box.office == office {
		delete(office.penpals, box)
		for _, penpals := range office.penpals {
			for _, boxes := range penpals {
				delete(boxes, box)
			}
		}
		box.office = nil
	}
}

func NewPostOffice() *PostOffice {
	return &PostOffice{
		mu:      &sync.RWMutex{},
		penpals: make(map[*MailBox]map[MessageType]map[*MailBox]struct{}),
	}
}

type MailBox struct {
	office *PostOffice
}

func (box *MailBox) PostMessage(msg Message) {
	if box.office != nil {
		box.office.RecieveMessage(box, msg)
	}
}

func (box *MailBox) RecieveMessage(msg Message) {

}

func (box *MailBox) Remove() {
	if box.office != nil {
		box.office.Remove(box)
	}
}
