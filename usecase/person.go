package usecase

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"

	"github.com/pkg/errors"

	"github.com/ars0915/matching-system/constant"
	"github.com/ars0915/matching-system/entity"
	"github.com/ars0915/matching-system/internal/tree"
)

func (h *PersonHandler) GenerateNextID() uint64 {
	return atomic.AddUint64(h.id, 1)
}

func (h *PersonHandler) AddPerson(p entity.Person) (entity.Person, error) {
	var err error

	p.ID = h.GenerateNextID()
	switch p.Gender {
	case constant.GenderMale:
		err = h.boys.AddPerson(&p)
	case constant.GenderFemale:
		err = h.girls.AddPerson(&p)
	}

	return p, err
}

func (h *PersonHandler) findPerson(id uint64) (entity.Person, error) {
	if p, exist := h.boys.FindByID(id); exist {
		return *p, nil
	}
	if p, exist := h.girls.FindByID(id); exist {
		return *p, nil
	}

	return entity.Person{}, ErrorPersonNotFound
}

func (h *PersonHandler) RemovePerson(ctx context.Context, id uint64) error {
	person, err := h.findPerson(id)
	if err != nil {
		return err
	}

	switch person.Gender {
	case constant.GenderMale:
		err = h.boys.RemovePerson(id)
	case constant.GenderFemale:
		err = h.girls.RemovePerson(id)
	}

	if err != nil {
		if errors.Is(err, tree.ErrorPersonNotFound) {
			return ErrorPersonNotFound
		}
		return err
	}

	return nil
}

func (h *PersonHandler) QuerySinglePeople(ctx context.Context, id uint64, num int) ([]entity.Person, error) {
	var result []entity.Person

	person, err := h.findPerson(id)
	if err != nil {
		return nil, err
	}

	switch person.Gender {
	case constant.GenderMale:
		result = h.girls.QueryByHeight(0, person.Height)
	case constant.GenderFemale:
		result = h.boys.QueryByHeight(person.Height, math.MaxFloat64)
	}

	if len(result) > num {
		return result[:num], nil
	}
	return result, nil
}

func (h *PersonHandler) AddPersonAndFindMatch(ctx context.Context, p entity.Person) ([]entity.Person, error) {
	p, err := h.AddPerson(p)
	if err != nil {
		return nil, err
	}
	return h.QuerySinglePeople(ctx, p.ID, 1)
}

func (h *PersonHandler) Print() {
	fmt.Println("boys:")
	h.boys.Print()

	fmt.Println("girls:")
	h.girls.Print()
}
