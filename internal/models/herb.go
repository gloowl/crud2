package models

import (
	"fmt"
	"strings"
	"time"
)

// Herb - трава
type Herb struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	LatinName   string    `json:"latin_name"`
	Description string    `json:"description"`
	IsPoisonous bool      `json:"is_poisonous"`
	ImagePath   string    `json:"image_path"`
	CreatedAt   time.Time `json:"created_at"`
}

func (h *Herb) String() string {
	poisonous := "Нет"
	if h.IsPoisonous {
		poisonous = "ДА!"
	}

	return fmt.Sprintf(`
ID: %d
Название: %s
Латинское название: %s
Описание: %s
Ядовито: %s
Изображение: %s
Создано: %s`,
		h.ID,
		h.Name,
		h.LatinName,
		truncateString(h.Description, 100),
		poisonous,
		h.ImagePath,
		h.CreatedAt.Format("2006-01-02 15:04:05"),
	)
}

func (h *Herb) Validate() error {
	if strings.TrimSpace(h.Name) == "" {
		return fmt.Errorf("название травы не может быть пустым")
	}

	if len(h.Name) < 2 {
		return fmt.Errorf("название травы должно содержать минимум 2 символа")
	}

	if len(h.Name) > 255 {
		return fmt.Errorf("название травы не должно превышать 255 символов")
	}

	if h.LatinName != "" && len(h.LatinName) > 255 {
		return fmt.Errorf("латинское название не должно превышать 255 символов")
	}

	if h.ImagePath != "" && len(h.ImagePath) > 500 {
		return fmt.Errorf("путь к изображению не должен превышать 500 символов")
	}

	return nil
}

// TableHeader returns the table header for herbs
func (h *Herb) TableHeader() string {
	return fmt.Sprintf("%-4s %-20s %-25s %-8s %-15s",
		"ID", "Название", "Латинское название", "Ядовито", "Создано")
}

// TableRow returns a formatted table row for the herb
func (h *Herb) TableRow() string {
	poisonous := "Нет"
	if h.IsPoisonous {
		poisonous = "ДА!"
	}

	return fmt.Sprintf("%-4d %-20s %-25s %-8s %-15s",
		h.ID,
		truncateString(h.Name, 20),
		truncateString(h.LatinName, 25),
		poisonous,
		h.CreatedAt.Format("2006-01-02"),
	)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
