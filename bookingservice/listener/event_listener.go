package listener

import (
	"log"

	"github.com/GO_NATIVE/contracts"
	"github.com/GO_NATIVE/lib/msgqueue"
	"github.com/GO_NATIVE/lib/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventProcessor struct {
	EventListener msgqueue.EventListener
	Database      persistence.DatabaseHandler
}

func (p *EventProcessor) ProcessEvents() error {
	log.Println("Listening for events...")

	received, errors, err := p.EventListener.Listen("event.created")
	if err != nil {
		return err
	}

	for {
		select {
		case evt := <-received:
			p.handleEvent(evt)
		case err := <-errors:
			log.Printf("received error while processing event: %s", &err)
		}
	}
}

func (p *EventProcessor) handleEvent(event msgqueue.Event) {
	switch e := event.(type) {
	case *contracts.EventCreatedEvent:
		log.Printf("event %s created: %s", e.ID, e)
		obj_id, err := primitive.ObjectIDFromHex(e.ID)
		if err != nil {
			panic(err)
		}
		p.Database.AddEvent(persistence.Event{ID: obj_id})
	case *contracts.LocationCreatedEvent:
		log.Printf("location %s created: %s", e.ID, e)
		obj_id, err := primitive.ObjectIDFromHex(e.ID)
		if err != nil {
			panic(err)
		}
		p.Database.AddLocation(persistence.Location{ID: obj_id})
	default:
		log.Printf("unknown event %t", e)
	}
}
