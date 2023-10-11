package repositories

import (
	"database/sql"
	"devbook_api/src/models"
)

type Posts struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *Posts {
	return &Posts{db}
}

func (repository Posts) Create(userId uint64, post models.Post) (uint64, error) {
	statement, err := repository.db.Prepare("insert into posts (title,content,author_id) values (?,?,?)")
	if err != nil {
		return 0, err
	}

	defer statement.Close()

	result, err := statement.Exec(post.Title, post.Content, post.AuthorId)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastId), nil
}

func (repository Posts) GetPostById(id uint64) (models.Post, error) {
	lines, err := repository.db.Query(`
	select p.*, u.nick from
	posts p inner join users u
	on u.id = p.author_id where p.id = ?
	`, id)
	if err != nil {
		return models.Post{}, err
	}
	var post models.Post
	if lines.Next() {

		err := lines.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		)
		if err != nil {
			return models.Post{}, err
		}
	}

	return post, nil
}

func (repository Posts) GetPosts(id uint64) ([]models.Post, error) {
	lines, err := repository.db.Query(`
	select distinct p.*, u.nick from posts p 
	join users u on u.id = p.author_id 
	left join followers s on p.author_id = s.user_id 
	where u.id = ? or s.follower_id = ?
	order by 1 desc`,
		id, id)
	if err != nil {
		return nil, err
	}
	var posts []models.Post
	for lines.Next() {
		var post models.Post
		err := lines.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (repository Posts) Update(id uint64, post models.Post) error {
	statement, err := repository.db.Prepare("update posts set title = ?, content = ? where id = ?")
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(post.Title, post.Content, id)
	if err != nil {
		return err
	}

	return nil
}

func (repository Posts) Delete(id uint64) error {
	statement, err := repository.db.Prepare("delete from posts where id = ?")
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (repository Posts) GetPostsByUser(id uint64) ([]models.Post, error) {
	lines, err := repository.db.Query(`
	select p.*, u.nick from
	posts p inner join users u
	on u.id = p.author_id 
	where p.author_id = ?
	order by 1 desc
	`, id)
	if err != nil {
		return nil, err
	}
	var posts []models.Post
	for lines.Next() {
		var post models.Post
		err := lines.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
func (repository Posts) Like(id uint64) error {
	statement, err := repository.db.Prepare("update posts set likes = likes + 1 where id = ?")
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (repository Posts) Dislike(id uint64) error {
	statement, err := repository.db.Prepare(`
	update posts set likes =
	CASE 
		WHEN likes > 0 THEN likes - 1
		ELSE likes 
	END 
	where id = ?
	`)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
