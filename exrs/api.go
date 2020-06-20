package exrs

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/maxvw8/exercise_lib/exrs/storage"
	pbexrs "github.com/maxvw8/exercise_lib/pbexrs/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

//TODO: Unit test this stuff!!!
//TODO: Manage grpc error codes https://grpc.io/docs/guides/error/

//API asd
type API struct {
	storage.ExerciseStorage
}

//Server creates a new instance of Exercise API
func Server(repo storage.ExerciseStorage) (*API, error) {
	s := &API{repo}
	return s, nil
}

//CreateExercise creates an exercise
func (s *API) CreateExercise(ctx context.Context, req *pbexrs.CreateExerciseRequest) (*pbexrs.Exercise, error) {
	log := ctxzap.Extract(ctx).Sugar()
	e := MarshallExercise(req.Exercise)
	log.Debugf("creating exercise %v", e)
	r, err := s.ExerciseStorage.Create(e)
	if err != nil {
		log.Warnf("failed creating new exercise %v. Error was %v", e, err)
		return &pbexrs.Exercise{}, err
	}
	log.Debugf("created exercise %v and error %v", e, err)
	return UnmarshallExercise(r), err
}

//GetExercise reads an exercise from the repository matching the id
func (s *API) GetExercise(ctx context.Context, req *pbexrs.GetExerciseRequest) (*pbexrs.Exercise, error) {
	log := ctxzap.Extract(ctx).Sugar()
	log.Debugf("getting exercise id %v", req.Id)
	r, err := s.ExerciseStorage.Read(req.Id)
	if err != nil {
		log.Warnf("could not find exercise with id %v. Error was %v", req.Id, err)
		return &pbexrs.Exercise{}, err
	}
	log.Debugf("found exercise %v", r)
	return UnmarshallExercise(r), err
}

//UpdateExercise updates an existing record, changing everything but the exercise id
func (s *API) UpdateExercise(ctx context.Context, req *pbexrs.UpdateRequest) (*pbexrs.Exercise, error) {
	log := ctxzap.Extract(ctx).Sugar()
	e := MarshallExercise(req.Exercise)
	log.Debugf("updating exercise with id %v", req.GetId())
	r, err := s.ExerciseStorage.Update(req.GetId(), e)
	if err != nil {
		log.Warnf("could not update exercise %v. Error was %v", req.GetId(), err)
		return &pbexrs.Exercise{}, err
	}
	log.Debugf("updated exercise %v", req.GetId())
	return UnmarshallExercise(r), err
}

//DeleteExercise deletes an exercise by Id, if doesnt exists, returns an error
func (s *API) DeleteExercise(ctx context.Context, req *pbexrs.DeleteRequest) (*empty.Empty, error) {
	log := ctxzap.Extract(ctx).Sugar()
	log.Debugf("deleting exercise with id %v", req.GetId())
	_, err := s.ExerciseStorage.Delete(req.GetId())
	if err != nil {
		log.Warnf("failed to delete exercise with id %v", req.GetId())
		return &emptypb.Empty{}, err
	}
	return &empty.Empty{}, err
}

//ListExercises returns a paged list of exercises
func (s *API) ListExercises(ctx context.Context, req *pbexrs.ListExercisesRequest) (*pbexrs.ListExercisesResponse, error) {
	log := ctxzap.Extract(ctx).Sugar()
	log.Debugf("[Request] Listing exercises")
	l, err := s.ExerciseStorage.List()
	if err != nil {
		log.Warnf("failed to get list of exercises")
		return &pbexrs.ListExercisesResponse{}, err
	}
	ul := UnmarshallExerciseList(l)
	if ul == nil { //in case returned list is nil, this funciton never returns nil
		ul = []*pbexrs.Exercise{}
	}
	log.Debugf("[Response] list %v with error %v", l, err)
	return &pbexrs.ListExercisesResponse{Exercises: ul}, err
}

//UnmarshallExerciseList converts a list of storage Exercises into transport layer exercises
func UnmarshallExerciseList(l []*storage.Exercise) []*pbexrs.Exercise {
	if l == nil {
		return nil
	}
	vsm := make([]*pbexrs.Exercise, len(l))
	for i, v := range l {
		vsm[i] = UnmarshallExercise(v)
	}
	return vsm
}

//MarshallExercise converts a storage layer exercise struct transport layer exercise into an
func MarshallExercise(e *pbexrs.Exercise) *storage.Exercise {
	if e == nil {
		return nil
	}
	return &storage.Exercise{Id: e.Id,
		Name:         e.Name,
		Kind:         e.Kind,
		Categories:   e.Categories,
		Muscles:      e.Muscles,
		MuscleGroups: e.MuscleGroups,
		Images:       e.Images,
		Videos:       e.Videos,
	}
}

//UnmarshallExercise converts a transport layer exercise into an storage layer exercise struct
func UnmarshallExercise(e *storage.Exercise) *pbexrs.Exercise {
	if e == nil {
		return nil
	}
	return &pbexrs.Exercise{Id: e.Id,
		Name:         e.Name,
		Kind:         e.Kind,
		Categories:   e.Categories,
		Muscles:      e.Muscles,
		MuscleGroups: e.MuscleGroups,
		Images:       e.Images,
		Videos:       e.Videos,
	}
}
