package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/GO_NATIVE/lib/msgqueue"
	"github.com/GO_NATIVE/lib/persistence"
)

func ServeAPI(endpoint, tlsendpoint string, databasehandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) (chan error, chan error) {
	handler := NewEventHandler(databasehandler, eventEmitter)
	r := mux.NewRouter()
	httpErrChan := make(chan error)
	httptlsErrChan := make(chan error)

	eventsRouter := r.PathPrefix("/events").Subrouter()
	eventsRouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)
	eventsRouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)
	go func() { httpErrChan <- http.ListenAndServe(endpoint, r) }()
	go func() { httptlsErrChan <- http.ListenAndServeTLS(tlsendpoint, "cert.pem", "key.pem", r) }()
	return httpErrChan, httptlsErrChan
}
