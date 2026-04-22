package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	redis "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mskKandula/oes/api/model"
)

type commonMySQLRepository struct {
	MySQLDB *sql.DB
	Redis   *redis.Client
}

func NewCommonMySQLRepository(rc *RepositoryConfig) model.CommonRepository {
	return &commonMySQLRepository{
		MySQLDB: rc.MySQLDB,
		Redis:   rc.Redis,
	}
}

func (cs *commonMySQLRepository) LoginUser(ctx context.Context, userLogin model.UserLogin) (int, string, string, string, error) {

	var (
		id       int
		password string
		clientId string
		userType string
	)

	query := `
		SELECT id, password, clientId, 'Examiner' AS userType
		FROM Examiners
		WHERE email = ?

		UNION ALL

		SELECT id, password, clientId, 'Student' AS userType
		FROM Students
		WHERE email = ?
		LIMIT 1`

	row := cs.MySQLDB.QueryRowContext(ctx, query, userLogin.Email, userLogin.Email)

	err := row.Scan(&id, &password, &clientId, &userType)

	if err != nil {
		return 0, "", "", "", err
	}

	return id, userType, password, clientId, nil
}

func (cs *commonMySQLRepository) ReadRoutes(ctx context.Context, userId int, userType string) ([]model.Route, error) {

	var routes []model.Route

	// --- Cache read ---
	cacheKey := fmt.Sprintf("routes:%s", userType)

	val, err := cs.Redis.Get(ctx, cacheKey).Bytes()
	if err != nil {
		log.Println(err)
	} else {
		if val != nil {
			json.Unmarshal(val, &routes)
			return routes, nil
		}
	}

	// --- MySQL fallback ---
	rows, err := cs.MySQLDB.QueryContext(ctx, `SELECT m.id,m.name,m.url,m.description FROM Role r
    INNER JOIN UserRole ur ON r.id = ur.roleId
	INNER JOIN RoleMenu rm ON ur.roleId = rm.roleId
	INNER JOIN Menu m ON rm.menuId = m.id
	WHERE ur.userId=? AND r.name=? ORDER BY m.id;`, userId, userType)

	if err != nil {
		return routes, err
	}

	defer rows.Close()

	for rows.Next() {
		var route model.Route

		if err := rows.Scan(&route.Id, &route.Name, &route.Url, &route.Description); err != nil {
			return routes, err
		}
		routes = append(routes, route)
	}

	// --- Cache write ---
	jsonData, err := json.Marshal(routes)
	if err != nil {
		log.Println(err)
	} else {
		err = cs.Redis.Set(ctx, cacheKey, jsonData, 24*time.Hour).Err()
		if err != nil {
			log.Println(err)
		}
	}

	return routes, nil
}

func (cs *commonMySQLRepository) ReadVideos(ctx context.Context, clientId string) ([]model.Video, error) {
	var (
		videos []model.Video
	)

	val, err := cs.Redis.Get(ctx, clientId).Bytes()

	if err != nil {
		log.Println(err)
	} else {
		if val != nil {
			json.Unmarshal(val, &videos)

			return videos, nil
		}
	}

	rows, err := cs.MySQLDB.QueryContext(ctx, `SELECT name, videoUrl,thumbnailPath,description from VideoContent where clientId = ?`, clientId)
	if err != nil {
		return videos, err
	}

	defer rows.Close()

	for rows.Next() {
		var video model.Video

		if err := rows.Scan(&video.Name, &video.VideoUrl, &video.ThumbnailPath, &video.Description); err != nil {
			return videos, err
		}
		videos = append(videos, video)
	}

	json, err := json.Marshal(videos)

	if err != nil {
		log.Println(err)
	} else {
		err = cs.Redis.Set(ctx, clientId, json, 15*time.Minute).Err()
		if err != nil {
			log.Println(err)
		}
	}
	return videos, nil
}
