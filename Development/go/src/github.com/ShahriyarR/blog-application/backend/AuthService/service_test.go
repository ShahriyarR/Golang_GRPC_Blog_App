package main

import (
	"blog-application/global"
	"blog-application/proto"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func Test_authServer_Login(t *testing.T) {
	global.ConnectToTestDB()
	pw, _ := bcrypt.GenerateFromPassword([]byte("example"), bcrypt.DefaultCost)
	global.DB.Collection("user").InsertOne(context.Background(), global.User{ID: primitive.NewObjectID(),
																					Email: "test@gmail.com",
																					UserName: "Carl",
																					Password: string(pw)})
	server := authServer{}
	_, err := server.Login(context.Background(), &proto.LoginRequest{Login: "test@gmail.com", Password: "example"})
	if err != nil {
		t.Error("1: An error was returned: ", err.Error())
	}

	_, err = server.Login(context.Background(), &proto.LoginRequest{Login: "something", Password: "something"})
	if err == nil {
		t.Error("2: Error was nil")
	}

	_, err = server.Login(context.Background(), &proto.LoginRequest{Login: "Carl", Password: "example"})
	if err != nil {
		t.Error("3: An error was returned: ", err.Error())
	}
}

func Test_authServer_UserNameUsed(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{UserName: "Carl"})
	server := authServer{}
	res, err := server.UserNameUsed(context.Background(), &proto.UserNameUsedRequest{UserName: "Carlo"})
	if err != nil {
		t.Error("1: An error was returned: ", err.Error())

	}
	if res.GetUsed() {
		t.Error("1: wrong result")
	}

	res, err = server.UserNameUsed(context.Background(), &proto.UserNameUsedRequest{UserName: "Carl"})
	if err != nil {
		t.Error("2: An error was returned: ", err.Error())
	}

	if !res.GetUsed() {
		t.Error("2: wrong result")
	}
}

func Test_authServer_EmailUsed(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Email: "carl@gmail.com"})
	server := authServer{}
	res, err := server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: "carlo@gmail.com"})
	if err != nil {
		t.Error("1: An error was returned: ", err.Error())

	}
	if res.GetUsed() {
		t.Error("1: wrong result")
	}

	res, err = server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: "carl@gmail.com"})
	if err != nil {
		t.Error("2: An error was returned: ", err.Error())
	}

	if !res.GetUsed() {
		t.Error("2: wrong result")
	}
}

func Test_authServer_SignUp(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{UserName: "carl", Email: "carl@gmail.com"})
	server := authServer{}

	_, err := server.SignUp(context.Background(), &proto.SignUpRequest{UserName: "carl", Email: "example@gmail.com", Password: "examplestring"})
	if err.Error() != "username is used" {
		t.Error("1: No or the wrong Error was returned")
	}

	_, err = server.SignUp(context.Background(), &proto.SignUpRequest{UserName: "example", Email: "carl@gmail.com", Password: "examplestring"})
	if err.Error() != "email is used" {
		t.Error("2: No or the wrong Error was returned")
	}

	_, err = server.SignUp(context.Background(), &proto.SignUpRequest{UserName: "example", Email: "example@gmail.com", Password: "examplestring"})
	if err != nil {
		t.Error("3: an error was returned")
	}

	_, err = server.SignUp(context.Background(), &proto.SignUpRequest{UserName: "example", Email: "example@gmail.com", Password: "exam"})
	if err.Error() != "validation failed" {
		t.Error("4: No or the wrong Error was returned")
	}

}

func Test_authServer_AuthUser(t *testing.T) {
	server := authServer{}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoie1wiSURcIjpcIjVmMTczYTFiNWQ0MzFlY2U0OGE0ZGYxYVwiLFwiVXNlck5hbWVcIjpcIkNhcmxcIixcIkVtYWlsXCI6XCJ0ZXN0QGdtYWlsLmNvbVwiLFwiUGFzc3dvcmRcIjpcIiQyYSQxMCRMcjRqNFlKb3RkRmxBZDBBL2tJTlR1M2paQm42NUpLRjEzbC9PZEVjYmc4YWpJM2V5Yk03aVwifSJ9.FUKL0keIlIfx8-6QFg38P7aUzkL3fXCKl4ReJae9fdE"
	res, err := server.AuthUser(context.Background(), &proto.AuthUserRequest{Token: token})

	if err != nil {
		t.Error("an error was returned")
	}

	if res.GetID() != "5f173a1b5d431ece48a4df1a" || res.GetUserName() != "Carl" || res.GetEmail() != "test@gmail.com" {
		t.Error("wrong result returned: ", res)
	}


}