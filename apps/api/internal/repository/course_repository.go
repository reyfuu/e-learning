package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"learn-platform/backend/internal/model"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) CreateCourse(ctx context.Context, course *model.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *CourseRepository) GetCoursesByUserID(ctx context.Context, userID string) ([]model.Course, error) {
	var courses []model.Course
	// Preload Sections and their Lessons ordered by position
	err := r.db.WithContext(ctx).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Sections.Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("instructor_id = ?", userID).
		Find(&courses).Error
	return courses, err
}

func (r *CourseRepository) GetCourseByID(ctx context.Context, id string, userID string) (*model.Course, error) {
	var course model.Course
	err := r.db.WithContext(ctx).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Sections.Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("id = ? AND instructor_id = ?", id, userID).
		First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *CourseRepository) UpdateCourse(ctx context.Context, course *model.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

// Section & Lesson Methods for Course Curriculum

func (r *CourseRepository) CreateSection(ctx context.Context, section *model.Section) error {
	return r.db.WithContext(ctx).Create(section).Error
}

func (r *CourseRepository) GetSectionsByCourseID(ctx context.Context, courseID string) ([]model.Section, error) {
	var sections []model.Section
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("position ASC").
		Find(&sections).Error
	return sections, err
}

func (r *CourseRepository) CreateLesson(ctx context.Context, lesson *model.Lesson) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}

func (r *CourseRepository) GetLessonByID(ctx context.Context, id string) (*model.Lesson, error) {
	var lesson model.Lesson
	err := r.db.WithContext(ctx).First(&lesson, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (r *CourseRepository) UpdateLesson(ctx context.Context, lesson *model.Lesson) error {
	return r.db.WithContext(ctx).Save(lesson).Error
}

// Lesson Progress Methods

func (r *CourseRepository) GetLessonProgress(ctx context.Context, userID, lessonID string) (*model.LessonProgress, error) {
	var progress model.LessonProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND lesson_id = ?", userID, lessonID).
		First(&progress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found so we can easily handle it
		}
		return nil, err
	}
	return &progress, nil
}

func (r *CourseRepository) SaveLessonProgress(ctx context.Context, progress *model.LessonProgress) error {
	return r.db.WithContext(ctx).Save(progress).Error
}

func (r *CourseRepository) GetUserLessonProgresses(ctx context.Context, userID string) ([]model.LessonProgress, error) {
	var progresses []model.LessonProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&progresses).Error
	return progresses, err
}

func (r *CourseRepository) CountLessonsBySectionID(ctx context.Context, sectionID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Lesson{}).
		Where("section_id = ?", sectionID).
		Count(&count).Error
	return count, err
}

func (r *CourseRepository) GetSectionByID(ctx context.Context, id string) (*model.Section, error) {
	var section model.Section
	err := r.db.WithContext(ctx).First(&section, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

// GetAllPublishedCourses returns all courses (for the student catalog).
func (r *CourseRepository) GetAllPublishedCourses(ctx context.Context) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.WithContext(ctx).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Sections.Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Instructor").
		Find(&courses).Error
	return courses, err
}

// GetCourseByIDPublic fetches a single course without the instructor_id filter.
func (r *CourseRepository) GetCourseByIDPublic(ctx context.Context, id string) (*model.Course, error) {
	var course model.Course
	err := r.db.WithContext(ctx).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Sections.Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Instructor").
		First(&course, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// EnrollUser creates an enrollment record.
func (r *CourseRepository) EnrollUser(ctx context.Context, enrollment *model.Enrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}

// IsEnrolled checks whether a user is already enrolled in a course.
func (r *CourseRepository) IsEnrolled(ctx context.Context, userID, courseID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Enrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).Error
	return count > 0, err
}

// GetEnrolledCourses returns all courses a student has enrolled in.
func (r *CourseRepository) GetEnrolledCourses(ctx context.Context, userID string) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.WithContext(ctx).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Sections.Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Instructor").
		Joins("JOIN enrollments ON enrollments.course_id = courses.id").
		Where("enrollments.user_id = ? AND enrollments.deleted_at IS NULL", userID).
		Find(&courses).Error
	return courses, err
}

