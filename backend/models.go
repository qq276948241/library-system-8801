package main

import "time"

type User struct {
	ID           int64  `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"-"`
	Name         string `db:"name" json:"name"`
	Role         string `db:"role" json:"role"`
	Email        string `db:"email" json:"email"`
	CreatedAt    string `db:"created_at" json:"created_at"`
}

type Book struct {
	ID              int64  `db:"id" json:"id"`
	Title           string `db:"title" json:"title"`
	Author          string `db:"author" json:"author"`
	ISBN            string `db:"isbn" json:"isbn"`
	Category        string `db:"category" json:"category"`
	Description     string `db:"description" json:"description"`
	CoverColor      string `db:"cover_color" json:"cover_color"`
	Publisher       string `db:"publisher" json:"publisher"`
	PublishedYear   int    `db:"published_year" json:"published_year"`
	TotalCopies     int    `db:"total_copies" json:"total_copies"`
	AvailableCopies int    `db:"available_copies" json:"available_copies"`
	CreatedAt       string `db:"created_at" json:"created_at"`
}

type Borrow struct {
	ID         int64  `db:"id" json:"id"`
	UserID     int64  `db:"user_id" json:"user_id"`
	BookID     int64  `db:"book_id" json:"book_id"`
	BorrowDate string `db:"borrow_date" json:"borrow_date"`
	DueDate    string `db:"due_date" json:"due_date"`
	ReturnDate string `db:"return_date" json:"return_date"`
	Status     string `db:"status" json:"status"`
	User       *User  `json:"user,omitempty"`
	Book       *Book  `json:"book,omitempty"`
}

type CategoryCount struct {
	Category string `db:"category" json:"category"`
	Count    int64  `db:"count" json:"count"`
}

type Stats struct {
	TotalBooks      int64           `db:"total_books" json:"total_books"`
	TotalCopies     int64           `db:"total_copies" json:"total_copies"`
	AvailableCopies int64           `db:"available_copies" json:"available_copies"`
	TotalUsers      int64           `db:"total_users" json:"total_users"`
	ActiveBorrows   int64           `db:"active_borrows" json:"active_borrows"`
	OverdueBorrows  int64           `db:"overdue_borrows" json:"overdue_borrows"`
	Categories      []CategoryCount `json:"categories"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

type loginResult struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func normTime(s string) string {
	if s == "" {
		return ""
	}
	layouts := []string{
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.UTC().Format(time.RFC3339)
		}
	}
	return s
}
