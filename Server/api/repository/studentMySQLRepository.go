package repository

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mskKandula/oes/api/model"
)

const (
	createStudentQuery  = "INSERT INTO Students(name, email, mobileNo, password, clientId) VALUES(?,?,?,?,?)"
	createUserRoleQuery = "INSERT INTO UserRole(userId, roleId) VALUES(?,?)"
	readAllStudentsQuery = "SELECT id, name, email, mobileNo FROM Students WHERE clientId=?"
)

type studentMySQLRepository struct {
	MySQLDB *sql.DB
}

func NewStudentMySQLRepository(rc *RepositoryConfig) model.StudentRepository {
	return &studentMySQLRepository{
		MySQLDB: rc.MySQLDB,
	}
}

func (sr *studentMySQLRepository) Create(ctx context.Context, student *model.Student) error {

	return withTransactionContext(ctx, sr.MySQLDB, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx,
			createStudentQuery,
			student.Name, student.Email, student.Mobile, student.Password, student.ClientId,
		)
		if err != nil {
			return err
		}

		lId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx,
			createUserRoleQuery,
			lId, 2,
		)
		return err
	})
}

func (sr *studentMySQLRepository) ReadAll(ctx context.Context, clientId string) ([]model.Student, error) {
	var students []model.Student

	rows, err := sr.MySQLDB.QueryContext(ctx, readAllStudentsQuery, clientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var student model.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Mobile); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
