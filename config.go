package main

import "github.com/gofiber/storage/dynamodb/v2"

func InitializeDB(config dynamodb.Config) dynamodb.Storage{
	// Initialize dynamodb
	store := dynamodb.New(config)
	return *store
}