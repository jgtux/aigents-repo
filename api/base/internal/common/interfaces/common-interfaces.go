package interfaces


type Common[T any] interface {
	Create(data *T) func()
	GetByID(data *T) func()
	Fetch(limit, offset int) ([]T, func())
	Update(data *T) func()
	Delete(data *T) func()
}
