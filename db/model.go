package db

import (
	"context"
	"time"

	"github.com/pawanverma1337/atlan-challenge/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// M global variable to store the struct
var M mongoStore

type mongoStore struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

// File struct to store the data.
type File struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	FileUploadLength    int                `bson:"file_length"`
	FileUploadOffset    int                `bson:"file_offset"`
	FileUploadComplete  bool               `bson:"file_complete"`
	FileUploadTerminate bool               `bson:"file_terminate"`
}

// CreateFile function
func (f *File) CreateFile() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = M.Collection.InsertOne(ctx, *f)
	return err
}

// UpdateFile func
func (f *File) UpdateFile(update bson.D) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = M.Collection.UpdateOne(ctx, bson.M{"_id": (*f).ID}, update)
	return err
}

// GetFile func
func (f *File) GetFile() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = M.Collection.FindOne(ctx, bson.D{{"_id", (*f).ID}}).Decode(f)
	return err
}

// Connect function to connect to the mongo
// db database and print the log information
// to the terminal.
func Connect() {
	var err error

	M.Client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/test"))
	util.CheckError(err)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = M.Client.Connect(ctx)
	util.CheckError(err)
}
