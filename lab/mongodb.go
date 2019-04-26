package main
import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
)

func main(){
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://monitor:76shG8jsd87dsd@52.220.171.16:27018?authSource=cybex_gateway"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	collection := client.Database("cybex_gateway").Collection("numbers")
}