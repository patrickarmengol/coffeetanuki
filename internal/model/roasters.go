package model

import (
	"fmt"
	"time"

	"github.com/patrickarmengol/somethingsomethingcoffee/internal/validator"
)

// passed from handler to service
type RoasterCreateInput struct {
	Name        string `form:"name"`
	Description string `form:"description"`
	Website     string `form:"website"`
	Location    string `form:"location"`

	validator.Validator `form:"-"`
}

func (i *RoasterCreateInput) Validate() {
	i.CheckField(validator.NotBlank(i.Name), "name", "this field cannot be blank")
	i.CheckField(validator.MaxChars(i.Name, 50), "name", "this field must have at most 50 characters")
	i.CheckField(validator.MaxChars(i.Description, 300), "description", "this field must have at most 300 characters")
	i.CheckField(validator.IsURL(i.Website), "website", "this field must be a valid URL")
	i.CheckField(validator.NotBlank(i.Location), "location", "this field cannot be blank")
}

func (i *RoasterCreateInput) ToParams() *RoasterCreateParams {
	return &RoasterCreateParams{
		Name:        i.Name,
		Description: i.Description,
		Website:     i.Website,
		Location:    i.Website,
	}
}

// passed from service to repository
type RoasterCreateParams struct {
	Name        string
	Description string
	Website     string
	Location    string
}

// passed from handler to service
type RoasterEditInput struct {
	ID          int64  `form:"-"` // parsed from URL param
	Name        string `form:"name"`
	Description string `form:"description"`
	Website     string `form:"website"`
	Location    string `form:"location"`

	validator.Validator `form:"-"`
}

func (i *RoasterEditInput) Validate() {
	i.CheckField(i.ID > 0, "id", "this field must be greater than 0")
	i.CheckField(validator.NotBlank(i.Name), "name", "this field cannot be blank")
	i.CheckField(validator.MaxChars(i.Name, 50), "name", "this field must have at most 50 characters")
	i.CheckField(validator.MaxChars(i.Description, 300), "description", "this field must have at most 300 characters")
	i.CheckField(validator.IsURL(i.Website), "website", "this field must be a valid URL")
	i.CheckField(validator.NotBlank(i.Location), "location", "this field cannot be blank")
}

func (i *RoasterEditInput) ToParams() *RoasterEditParams {
	return &RoasterEditParams{
		ID:          i.ID,
		Name:        i.Name,
		Description: i.Description,
		Website:     i.Website,
		Location:    i.Website,
	}
}

// passed from service to repository
type RoasterEditParams struct {
	ID          int64
	Name        string
	Description string
	Website     string
	Location    string
}

// returned from repository to service
type RoasterDB struct {
	ID          int64
	Name        string
	Description string
	Website     string
	Location    string
	CreatedAt   time.Time
	Version     int

	Beans []*BeanDB
}

func (m *RoasterDB) ToResponse() *RoasterResponse {
	r := &RoasterResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Website:     m.Website,
		Location:    m.Location,
	}
	if m.Beans != nil {
		beans := []*BeanResponse{}
		for _, b := range m.Beans {
			beans = append(beans, b.ToResponse())
		}
		r.Beans = beans
	}

	return r
}

// returned from service to handler
type RoasterResponse struct {
	ID          int64
	Name        string
	Description string
	Website     string
	Location    string

	Beans []*BeanResponse
}

func (r *RoasterResponse) ToEditInput() *RoasterEditInput {
	return &RoasterEditInput{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Website:     r.Website,
		Location:    r.Location,
	}
}

type RoasterFilterInput struct {
	Term string `form:"term"`
	Sort string `form:"sort"`

	// PageNum  int
	// PageSize int

	validator.Validator
}

func (i *RoasterFilterInput) Validate() {
	i.CheckField(validator.MaxChars(i.Term, 50), "term", "this field must be at most 50 characters")
	i.CheckField(validator.NotBlank(i.Sort), "sort", "this field must not be empty")
	i.CheckField(validator.PermittedValue(i.Sort, roasterSortBys...), "sort", fmt.Sprintf("this field must be in one of %v", roasterSortBys))
}

func (i *RoasterFilterInput) ToParams() *RoasterFilterParams {
	p := &RoasterFilterParams{
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

type RoasterFilterParams struct {
	SearchTerm string
	SortField  string
	SortDir    string
}

var roasterSortBys = []string{
	SortByIDAsc,
	SortByIDDesc,
	SortByNameAsc,
	SortByNameDesc,
}
