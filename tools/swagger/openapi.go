package swagger

import "github.com/daodao97/egin/lib"

type Api struct {
	Path        string      `json:"-"`
	Method      string      `json:"-" default:"GET"`
	Tags        []string    `json:"tags" default:"[]"`
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	OperationId string      `json:"operationId"`
	Produces    []string    `json:"produces" default:"[]"`
	Parameters  []Parameter `json:"parameters" default:"[]"`
	Response    []Response  `json:"response"`
}

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Type        string `json:"type"`
}

type Response struct {
	Code        int
	Description string
	Schema      Schema
}

type Schema struct {
	Ref string `json:"$ref"`
}

type Path string
type Method string

type Paths map[Path]map[Method]Api

type Tags struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Swagger struct {
	Swagger string `default:"2.0" json:"swagger"`
	Info    struct {
		Description    string `json:"description" default:"egin swagger"`
		Version        string `json:"version"`
		Title          string `json:"title"`
		TermsOfService string `json:"termsOfService"`
		Contact        struct {
			Email string `json:"email"`
		} `json:"contact"`
		License struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"license"`
	} `json:"info"`
	Host     string   `json:"host,omitempty"`
	BasePath string   `json:"basePath,omitempty"`
	Tags     []Tags   `json:"tags" default:"[]"`
	Schemes  []string `json:"schemes" default:"[]"`
	Paths    Paths    `json:"paths"`
}

func NewSwagger() *Swagger {
	s := &Swagger{}
	lib.MustSet(s)
	return s
}
