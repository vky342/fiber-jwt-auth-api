package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"context"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb" 
	"github.com/gofiber/storage/dynamodb/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
        
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


func ParseAndVerifyToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
}



func ExtractUserFromToken(token *jwt.Token) (*User, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, usernameOk := claims["username"].(string)
		password, passwordOk := claims["password"].(string)

		if !usernameOk {
			return nil, fmt.Errorf("username claim missing or not a string")
		}
		if !passwordOk {
			password = "" // Optional: Handle missing password gracefully
		}

		return &User{Username: username, Password: password}, nil
	}
	return nil, fmt.Errorf("invalid token")
}


func AuthMiddleware(c *fiber.Ctx) (*jwt.Token, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Extract token from "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization format"})
	}

	tokenString := tokenParts[1];

	parsed_token, err := ParseAndVerifyToken(tokenString)

	if err != nil{
		return  nil,c.Status(fiber.StatusNonAuthoritativeInformation).JSON(fiber.Map{
			"error": "Parsing and verifying failed",
		})
	}

	return parsed_token, nil
}

func listPrimaryKeys(store *dynamodb.Storage, tableName string, primaryKeyName string) ([]string, error) {

	ctx := context.Background() 

	input := &awsdynamodb.ExecuteStatementInput{
			Statement: aws.String(fmt.Sprintf("SELECT %s FROM \"%s\"", primaryKeyName, tableName)),
	}

	output, err := store.Conn().ExecuteStatement(ctx, input)

	if err != nil {
			return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	Keys := make([]string, 0)
	
	for _, key := range output.Items {
		
		primaryKeyValue, ok := key[primaryKeyName].(*types.AttributeValueMemberS) // Assuming ID is a string
		if ok {
				//fmt.Println("Primary Key Value:", primaryKeyValue.Value)
				Keys = append(Keys, primaryKeyValue.Value)
		}
	}

	return Keys, nil
}