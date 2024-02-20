package repositories

import (
	"database/sql"
	"devbook_api/src/models"
	"errors"
	"fmt"
)

type Users struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *Users {
	return &Users{db}
}

func (repository Users) Create(user models.User) (uint64, error) {
	statement, err := repository.db.Prepare(
		"insert into users (name, nick, email, password) values (?,?,?,?)",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(user.Name, user.Nick, user.Email, user.Password)

	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(lastId), nil
}

func (repository Users) Get(nameOrNick string) ([]models.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick) // %nameOrNick%

	lines, err := repository.db.Query(
		"select id, name, nick, email, createdAt from users where name LIKE ? or nick LIKE ?",
		nameOrNick, nameOrNick,
	)
	if err != nil {
		return nil, err
	}

	defer lines.Close()

	var users []models.User

	for lines.Next() {
		var user models.User

		if err = lines.Scan(
			&user.Id,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil

}

func (repository Users) GetById(userId int64) (models.User, error) {

	lines, err := repository.db.Query(
		"select id, name, nick, email, createdAt from users where id = ?",
		userId,
	)
	if err != nil {
		return models.User{}, err
	}

	defer lines.Close()

	var user models.User

	if lines.Next() {
		if err = lines.Scan(
			&user.Id,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return models.User{}, err
		}

	} else {
		return models.User{}, errors.New("user not found with this id")
	}

	return user, nil

}

func (repository Users) Update(userId uint64, user models.User) error {
	statement, err := repository.db.Prepare(
		"update users set name = ?, nick = ?, email = ? where id = ?",
	)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(user.Name, user.Nick, user.Email, userId)
	if err != nil {
		return err
	}

	return nil
}

func (repository Users) Delete(userId uint64) error {
	statement, err := repository.db.Prepare(
		"delete from users where id = ?",
	)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(userId)
	if err != nil {
		return err
	}

	return nil
}

func (repository Users) SearchByEmail(email string) (models.User, error) {
	line, err := repository.db.Query("select id, password from users where email = ?", email)
	if err != nil {
		return models.User{}, err
	}
	defer line.Close()

	var user models.User
	if line.Next() {
		err := line.Scan(
			&user.Id,
			&user.Password,
		)
		if err != nil {
			return models.User{}, err
		}
	} else {
		return models.User{}, errors.New("user not found with this email")
	}

	return user, nil
}

func (repository Users) FollowUser(followerId, userId uint64) error {
	statement, err := repository.db.Prepare(
		"insert ignore into followers (user_id, follower_id) values (?, ?)",
	)

	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(followerId, userId)

	if err != nil {
		return err
	}

	return nil
}

func (repository Users) UnfollowUser(followerId, userId uint64) error {
	statement, err := repository.db.Prepare(
		"delete from followers where user_id = ? and follower_id = ?",
	)

	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(followerId, userId)

	if err != nil {
		return err
	}

	return nil
}

func (repository Users) GetFollowers(userId uint64) ([]models.User, error) {
	lines, err := repository.db.Query(
		`
		select u.id, u.name, u.nick, u.email, u.createdAt
		from users u inner join followers s on u.id = s.follower_id where s.user_id = ?
		`, userId,
	)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var users []models.User

	for lines.Next() {
		var user models.User

		if err = lines.Scan(
			&user.Id,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)

	}
	return users, nil
}

func (repository Users) GetFollowings(userId uint64) ([]models.User, error) {
	lines, err := repository.db.Query(
		`
		select u.id, u.name, u.nick, u.email, u.createdAt
		from users u inner join followers s on u.id = s.user_id where s.follower_id = ?
		`, userId,
	)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var users []models.User

	for lines.Next() {
		var user models.User

		if err = lines.Scan(
			&user.Id,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)

	}
	return users, nil
}

func (repository Users) SearchPassword(id uint64) (string, error) {
	line, err := repository.db.Query("select password from users where id = ?", id)
	if err != nil {
		return "", err
	}
	defer line.Close()

	var password string
	if line.Next() {
		err := line.Scan(
			&password,
		)
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("user not found with this id")
	}

	return password, nil
}

func (repository Users) UpdatePassword(userId uint64, password []byte) error {
	statement, err := repository.db.Prepare(
		"update users set password = ? where id = ?",
	)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(password, userId)
	if err != nil {
		return err
	}

	return nil
}
