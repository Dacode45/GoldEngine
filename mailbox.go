package goldengine

import (
	"log"
	"sync"
	"sync/atomic"
)

//MessageType : String Defining the Message
type MessageType string

//Message : Contains two fields, MessageType and the pointer to
//the object you are dealing with. Used for RPC.
type Message struct {
	Message MessageType
	Content interface{}
}

//Address : uint32 unique to every box
type Address uint32

var runningAddress uint32 = 1

//PostOffice : Manages MailBoxs
type PostOffice struct {
	penpals   map[Address]map[MessageType]map[Address]struct{}
	mailboxes map[Address]MailBox
	logger    *log.Logger
	mu        sync.RWMutex
}

//Broadcast : Send a message to every box
func (office *PostOffice) Broadcast(msg Message) {
	if office.logger != nil {
		office.logger.Printf("Office %p Broadcasting Message: %s\n", office, msg.Message)
	}
	for _, box := range office.mailboxes {
		box.RecieveMessage(msg)
	}
}

//RecieveMessage : Called whenever a MailBox post. Sends this Message out to
//Subscribed MailBoxs
func (office *PostOffice) RecieveMessage(sender Address, msg Message) {
	if office.logger != nil {
		office.logger.Printf("Office %p Recieved Message: %s from Address: %v\n", office, msg.Message, sender)
	}
	boxes := office.getSubscribers(sender, msg.Message)
	for address := range boxes {
		box, ok := office.mailboxes[address]
		if ok {
			box.RecieveMessage(msg)
		}
	}
}

//GetSubscribers : Get all subscribers to this mailbox for a specefic message
func (office *PostOffice) getSubscribers(sender Address, msg MessageType) map[Address]struct{} {
	office.mu.RLock()
	defer office.mu.RUnlock()
	penpals, ok := office.penpals[sender]
	if !ok || penpals == nil {
		return nil
	}
	boxes, _ := penpals[msg]
	return boxes
}

//Subscribe : Tells a MailBox to Listen to a specific message from this mailbox
func (office *PostOffice) Subscribe(sender Address, reciever Address, msg MessageType) {
	office.mu.Lock()
	defer office.mu.Unlock()
	penpals, ok := office.penpals[sender]
	if !ok || penpals == nil {
		return
	}
	boxes, ok := penpals[msg]
	if !ok || boxes == nil {
		penpals[msg] = make(map[Address]struct{})
	}
	penpals[msg][reciever] = struct{}{}
}

//UnSubscribe : Removes Subscriber for a given message and MailBox
func (office *PostOffice) UnSubscribe(sender Address, reciever Address, msg MessageType) {
	boxes := office.getSubscribers(sender, msg)
	office.mu.Lock()
	defer office.mu.Unlock()
	delete(boxes, reciever)

}

//IsSubscribed : See if a given mailbox is the subscriber for a message on another
func (office *PostOffice) IsSubscribed(sender Address, reciever Address, msg MessageType) bool {
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

//Add : Adds a MailBox to this office
func (office *PostOffice) Add(box MailBox) {
	box.Remove()
	office.mu.Lock()
	defer office.mu.Unlock()
	id := runningAddress
	atomic.AddUint32(&runningAddress, 1)
	box.SetAddress(Address(id))
	office.penpals[Address(id)] = make(map[MessageType]map[Address]struct{})
	office.mailboxes[Address(id)] = box
	box.SetOffice(office)
}

//Remove : PostOffice will never send this MailBox another message
func (office *PostOffice) Remove(box Address) {
	mailbox, ok := office.mailboxes[box]
	if ok {
		mailbox.Remove()
		office.mu.Lock()
		defer office.mu.Unlock()
		delete(office.penpals, box)
		for _, penpals := range office.penpals {
			for _, boxes := range penpals {
				delete(boxes, box)
			}
		}
		delete(office.mailboxes, box)
	}

}

//NewPostOffice : Creates a NewPostOffice
func NewPostOffice() *PostOffice {
	return &PostOffice{
		mu:        sync.RWMutex{},
		penpals:   make(map[Address]map[MessageType]map[Address]struct{}),
		mailboxes: make(map[Address]MailBox),
	}
}

//MailBox : Interface for accepting and recieve message
type MailBox interface {
	PostMessage(msg Message)
	RecieveMessage(msg Message)
	Remove()
	SetOffice(office *PostOffice)
	GetOffice() *PostOffice
	GetAddress() Address
	SetAddress(Address)
}

//BasicMailBox : Struct that accepts and recieves messages
type BasicMailBox struct {
	address Address
	office  *PostOffice
}

//SetAddress : Sets the Address of Mailbox
func (box *BasicMailBox) SetAddress(a Address) {
	box.address = a
}

//GetAddress : Gets the Address of Mailbox
func (box *BasicMailBox) GetAddress() Address {
	return box.address
}

//PostMessage : Send Out a Message
func (box *BasicMailBox) PostMessage(msg Message) {
	if box.office != nil {
		box.office.RecieveMessage(box.GetAddress(), msg)
	}
}

//RecieveMessage : Recieve a Message. Should Override
func (box *BasicMailBox) RecieveMessage(msg Message) {
}

//Remove : mailbox removes itself from its office
func (box *BasicMailBox) Remove() {
	box.office = nil
}

//GetOffice : Returns the mailbox office
func (box *BasicMailBox) GetOffice() *PostOffice {
	return box.office
}

//SetOffice : Sets this Boxes Office
func (box *BasicMailBox) SetOffice(office *PostOffice) {
	box.office = office
}
