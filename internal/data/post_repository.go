package data

import (
	"context"

	"github.com/lucabecci/golang-blog-psql.git/pkg/post"
)

type PostRespository struct {
	Data *Data
}

func (pr *PostRespository) GetAll(ctx context.Context) ([]post.Post, error) {
	q := `
    SELECT id, body, user_id, created_at, updated_at
        FROM posts;
    `

	rows, err := pr.Data.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []post.Post
	for rows.Next() {
		var p post.Post
		rows.Scan(&p.ID, &p.Body, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
		posts = append(posts, p)
	}

	return posts, nil
}

func (pr *PostRespository) GetOne(ctx context.Context, id uint) (post.Post, error) {
	q := `
    SELECT id, body, user_id, created_at, updated_at
        FROM posts WHERE id = $1;
    `

	row := pr.Data.DB.QueryRowContext(ctx, q, id)

	var p post.Post
	err := row.Scan(&p.ID, &p.Body, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return post.Post{}, err
	}

	return p, nil
}

func (pr *PostRespository) GetByUser(ctx context.Context, userID uint) ([]post.Post, error) {
	q := `
    SELECT id, body, user_id, created_at, updated_at
        FROM posts
        WHERE user_id = $1;
    `

	rows, err := pr.Data.DB.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []post.Post
	for rows.Next() {
		var p post.Post
		rows.Scan(&p.ID, &p.Body, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
		posts = append(posts, p)
	}

	return posts, nil
}
