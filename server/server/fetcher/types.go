package fetcher

type Fetcher struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	ImageName   string `bson:"image_name" json:"image_name"`
	ID          string `bson:"id" json:"id"`
}
