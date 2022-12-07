// Copyright 2022 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
