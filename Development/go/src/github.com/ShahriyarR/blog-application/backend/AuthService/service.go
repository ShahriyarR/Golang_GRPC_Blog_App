package main

import (
	"blog-application/global"
	"blog-application/proto"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"log"
	"net"
	"regexp"
	"time"
)

type authServer struct{}

func (authServer) Login(_ context.Context, in *proto.LoginRequest) (*proto.AuthResponse, error){
	login, password := in.GetLogin(), in.GetPassword()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var user global.User
	global.DB.Collection("user").FindOne(ctx, bson.M{"$or": []bson.M{bson.M{"username": login}, bson.M{"email": login}}}).Decode(&user)
	if user == global.NilUser {
		return &proto.AuthResponse{}, errors.New("wrong login credentials provided")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return &proto.AuthResponse{}, errors.New("wrong login credentials provided")
	}

	return &proto.AuthResponse{Token: user.GetToken()}, nil

}

var server authServer

func (authServer) SignUp(_ context.Context, in *proto.SignUpRequest) (*proto.AuthResponse, error) {
	userName, email, password := in.GetUserName(), in.GetEmail(), in.GetPassword()
	match, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email)
	if len(userName) < 4 || len(userName) > 20 ||
		len(email) > 35  || len(password) < 8 || len(password) > 128 || !match{
		return &proto.AuthResponse{}, errors.New("validation failed")
	}

	res, err := server.UserNameUsed(context.Background(), &proto.UserNameUsedRequest{UserName: userName})
	if err != nil {
		log.Println("Error returned from UserNameUsed: ", err.Error())
		return &proto.AuthResponse{}, errors.New("something went wrong")
	}

	if res.GetUsed() {
		return &proto.AuthResponse{}, errors.New("username is used")
	}

	res, err = server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: email})
	if err != nil {
		log.Println("Error returned from UserNameUsed: ", err.Error())
		return &proto.AuthResponse{}, errors.New("something went wrong")
	}

	if res.GetUsed() {
		return &proto.AuthResponse{}, errors.New("email is used")
	}
	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	newUser := global.User{ID: primitive.NewObjectID(), UserName: userName, Email: email, Password: string(pw)}

	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	_, err = global.DB.Collection("user").InsertOne(ctx, newUser)
	if err != nil {
		log.Println("Error inserting newUser: ", err.Error())
		return &proto.AuthResponse{}, errors.New("something went wrong")
	}
	return &proto.AuthResponse{Token: newUser.GetToken()}, nil
}

func (authServer) UserNameUsed(_ context.Context, in *proto.UserNameUsedRequest) (*proto.UsedResponse, error) {
	userName := in.GetUserName()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var result global.User
	global.DB.Collection("user").FindOne(ctx, bson.M{"username": userName}).Decode(&result)
	return &proto.UsedResponse{Used: result != global.NilUser}, nil
}

func (authServer) EmailUsed(_ context.Context, in *proto.EmailUsedRequest) (*proto.UsedResponse, error) {
	email := in.GetEmail()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var result global.User
	global.DB.Collection("user").FindOne(ctx, bson.M{"email": email}).Decode(&result)
	return &proto.UsedResponse{Used: result != global.NilUser}, nil
}
func main() {
	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, authServer{})
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("Error creating listener: ", err.Error())
	}
	server.Serve(listener)
}
