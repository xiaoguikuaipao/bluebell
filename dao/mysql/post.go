package mysql

import (
	"strings"

	"web_app/models"

	"github.com/jmoiron/sqlx"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(
		post_id, title, content, author_id, community_id)
		values (?, ?, ?, ?, ?)
		`
	if _, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID); err != nil {
		return err
	}
	return err
}

func GetPostByID(pid int64) (*models.Post, error) {
	p := new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id from post where post_id = ?`
	err := db.Get(p, sqlStr, pid)
	if err != nil {
		return nil, err
	}
	return p, err
}

func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, created_time 
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)
	`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	postList = make([]*models.Post, 0, len(ids))
	err = db.Select(&postList, query, args...)
	return
}

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
					from post 
					order by create_time
					desc 
					limit ?, ?`
	posts = make([]*models.Post, 0, 3)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}
