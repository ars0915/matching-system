package usecase

import (
	"context"
	"math"
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
	s.boys.EXPECT().AddPerson(ctest.DiffWrapper(&person)).Return(nil)
	person, err := s.h.AddPerson(person)
	assert.Nil(s.T(), err)

	s.boys.EXPECT().FindByID(person.ID).Return(&person, true)
	s.boys.EXPECT().RemovePerson(person.ID).Return(nil)

	err = s.h.RemovePerson(context.Background(), person.ID)
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

	s.boys.EXPECT().FindByID(uint64(1)).Return(&people[0], true)
	s.boys.EXPECT().FindByID(uint64(2)).Return(&people[1], true)

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

	s.boys.EXPECT().FindByID(uint64(1)).Return(&people[0], true)
	s.boys.EXPECT().FindByID(uint64(2)).Return(nil, false)
	s.girls.EXPECT().FindByID(uint64(2)).Return(&people[1], true)

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

	s.boys.EXPECT().FindByID(uint64(1)).Return(&people[0], true)
	s.boys.EXPECT().FindByID(uint64(2)).Return(nil, false)
	s.girls.EXPECT().FindByID(uint64(2)).Return(&people[1], true)

	err := s.h.Match(context.Background(), 1, 2)
	assert.Equal(s.T(), ErrorWantedDateLimit, err)
}
