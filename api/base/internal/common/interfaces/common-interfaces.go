package interfaces


type Errfunc func()

type Common[T any] interface {
	Create(data *T) Errfunc
	GetByID(data *T) Errfunc
	Fetch(limit, offset int) ([]T, Errfunc)
	Update(data *T) Errfunc
	Delete(data *T) Errfunc
}
