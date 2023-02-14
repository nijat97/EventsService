package mongolayer

import (
	"context"
	"fmt"
	"time"

	"github.com/GO_NATIVE/lib/persistence"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB        = "myevents"
	USERS     = "users"
	EVENTS    = "events"
	LOCATIONS = "locations"
)

type MongoDBLayer struct {
	client *mongo.Client
}

func NewMongoDBLayer(connection string) (*MongoDBLayer, error) {
	clt, err := mongo.NewClient(options.Client().ApplyURI(connection))
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	err = clt.Connect(ctx)
	if err != nil {
		fmt.Printf("Error occured: ", err)
		return nil, err
	}

	return &MongoDBLayer{
		client: clt,
	}, err
}

func (mgoLayer *MongoDBLayer) AddUser(u persistence.User) ([]byte, error) {
	u.ID = primitive.NewObjectID()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := mgoLayer.client.Database(DB).Collection(USERS).InsertOne(ctx, u)
	return []byte(u.ID.Hex()), err
}

func (mgoLayer *MongoDBLayer) AddLocation(l persistence.Location) (persistence.Location, error) {
	l.ID = primitive.NewObjectID()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := mgoLayer.client.Database(DB).Collection(LOCATIONS).InsertOne(ctx, l)
	return l, err
}

func (mgoLayer *MongoDBLayer) AddEvent(e persistence.Event) ([]byte, error) {
	if !e.ID.IsZero() {
		e.ID = primitive.NewObjectID()
	}

	if !e.Location.ID.IsZero() {
		e.Location.ID = primitive.NewObjectID()
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := mgoLayer.client.Database(DB).Collection(EVENTS).InsertOne(ctx, e)
	byteSlice := []byte(e.ID.Hex())
	return byteSlice, err
}

func (mgoLayer *MongoDBLayer) AddBookingForUser(id []byte, bk persistence.Booking) error {
	obj_id, err := primitive.ObjectIDFromHex(string(id))
	filter := bson.M{"_id": obj_id}
	update := bson.M{"$addToSet": bson.M{"bookings": []persistence.Booking{bk}}}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = mgoLayer.client.Database(DB).Collection(USERS).UpdateOne(ctx, filter, update)
	return err
}

func (mgoLayer *MongoDBLayer) FindUser(first string, last string) (persistence.User, error) {
	u := persistence.User{}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := mgoLayer.client.Database(DB).Collection(USERS).FindOne(ctx, bson.M{"first": first, "last": last}).Decode(&u)
	return u, err
}

func (mgoLayer *MongoDBLayer) FindBookingsForUser(id []byte) ([]persistence.Booking, error) {
	u := persistence.User{}
	obj_id, err := primitive.ObjectIDFromHex(string(id))
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = mgoLayer.client.Database(DB).Collection(USERS).FindOne(ctx, obj_id).Decode(&u)
	return u.Bookings, err
}

// we pass id as []byte instead of bson object cause we want to keep FindEvent in lib/persistence.go stays as generic as possible.
// an interface{} can also be used instead of byte. in Go, it can be converted to any other type
func (mgoLayer *MongoDBLayer) FindEvent(id []byte) (persistence.Event, error) {
	e := persistence.Event{}
	obj_id, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		e.Location.Name = string(id)
	}
	filter := bson.M{"_id": obj_id}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = mgoLayer.client.Database(DB).Collection(EVENTS).FindOne(ctx, filter).Decode(&e)
	if err == mongo.ErrNoDocuments {
		e.ID = obj_id
		e.Name = "not found"
	}
	return e, err
}

// bson.M{} is a type(map) to represent query parameters. bson.M{"field_name": value}
func (mgoLayer *MongoDBLayer) FindEventByName(name string) (persistence.Event, error) {
	e := persistence.Event{}
	filter := bson.M{"name": name}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := mgoLayer.client.Database(DB).Collection(EVENTS).FindOne(ctx, filter).Decode(&e)
	return e, err
}

func (mgoLayer *MongoDBLayer) FindAllAvailableEvents() ([]persistence.Event, error) {
	events := []persistence.Event{}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := mgoLayer.client.Database(DB).Collection(EVENTS).Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("Error:", err)
	}
	cursor.All(ctx, &events)
	return events, err
}

func (mgoLayer *MongoDBLayer) FindLocation(id string) (persistence.Location, error) {
	location := persistence.Location{}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := mgoLayer.client.Database(DB).Collection(LOCATIONS).FindOne(ctx, bson.M{"_id": id}).Decode(&location)
	return location, err
}

func (mgoLayer *MongoDBLayer) FindAllLocations() ([]persistence.Location, error) {
	locations := []persistence.Location{}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := mgoLayer.client.Database(DB).Collection(LOCATIONS).Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("Error:", err)
	}
	cursor.All(ctx, &locations)
	return locations, err
}
