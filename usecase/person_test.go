package usecase

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/ars0915/matching-system/constant"
	"github.com/ars0915/matching-system/entity"
	mocks "github.com/ars0915/matching-system/internal/mocks/tree"
	ctest "github.com/ars0915/matching-system/util/cTest"
	"github.com/ars0915/matching-system/util/cTypes"
)

type personTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller

	h     *PersonHandler
	boys  *mocks.MockTree
	girls *mocks.MockTree
}

func Test_personTestSuite(t *testing.T) {
	suite.Run(t, &personTestSuite{})
}

func (s *personTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.boys = mocks.NewMockTree(s.ctrl)
	s.girls = mocks.NewMockTree(s.ctrl)
	s.h = NewPersonHandler(s.boys, s.girls)
}

func (s *personTestSuite) TearDownTest(t *testing.T) {
	defer s.ctrl.Finish()
}

func (s *personTestSuite) initPeople(people []entity.Person) {
	for _, person := range people {
		if person.Gender == constant.GenderMale {
			s.boys.EXPECT().AddPerson(gomock.Any()).Return(nil)
		} else {
			s.girls.EXPECT().AddPerson(gomock.Any()).Return(nil)
		}
		_, err := s.h.AddPerson(person)
		assert.Nil(s.T(), err)
	}
}

func (s *personTestSuite) Test_AddPerson() {
	person := entity.Person{
		ID:          1,
		Name:        "a",
		Height:      150,
		Gender:      "male",
		WantedDates: cTypes.Uint64(2),
	}

	s.boys.EXPECT().AddPerson(ctest.DiffWrapper(&person)).Return(nil)

	actualPerson, err := s.h.AddPerson(person)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), person, actualPerson)
}

func (s *personTestSuite) Test_RemovePerson() {
	person := entity.Person{
		ID:          1,
		Name:        "a",
		Height:      150,
		Gender:      "male",
		WantedDates: cTypes.Uint64(2),
	}
	s.initPeople([]entity.Person{person})

	s.boys.EXPECT().FindByID(person.ID).Return(&person, true)
	s.boys.EXPECT().RemovePerson(person.ID).Return(nil)

	err := s.h.RemovePerson(context.Background(), person.ID)
	assert.Nil(s.T(), err)
}

func (s *personTestSuite) Test_QuerySinglePeople() {
	people := []entity.Person{
		{
			ID:          1,
			Name:        "a",
			Height:      151,
			Gender:      "male",
			WantedDates: cTypes.Uint64(2),
		},
		{
			ID:          2,
			Name:        "b",
			Height:      152,
			Gender:      "male",
			WantedDates: cTypes.Uint64(2),
		},
		{
			ID:          3,
			Name:        "c",
			Height:      153,
			Gender:      "male",
			WantedDates: cTypes.Uint64(2),
		},
		{
			ID:          4,
			Name:        "d",
			Height:      150,
			Gender:      "female",
			WantedDates: cTypes.Uint64(2),
		},
	}
	s.initPeople(people)

	targetPerson := people[3]
	s.boys.EXPECT().FindByID(targetPerson.ID).Return(nil, false)
	s.girls.EXPECT().FindByID(targetPerson.ID).Return(&targetPerson, true)
	s.boys.EXPECT().QueryByHeight(targetPerson.Height, math.MaxFloat64).Return(people[:3])

	gotPeople, err := s.h.QuerySinglePeople(context.Background(), targetPerson.ID, 2)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(gotPeople))
}

func (s *personTestSuite) Test_MatchSameGender() {
	people := []entity.Person{
		{
			ID:          1,
			Name:        "a",
			Height:      151,
			Gender:      "male",
			WantedDates: cTypes.Uint64(2),
		},
		{
			ID:          2,
			Name:        "b",
			Height:      152,
			Gender:      "male",
			WantedDates: cTypes.Uint64(2),
		},
	}
	s.initPeople(people)

	s.boys.EXPECT().FindByID(people[0].ID).Return(&people[0], true)
	s.boys.EXPECT().FindByID(people[1].ID).Return(&people[1], true)

	err := s.h.Match(context.Background(), 1, 2)
	assert.Equal(s.T(), ErrorMatchSameGender, err)
}

