package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/huandu/go-sqlbuilder"

	mooc "github.com/AguilaMike/go-gin-hexagonal/internal"
)

// CourseRepository is a MySQL mooc.CourseRepository implementation.
type CourseRepository struct {
	db *sql.DB
}

// NewCourseRepository initializes a MySQL-based implementation of mooc.CourseRepository.
func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

// GetAll implements the mooc.CourseRepository interface.
func (r *CourseRepository) GetAll(ctx context.Context) ([]mooc.Course, error) {
	courseSQLStruct := sqlbuilder.NewStruct(new(sqlCourse))
	query, args := courseSQLStruct.SelectFrom(sqlCourseTable).Build()

	rows, _ := r.db.QueryContext(ctx, query, args...)
	defer rows.Close()

	var courses []mooc.Course
	for rows.Next() {
		var courseDB sqlCourse
		if err := rows.Scan(courseSQLStruct.Addr(&courseDB)...); err != nil {
			return nil, fmt.Errorf("error trying to scan course: %w", err)
		}
		course, err := mooc.NewCourse(courseDB.ID, courseDB.Name, courseDB.Duration)
		if err != nil {
			return nil, fmt.Errorf("error trying to create course: %w", err)
		}
		courses = append(courses, course)
	}
	return courses, nil
}

// GetByID implements the mooc.CourseRepository interface.
func (r *CourseRepository) GetByID(ctx context.Context, id string) (*mooc.Course, error) {
	courseSQLStruct := sqlbuilder.NewStruct(new(sqlCourse))
	sb := courseSQLStruct.SelectFrom(sqlCourseTable)
	sb.Where(sb.Equal("id", id))
	query, args := sb.Build()

	rows, _ := r.db.QueryContext(ctx, query, args...)
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("course with id %s not found", id)
	}

	var courseDB sqlCourse
	if err := rows.Scan(courseSQLStruct.Addr(&courseDB)...); err != nil {
		return nil, fmt.Errorf("error trying to scan course: %w", err)
	}
	course, err := mooc.NewCourse(courseDB.ID, courseDB.Name, courseDB.Duration)
	if err != nil {
		return nil, fmt.Errorf("error trying to create course: %w", err)
	}
	return &course, nil
}

// Save implements the mooc.CourseRepository interface.
func (r *CourseRepository) Save(ctx context.Context, course mooc.Course) error {
	courseSQLStruct := sqlbuilder.NewStruct(new(sqlCourse))
	query, args := courseSQLStruct.InsertInto(sqlCourseTable, sqlCourse{
		ID:       course.ID().String(),
		Name:     course.Name().String(),
		Duration: course.Duration().String(),
	}).Build()

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist course on database: %v", err)
	}

	return nil
}
