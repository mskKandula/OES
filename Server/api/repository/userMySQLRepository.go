package repository

import (
	"context"
	"database/sql"
	"time"

	redis "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/mskKandula/oes/api/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type userMySQLRepository struct {
	MySQLDB  *sql.DB
	RabbitMQ *amqp.Channel
	Queue    amqp.Queue
	Redis    *redis.Client
}

func NewUserMySQLRepository(rc *RepositoryConfig) model.UserRepository {
	return &userMySQLRepository{
		MySQLDB:  rc.MySQLDB,
		RabbitMQ: rc.RabbitMQ,
		Queue:    rc.Queue,
		Redis:    rc.Redis,
	}
}

func (ur *userMySQLRepository) Create(ctx context.Context, user model.User, password string) error {

	return withTransactionContext(ctx, ur.MySQLDB, func(tx *sql.Tx) error {
		id := uuid.New().String()

		insertExaminer := "INSERT INTO Examiners(name, age, email, mobileNo, password, clientId) VALUES(?,?,?,?,?,?)"
		result, err := tx.ExecContext(ctx, insertExaminer,
			user.Name, user.Age, user.Email, user.MobileNo, password, id,
		)
		if err != nil {
			return err
		}

		lId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		insertUserRole := "INSERT INTO UserRole(userId, roleId) VALUES(?,?)"
		_, err = tx.ExecContext(ctx, insertUserRole, lId, 1)
		return err
	})
}

func (ur *userMySQLRepository) CreateVideo(ctx context.Context, fileName, videoUrl, imagePath, clientId, dstpath string) error {

	err := withTransactionContext(ctx, ur.MySQLDB, func(tx *sql.Tx) error {
		insertVideo := "INSERT INTO VideoContent(name, videoUrl, thumbnailPath, contentType, description, clientId) VALUES(?,?,?,?,?,?)"
		_, err := tx.ExecContext(ctx, insertVideo,
			fileName, videoUrl, imagePath, "video/mp4", "Sample Video", clientId,
		)
		if err != nil {
			return err
		}

		// Send message to queue using the caller's context so cancellation/deadlines propagate
		return ur.RabbitMQ.PublishWithContext(ctx,
			"",            // exchange
			ur.Queue.Name, // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Timestamp:    time.Now(),
				Body:         []byte(dstpath),
			})
	})

	ur.Redis.Del(ctx, clientId) // Invalidate the video cache to prevent stale data
	return err
}

func (ur *userMySQLRepository) ExamCreation(ctx context.Context, clientId, examName, examType string) (int64, error) {

	insertExam := "INSERT INTO Exams(clientId, name, type) VALUES(?,?,?)"
	result, err := ur.MySQLDB.ExecContext(ctx, insertExam, clientId, examName, examType)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
