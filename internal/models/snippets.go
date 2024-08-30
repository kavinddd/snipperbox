package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	query := `
	INSERT INTO snippets (title, content, created, expires) 
		VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
	`

	result, err := m.DB.Exec(query, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `
	SELECT id, title, content, created, expires 
	FROM snippets 
	WHERE id = ?
		AND expires > UTC_TIMESTAMP()
	`

	row := m.DB.QueryRow(query, id)

	s := &Snippet{}

	err := row.Scan(
		&s.ID,
		&s.Title,
		&s.Content,
		&s.Created,
		&s.Expires,
	)

	if err != nil {
		return nil, ErrNoRecord
	}

	return s, nil

}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	query := `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10	
	`

	rows, err := m.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := make([]*Snippet, 0)

	for rows.Next() {
		tempSnippet := &Snippet{}
		err = rows.Scan(
			&tempSnippet.ID,
			&tempSnippet.Title,
			&tempSnippet.Content,
			&tempSnippet.Created,
			&tempSnippet.Expires,
		)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, tempSnippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
