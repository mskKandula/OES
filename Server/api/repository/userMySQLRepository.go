package repository

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/mskKandula/oes/api/model"
)

type userMySQLRepository struct {
	MySQLDB *sql.DB
}

func NewUserMySQLRepository(rc *RepositoryConfig) model.UserRepository {
	return &userMySQLRepository{
		MySQLDB: rc.MySQLDB,
	}
}

func (ur *userMySQLRepository) Create(user model.User, password string) error {
	query, err := ur.MySQLDB.Prepare("INSERT INTO Examiners(name, age, email, mobileNo, password,clientId) VALUES(?,?,?,?,?)")

	if err != nil {
		return err
	}

	id := uuid.New()

	result, err := query.Exec(user.Name, user.Age, user.Email, user.MobileNo, password, strings.Replace(id.String(), "-", "", -1))

	if err != nil {
		return err
	}

	lId, _ := result.LastInsertId()

	query, err = ur.MySQLDB.Prepare("INSERT INTO UserRole(userId, roleId) VALUES(?,?)")

	if err != nil {
		return err
	}

	query.Exec(lId, 1)

	return nil
}

func (ur *userMySQLRepository) CreateVideo(fileName, videoUrl, imagePath, clientId string) error {
	query, err := ur.MySQLDB.Prepare("INSERT INTO VideoContent(name, videoUrl,thumbnailPath,contentType,description,clientId) VALUES(?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = query.Exec(fileName, videoUrl, imagePath, "video/mp4", "Sample Video", clientId)
	if err != nil {
		return err
	}
	return nil
}
