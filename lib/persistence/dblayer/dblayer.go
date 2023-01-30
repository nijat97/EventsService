//build a database layer package that acts as the gateway to the
//persistence layer in our microservice. The package will utilize the factory design pattern by
//implementing a factory function.

package dblayer

import (
	"github.com/GO_NATIVE/lib/persistence"
	"github.com/GO_NATIVE/lib/persistence/mongolayer"
)

type DBTYPE string

const (
	MONGODB  DBTYPE = "mongodb"
	DYNAMODB DBTYPE = "dynamodb"
)

func NewPersistenceLayer(options DBTYPE, connection string) (persistence.DatabaseHandler, error) {
	switch options {
	case MONGODB:
		return mongolayer.NewMongoDBLayer(connection)
	}
	return nil, nil
}
