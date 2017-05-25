package memory

import (
	"github.com/PiotrPrzybylak/koop-wooch/domain"
	"strconv"
)

func NewCategoryRepository() domain.CategoryRepository {

	return &categoryRepository{}
}

type categoryRepository struct {
	ID         int
	categories []domain.Category
}

func (c categoryRepository) All() ([]domain.Category, error) {
	return c.categories, nil
}
func (c *categoryRepository) Create(category domain.Category) (string, error) {
	c.ID++
	category.ID = strconv.Itoa(c.ID)
	c.categories = append(c.categories, category)
	return category.ID, nil
}
