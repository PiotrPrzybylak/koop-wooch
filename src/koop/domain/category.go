package domain

type Category struct {
	ID   string
	Name string
}

type CategoryService interface {
	Create(category Category) (string, error)
	All() ([]Category, error)
}

type CategoryRepository interface {
	Create(category Category) (string, error)
	All() ([]Category, error)
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return categoryService{repo}
}

type categoryService struct {
	repo CategoryRepository
}

func (s categoryService) All() ([]Category, error) {
	return s.repo.All()
}
func (s categoryService) Create(category Category) (string, error) {
	return s.repo.Create(category)
}
