package tree

import (
	"math"
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"

	"github.com/ars0915/matching-system/entity"
)

func TestPersonTree_QueryByHeight(t *testing.T) {
	idMap := map[uint64]*entity.Person{
		1: {
			ID:     1,
			Name:   "1",
			Height: 150,
		},
		2: {
			ID:     2,
			Name:   "2",
			Height: 155,
		},
		3: {
			ID:     3,
			Name:   "3",
			Height: 155,
		},
		4: {
			ID:     4,
			Name:   "4",
			Height: 160,
		},
		5: {
			ID:     5,
			Name:   "5",
			Height: 170,
		},
	}

	pt := &PersonTree{
		tree:  redblacktree.NewWith(utils.Float64Comparator),
		idMap: idMap,
		mu:    sync.RWMutex{},
	}
	insertPersonToTree(pt, idMap)

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
		t.Run(tt.name, func(t *testing.T) {

			got := pt.QueryByHeight(tt.args.minHeight, tt.args.maxHeight)
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
