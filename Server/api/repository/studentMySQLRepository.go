package repository

import (
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

func (sr *studentMySQLRepository) Create(student *model.Student) error {
	query, err := sr.MySQLDB.Prepare("INSERT INTO Students(name, email, mobileNo, password) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	result, err := query.Exec(student.Name, student.Email, student.Mobile, student.Password)
	if err != nil {
		return err
	}

	lId, _ := result.LastInsertId()

	query, err = sr.MySQLDB.Prepare("INSERT INTO UserRole(userId, roleId) VALUES(?,?)")
	if err != nil {
		return err
	}

	query.Exec(lId, 2)

	return nil
}

func (sr *studentMySQLRepository) ReadAll() ([]model.Student, error) {
	var students []model.Student
	rows, err := sr.MySQLDB.Query(`SELECT name,email,mobileNo from Students`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var student model.Student

		if err := rows.Scan(&student.Name, &student.Email, &student.Mobile); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}
