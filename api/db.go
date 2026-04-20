package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

var client *supabase.Client

func initDB() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Erreur chargement .env:", err)
		return
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	client, err = supabase.NewClient(supabaseURL, supabaseKey, nil)

	if err != nil {
		fmt.Println("Erreur init Supabase:", err)
		return
	}
}
