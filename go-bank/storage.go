package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	DeleteAccount(int) error
	UpdateAcoount(*Account) error
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=mysecretpassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (db *PostgresStore) Init() error {
	return db.CreateAccountsTable()
}

func (db *PostgresStore) CreateAccountsTable() error {
	query := ` create table if not exists account( 
	id serial primary key,
	first_name varchar(50),
	last_name varchar(50),
	number serial,
	balance float8,
	created_at timestamp
	)
	`
	if _, err := db.db.Exec(query); err != nil {
		return err
	}
	return nil
}
func (db *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into account
	(first_name, last_name, number, balance, created_at)
	values ($1, $2, $3, $4, $5)`
	resp, err := db.db.Exec(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	return nil
}

func (db *PostgresStore) GetAccountById(id int) (*Account, error) {
	rows, err := db.db.Query("select * from account where id = $1", id)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, fmt.Errorf("Account with id %d not found", id)
	}
	account, err := scanIntoAccount(rows)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (db *PostgresStore) DeleteAccount(id int) error {
	_, err := db.db.Query("delete from account where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (db *PostgresStore) UpdateAcoount(acc *Account) error {
	return nil
}

func (db *PostgresStore) GetAccounts() ([]*Account, error) {
	accounts := []*Account{}
	rows, err := db.db.Query("select * from account")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}
