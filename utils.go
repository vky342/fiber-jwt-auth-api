package main

import (
	"os"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(user *User) (string,error){
	claims := jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	}

	//secret := "secret"

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(os.Getenv("SECRET"))) 

	if err != nil{
		return "",err
	}

	return token,nil
}