package main

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

// built following this guide https://dev.to/hlubek/metaprogramming-with-go-or-how-to-build-code-generators-that-parse-go-code-2k3j

func main() {
	if len(os.Args) != 3 {
		failErr(fmt.Errorf("expected exactly two arguments: <source type> and <table name>"))
	}
	sourceType := os.Args[1]
	tableName := os.Args[2]
	sourceTypePackage, sourceTypeName := splitSourceType(sourceType)

	pkg := loadPackage(sourceTypePackage)

	obj := pkg.Types.Scope().Lookup(sourceTypeName)
	if obj == nil {
		failErr(fmt.Errorf("%s not found in declared types of %s",
			sourceTypeName, pkg))
	}

	if _, ok := obj.(*types.TypeName); !ok {
		failErr(fmt.Errorf("%v is not a named type", obj))
	}

	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		failErr(fmt.Errorf("type %v is not a struct", obj))
	}

	f := jen.NewFile("generated")

	f.PackageComment("Code generated by generator, DO NOT EDIT.")

	err := os.MkdirAll("../generated", os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	goFile := os.Getenv("GOFILE")
	ext := filepath.Ext(goFile)
	baseFilename := goFile[0 : len(goFile)-len(ext)]

	targetFilename := "../generated/" + baseFilename + "." + strings.ToLower(sourceTypeName) + ".gen.go"
	generateCreate(f, tableName, sourceTypeName, structType, sourceTypePackage)
	generateUpdate(f, tableName, sourceTypeName, structType, sourceTypePackage)
	generateRead(f, tableName, sourceTypeName, structType, sourceTypePackage)
	generateDelete(f, tableName, sourceTypeName, structType)

	err = f.Save(targetFilename)
	if err != nil {
		failErr(err)
	}
}

func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		failErr(fmt.Errorf("loading packages for inspection: %v", err))
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0]
}

func splitSourceType(sourceType string) (string, string) {
	idx := strings.LastIndexByte(sourceType, '.')
	if idx == -1 {
		failErr(fmt.Errorf(`expected qualified type as "pkg/path.MyType"`))
	}
	sourceTypePackage := sourceType[0:idx]
	sourceTypeName := sourceType[idx+1:]
	return sourceTypePackage, sourceTypeName
}

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var structColPattern = regexp.MustCompile(`col:"([^"]+)"`)

func generateCreate(f *jen.File, tableName string, sourceTypeName string, structType *types.Struct, sourceTypePackage string) {
	//goPackage := os.Getenv("GOPACKAGE")
	cblck := []jen.Code{}
	cblck = append(cblck, jen.Id("args").Op(":=").Op("[]").Interface().Op("{"))
	sqlstr := fmt.Sprintf(`INSERT INTO %s (`, tableName)
	valuecols := ""
	valuesstr := ") VALUES ("
	questionmarks := ""
	end := ") ON CONFLICT DO NOTHING"
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		rawTagValue := structType.Tag(i)

		matches := structColPattern.FindStringSubmatch(rawTagValue)
		if matches == nil {
			continue
		}
		if len(valuecols) > 0 {
			// need a comma
			valuecols += ","
			questionmarks += ","
		}
		raw := strings.Split(matches[1], ",")
		col := raw[0]
		valuecols += col
		questionmarks += fmt.Sprintf("$%d", i+1)
		cblck = append(cblck, jen.Id("c").Dot(field.Name()).Op(","))
	}

	cblck = append(cblck, jen.Op("}"))
	cblck = append(cblck, jen.Id("sql").Op(":=").Op(fmt.Sprintf("\"%s%s%s%s%s\"", sqlstr, valuecols, valuesstr, questionmarks, end)))

	cblck = append(cblck, jen.Id("_").Op(",").Id("err").Op(":=").Id("pg").Dot("Exec").Params(
		jen.Id("ctx"), jen.Id("sql"), jen.Id("args").Op("..."),
	))
	cblck = append(cblck, jen.If(jen.Id("err").Op("!=").Nil()).Block(
		jen.Return(jen.Nil(), jen.Id("err")),
	))
	cblck = append(cblck, jen.Return(jen.Op("&").Id("c"), jen.Nil()))
	f.ImportName("github.com/jackc/pgx/v5/pgxpool", "pgxpool")

	f.Commentf("Insert new '%s' into table '%s' - duplicate ids is a no-op", sourceTypeName, tableName)
	f.Func().Id(sourceTypeName+"Create").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("pg").Op("*").Qual(
			"github.com/jackc/pgx/v5/pgxpool",
			"Pool",
		),
		jen.Id("c").Qual(sourceTypePackage, sourceTypeName),
	).Params(jen.Op("*").Qual(sourceTypePackage, sourceTypeName), jen.Error()).Block(cblck...)
}

