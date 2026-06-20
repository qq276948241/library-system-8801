package main

import "time"

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

func deriveStatus(dueDate string, returnDate *string) string {
	if returnDate != nil && *returnDate != "" {
		return "returned"
	}
	t, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		return "borrowed"
	}
	if t.Before(time.Now()) {
		return "overdue"
	}
	return "borrowed"
}

func keysInt64(m map[int64]bool) []int64 {
	ks := make([]int64, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func enrichBorrow(b *Borrow) {
	b.BorrowDate = normTime(b.BorrowDate)
	b.DueDate = normTime(b.DueDate)
	if b.ReturnDate != nil {
		rd := normTime(*b.ReturnDate)
		b.ReturnDate = &rd
	}
	b.Status = deriveStatus(b.DueDate, b.ReturnDate)
	if b.User != nil {
		b.User.CreatedAt = normTime(b.User.CreatedAt)
	}
	if b.Book != nil {
		b.Book.CreatedAt = normTime(b.Book.CreatedAt)
	}
}
