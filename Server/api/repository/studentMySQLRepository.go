package repository

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mskKandula/oes/api/model"
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

	err := withTransaction(sr.MySQLDB, func(tx *sql.Tx) error {
		query, err := tx.Prepare("INSERT INTO Students(name, email, mobileNo, password,clientId) VALUES(?,?,?,?,?)")
		if err != nil {
			return err
		}

		result, err := query.ExecContext(ctx, student.Name, student.Email, student.Mobile, student.Password, student.ClientId)
		if err != nil {
			return err
		}

		lId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		query, err = tx.Prepare("INSERT INTO UserRole(userId, roleId) VALUES(?,?)")
		if err != nil {
			return err
		}

		_, err = query.ExecContext(ctx, lId, 2)
		if err != nil {
			return err
		}

		return err
	})

	return err
}

func (sr *studentMySQLRepository) ReadAll(ctx context.Context, clientId string) ([]model.Student, error) {
	var students []model.Student
	rows, err := sr.MySQLDB.QueryContext(ctx, `SELECT id,name,email,mobileNo from Students where clientId=?`, clientId)
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

	return students, nil
}
