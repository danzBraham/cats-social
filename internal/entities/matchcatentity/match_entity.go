package matchcatentity

import "github.com/danzBraham/cats-social/internal/entities/catentity"

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

type IssuerDetail struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type GetMatchCatResponse struct {
	Id             string                   `json:"id"`
	IssuedBy       IssuerDetail             `json:"issuedBy"`
	MatchCatDetail catentity.GetCatResponse `json:"matchCatDetail"`
	UserCatDetail  catentity.GetCatResponse `json:"userCatDetail"`
	Message        string                   `json:"message"`
	CreatedAt      string                   `json:"createdAt"`
}

type ApproveMatchCatRequest struct {
	MatchId string `json:"matchId" validate:"required,len=26"`
}
