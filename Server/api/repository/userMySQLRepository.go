package repository

import (
	"context"
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/mskKandula/oes/api/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type userMySQLRepository struct {
	MySQLDB  *sql.DB
	RabbitMQ *amqp.Channel
	Queue    amqp.Queue
}

func NewUserMySQLRepository(rc *RepositoryConfig) model.UserRepository {
	return &userMySQLRepository{
		MySQLDB:  rc.MySQLDB,
		RabbitMQ: rc.RabbitMQ,
		Queue:    rc.Queue,
	}
}

func (ur *userMySQLRepository) Create(ctx context.Context, user model.User, password string) error {

	err := withTransaction(ur.MySQLDB, func(tx *sql.Tx) error {
		query, err := tx.Prepare("INSERT INTO Examiners(name, age, email, mobileNo, password,clientId) VALUES(?,?,?,?,?,?)")
		if err != nil {
			return err
		}

		id := uuid.New()

		cid := strings.Replace(id.String(), "-", "", -1)

		result, err := query.ExecContext(ctx, user.Name, user.Age, user.Email, user.MobileNo, password, cid)
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

		_, err = query.ExecContext(ctx, lId, 1)
		if err != nil {
			return err
		}

		return err
	})

	return err
}

func (ur *userMySQLRepository) CreateVideo(ctx context.Context, fileName, videoUrl, imagePath, clientId, dstpath string) error {

	err := withTransaction(ur.MySQLDB, func(tx *sql.Tx) error {
		// Insert into DB
		query, err := tx.Prepare("INSERT INTO VideoContent(name, videoUrl,thumbnailPath,contentType,description,clientId) VALUES(?,?,?,?,?,?)")
		if err != nil {
			return err
		}

		_, err = query.ExecContext(ctx, fileName, videoUrl, imagePath, "video/mp4", "Sample Video", clientId)
		if err != nil {
			return err
		}

		ctx := context.Background()
		// Send message to queue
		err = ur.RabbitMQ.PublishWithContext(ctx,
			"",            // exchange
			ur.Queue.Name, // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(dstpath),
			})

		if err != nil {
			return err
		}

		return err
	})

	return err
}

func (ur *userMySQLRepository) ExamCreation(ctx context.Context, clientId, examName, examType string) (int64, error) {

	// Insert into DB
	query, err := ur.MySQLDB.Prepare("INSERT INTO Exams(clientId,name,type) VALUES(?,?,?)")
	if err != nil {
		return 0, err
	}

	result, err := query.ExecContext(ctx, clientId, examName, examType)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
