package tree

import (
	"math"
	"reflect"
	"slices"
	"sort"
	"sync"
	"testing"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/ars0915/matching-system/entity"
	"github.com/ars0915/matching-system/util/cTypes"
)

type personTreeTestSuite struct {
	suite.Suite

	pt *PersonTree
}

func Test_personTreeTestSuite(t *testing.T) {
	suite.Run(t, &personTreeTestSuite{})
}

func (s *personTreeTestSuite) SetupTest() {
	idMap := map[uint64]*entity.Person{
		1: {
			ID:          1,
			Name:        "1",
			Height:      150,
			WantedDates: cTypes.Uint64(1),
		},
		2: {
			ID:          2,
			Name:        "2",
			Height:      155,
			WantedDates: cTypes.Uint64(2),
		},
		3: {
			ID:          3,
			Name:        "3",
			Height:      155,
			WantedDates: cTypes.Uint64(1),
		},
		4: {
			ID:          4,
			Name:        "4",
			Height:      160,
			WantedDates: cTypes.Uint64(1),
		},
		5: {
			ID:          5,
			Name:        "5",
			Height:      170,
			WantedDates: cTypes.Uint64(1),
		},
	}

	s.pt = &PersonTree{
		tree:  redblacktree.NewWith(utils.Float64Comparator),
		idMap: idMap,
		mu:    sync.RWMutex{},
	}
	insertPersonToTree(s.pt, idMap)
}

func insertPersonToTree(pt *PersonTree, idMap map[uint64]*entity.Person) {
	for _, p := range idMap {
		if value, found := pt.tree.Get(p.Height); found {
			people := value.([]uint64)
			pt.tree.Put(p.Height, append(people, p.ID))
			continue
		}
		pt.tree.Put(p.Height, []uint64{p.ID})
	}
}

func (s *personTreeTestSuite) Test_QueryByHeight() {
	type args struct {
		minHeight float64
		maxHeight float64
	}
	tests := []struct {
		name string
		args args
		want []uint64
	}{
		{
			name: "All people higher than minHeight",
			args: args{
				minHeight: 100,
				maxHeight: math.MaxFloat64,
			},
			want: []uint64{1, 2, 3, 4, 5},
		},
		{
			name: "No people higher than minHeight",
			args: args{
				minHeight: 200,
				maxHeight: math.MaxFloat64,
			},
			want: nil,
		},
		{
			name: "All people lower than maxHeight",
			args: args{
				minHeight: 0,
				maxHeight: 190,
			},
			want: []uint64{1, 2, 3, 4, 5},
		},
		{
			name: "All people higher than maxHeight",
			args: args{
				minHeight: 0,
				maxHeight: 100,
			},
			want: nil,
		},
		{
			name: "Success",
			args: args{
				minHeight: 151,
				maxHeight: 161,
			},
			want: []uint64{2, 3, 4},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			got := s.pt.QueryByHeight(tt.args.minHeight, tt.args.maxHeight)
			var gotIDs []uint64
			for _, person := range got {
				gotIDs = append(gotIDs, person.ID)
			}
			sort.Slice(gotIDs, func(i, j int) bool { return gotIDs[i] < gotIDs[j] })

			if !reflect.DeepEqual(gotIDs, tt.want) {
				t.Errorf("QueryByHeight() = %v, want %v", gotIDs, tt.want)
			}
		})
	}
}

func (s *personTreeTestSuite) Test_AddPerson() {
	person := entity.Person{
		ID:          99999,
		Name:        "New Person",
		Height:      999,
		Gender:      "male",
		WantedDates: cTypes.Uint64(3),
	}
	tests := []struct {
		name    string
		p       entity.Person
		wantErr error
	}{
		{
			"Success",
			person,
			nil,
		},
		{
			"Exist",
			person,
			ErrorPersonExist,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.pt.AddPerson(&tt.p)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			s.checkPersonExist(tt.p)
		})
	}
}

func (s *personTreeTestSuite) checkPersonExist(p entity.Person) {
	s.pt.mu.RLock()
	defer s.pt.mu.RUnlock()

	// check in idMap
	gotPerson, exist := s.pt.idMap[p.ID]
	assert.True(s.T(), exist, "person should be found in idMap")
	assert.Truef(s.T(), reflect.DeepEqual(*gotPerson, p), "person different, got = %v, want = %v", *gotPerson, p)

	// check in tree
	value, exist := s.pt.tree.Get(p.Height)
	assert.True(s.T(), exist, "height should be found in tree")
	gotIDs := value.([]uint64)
	assert.True(s.T(), slices.Contains(gotIDs, p.ID), "person id should be found in node")
}
