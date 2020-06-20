// +build unit

package exrs

import (
	"testing"

	"github.com/maxvw8/exercise_lib/exrs/storage"
	pbexrs "github.com/maxvw8/exercise_lib/pbexrs/v1"
	"github.com/stretchr/testify/assert"
)

func TestMarshall(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    *pbexrs.Exercise
		Expected *storage.Exercise
	}{
		{
			Name:     "nil",
			Input:    nil,
			Expected: nil,
		},
		{
			Name: "full instance",
			Input: &pbexrs.Exercise{
				Id:           "id",
				Name:         "name",
				Category:     []string{"category"},
				Kind:         "kind",
				Images:       []string{"images"},
				Videos:       []string{"yutub"},
				Muscles:      []string{"muscles"},
				MuscleGroups: []string{"muscle groups"},
			},
			Expected: &storage.Exercise{
				Id:           "id",
				Name:         "name",
				Category:     []string{"category"},
				Kind:         "kind",
				Images:       []string{"images"},
				Videos:       []string{"yutub"},
				Muscles:      []string{"muscles"},
				MuscleGroups: []string{"muscle groups"},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc //capturing test case issue with parallel execution
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			assert.EqualValues(t, tc.Expected, MarshallExercise(tc.Input))
		})
	}
}

func TestUnarshall(t *testing.T) {
	testCases := []struct {
		Name     string
		Expected *pbexrs.Exercise
		Input    *storage.Exercise
	}{
		{
			Name:     "nil",
			Input:    nil,
			Expected: nil,
		},
		{
			Name: "full instance",
			Input: &storage.Exercise{
				Id:           "id",
				Name:         "name",
				Category:     []string{"category"},
				Kind:         "kind",
				Images:       []string{"images"},
				Videos:       []string{"yutub"},
				Muscles:      []string{"muscles"},
				MuscleGroups: []string{"muscle groups"},
			},
			Expected: &pbexrs.Exercise{
				Id:           "id",
				Name:         "name",
				Category:     []string{"category"},
				Kind:         "kind",
				Images:       []string{"images"},
				Videos:       []string{"yutub"},
				Muscles:      []string{"muscles"},
				MuscleGroups: []string{"muscle groups"},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc //capturing test case issue with parallel execution
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			assert.EqualValues(t, tc.Expected, UnmarshallExercise(tc.Input))
		})
	}
}

func TestUnmarshallList(t *testing.T) {
	storageSample := &storage.Exercise{
		Id:           "id",
		Name:         "name",
		Category:     []string{"category"},
		Kind:         "kind",
		Images:       []string{"images"},
		Videos:       []string{"yutub"},
		Muscles:      []string{"muscles"},
		MuscleGroups: []string{"muscle groups"},
	}
	pbexrsSample := &pbexrs.Exercise{
		Id:           "id",
		Name:         "name",
		Category:     []string{"category"},
		Kind:         "kind",
		Images:       []string{"images"},
		Videos:       []string{"yutub"},
		Muscles:      []string{"muscles"},
		MuscleGroups: []string{"muscle groups"},
	}
	testCases := []struct {
		Name     string
		Input    []*storage.Exercise
		Expected []*pbexrs.Exercise
	}{
		{
			Name:     "nil",
			Input:    nil,
			Expected: nil, //never return
		},
		{
			Name:     "empty list",
			Input:    []*storage.Exercise{},
			Expected: []*pbexrs.Exercise{},
		},
		{
			Name:     "one element",
			Input:    []*storage.Exercise{storageSample},
			Expected: []*pbexrs.Exercise{pbexrsSample},
		},
		{
			Name:     "many element",
			Input:    []*storage.Exercise{storageSample, storageSample, storageSample},
			Expected: []*pbexrs.Exercise{pbexrsSample, pbexrsSample, pbexrsSample},
		},
	}

	for _, tc := range testCases {
		tc := tc //capturing test case issue with parallel execution
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			assert.EqualValues(t, tc.Expected, UnmarshallExerciseList(tc.Input))
		})
	}
}
