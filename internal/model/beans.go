package model

import (
	"fmt"
	"time"

	"github.com/patrickarmengol/somethingsomethingcoffee/internal/validator"
)

// passed from handler to service
// gets validated in service
type BeanCreateInput struct {
	Name       string         `form:"name"`
	RoastLevel RoastLevelEnum `form:"roast_level"`
	RoasterID  int64          `form:"roaster_id"`

	validator.Validator `form:"-"`
}

func (i *BeanCreateInput) Validate() {
	i.CheckField(validator.NotBlank(i.Name), "name", "this field cannot be blank")
	i.CheckField(validator.NotBlank(string(i.RoastLevel)), "roast_level", "this field cannot be blank")
	i.CheckField(validator.PermittedValue(i.RoastLevel, roastLevels...), "roast_level", fmt.Sprintf("this field must be one of %v", roastLevels))
	i.CheckField(i.RoasterID > 0, "roaster_id", "this field must be greater than 0")
}

func (i *BeanCreateInput) ToParams() *BeanCreateParams {
	return &BeanCreateParams{
		Name:       i.Name,
		RoastLevel: i.RoastLevel,
		RoasterID:  i.RoasterID,
	}
}

// passed from service to repository
type BeanCreateParams struct {
	Name       string
	RoastLevel RoastLevelEnum
	RoasterID  int64
}

// passed from handler to service
// gets validated in service
type BeanEditInput struct {
	ID         int64          `form:"-"`
	Name       string         `form:"name"`
	RoastLevel RoastLevelEnum `form:"roast_level"`
	RoasterID  int64          `form:"roaster_id"`

	validator.Validator `form:"-"`
}

func (i *BeanEditInput) Validate() {
	i.CheckField(i.ID > 0, "id", "this field musts be greater than 0")
	i.CheckField(validator.NotBlank(i.Name), "name", "this field cannot be blank")
	i.CheckField(validator.NotBlank(string(i.RoastLevel)), "roast_level", "this field cannot be blank")
	i.CheckField(validator.PermittedValue(i.RoastLevel, roastLevels...), "roast_level", fmt.Sprintf("this field must be one of %v", roastLevels))
	i.CheckField(i.RoasterID > 0, "roaster_id", "this field must be greater than 0")
}

func (i *BeanEditInput) ToParams() *BeanEditParams {
	return &BeanEditParams{
		ID:         i.ID,
		Name:       i.Name,
		RoastLevel: i.RoastLevel,
		RoasterID:  i.RoasterID,
	}
}

// passed from service to repository
type BeanEditParams struct {
	ID         int64
	Name       string
	RoastLevel RoastLevelEnum
	RoasterID  int64
}

// returned from repository to service
type BeanDB struct {
	ID         int64
	Name       string
	RoastLevel RoastLevelEnum
	RoasterID  int64
	CreatedAt  time.Time
	Version    int

	Roaster *RoasterDB
}

func (m *BeanDB) ToResponse() *BeanResponse {
	r := &BeanResponse{
		ID:         m.ID,
		Name:       m.Name,
		RoastLevel: m.RoastLevel,
		RoasterID:  m.RoasterID,
	}
	if m.Roaster != nil {
		r.Roaster = m.Roaster.ToResponse()
	}
	return r
}

// returned from service to handler
type BeanResponse struct {
	ID         int64
	Name       string
	RoastLevel RoastLevelEnum
	RoasterID  int64

	Roaster *RoasterResponse
}

func (r *BeanResponse) ToEditInput() *BeanEditInput {
	return &BeanEditInput{
		ID:         r.ID,
		Name:       r.Name,
		RoastLevel: r.RoastLevel,
		RoasterID:  r.RoasterID,
	}
}

type BeanFilterInput struct {
	Term string `form:"term"`
	Sort string `form:"sort"`

	// PageNum  int
	// PageSize int

	validator.Validator
}

func (i *BeanFilterInput) Validate() {
	i.CheckField(validator.MaxChars(i.Term, 50), "term", "this field must be at most 50 characters")
	i.CheckField(validator.NotBlank(i.Sort), "sort", "this field must not be empty")
	i.CheckField(validator.PermittedValue(i.Sort, beanSortBys...), "sort", fmt.Sprintf("this field must be in one of %v", beanSortBys))
}

func (i *BeanFilterInput) ToParams() *BeanFilterParams {
	p := &BeanFilterParams{
		SearchTerm: i.Term,
	}
	// TODO: maybe use a map instead since sorts used by multiple filters
	switch i.Sort {
	case SortByIDAsc:
		p.SortField = "id"
		p.SortDir = "asc"
	case SortByIDDesc:
		p.SortField = "id"
		p.SortDir = "desc"
	case SortByNameAsc:
		p.SortField = "name"
		p.SortDir = "asc"
	case SortByNameDesc:
		p.SortField = "name"
		p.SortDir = "desc"
	default:
		// should never happen since input must be validated
		p.SortField = "id"
		p.SortDir = "asc"
	}
	return p
}

type BeanFilterParams struct {
	SearchTerm string
	SortField  string
	SortDir    string
}

// value models

type RoastLevelEnum string

const (
	RLLight       RoastLevelEnum = "light"
	RLMediumLight RoastLevelEnum = "medium-light"
	RLMedium      RoastLevelEnum = "medium"
	RLMediumDark  RoastLevelEnum = "medium-dark"
	RLDark        RoastLevelEnum = "dark"
)

var roastLevels = []RoastLevelEnum{
	RLLight,
	RLMediumLight,
	RLMedium,
	RLMediumDark,
	RLDark,
}

// TODO: move this to filter.go
const (
	SortByIDAsc    string = "id_asc"
	SortByIDDesc   string = "id_desc"
	SortByNameAsc  string = "name_asc"
	SortByNameDesc string = "name_desc"
)

var beanSortBys = []string{
	SortByIDAsc,
	SortByIDDesc,
	SortByNameAsc,
	SortByNameDesc,
}