func generateUpdate(f *jen.File, tableName string, sourceTypeName string, structType *types.Struct, sourceTypePackage string) {
	//goPackage := os.Getenv("GOPACKAGE")
	cblck := []jen.Code{}
	cblck = append(cblck, jen.Id("parameters").Op(":=").Op("[]").Interface().Op("{}"))
	sets := []string{}
	conditions := []string{}
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		rawTagValue := structType.Tag(i)

		matches := structColPattern.FindStringSubmatch(rawTagValue)
		if matches == nil {
			continue
		}
		raw := strings.Split(matches[1], ",")
		col := raw[0]
		isId := false
		if len(raw) > 1 && raw[1] == "primary" {
			isId = true
		}
		if isId {
			conditions = append(conditions, col+fmt.Sprintf("=$%d", i+1))
		} else {
			sets = append(sets, col+fmt.Sprintf("=$%d", i+1))
		}
		cblck = append(cblck, jen.Id("parameters").Op("=").Op("append(").Id("parameters").Op(",").Id("u").Dot(field.Name()).Op(")"))
	}

	cblck = append(cblck, jen.Id("sql").Op(":=").Op(fmt.Sprintf("\"%s%s%s\"",
		fmt.Sprintf(`UPDATE %s `, tableName),
		"SET "+strings.Join(sets, ","),
		" WHERE "+strings.Join(conditions, ","),
	)))

	cblck = append(cblck, jen.Id("_").Op(",").Id("err").Op(":=").Id("pg").Dot("Exec").Params(
		jen.Id("ctx"), jen.Id("sql"), jen.Id("parameters").Op("..."),
	))
	cblck = append(cblck, jen.If(jen.Id("err").Op("!=").Nil()).Block(
		jen.Return(jen.Nil(), jen.Id("err")),
	))
	cblck = append(cblck, jen.Return(jen.Op("&").Id("u"), jen.Nil()))
	f.ImportName("github.com/jackc/pgx/v5/pgxpool", "pgxpool")

	f.Commentf("Update '%s' in table '%s' based on id columns", sourceTypeName, tableName)
	f.Func().Id(sourceTypeName+"Update").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("pg").Op("*").Qual(
			"github.com/jackc/pgx/v5/pgxpool",
			"Pool",
		),
		jen.Id("u").Qual(sourceTypePackage, sourceTypeName),
	).Params(jen.Op("*").Qual(sourceTypePackage, sourceTypeName), jen.Error()).Block(cblck...)
}

