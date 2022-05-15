package users

import (
	"fmt"
	"go-ms-bookstore-user/datasource/postgres/users_db"
	"go-ms-bookstore-user/logger"
	"go-ms-bookstore-user/utils/errors"
	pgutils "go-ms-bookstore-user/utils/pg_utils"
)

const (
	indexUniqueEmail    = "pq: duplicate key value violates unique constraint"
	noUserError         = "sql: no rows in result set"
	queryUserById       = " id,first_name,last_name,email,date_created ,status FROM users WHERE id = $1;"
	queryInsertUser     = "INSERT INTO users (first_name, last_name, email, date_created, status, password) VALUES($1, $2, $3, $4, $5, $6) RETURNING id;"
	queryUpdateUser     = "UPDATE users SET first_name = $1, last_name =$2, email =$3 WHERE id =$4;"
	queryDeleteUser     = "DELETE FROM users WHERE id =$1;"
	querySearchByStatus = "Select id,first_name,last_name,email,date_created ,status FROM users WHERE status = $1;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUserById)

	if err != nil {
		logger.Error("Error while preparing the query", err)
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	err = stmt.QueryRow(user.Id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)

	if err != nil {
		logger.Error("Error while scan/query", err)
		return pgutils.ParseError(err)
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)

	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	fmt.Println(user)
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		return pgutils.ParseError(saveErr)
	}

	rowAffected, err := insertResult.RowsAffected()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("Error while saving user: %s", err.Error()))
	}

	fmt.Println("Row Affected", rowAffected)

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)

	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	fmt.Println(user)
	updateResult, updateErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if updateErr != nil {
		return pgutils.ParseError(updateErr)
	}

	rowAffected, err := updateResult.RowsAffected()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("Error while updating user: %s", err.Error()))
	}

	fmt.Println("Row Affected", rowAffected)

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)

	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	fmt.Println(user)
	deleteResult, deleteErr := stmt.Exec(user.Id)
	if deleteErr != nil {
		return pgutils.ParseError(deleteErr)
	}

	rowAffected, err := deleteResult.RowsAffected()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("Error while deleting user: %s", err.Error()))
	}

	fmt.Println("Row Affected", rowAffected)

	return nil
}

func (user *User) Search(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(querySearchByStatus)

	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, searchErr := stmt.Query(status)

	if searchErr != nil {
		return nil, pgutils.ParseError(searchErr)
	}

	defer rows.Close()

	var users []User = nil

	for rows.Next() {
		error := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
		if error != nil {
			return nil, pgutils.ParseError(error)
		}

		users = append(users, *user)
	}

	if len(users) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("No users found for status %s", status))
	}
	return users, nil
}
