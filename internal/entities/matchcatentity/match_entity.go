package matchcatentity

type Status string

const (
	Pending  Status = "pending"
	Approved Status = "approved"
	Rejected Status = "rejected"
)

type MatchCat struct {
	Id         string
	MatchCatId string
	UserCatId  string
	Message    string
	Status     Status
	IssuedBy   string
	IsDeleted  bool
	CreatedAt  string
	UpdatedAt  string
}

type CreateMatchCatRequest struct {
	MatchCatId string `json:"matchCatId" validate:"required,len=26"`
	UserCatId  string `json:"userCatId" validate:"required,len=26"`
	Message    string `json:"message" validate:"required,min=5,max=120"`
}
