package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"learn-platform/backend/internal/model"
	"learn-platform/backend/internal/repository"
)

type CourseService struct {
	repo *repository.CourseRepository
}

func NewCourseService(repo *repository.CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func slugify(s string) string {
	var re = regexp.MustCompile("[^a-z0-9]+")
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func (s *CourseService) CreateCourse(ctx context.Context, userID string, req *model.CreateCourseRequest) (*model.CourseResponse, error) {
	if req.Title == "" {
		return nil, errors.New("course title is required")
	}

	courseID := uuid.New().String()
	course := &model.Course{
		ID:        courseID,
		UserID:    userID,
		Title:     req.Title,
		Link:      req.Link,
		Icon:      req.Icon,
		Status:    "learning",
		Slug:      slugify(req.Title) + "-" + courseID[:8],
		Language:  "id",
	}

	if course.Icon == "" {
		course.Icon = "book"
	}

	if err := s.repo.CreateCourse(ctx, course); err != nil {
		return nil, err
	}

	return s.mapCourseToResponse(ctx, course, userID), nil
}

func (s *CourseService) GetUserCourses(ctx context.Context, userID string) ([]model.CourseResponse, error) {
	courses, err := s.repo.GetCoursesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []model.CourseResponse
	for _, c := range courses {
		responses = append(responses, *s.mapCourseToResponse(ctx, &c, userID))
	}
	
	if responses == nil {
		responses = make([]model.CourseResponse, 0)
	}
	return responses, nil
}

func (s *CourseService) GetCourse(ctx context.Context, courseID string, userID string) (*model.CourseResponse, error) {
	// Try owned course first; fall back to public (enrolled students)
	course, err := s.repo.GetCourseByIDPublic(ctx, courseID)
	if err != nil {
		return nil, errors.New("course not found")
	}
	return s.mapCourseToResponse(ctx, course, userID), nil
}

// GetAllCourses returns all courses for the catalog.
func (s *CourseService) GetAllCourses(ctx context.Context) ([]model.CourseResponse, error) {
	courses, err := s.repo.GetAllPublishedCourses(ctx)
	if err != nil {
		return nil, err
	}
	var responses []model.CourseResponse
	for _, c := range courses {
		resp := model.CourseResponse{
			ID:           c.ID,
			Title:        c.Title,
			Status:       c.Status,
			TotalModules: 0,
		}
		for _, sec := range c.Sections {
			resp.TotalModules += len(sec.Lessons)
		}
		if c.Instructor.Name != "" {
			resp.Link = c.Instructor.Name // reuse Link field to carry instructor name
		}
		responses = append(responses, resp)
	}
	if responses == nil {
		responses = make([]model.CourseResponse, 0)
	}
	return responses, nil
}

// EnrollCourse enrolls a student in a course.
func (s *CourseService) EnrollCourse(ctx context.Context, userID, courseID string) error {
	enrolled, err := s.repo.IsEnrolled(ctx, userID, courseID)
	if err != nil {
		return err
	}
	if enrolled {
		return errors.New("already enrolled")
	}
	enrollment := &model.Enrollment{
		UserID:   userID,
		CourseID: courseID,
	}
	return s.repo.EnrollUser(ctx, enrollment)
}

// GetEnrolledCourses returns courses a student has enrolled in.
func (s *CourseService) GetEnrolledCourses(ctx context.Context, userID string) ([]model.CourseResponse, error) {
	courses, err := s.repo.GetEnrolledCourses(ctx, userID)
	if err != nil {
		return nil, err
	}
	var responses []model.CourseResponse
	for _, c := range courses {
		responses = append(responses, *s.mapCourseToResponse(ctx, &c, userID))
	}
	if responses == nil {
		responses = make([]model.CourseResponse, 0)
	}
	return responses, nil
}


func (s *CourseService) CreateModule(ctx context.Context, userID string, courseID string, req *model.CreateModuleRequest) (*model.ModuleResponse, error) {
	if req.Title == "" {
		return nil, errors.New("module title is required")
	}

	// Verify course belongs to user
	_, err := s.repo.GetCourseByID(ctx, courseID, userID)
	if err != nil {
		return nil, errors.New("course not found or unauthorized")
	}

	// Get sections for course
	sections, err := s.repo.GetSectionsByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	var section model.Section
	if len(sections) == 0 {
		section = model.Section{
			CourseID: courseID,
			Title:    "Materi Utama",
			Position: 1,
		}
		if err := s.repo.CreateSection(ctx, &section); err != nil {
			return nil, err
		}
	} else {
		section = sections[0]
	}

	// Count lessons in section to determine next position
	count, err := s.repo.CountLessonsBySectionID(ctx, section.ID)
	if err != nil {
		return nil, err
	}
	nextPos := int(count) + 1

	durationSec := 1800 // Default 30 min = 1800 seconds
	lesson := &model.Lesson{
		SectionID:   section.ID,
		Title:       req.Title,
		ContentType: "video", // Default to video
		Position:    nextPos,
		Duration:    durationSec,
	}

	if err := s.repo.CreateLesson(ctx, lesson); err != nil {
		return nil, err
	}

	return &model.ModuleResponse{
		ID:          lesson.ID,
		Title:       lesson.Title,
		IsCompleted: false,
		Duration:    req.Duration,
	}, nil
}

func (s *CourseService) ToggleModuleCompletion(ctx context.Context, userID string, lessonID string) (*model.ModuleResponse, error) {
	lesson, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, errors.New("lesson not found")
	}

	// Fetch section to verify course owner
	section, err := s.repo.GetSectionByID(ctx, lesson.SectionID)
	if err != nil {
		return nil, errors.New("section not found")
	}

	_, err = s.repo.GetCourseByID(ctx, section.CourseID, userID)
	if err != nil {
		return nil, errors.New("unauthorized to modify this lesson")
	}

	progress, err := s.repo.GetLessonProgress(ctx, userID, lessonID)
	if err != nil {
		return nil, err
	}

	isCompleted := true
	if progress != nil {
		progress.IsCompleted = !progress.IsCompleted
		isCompleted = progress.IsCompleted
		if err := s.repo.SaveLessonProgress(ctx, progress); err != nil {
			return nil, err
		}
	} else {
		progress = &model.LessonProgress{
			UserID:      userID,
			LessonID:    lessonID,
			IsCompleted: true,
		}
		if err := s.repo.SaveLessonProgress(ctx, progress); err != nil {
			return nil, err
		}
	}

	durStr := "30 min"
	if lesson.Duration > 0 {
		durStr = fmt.Sprintf("%d min", lesson.Duration/60)
	}

	return &model.ModuleResponse{
		ID:          lesson.ID,
		Title:       lesson.Title,
		IsCompleted: isCompleted,
		Duration:    durStr,
	}, nil
}

func (s *CourseService) UpdateCourseStatus(ctx context.Context, userID string, courseID string, req *model.UpdateCourseStatusRequest) (*model.CourseResponse, error) {
	// Accept old states too
	if req.Status != "learning" && req.Status != "completed" && req.Status != "cancelled" &&
		req.Status != "draft" && req.Status != "pending" && req.Status != "published" && req.Status != "archived" {
		return nil, errors.New("invalid status")
	}

	course, err := s.repo.GetCourseByID(ctx, courseID, userID)
	if err != nil {
		return nil, errors.New("course not found or unauthorized")
	}

	course.Status = req.Status
	if err := s.repo.UpdateCourse(ctx, course); err != nil {
		return nil, err
	}

	return s.mapCourseToResponse(ctx, course, userID), nil
}

func (s *CourseService) mapCourseToResponse(ctx context.Context, c *model.Course, userID string) *model.CourseResponse {
	progresses, _ := s.repo.GetUserLessonProgresses(ctx, userID)
	completedMap := make(map[string]bool)
	for _, p := range progresses {
		if p.IsCompleted {
			completedMap[p.LessonID] = true
		}
	}

	var completed int
	var moduleResponses []model.ModuleResponse

	for _, sec := range c.Sections {
		for _, l := range sec.Lessons {
			isComp := completedMap[l.ID]
			if isComp {
				completed++
			}
			durStr := "30 min"
			if l.Duration > 0 {
				durStr = fmt.Sprintf("%d min", l.Duration/60)
			}
			moduleResponses = append(moduleResponses, model.ModuleResponse{
				ID:          l.ID,
				Title:       l.Title,
				IsCompleted: isComp,
				Duration:    durStr,
			})
		}
	}

	if moduleResponses == nil {
		moduleResponses = make([]model.ModuleResponse, 0)
	}

	progress := 0
	total := len(moduleResponses)
	if total > 0 {
		progress = (completed * 100) / total
	}

	return &model.CourseResponse{
		ID:               c.ID,
		Title:            c.Title,
		Link:             c.Link,
		Icon:             c.Icon,
		Status:           c.Status,
		Progress:         progress,
		TotalModules:     total,
		CompletedModules: completed,
		Modules:          moduleResponses,
	}
}
