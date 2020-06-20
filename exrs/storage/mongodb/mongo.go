package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/maxvw8/exercise_lib/exrs/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const colName = "exercises"

//Storage manages all interactions to the collection
type Storage struct {
	*mongo.Collection
	client *mongo.Client
}

//New creates an instance based on a DB connection
func New(database string) (*Storage, error) {
	var client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo db client. Error %v", err)
	}
	//creating connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database. Error %v", err)
	}
	// Check the connection
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database. Error %v", err)
	}
	db := client.Database(database)
	//init collection
	col := db.Collection(colName)
	return &Storage{col, client}, nil
}

//Provide CRUD

//Create a new Exercise
func (lib *Storage) Create(e *storage.Exercise) (*storage.Exercise, error) {
	r, err := lib.InsertOne(context.Background(), e)
	if err != nil {
		return nil, fmt.Errorf("failed to create new exercise %v. Error was %v", e, err)
	}
	e.Id = r.InsertedID.(primitive.ObjectID).Hex()
	return e, nil
}

//Read an exercise by id
func (lib *Storage) Read(hID string) (*storage.Exercise, error) {
	id, err := primitive.ObjectIDFromHex(hID)
	if err != nil {
		return nil, fmt.Errorf("unparseable id %v. Error was %v", id, err)
	}
	filter := bson.D{{Key: "_id", Value: id}}
	var exe *storage.Exercise
	err = lib.FindOne(context.Background(), filter).Decode(&exe)
	if err != nil {
		return nil, fmt.Errorf("could not find record by id %v. Error was %v", id, err)
	}
	return exe, nil
}

//Update an Exercise
func (lib *Storage) Update(id string, e *storage.Exercise) (*storage.Exercise, error) {
	idh, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unparseable id %v. Error was %v", idh, err)
	}
	filter := bson.D{{Key: "_id", Value: idh}}
	var updated *storage.Exercise
	err = lib.FindOneAndUpdate(context.Background(),
		filter,
		bson.M{"$set": e},
		options.FindOneAndUpdate().SetReturnDocument(1)).Decode(&updated)
	if err != nil {
		return nil, fmt.Errorf("could not update record %v. Error was %v", e, err)
	}
	return updated, nil
}

//Delete an Exercise
func (lib *Storage) Delete(id string) (bool, error) {
	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("unparseable id %v. Error was %v", pid, err)
	}
	filter := bson.D{{Key: "_id", Value: pid}}
	_, err = lib.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, fmt.Errorf("could not find record by id %v. Error was %v", pid, err)
	}
	return true, nil
}

//List obtains all the exercises
func (lib *Storage) List() ([]*storage.Exercise, error) {
	cursor, err := lib.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("could not find records. %v", err)
	}
	var exes []*storage.Exercise
	if err = cursor.All(context.Background(), &exes); err != nil {
		return nil, fmt.Errorf("could not parse records. %v", err)
	}
	return exes, nil
}

//older api

// GetByType returns a list of exercises matching the type or error in db connection issues
func (lib *Storage) GetByType(t string) ([]storage.Exercise, error) {
	fmt.Printf("filtering by: %s\n", t)
	filter := bson.D{{Key: "type", Value: t}}
	cursor, err := lib.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("could not find records. %v", err)
	}
	var exes []storage.Exercise
	if err = cursor.All(context.Background(), &exes); err != nil {
		return nil, fmt.Errorf("could not read records. %v", err)
	}
	return exes, nil
}

//GetByName returns an Exercise by a given name, nil if not found. Error in db connection issues
func (lib *Storage) GetByName(name string) (*storage.Exercise, error) {
	filter := bson.D{{Key: "name", Value: name}}
	var exe *storage.Exercise
	err := lib.FindOne(context.Background(), filter).Decode(&exe)
	if err != nil {
		return nil, fmt.Errorf("could not find record by name %s. Error was %v", name, err)
	}
	return exe, nil
}

//GetByMuscleGroup obtains all the excersices by a given muscle group
func (lib *Storage) GetByMuscleGroup(mg string) ([]storage.Exercise, error) {
	filter := bson.D{{Key: "muscle_groups", Value: mg}}
	cursor, err := lib.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("could not find records. %v", err)
	}
	var exes []storage.Exercise
	if err = cursor.All(context.Background(), &exes); err != nil {
		return nil, fmt.Errorf("could not read records. %v", err)
	}
	return exes, nil
}

//AddExercise add a new exercise to the database
func (lib *Storage) AddExercise(exe *storage.Exercise) (*storage.Exercise, error) {
	exists, err := lib.exists(exe)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("exercise with name '%s', already exists", exe.Name)
	}
	_, err = lib.InsertOne(context.Background(), exe)
	if err != nil {
		return nil, fmt.Errorf("failed to create new exercise %v. Error was %v", exe, err)
	}
	return exe, nil
}

//Update an exercise. error in case of db issue
// func (lib *Lib) Update(exe *Exercise) error {
// 	exists, err := lib.exists(exe)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("exercise with name '%s' does not exists", exe.GetName())
// 	}
// 	filter := bson.D{{"name", exe.GetName()}}
// 	_, err = lib.UpdateOne(context.TODO(), filter, exe)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return nil
// }

func (lib *Storage) exists(exe *storage.Exercise) (bool, error) {
	e, err := lib.GetByName(exe.Name)
	if err != nil {
		return false, fmt.Errorf("could not find record. %v", err)
	}
	return e != nil, nil
}

// Close the connection to the mongo server
func (lib *Storage) Close() error {
	ctx := context.Background()
	err := lib.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("failed to close connection to database. Error %v", err)
	}
	return nil
}
