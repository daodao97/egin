package swagger

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/daodao97/egin/lib"

	"github.com/daodao97/egin-tools/parser"
)

var (
	matchController = regexp.MustCompile(`@Controller .*`)
	matchAnyApi     = regexp.MustCompile(`@AnyApi .*`)
	matchGetApi     = regexp.MustCompile(`@GetApi .*`)
	matchPostApi    = regexp.MustCompile(`@PostApi .*`)
	matchPutApi     = regexp.MustCompile(`@PutApi .*`)
	matchDeleteApi  = regexp.MustCompile(`@DeleteApi .*`)
	matchDesc       = regexp.MustCompile(`@Desc .*`)
	matchTag        = regexp.MustCompile(`@Tag .*`)
	matchSummary    = regexp.MustCompile(`@Summary .*`)
	matchParams     = regexp.MustCompile(`@Params .*`)
)

type Controller struct {
	Tag  string
	Desc string
}

func transController(info parser.StructInfo) (c Controller) {
	c.Tag = info.Name
	for _, v := range info.Doc {
		if matchController.MatchString(v) {
			token := explode(v)
			if token[1] != "" {
				c.Tag = token[1]
			}
			if len(token) > 2 && token[2] != "" {
				c.Desc = token[2]
			}

			continue
		}
	}
	return c
}

func transApi(sf parser.StructFunc, info []parser.StructInfo) (api Api, err error) {
	if sf.Doc == nil {
		return api, errors.New("func doc not found")
	}
	_ = lib.Set(&api)
	for _, v := range sf.Doc {
		if matchAnyApi.MatchString(v) {
			api.Method = "ANY"
			api.Path = explode(v)[1]
			continue
		}
		if matchGetApi.MatchString(v) {
			api.Method = "GET"
			api.Path = explode(v)[1]
			continue
		}
		if matchPostApi.MatchString(v) {
			api.Method = "POST"
			api.Path = explode(v)[1]
			continue
		}
		if matchPutApi.MatchString(v) {
			api.Method = "PUT"
			api.Path = explode(v)[1]
			continue
		}
		if matchDeleteApi.MatchString(v) {
			api.Method = "DELETE"
			api.Path = explode(v)[1]
			continue
		}
		if matchDesc.MatchString(v) {
			api.Description = strings.TrimPrefix(v, "@Desc")
			continue
		}
		if matchTag.MatchString(v) {
			api.Tags = explode(v)[1:]
			continue
		}
		if matchSummary.MatchString(v) {
			api.Summary = strings.TrimPrefix(v, "@Summary")
			continue
		}
		if matchParams.MatchString(v) {
			p := explode(v)[1]
			si, err := filterByName(info, p)
			if err != nil {
				fmt.Println("not found " + p)
				continue
			}
			api.Parameters = transParams(si.Fields)
			continue
		}
	}
	return api, nil
}

func explode(str string) []string {
	return strings.Split(str, " ")
}

func filterByName(info []parser.StructInfo, name string) (result parser.StructInfo, err error) {
	for _, v := range info {
		if v.Name == name {
			return v, nil
		}
	}
	return result, errors.New("not found")
}

func transParams(fields []parser.StructField) (ps []Parameter) {
	for _, v := range fields {
		ps = append(ps, transParam(v))
	}
	return ps
}

func transParam(field parser.StructField) (param Parameter) {
	param.Type = field.Type
	if name, ok := field.Tags["json"]; ok {
		param.Name = name
	}
	if label, ok := field.Tags["label"]; ok {
		param.Description = label
	}
	if where, ok := field.Tags["in"]; ok {
		param.In = where
	}
	if binding, ok := field.Tags["binding"]; ok {
		rules := strings.Split(binding, ",")
		if _, ok := lib.Find(rules, "required"); ok {
			param.Required = true
		}
	}
	return param
}

func Filter(info []parser.StructInfo) (paths Paths, tags []Tags) {
	paths = make(Paths)
	for _, v := range info {
		c := transController(v)
		if v.Funcs != nil {
			for _, f := range v.Funcs {
				api, err := transApi(f, info)
				if err == nil {
					api.Tags = []string{c.Tag}
					path := Path(api.Path)
					method := Method(strings.ToLower(api.Method))
					if _, ok := paths[path]; !ok {
						paths[path] = make(map[Method]Api)
					}
					paths[path][method] = api
				}
			}
			tags = append(tags, Tags{Name: c.Tag, Description: c.Desc})
		}
	}
	return paths, tags
}
