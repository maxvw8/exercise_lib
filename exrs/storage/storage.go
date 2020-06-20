package storage

//ExerciseStorage defines crud for exercise
type ExerciseStorage interface {
	Create(*Exercise) (*Exercise, error)
	Delete(string) (bool, error)
	Read(string) (*Exercise, error)
	Update(string, *Exercise) (*Exercise, error)
	List() ([]*Exercise, error)
}

//Exercise type stored on database
type Exercise struct {
	Id           string   `bson:"_id,omitempty"`
	Name         string   `bson:"name,omitempty"`
	Kind         string   `bson:"kind,omitempty"`
	Categories   []string `bson:"category,omitempty"`
	Muscles      []string `bson:"muscles,omitempty"`
	MuscleGroups []string `bson:"muscle_groups,omitempty"`
	Images       []string `bson:"images,omitempty"`
	Videos       []string `bson:"videos,omitempty"`
}