func generateRead(f *jen.File, tableName string, sourceTypeName string, structType *types.Struct, sourceTypePackage string) {
	//goPackage := os.Getenv("GOPACKAGE")

	f.ImportName("github.com/jackc/pgx/v5/pgxpool", "pgxpool")
	f.ImportAlias("github.com/jackc/pgx/v5", "pgx")
	identifiers := []jen.Code{}
	cblck := []jen.Code{jen.Id("r").Op(":=").Qual(sourceTypePackage, sourceTypeName).Op("{}")}
	cblck = append(cblck, jen.Id("identifiers").Op(":=").Op("[]").Interface().Op("{}"))
	selectStr := "SELECT "
	whereStr := " WHERE "
	paramCount := 1
	scanParams := []jen.Code{}
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		rawTagValue := structType.Tag(i)
		//sqlStr := "SELECT id, owner_id, owner_type, name, keys FROM dataset_metadata WHERE id = $1"
		matches := structColPattern.FindStringSubmatch(rawTagValue)
		if matches == nil {
			continue
		}
		raw := strings.Split(matches[1], ",")
		col := raw[0]
		isId := false
		if len(raw) > 1 && raw[1] == "primary" {
			isId = true
		}
		if isId {
			identifiers = append(identifiers, jen.Id(field.Name()).Qual("", field.Type().String()))
			if paramCount > 1 {
				whereStr += ", "
			}
			whereStr += fmt.Sprintf("%s = $%d", col, paramCount)
			cblck = append(cblck, jen.Id("identifiers").Op("=").Op("append(").Id("identifiers").Op(",").Id(field.Name()).Op(")"))
			paramCount += 1
		}
		if i > 0 {
			selectStr += ", "
		}
		// 	err := res.Scan(&r.Id, &r.OwnerId, &r.OwnerType, &r.Name, &r.Keys)

		scanParams = append(scanParams, jen.Op("&").Id("r").Dot(field.Name()))
		selectStr += col

	}

	cblck = append(cblck, jen.Id("sql").Op(":=").Op(fmt.Sprintf("\"%s FROM %s %s\"", selectStr, tableName, whereStr)))

	cblck = append(cblck, jen.Id("res").Op(":=").Id("pg").Dot("QueryRow").Params(
		jen.Id("ctx"), jen.Id("sql"), jen.Id("identifiers").Op("..."),
	))

	cblck = append(cblck, jen.Id("err").Op(":=").Id("res").Dot("Scan").Params(scanParams...))

	cblck = append(cblck, jen.If(jen.Id("err").Op("==").Qual("github.com/jackc/pgx/v5", "ErrNoRows")).Block(
		jen.Return(jen.Nil(), jen.Nil()),
	))

	cblck = append(cblck, jen.If(jen.Id("err").Op("!=").Nil()).Block(
		jen.Return(jen.Nil(), jen.Id("err")),
	))

	cblck = append(cblck, jen.Return(jen.Op("&").Id("r"), jen.Nil()))

	inputs := []jen.Code{jen.Id("ctx").Qual("context", "Context"), jen.Id("pg").Op("*").Qual(
		"github.com/jackc/pgx/v5/pgxpool",
		"Pool",
	)}

	inputs = append(inputs, identifiers...)

	f.Commentf("Read '%s' in table '%s' based on id columns - nil but no error if not found", sourceTypeName, tableName)
	f.Func().Id(sourceTypeName+"Read").Params(
		inputs...,
	).Params(jen.Op("*").Qual(sourceTypePackage, sourceTypeName), jen.Error()).Block(cblck...)
}

func generateDelete(f *jen.File, tableName string, sourceTypeName string, structType *types.Struct) {
	f.ImportName("github.com/jackc/pgx/v5/pgxpool", "pgxpool")
	identifiers := []jen.Code{}
	cblck := []jen.Code{jen.Id("identifiers").Op(":=").Op("[]").Interface().Op("{}")}
	deleteStr := fmt.Sprintf("DELETE FROM %s WHERE ", tableName)
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		rawTagValue := structType.Tag(i)
		matches := structColPattern.FindStringSubmatch(rawTagValue)
		if matches == nil {
			continue
		}
		raw := strings.Split(matches[1], ",")
		col := raw[0]
		isId := false
		if len(raw) > 1 && raw[1] == "primary" {
			isId = true
		}
		if isId {
			identifiers = append(identifiers, jen.Id(field.Name()).Qual("", field.Type().String()))
			if i > 0 {
				deleteStr += ", "
			}
			deleteStr += fmt.Sprintf("%s = $%d", col, i+1)
			cblck = append(cblck, jen.Id("identifiers").Op("=").Op("append(").Id("identifiers").Op(",").Id(field.Name()).Op(")"))
		}

	}

	cblck = append(cblck, jen.Id("sql").Op(":=").Op(fmt.Sprintf("\"%s\"", deleteStr)))

	cblck = append(cblck, jen.Id("_").Op(",").Id("err").Op(":=").Id("pg").Dot("Exec").Params(
		jen.Id("ctx"), jen.Id("sql"), jen.Id("identifiers").Op("..."),
	))

	cblck = append(cblck, jen.If(jen.Id("err").Op("!=").Nil()).Block(
		jen.Return(jen.Id("err")),
	))

	cblck = append(cblck, jen.Return(jen.Nil()))

	inputs := []jen.Code{jen.Id("ctx").Qual("context", "Context"), jen.Id("pg").Op("*").Qual(
		"github.com/jackc/pgx/v5/pgxpool",
		"Pool",
	)}

	inputs = append(inputs, identifiers...)

	f.Commentf("Deletes '%s' in table '%s' based on id columns - no error if not found", sourceTypeName, tableName)
	f.Func().Id(sourceTypeName + "Delete").Params(
		inputs...,
	).Params(jen.Error()).Block(cblck...)
}
