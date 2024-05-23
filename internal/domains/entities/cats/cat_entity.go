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
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Race        Race     `json:"race"`
	Sex         Sex      `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	Description string   `json:"description"`
	ImageUrls   []string `json:"imageUrls"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

type AddCatRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=30"`
	Race        Race     `json:"race" validate:"required,oneof='Persian' 'Maine Coon' 'Siamese' 'Ragdoll' 'Bengal' 'Sphynx' 'British Shorthair' 'Abyssinian' 'Scottish Fold' 'Birman'"`
	Sex         Sex      `json:"sex" validate:"required,oneof='male' 'female'"`
	AgeInMonth  int      `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string   `json:"description" validate:"required,min=1,max=200"`
	ImageUrls   []string `json:"imageUrls" validate:"required,dive,required,imageurl"`
}

type AddCatResponse struct {
	ID        string `json:"name"`
	CreatedAt string `json:"createdAt"`
}
