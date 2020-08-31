/*
Copyright 2020 The Netease Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package curveservice

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type CurveVolumeStatus string

type CurveVolumeDetail struct {
	Id         string            `param:"id"`
	ParentId   string            `param:"parentid"`
	FileType   string            `param:"filetype"`
	LengthGiB  int               `param:"length(GB)"`
	CreateTime string            `param:"createtime"`
	User       string            `param:"user"`
	FileName   string            `param:"filename"`
	FileStatus CurveVolumeStatus `param:"fileStatus"`
}

// Parse the output of 'curve stat':
// id: 39007
// parentid: 39005
// filetype: INODE_PAGEFILE
// length(GB): 10
// createtime: 2020-08-07 10:51:52
// user: k8s
// filename: pvc-ce482926-91d8-11ea-bf6e-fa163e23ce53a
// fileStatus: Created
func simpleParseVolumeDetail(info []byte) (*CurveVolumeDetail, error) {
	infoMap := make(map[string]string)
	for _, line := range strings.Split(string(info), "\n") {
		if !strings.Contains(line, ": ") {
			continue
		}
		kvSlice := strings.SplitN(line, ": ", 2)
		if strings.TrimSpace(kvSlice[0]) != "" {
			infoMap[strings.TrimSpace(kvSlice[0])] = strings.TrimSpace(kvSlice[1])
		}
	}

	volDetail := &CurveVolumeDetail{}
	rv := reflect.ValueOf(volDetail).Elem()
	numField := rv.NumField()
	for i := 0; i < numField; i++ {
		tag := rv.Type().Field(i).Tag.Get("param")
		if tag == "" {
			continue
		}
		value, ok := infoMap[tag]
		if !ok || value == "" {
			continue
		}
		field := rv.Field(i)
		if !field.CanSet() {
			continue
		}
		switch field.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return nil, err
			}
			field.SetBool(v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(value, 10, 0)
			if err != nil {
				return nil, err
			}
			field.SetInt(v)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.ParseUint(value, 10, 0)
			if err != nil {
				return nil, err
			}
			field.SetUint(v)
		case reflect.Float32, reflect.Float64:
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			field.SetFloat(v)
		case reflect.String:
			field.SetString(value)
		default:
			return nil, fmt.Errorf("can not support parse type %v", field.Kind())
		}
	}

	return volDetail, nil
}
