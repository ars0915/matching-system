package usecase

import (
	"context"
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

func (h *PersonHandler) findPerson(id uint64) (*entity.Person, error) {
	if p, exist := h.boys.FindByID(id); exist {
		return p, nil
	}
	if p, exist := h.girls.FindByID(id); exist {
		return p, nil
	}

	return nil, ErrorPersonNotFound
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

func (h *PersonHandler) Match(ctx context.Context, id1, id2 uint64) error {
	person1, err := h.findPerson(id1)
	if err != nil {
		return err
	}

	person2, err := h.findPerson(id2)
	if err != nil {
		return err
	}

	if person1.Gender == person2.Gender {
		return ErrorMatchSameGender
	}

	// Check height
	if (person1.Gender == constant.GenderMale && person1.Height < person2.Height) ||
		(person1.Gender == constant.GenderFemale && person1.Height > person2.Height) {
		return ErrorHeightCheckFailed
	}

	if !h.tryMatch(person1, person2) {
		return ErrorWantedDateLimit
	}

	return nil
}

func (h *PersonHandler) tryMatch(person1, person2 *entity.Person) bool {
	// Atomically decrement wanted dates and check if either person has exhausted their dates
	if !h.decrementWantedDate(person1) || !h.decrementWantedDate(person2) {
		return false
	}

	// Remove from the system if any person's dates reach 0
	h.removeIfExhausted(person1)
	h.removeIfExhausted(person2)

	return true
}

func (h *PersonHandler) decrementWantedDate(person *entity.Person) bool {
	current := atomic.LoadUint64(person.WantedDates)
	if current == 0 {
		return false
	}

	if atomic.CompareAndSwapUint64(person.WantedDates, current, current-1) {
		return true
	}
	return false
}

func (h *PersonHandler) removeIfExhausted(person *entity.Person) {
	if atomic.LoadUint64(person.WantedDates) == 0 {
		// Remove from the appropriate gender group
		// Ignore the error because person has already been removed
		switch person.Gender {
		case constant.GenderMale:
			_ = h.boys.RemovePerson(person.ID)
		case constant.GenderFemale:
			_ = h.girls.RemovePerson(person.ID)
		}
	}
}
