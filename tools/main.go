package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/pkg/errors"

	"github.com/daodao97/egin/lib"

	"github.com/daodao97/egin-tools/asset"

	"github.com/daodao97/egin-tools/gen"
	"github.com/daodao97/egin-tools/parser"
	"github.com/daodao97/egin-tools/swagger"
)

var genDoc = flag.Bool("swagger", false, "是否生成 swagger.json 文件, 默认否")
var startUi = flag.Bool("ui", false, "是否开启 swagger ui 的 http 服务")
var uiPort = flag.String("ui-port", "8000", "swagger ui的监听端口")
var genRoute = flag.Bool("route", false, "是否生成文件对应的路由文件")
var genModelMode = flag.Bool("model", false, "根据mysql表结构生成数据模型")
var connection = flag.String("connection", "default", "数据库丽连接名")
var database = flag.String("database", "", "数据库名")
var genCtrl = flag.Bool("controller", false, "创建控制器")
var table = flag.String("table", "", "表名")
var apidoc interface{}

// go:generate go-bindata-assetfs -o=asset/asset.go -pkg=asset ui/...
func main() {
	flag.Parse()

	if *genDoc {
		genSwagger()
		if *startUi {
			ui()
		}
	}

	if *genRoute {
		genRouter()
	}

	if *genModelMode {
		genModel()
	}

	if *genCtrl {
		genController()
	}
}

func ui() {

	files := assetfs.AssetFS{
		Asset:    asset.Asset,
		AssetDir: asset.AssetDir,
		Prefix:   "ui", // 访问文件1.html = > 访问文件 www/1.html
	}

	http.Handle("/", http.FileServer(&files))
	http.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		js, err := json.Marshal(apidoc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	})
	fmt.Println("SwaggerUI已启动, 使用 http://localhost:" + *uiPort + " 打开")
	err := http.ListenAndServe(":"+*uiPort, nil)
	onErr(err)
}

func genSwagger() {

	openApi := swagger.NewSwagger()

	AllPath := make(swagger.Paths)
	var AllTags []swagger.Tags

	lib.RecursiveDir("controller", func(filePath string) {
		structInfo, err := parser.FileStructInfo(filePath)
		onErr(err)
		paths, tags := swagger.Filter(structInfo)
		for p, v := range paths {
			if _, ok := AllPath[p]; !ok {
				AllPath[p] = make(map[swagger.Method]swagger.Api)
			}
			for m, api := range v {
				AllPath[p][m] = api
			}
		}
		AllTags = append(AllTags, tags...)
	})

	openApi.Paths = AllPath
	openApi.Tags = AllTags
	apidoc = openApi
}

func genRouter() {
	lib.RecursiveDir("controller", func(filePath string) {
		structInfo, err := parser.FileStructInfo(filePath)
		onErr(err)
		varsInfo, err := parser.FileVarInfo(filePath)
		onErr(err)
		fmt.Println(filePath)
		gen.MakeRouteFile(structInfo, varsInfo)
	})
	gen.MakeRouteExport()
}

func genModel() {
	if *database != "" && *table != "" {
		tableInfo := gen.GetTableInfo(*connection, *database, *table)
		gen.MakeModel(*connection, *database, tableInfo)
	}
	if *database != "" && *table == "" {
		tables := gen.GetDbAllTable(*connection, *database)
		for _, t := range tables {
			gen.MakeModel(*connection, *database, t)
		}
	}
}

func genController() {
	if *table == "" {
		onErr(errors.New("need -table ***"))
	}
	gen.MakeController(*table, "")
}

func onErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
