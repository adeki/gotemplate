package template

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

const basePath = "path/to/templates-dir/"

func New(pathList ...string) *template.Template {
	tmpl := template.New("").Funcs(funcMap())
	for _, path := range pathList {
		tmpl = template.Must(tmpl.ParseGlob(basePath + path))
	}
	return tmpl
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"atoi":        atoi,
		"comment":     comment,
		"diff":        diff,
		"encode_json": encodeJson,
		"ftime":       ftime,
		"incr":        incr,
		"itoa":        itoa,
		"join":        join,
		"now":         now,
		"raw":         raw,
		"repeat":      repeat,
		"rmdr":        rmdr,
		"sum":         sum,
		"split":       split,
		"prod":        prod,
		"qtnt":        qtnt,
	}
}

//
// template functions
//

func raw(s string) template.HTML {
	return template.HTML(s)
}

func incr(i int) int {
	return i + 1
}

func sum(a, b int) int {
	return a + b
}

func diff(a, b int) int {
	return a - b
}

func prod(a, b int) int {
	return a * b
}

func qtnt(a, b int) int {
	return a / b
}

func rmdr(a, b int) int {
	return a % b
}

func repeat(start, end int) []int {
	limit := end - start + 1
	res := make([]int, limit)
	for i := 0; i < limit; i++ {
		res[i] = start + i
	}
	return res
}

func now() time.Time {
	return time.Now()
}

func ftime(fmt string, t time.Time) string {
	return t.Format(fmt)
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func join(a []string, sep string) string {
	return strings.Join(a, sep)
}

func encodeJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func comment(s string) template.HTML {
	return raw(fmt.Sprint("<!--", s, "-->"))
}

func split(sep, target string) []string {
	return strings.Split(target, sep)
}
