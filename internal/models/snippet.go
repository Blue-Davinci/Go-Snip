package models

import (
	"database/sql"
	"errors"
	"time"
)

// SnippetModelInterface is an interface that defines the methods that a SnippetModel must implement
// We're going to use it better for testing as well.
type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// Use the Exec() method on the embedded connection pool to execute the
	// statement.
	result, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, nil
	}
	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the our table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// Create empty struct pointer to hold returned values
	snip := &Snippet{}
	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement and Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct.
	err := m.DB.QueryRow(query, id).Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return snip, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//create an empty slice of snippet structs
	snippets := []*Snippet{}
	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		snip := &Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created.
		err = rows.Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)
		if err != nil {
			return nil, err
		}
		//append the new snipet struct
		snippets = append(snippets, snip)
	}
	//check if err occured during rows operation
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
