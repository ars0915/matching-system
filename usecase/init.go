package usecase

import "github.com/ars0915/matching-system/repo"

func InitHandler(db repo.App) Handler {
	task := NewTaskHandler(db)

	h := newHandler(
		WithTask(task),
	)

	return h
}