func (s *personTestSuite) Test_MatchHeightCheckFail() {
	people := []entity.Person{
		{
			ID:          1,
			Name:        "a",
			Height:      150,
			Gender:      "male",
			WantedDates: cTypes.Uint64(2),
		},
		{
			ID:          2,
			Name:        "b",
			Height:      152,
			Gender:      "female",
			WantedDates: cTypes.Uint64(2),
		},
	}
	s.initPeople(people)

	s.boys.EXPECT().FindByID(people[0].ID).Return(&people[0], true)
	s.boys.EXPECT().FindByID(people[1].ID).Return(nil, false)
	s.girls.EXPECT().FindByID(people[1].ID).Return(&people[1], true)

	err := s.h.Match(context.Background(), 1, 2)
	assert.Equal(s.T(), ErrorHeightCheckFailed, err)
}

func (s *personTestSuite) Test_MatchDecrementWantedDatesFail() {
	people := []entity.Person{
		{
			ID:          1,
			Name:        "a",
			Height:      170,
			Gender:      "male",
			WantedDates: cTypes.Uint64(0),
		},
		{
			ID:          2,
			Name:        "b",
			Height:      152,
			Gender:      "female",
			WantedDates: cTypes.Uint64(2),
		},
	}
	s.initPeople(people)

	s.boys.EXPECT().FindByID(people[0].ID).Return(&people[0], true)
	s.boys.EXPECT().FindByID(people[1].ID).Return(nil, false)
	s.girls.EXPECT().FindByID(people[1].ID).Return(&people[1], true)

	err := s.h.Match(context.Background(), 1, 2)
	assert.Equal(s.T(), ErrorWantedDateLimit, err)
}

func (s *personTestSuite) Test_MatchSuccess() {
	people := []entity.Person{
		{
			ID:          1,
			Name:        "a",
			Height:      170,
			Gender:      "male",
			WantedDates: cTypes.Uint64(2),
		},
		{
			ID:          2,
			Name:        "b",
			Height:      152,
			Gender:      "female",
			WantedDates: cTypes.Uint64(2),
		},
	}
	s.initPeople(people)

	s.boys.EXPECT().FindByID(people[0].ID).Return(&people[0], true)
	s.boys.EXPECT().FindByID(people[1].ID).Return(nil, false)
	s.girls.EXPECT().FindByID(people[1].ID).Return(&people[1], true)

	err := s.h.Match(context.Background(), 1, 2)
	assert.Nil(s.T(), err)
}

func (s *personTestSuite) Test_MatchWithRemovePerson() {
	people := []entity.Person{
		{
			ID:          1,
			Name:        "a",
			Height:      170,
			Gender:      "male",
			WantedDates: cTypes.Uint64(1),
		},
		{
			ID:          2,
			Name:        "b",
			Height:      152,
			Gender:      "female",
			WantedDates: cTypes.Uint64(2),
		},
	}
	s.initPeople(people)

	s.boys.EXPECT().FindByID(people[0].ID).Return(&people[0], true)
	s.boys.EXPECT().FindByID(people[1].ID).Return(nil, false)
	s.girls.EXPECT().FindByID(people[1].ID).Return(&people[1], true)

	s.boys.EXPECT().RemovePerson(people[0].ID).Return(nil)

	err := s.h.Match(context.Background(), 1, 2)
	assert.Nil(s.T(), err)
}

func (s *personTestSuite) Test_Concurrent_DecrementWantedDate() {
	person := entity.Person{
		ID:          1,
		Name:        "a",
		Height:      170,
		Gender:      "male",
		WantedDates: cTypes.Uint64(5),
	}

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.h.decrementWantedDate(&person)
		}()
	}

	wg.Wait()

	actualWantedDates := atomic.LoadUint64(person.WantedDates)
	assert.Equal(s.T(), uint64(3), actualWantedDates)
}
