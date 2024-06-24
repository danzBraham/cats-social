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
