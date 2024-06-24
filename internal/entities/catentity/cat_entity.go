package catentity

type Race string

const (
	Persian          Race = "Persian"
	MaineCoon        Race = "Maine Coon"
	Siamese          Race = "Siamese"
	Ragdoll          Race = "Ragdoll"
	Bengal           Race = "Bengal"
	Sphynx           Race = "Sphynx"
	BritishShorthair Race = "British Shorthair"
	Abyssinian       Race = "Abyssinian"
	ScottishFold     Race = "Scottish Fold"
	Birman           Race = "Birman"
)

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
)

type Cat struct {
	Id          string
	Name        string
	Race        Race
	Sex         Sex
	AgeInMonth  int
	Description string
	ImageUrls   []string
	HasMatched  bool
	OwnerId     string
	CreatedAt   string
	UpdatedAt   string
}

type CreateCatRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=30"`
	Race        Race     `json:"race" validate:"required,oneof='Persian' 'Maine Coon' 'Siamese' 'Ragdoll' 'Bengal' 'Sphynx' 'British Shorthair' 'Abyssinian' 'Scottish Fold' 'Birman'"`
	Sex         Sex      `json:"sex" validate:"required,oneof='male' 'female'"`
	AgeInMonth  int      `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string   `json:"description" validate:"required,min=1,max=200"`
	ImageUrls   []string `json:"imageUrls" validate:"required,min=1,dive,required,http_url"`
}

type CreateCatResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type CatQueryParams struct {
	Id         string
	Limit      int
	Offset     int
	Race       Race
	Sex        Sex
	HasMatched bool
	AgeInMonth string
	Owned      bool
	Search     string
}

type GetCatResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Race        Race     `json:"race"`
	Sex         Sex      `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	Description string   `json:"description"`
	ImageUrls   []string `json:"imageUrls"`
	HasMatched  bool     `json:"hasMatched"`
	CreatedAt   string   `json:"createdAt"`
}
