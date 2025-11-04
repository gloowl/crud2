package repository

import (
	"crud2/internal/models"
	"database/sql"
	"fmt"
)

type HerbRepository struct {
	db *sql.DB
}

func NewHerbRepository(db *sql.DB) *HerbRepository {
	return &HerbRepository{db: db}
}

// Create adds a new herb to the database
func (r *HerbRepository) Create(herb *models.Herb) error {
	if err := herb.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO herbs (name, latin_name, description, is_poisonous, image_path) 
		VALUES (?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query, herb.Name, herb.LatinName, herb.Description, herb.IsPoisonous, herb.ImagePath)

	// err := r.db.QueryRow(query, herb.Name, herb.LatinName, herb.Description, herb.IsPoisonous, herb.ImagePath).Scan(&herb.ID, &herb.CreatedAt)

	if err != nil {
		return fmt.Errorf("ошибка создания травы: %v", err)
	}

	id, _ := res.LastInsertId()
	herb.ID = int(id)

	return nil
}

// GetByID retrieves a herb by its ID
func (r *HerbRepository) GetByID(id int) (*models.Herb, error) {
	herb := &models.Herb{}
	query := `
		SELECT id, name, latin_name, description, is_poisonous, image_path, created_at
		FROM herbs 
		WHERE id = ?;`

	err := r.db.QueryRow(query, id).Scan(
		&herb.ID, &herb.Name, &herb.LatinName, &herb.Description,
		&herb.IsPoisonous, &herb.ImagePath, &herb.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("трава с ID %d не найдена", id)
		}
		return nil, fmt.Errorf("ошибка получения травы: %v", err)
	}
	return herb, nil
}

// GetAll retrieves all herbs
func (r *HerbRepository) GetAll() ([]models.Herb, error) {
	query := `
		SELECT id, name, latin_name, description, is_poisonous, image_path, created_at
		FROM herbs 
		ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка трав: %v", err)
	}
	defer rows.Close()

	var herbs []models.Herb
	for rows.Next() {
		herb := models.Herb{}
		err := rows.Scan(&herb.ID, &herb.Name, &herb.LatinName, &herb.Description,
			&herb.IsPoisonous, &herb.ImagePath, &herb.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования травы: %v", err)
		}
		herbs = append(herbs, herb)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по травам: %v", err)
	}

	return herbs, nil
}

// Update modifies an existing herb
func (r *HerbRepository) Update(herb *models.Herb) error {
	if err := herb.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE herbs 
		SET name = ?, latin_name = ?, description = ?, 
		    is_poisonous = ?, image_path = ?
		WHERE id = ?;`

	result, err := r.db.Exec(query, herb.Name, herb.LatinName,
		herb.Description, herb.IsPoisonous, herb.ImagePath, herb.ID)

	if err != nil {
		return fmt.Errorf("ошибка обновления травы: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества затронутых строк: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("трава с ID %d не найдена", herb.ID)
	}

	return nil
}

// Delete removes a herb from the database
func (r *HerbRepository) Delete(id int) error {
	query := `DELETE FROM herbs WHERE id = ?;`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления травы: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества затронутых строк: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("трава с ID %d не найдена", id)
	}

	return nil
}

// Search finds herbs by name (case-insensitive partial match)
func (r *HerbRepository) Search(name string) ([]models.Herb, error) {
	query := `
		SELECT id, name, latin_name, description, is_poisonous, image_path, created_at
		FROM herbs 
		WHERE LOWER(name) LIKE LOWER($1) OR LOWER(latin_name) LIKE LOWER($1)
		ORDER BY name;`

	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска трав: %v", err)
	}
	defer rows.Close()

	var herbs []models.Herb
	for rows.Next() {
		herb := models.Herb{}
		err := rows.Scan(&herb.ID, &herb.Name, &herb.LatinName, &herb.Description,
			&herb.IsPoisonous, &herb.ImagePath, &herb.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования травы: %v", err)
		}
		herbs = append(herbs, herb)
	}

	return herbs, rows.Err()
}

// GetPoisonous retrieves all poisonous herbs
func (r *HerbRepository) GetPoisonous() ([]models.Herb, error) {
	query := `
		SELECT id, name, latin_name, description, is_poisonous, image_path, created_at
		FROM herbs 
		WHERE is_poisonous = true
		ORDER BY name;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ядовитых трав: %v", err)
	}
	defer rows.Close()

	var herbs []models.Herb
	for rows.Next() {
		herb := models.Herb{}
		err := rows.Scan(&herb.ID, &herb.Name, &herb.LatinName, &herb.Description,
			&herb.IsPoisonous, &herb.ImagePath, &herb.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования травы: %v", err)
		}
		herbs = append(herbs, herb)
	}

	return herbs, rows.Err()
}
