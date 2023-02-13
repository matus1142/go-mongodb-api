package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id       int32  `bson:"id"`
	Fname    string `bson:"fname"`
	Lname    string `bson:"lname"`
	Username string `bson:"username"`
	Email    string `bson:"email"`
	Avatar   string `bson:"avatar"`
}

func GetHelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}

func GetUsers(c *gin.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	checkErr(err)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.D{}
	err = client.Connect(ctx)
	checkErr(err)
	// เรียกดู database
	collection := client.Database("mydb").Collection("users")
	// เรียกดู collection
	cursor, err := collection.Find(context.TODO(), filter)
	checkErr(err)
	for cursor.Next(ctx) {
		var user bson.M
		if err = cursor.Decode(&user); err != nil {
			panic(err)
		}
		output, err := json.MarshalIndent(user, "", "    ")
		checkErr(err)
		fmt.Printf("%s\n", output)
	}
	defer client.Disconnect(ctx)
}

func GetUserId(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	checkErr(err)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{"id": id}
	err = client.Connect(ctx)
	checkErr(err)
	// เรียกดู database
	collection := client.Database("mydb").Collection("users")
	// เรียกดู collection
	var user bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	checkErr(err)
	output, err := json.MarshalIndent(user, "", "    ")
	checkErr(err)
	fmt.Printf("%s\n", output)
	defer client.Disconnect(ctx)
}

func AddUser(c *gin.Context) {
	var u User
	c.BindJSON(&u)
	c.JSON(http.StatusOK, gin.H{"result": u})

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	checkErr(err)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	checkErr(err)
	// เรียกดู database
	collection := client.Database("mydb").Collection("users")
	// newUser := User{Id: 99, Fname: "Ftest", Lname: "Ltest", Username: "Utest", Email: "Etest", Avatar: "Atest"}
	newUser := User{Id: u.Id, Fname: u.Fname, Lname: u.Lname, Username: u.Username, Email: u.Email, Avatar: u.Avatar}
	result, err := collection.InsertOne(context.TODO(), newUser)
	checkErr(err)
	fmt.Printf("Document inserted with ID: %s\n", result.InsertedID)
	defer client.Disconnect(ctx)
}

func UpdateUserData(c *gin.Context) {
	var u User
	c.BindJSON(&u)
	c.JSON(http.StatusOK, gin.H{"result": u})
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	checkErr(err)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	checkErr(err)
	// เรียกดู database
	collection := client.Database("mydb").Collection("users")
	filter := bson.M{"id": u.Id}
	update := bson.D{{"$set", bson.D{
		{"fname", u.Fname},
		{"lname", u.Lname},
		{"username", u.Username},
		{"email", u.Email},
		{"avatar", u.Avatar},
	}}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	checkErr(err)
	fmt.Printf("Documents updated: %v\n", result.ModifiedCount)
	defer client.Disconnect(ctx)
}

func DeleteUser(c *gin.Context) {
	var u User
	c.BindJSON(&u)
	c.JSON(http.StatusOK, gin.H{"result": u})
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	checkErr(err)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	checkErr(err)
	// เรียกดู database
	collection := client.Database("mydb").Collection("users")
	filter := bson.M{"id": u.Id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	checkErr(err)
	fmt.Printf("Documents deleted: %d\n", result.DeletedCount)
	defer client.Disconnect(ctx)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()
	r.GET("/", GetHelloWorld)
	r.GET("/users", GetUsers)
	r.GET("/users/:id", GetUserId)
	r.POST("/users/create", AddUser)
	r.PUT("/users/update", UpdateUserData)
	r.DELETE("/users/delete", DeleteUser)
	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
