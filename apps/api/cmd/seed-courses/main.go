package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"learn-platform/backend/internal/seed"
)

func main() {
	apply := flag.Bool("apply", false, "insert missing courses; default is a read-only dry run")
	flag.Parse()
	dsn := os.Getenv("DATABASE_URL")
	ownerID := os.Getenv("COURSE_SEED_INSTRUCTOR_ID")
	if dsn == "" || ownerID == "" {
		log.Fatal("DATABASE_URL and COURSE_SEED_INSTRUCTOR_ID are required")
	}
	courses, err := seed.Catalog()
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	summary, err := seed.SeedCourses(context.Background(), db, ownerID, courses, *apply)
	if err != nil {
		log.Fatal(err)
	}
	mode := "dry-run"
	if *apply {
		mode = "applied"
	}
	fmt.Printf("%s: created=%d skipped=%d total=%d\n", mode, summary.Created, summary.Skipped, len(courses))
}
