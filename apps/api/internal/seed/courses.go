package seed

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"learn-platform/backend/internal/model"
)

//go:embed technology-courses.json
var files embed.FS

type LessonSpec struct {
	Title           string `json:"title"`
	DurationSeconds int    `json:"duration_seconds"`
	Content         string `json:"content"`
}

type SectionSpec struct {
	Title   string       `json:"title"`
	Lessons []LessonSpec `json:"lessons"`
}

type CourseSpec struct {
	Slug          string        `json:"slug"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	Outcomes      []string      `json:"outcomes"`
	Prerequisites []string      `json:"prerequisites"`
	Sections      []SectionSpec `json:"sections"`
}

type SeedSummary struct {
	Created int
	Skipped int
}

func Catalog() ([]CourseSpec, error) {
	data, err := files.ReadFile("technology-courses.json")
	if err != nil {
		return nil, err
	}
	var courses []CourseSpec
	if err := json.Unmarshal(data, &courses); err != nil {
		return nil, err
	}
	if err := validate(courses); err != nil {
		return nil, err
	}
	return courses, nil
}

func validate(courses []CourseSpec) error {
	if len(courses) != 56 {
		return fmt.Errorf("catalog must contain exactly 56 courses; got %d", len(courses))
	}
	slugs := make(map[string]bool, len(courses))
	for _, course := range courses {
		if strings.TrimSpace(course.Slug) == "" || strings.TrimSpace(course.Title) == "" || strings.TrimSpace(course.Description) == "" {
			return errors.New("course slug, title, and description are required")
		}
		if slugs[course.Slug] {
			return fmt.Errorf("duplicate course slug %q", course.Slug)
		}
		slugs[course.Slug] = true
		if len(course.Outcomes) == 0 || len(course.Prerequisites) == 0 || len(course.Sections) != 3 {
			return fmt.Errorf("course %q must include outcomes, prerequisites, and three sections", course.Slug)
		}
		for _, item := range append(append([]string{}, course.Outcomes...), course.Prerequisites...) {
			if strings.TrimSpace(item) == "" {
				return fmt.Errorf("course %q outcomes and prerequisites cannot be empty", course.Slug)
			}
		}
		for _, section := range course.Sections {
			if strings.TrimSpace(section.Title) == "" || len(section.Lessons) != 3 {
				return fmt.Errorf("course %q sections need a title and three lessons", course.Slug)
			}
			for _, lesson := range section.Lessons {
				if strings.TrimSpace(lesson.Title) == "" || strings.TrimSpace(lesson.Content) == "" || lesson.DurationSeconds <= 0 {
					return fmt.Errorf("course %q lessons need a title, content, and positive duration", course.Slug)
				}
			}
		}
	}
	return nil
}

func SeedCourses(ctx context.Context, db *gorm.DB, instructorID string, courses []CourseSpec, apply bool) (SeedSummary, error) {
	if err := validate(courses); err != nil {
		return SeedSummary{}, err
	}
	if strings.TrimSpace(instructorID) == "" {
		return SeedSummary{}, errors.New("COURSE_SEED_INSTRUCTOR_ID is required")
	}
	var instructor model.User
	if err := db.WithContext(ctx).First(&instructor, "id = ?", instructorID).Error; err != nil {
		return SeedSummary{}, fmt.Errorf("load seed instructor: %w", err)
	}
	if instructor.Role != "instructor" {
		return SeedSummary{}, errors.New("seed owner must have the instructor role")
	}

	var summary SeedSummary
	work := func(tx *gorm.DB) error {
		for _, spec := range courses {
			var count int64
			if err := tx.Model(&model.Course{}).Where("slug = ?", spec.Slug).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				summary.Skipped++
				continue
			}
			if !apply {
				summary.Created++
				continue
			}
			description := spec.Description + "\n\nYang akan dipelajari:\n- " + strings.Join(spec.Outcomes, "\n- ") + "\n\nPrasyarat:\n- " + strings.Join(spec.Prerequisites, "\n- ")
			course := model.Course{UserID: instructorID, Title: spec.Title, Slug: spec.Slug, Description: description, Price: 0, Level: "beginner", Status: "published", Language: "id", Icon: "book"}
			result := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoNothing: true}).Create(&course)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				summary.Skipped++
				continue
			}
			for sectionIndex, sectionSpec := range spec.Sections {
				section := model.Section{CourseID: course.ID, Title: sectionSpec.Title, Position: sectionIndex + 1}
				if err := tx.Create(&section).Error; err != nil {
					return err
				}
				for lessonIndex, lessonSpec := range sectionSpec.Lessons {
					content := lessonSpec.Content
					lesson := model.Lesson{SectionID: section.ID, Title: lessonSpec.Title, ContentType: "text", Content: &content, Duration: lessonSpec.DurationSeconds, Position: lessonIndex + 1}
					if err := tx.Create(&lesson).Error; err != nil {
						return err
					}
				}
			}
			summary.Created++
		}
		return nil
	}
	if !apply {
		err := work(db.WithContext(ctx))
		return summary, err
	}
	err := db.WithContext(ctx).Transaction(work)
	return summary, err
}
