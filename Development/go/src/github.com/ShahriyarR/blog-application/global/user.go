package global

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NilUser is a nil value for a user
var NilUser User

// User is the default user struct
type User struct{
	ID primitive.ObjectID `bson:"_id"`
	UserName string 	  `bson:"username"`
	Email string 	      `bson:"email"`
	Password string       `bson:"password"`
}

// GetToken returns the User's JWT
func (u User) GetToken() string {
	byteSlc, _ := json.Marshal(u)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": string(byteSlc),
	})
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}