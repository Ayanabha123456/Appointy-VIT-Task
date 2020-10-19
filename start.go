package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"encoding/json"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)
type Participant struct{
	Name string `json:"name"`
	Email string `json:"email"`
	RSVP string `json:"rsvp"`
}
type Meeting struct{
	Ide string `json:"_id"`
	Title string `json:"title"`
	Participants []Participant `json:"participants"`
	Start string `json:"start"`
	End string `json:"end"`
	Timestamp string `json:"timestamp"`
}
func scheduleMeeting(rw http.ResponseWriter,req *http.Request){
	var m Meeting
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		http.Error(rw,err.Error(),http.StatusBadRequest)
		return
	}
	collection:=client.Database("Appointy").Collection("meetings")
	result,err:=collection.InsertOne(context.TODO(),m);
	if(err!=nil){
		log.Fatal(err)
	}
	fmt.Println(result)
	json.NewEncoder(rw).Encode(m)
}
func getMeeting(rw http.ResponseWriter, req *http.Request){
	query:=req.URL.Query()
	id:=query.Get("id")
	var mres Meeting
	findOptions:=options.Find()
	findOptions.SetLimit(1)
	collection:=client.Database("Appointy").Collection("meetings")
	cur,err:=collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		log.Fatal(err)
	}
	for cur.Next(context.TODO()){
		var m Meeting
		err:=cur.Decode(&m)
		if err!=nil{
			log.Fatal(err)
		}
		if id==m.Ide{
			mres=m
			break;
		}

	}
	cur.Close(context.TODO())
	json.NewEncoder(rw).Encode(mres)
}
func listMeetings(rw http.ResponseWriter,req *http.Request){
	query:=req.URL.Query()
	sta:=query.Get("start")
	end:=query.Get("end")
	var meetings []Meeting
	findOptions:=options.Find()
	collection:=client.Database("Appointy").Collection("meetings")
	cur,err:=collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		log.Fatal(err)
	}
	for cur.Next(context.TODO()){
		var m Meeting
		err:=cur.Decode(&m)
		if err!=nil{
			log.Fatal(err)
		}
		if sta==m.Start && end==m.End{
			meetings=append(meetings,m)
		}

	}
	cur.Close(context.TODO())
	json.NewEncoder(rw).Encode(meetings)
}
func listParticipantMeetings(rw http.ResponseWriter,req *http.Request){
	query:=req.URL.Query()
	par:=query.Get("participant")
	var meetings []Meeting
	findOptions:=options.Find()
	collection:=client.Database("Appointy").Collection("meetings")
	cur,err:=collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		log.Fatal(err)
	}
	for cur.Next(context.TODO()){
		var m Meeting
		err:=cur.Decode(&m)
		if err!=nil{
			log.Fatal(err)
		}
		p:=m.Participants
		for i:=0;i<len(p);i++{
			if  par==p[i].Email{
				meetings=append(meetings,m)
			}
		}
		cur.Close(context.TODO())
		json.NewEncoder(rw).Encode(meetings)
	}

}
func chooseHandler(rw http.ResponseWriter,req *http.Request){
	if req.Method == http.MethodGet{
		query:=req.URL.Query()
		px:=query.Get("participant")
		if px!=""{
			listParticipantMeetings(rw,req)
		}else{
			listMeetings(rw,req)
		}
		
	}else{
		scheduleMeeting(rw,req)
	}
}
var client *mongo.Client
func main(){
	fmt.Println("Starting app...")
	client, err:=mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx,_:=context.WithTimeout(context.Background(),10*time.Second)
	err=client.Connect(ctx)
	if err != nil {
		fmt.Println("Connection error!!")
		return
	}
	mux:=http.NewServeMux()
	mux.HandleFunc("/meeting",getMeeting)
	mux.HandleFunc("/meetings",http.HandlerFunc(chooseHandler))
	http.ListenAndServe(":8081",mux)
}