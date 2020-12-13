package mysql

import (
	"database/sql"
	"errors"

	"github.com/DeviaVir/servente/pkg/models"
)

type ServiceModel struct {
	DB *sql.DB
}

func (m *ServiceModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO services (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *ServiceModel) Get(id int) (*models.Service, error) {
	stmt := `SELECT id, title, content, created, expires FROM services
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Service{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *ServiceModel) Latest() ([]*models.Service, error) {
	stmt := `SELECT id, title, content, created, expires FROM services
		WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := []*models.Service{}

	for rows.Next() {
		s := &models.Service{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		services = append(services, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}
