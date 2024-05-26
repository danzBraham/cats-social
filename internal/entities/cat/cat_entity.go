package cat_entity

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
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Race        Race     `json:"race"`
	Sex         Sex      `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	Description string   `json:"description"`
	ImageUrls   []string `json:"imageUrls"`
	HasMatched  bool     `json:"has_matched"`
	OwnerId     string   `json:"ownerId"`
}

type AddCatRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=30"`
	Race        Race     `json:"race" validate:"required,oneof='Persian' 'Maine Coon' 'Siamese' 'Ragdoll' 'Bengal' 'Sphynx' 'British Shorthair' 'Abyssinian' 'Scottish Fold' 'Birman'"`
	Sex         Sex      `json:"sex" validate:"required,oneof='male' 'female'"`
	AgeInMonth  int      `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string   `json:"description" validate:"required,min=1,max=200"`
	ImageUrls   []string `json:"imageUrls" validate:"required,min=1,dive,required,imageurl"`
	OwnerId     string   `validate:"required,len=26"`
}

type AddCatResponse struct {
	Id        string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type CatQueryParams struct {
	Id         string
	Limit      string
	Offset     string
	Race       string
	Sex        string
	HasMatched string
	AgeInMonth string
	Owned      string
	Search     string
}

type GetCatReponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Race        Race     `json:"race"`
	Sex         Sex      `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	ImageUrls   []string `json:"imageUrls"`
	Description string   `json:"description"`
	HasMatched  bool     `json:"hasMatched"`
	CreatedAt   string   `json:"createdAt"`
}

type UpdateCatRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=30"`
	Race        Race     `json:"race" validate:"required,oneof='Persian' 'Maine Coon' 'Siamese' 'Ragdoll' 'Bengal' 'Sphynx' 'British Shorthair' 'Abyssinian' 'Scottish Fold' 'Birman'"`
	Sex         Sex      `json:"sex" validate:"required,oneof='male' 'female'"`
	AgeInMonth  int      `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string   `json:"description" validate:"required,min=1,max=200"`
	ImageUrls   []string `json:"imageUrls" validate:"required,min=1,dive,required,imageurl"`
}
