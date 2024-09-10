// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/ars0915/matching-system/entity"
)

type (
	Handler interface {
		Person
	}
)

type (
	Person interface {
		AddPersonAndFindMatch(ctx context.Context, p entity.Person) ([]entity.Person, error)
		RemovePerson(ctx context.Context, id uint64) error
		QuerySinglePeople(ctx context.Context, id uint64, num int) ([]entity.Person, error)
		Match(ctx context.Context, id1, id2 uint64) error
	}
)
