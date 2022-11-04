package repository

import (
	"context"
	"database/sql"
	"encoding/json"
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

func (cs *commonMySQLRepository) LoginUser(userLogin model.UserLogin) (int, string, string, error) {

	var (
		id       int
		password string
		userType string = "User"
	)

	row := cs.MySQLDB.QueryRow("select id,password from Users where email=?", userLogin.Email)

	err := row.Scan(&id, &password)

	if err != nil {
		if err == sql.ErrNoRows {
			row := cs.MySQLDB.QueryRow("select id,password from Students where email=?", userLogin.Email)

			err = row.Scan(&id, &password)

			if err != nil {
				if err == sql.ErrNoRows {
					return 0, "", "", err
				}
			}
			userType = "Student"
		}
	}
	return id, userType, password, nil
}

func (cs *commonMySQLRepository) ReadRoutes(userId int) ([]model.Route, error) {

	var routes []model.Route

	// rows, err := Db.Query(`select * from menu where id in(
	// 	select menuId from roleMenu where roleId =(
	// 	select roleId from userRole where userId=(
	// 	select id from examiner where email=?
	// 	)))`, email)

	// if email != "admin@example.org" {
	// 	val = 2
	// } else {
	// 	val = 1
	// }

	// rows, err := Db.Query(`select * from menu where id in(
	// 	select menuId from roleMenu where roleId =(
	// 	select roleId from userRole where userId=?))`, val)

	rows, err := cs.MySQLDB.Query(`SELECT m.id,m.name,m.url,m.description FROM UserRole ur
	INNER JOIN RoleMenu rm ON ur.roleId = rm.roleId
	INNER JOIN Menu m ON rm.menuId = m.id
	where ur.userId=?`, userId)

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

	return routes, nil
}

func (cs *commonMySQLRepository) ReadVideos() ([]model.Video, error) {
	var (
		videos []model.Video
		ctx    = context.Background()
	)
	const id string = "videoData"

	val, err := cs.Redis.Get(ctx, id).Bytes()

	if err != nil {
		log.Println(err)
	} else {
		if val != nil {
			json.Unmarshal(val, &videos)

			return videos, nil
		}
	}

	rows, err := cs.MySQLDB.Query(`SELECT name, videoUrl,thumbnailPath,description from VideoContent`)
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
		err = cs.Redis.Set(ctx, id, json, 15*time.Minute).Err()
		if err != nil {
			log.Println(err)
		}
	}
	return videos, nil
}
