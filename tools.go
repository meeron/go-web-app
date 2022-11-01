package main

import (
	"fmt"
	"os"
	"web-app/database"
)

func parseArguments() bool {
	argLength := len(os.Args[1:])
	if argLength >= 1 {
		if os.Args[1] == "migrate" {
			migrateDatabase()
			return true
		}

		panic(fmt.Sprintf("Unsupported argument: %v", os.Args[1]))
	}

	return false
}

func migrateDatabase() {
	fmt.Println("Migrating database...")

	err := database.MigrateDb()
	if err != nil {
		panic(err)
	}

	fmt.Println("Database migration done.")
}
