package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/GO_NATIVE/lib/msgqueue"
	"github.com/GO_NATIVE/lib/persistence"
)

func ServeAPI(endpoint string, databasehandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) (chan error, chan error) {
	handler := NewEventHandler(databasehandler, eventEmitter)
	r := mux.NewRouter()
	httpErrChan := make(chan error)
	httptlsErrChan := make(chan error)

	eventsRouter := r.PathPrefix("/events").Subrouter()
	eventsRouter.Methods("POST").Path("/{eventID}/bookings").HandlerFunc(handler.newEventHandler)
	go func() { httpErrChan <- http.ListenAndServe(endpoint, r) }()
	return httpErrChan, httptlsErrChan
}
