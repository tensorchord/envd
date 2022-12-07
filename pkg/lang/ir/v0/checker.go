package v0

import (
	"reflect"
	"regexp"
	"strings"
)

func (g generalGraph) GetDepsFiles(deps []string) []string {
	tHandle := reflect.TypeOf(g)
	vHandle := reflect.ValueOf(g)
	deps = searchFileInGraph(tHandle, vHandle, deps)
	return deps
}

// Match all filed in ir.Graph with the given keyword
func likeFileFiled(str string) bool {
	nameKeyword := []string{
		"File",
		"Path",
		"Wheels",
	}
	if len(nameKeyword) == 0 {
		return true
	}
	re := regexp.MustCompile(strings.Join(nameKeyword, "|"))
	return re.MatchString(str)
}

// search all files in Graph
func searchFileInGraph(tHandle reflect.Type, vHandle reflect.Value, deps []string) []string {
	for i := 0; i < vHandle.NumField(); i++ {
		v := vHandle.Field(i)
		if v.Type().Kind() == reflect.Struct {
			t := v.Type()
			deps = searchFileInGraph(t, v, deps)
		} else if v.Type().Kind() == reflect.Ptr {
			if v.Type().Elem().Kind() == reflect.Struct {
				if v.Elem().CanAddr() {
					t := v.Type().Elem()
					deps = searchFileInGraph(t, v.Elem(), deps)
				}
			}
		} else {
			t := tHandle.Field(i)
			fieldName := t.Name
			if likeFileFiled(fieldName) {
				typeName := t.Type.String()
				if v.Interface() != nil {
					if typeName == "string" {
						deps = append(deps, v.Interface().(string))
					}
					if typeName == "*string" {
						deps = append(deps, *(v.Interface().(*string)))
					}
					if typeName == "[]string" {
						filesList := v.Interface().([]string)
						deps = append(deps, filesList...)
					}
				}
			}
		}
	}
	return deps
}
