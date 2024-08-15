package types

import "errors"

var (
	ErrEmptyTitle         = errors.New("заголовок не может быть пустым")
	ErrInvalidTitle       = errors.New("длина заголовка должна превышать 200 символов")
	ErrInvalidDescription = errors.New("длина описания не должна превышать 1000 симоволов")
	ErrInvalidPhotos      = errors.New("нельзя загрузить больше чем 3 ссылки на фото")
	ErrInvalidPageNumber  = errors.New("невалидный номер страницы")
	ErrInvalidUserId      = errors.New("невалидный идентификатор пользователя")
	ErrInvalidPriceSort   = errors.New("невалидный параметр сортировки по цене")
	ErrInvalidDateSort    = errors.New("невалидный параметр сортировки по дате")
	ErrInvalidId          = errors.New("невалидный идентификатор (id) объявления")
	ErrInvalidToken       = errors.New("invalid token")
	ErrFailedIsAdminCheck = errors.New("failed to check if user is admin")
	ErrAdNotFound         = errors.New("объявление не найдено")
)

var (
	KeyUser  = "user"
	KeyError = "error"
)
