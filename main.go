package main

import (
	"bytes"
	"cuelang.org/go/cue"
	"fmt"
	"github.com/oam-dev/kubevela/apis/types"
	velacue "github.com/oam-dev/kubevela/pkg/cue"
	"github.com/oam-dev/kubevela/pkg/cue/packages"
	"github.com/oam-dev/kubevela/pkg/utils/common"
	"github.com/pkg/errors"
	"strings"
)

func main() {
	// Create a new cue interpreter.
	c := common.Args{}
	config, err := c.GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	pd, err := packages.NewPackageDiscover(config)
	if err != nil {
		fmt.Println(err)
	}
	value, err := common.GetCUEParameterValue(annotations, pd)
	if err != nil {
		fmt.Println(err)
	}
	parameterStructs := GeneratorParameterStructs(value)
	printParamGosStruct(parameterStructs)
}

type StructParameter struct {
	types.Parameter
	// GoType is the same to parameter.Type but can be print in Go
	GoType string
	Fields []Field
}

type Field struct {
	Name   string
	Type   string
	GoType string
}

var structs []StructParameter

func GeneratorParameterStructs(param cue.Value) []StructParameter {
	structs = []StructParameter{}
	err := parseParameters(param, "Parameter")
	if err != nil {
		fmt.Println(err)
	}
	return structs
}

func NewStructParameter() StructParameter {
	return StructParameter{
		Parameter: types.Parameter{},
		GoType:    "",
		Fields:    []Field{},
	}
}

func parseParameters(paraValue cue.Value, paramKey string) error {
	param := NewStructParameter()
	param.Name = paramKey
	param.Type = paraValue.IncompleteKind()
	param.Short, param.Usage, param.Alias, param.Ignore = velacue.RetrieveComments(paraValue)
	if def, ok := paraValue.Default(); ok && def.IsConcrete() {
		param.Default = velacue.GetDefault(def)
	}

	switch param.Type {
	case cue.StructKind:
		arguments, err := paraValue.Struct()
		if err != nil {
			return fmt.Errorf("augument not as struct %w", err)
		}
		if arguments.Len() == 0 {
			var SubParam StructParameter
			SubParam.Name = "-"
			SubParam.Required = true
			tl := paraValue.Template()
			if tl != nil { // is map type
				kind, err := trimIncompleteKind(tl("").IncompleteKind().String())
				if err != nil {
					return errors.Wrap(err, "invalid parameter kind")
				}
				SubParam.GoType = kind
				// TODO: better way to name to subParam
				SubParam.Name = "Item"
				param.GoType = fmt.Sprintf("map[string]%s", SubParam.Name)
				structs = append(structs, SubParam)
			}
		}
		for i := 0; i < arguments.Len(); i++ {
			var subParam Field
			fi := arguments.Field(i)
			if fi.IsDefinition {
				continue
			}
			val := fi.Value
			name := fi.Name
			subParam.Name = name
			switch val.IncompleteKind() {
			case cue.StructKind:
				if subField, err := val.Struct(); err == nil && subField.Len() == 0 { // err cannot be not nil,so ignore it
					if mapValue, ok := val.Elem(); ok {
						// In the future we could recursively call to support complex map-value(struct or list)
						subParam.GoType = fmt.Sprintf("map[string]%s", mapValue.IncompleteKind().String())
					} else {
						return fmt.Errorf("failed to got Map kind from %s", subParam.Name)
					}
				} else {
					if err := parseParameters(val, name); err != nil {
						return err
					}
					subParam.GoType = dm.FieldName(name)
				}
			case cue.ListKind:
				elem, success := val.Elem()
				if !success {
					// fail to get elements, use the value of ListKind to be the type
					subParam.GoType = val.IncompleteKind().String()
					break
				}
				switch elem.Kind() {
				case cue.StructKind:
					subParam.GoType = fmt.Sprintf("[]%s", dm.FieldName(name))
					if err := parseParameters(elem, name); err != nil {
						return err
					}
				default:
					subParam.GoType = fmt.Sprintf("[]%s", elem.IncompleteKind().String())
				}
			default:
				subParam.GoType = val.IncompleteKind().String()
			}
			param.Fields = append(param.Fields, Field{
				Name:   subParam.Name,
				GoType: subParam.GoType,
			})
		}
	}
	structs = append(structs, param)
	return nil
}

func printParamGosStruct(parameters []StructParameter) {
	fmt.Printf("package main\n\n")
	for _, parameter := range parameters {
		if parameter.Usage == "" {
			parameter.Usage = "-"
		}
		fmt.Printf("// %s %s\n", dm.FieldName(parameter.Name), parameter.Usage)
		printField(parameter)
	}
}

func printField(param StructParameter) {
	buffer := &bytes.Buffer{}
	fieldName := dm.FieldName(param.Name)
	switch param.Type {
	case cue.StructKind:
		// struct in cue can be map or struct
		if strings.HasPrefix(param.GoType, "map[string]") {
			fmt.Fprintf(buffer, "type %s %s", fieldName, param.GoType)
		} else {
			fmt.Fprintf(buffer, "type %s struct {\n", fieldName)
			for _, f := range param.Fields {
				fmt.Fprintf(buffer, "    %s %s `json:\"%s\"`\n", dm.FieldName(f.Name), f.GoType, f.Name)
			}

			fmt.Fprintf(buffer, "}\n")
		}
	case cue.StringKind:
		fmt.Fprintf(buffer, "type %s string\n", fieldName)
	case cue.IntKind:
		fmt.Fprintf(buffer, "type %s int\n", fieldName)
	case cue.BoolKind:
		fmt.Fprintf(buffer, "type %s bool\n", fieldName)
	case cue.FloatKind:
		fmt.Fprintf(buffer, "type %s float64\n", fieldName)
	case cue.ListKind:
		fmt.Fprintf(buffer, "type %s []%s\n", fieldName, param.GoType)
	default:
		fmt.Fprintf(buffer, "type %s %s\n", fieldName, param.GoType)
	}
	//source, err := format.Source(buffer.Bytes())
	//if err != nil {
	//	fmt.Println("Failed to format source:", err)
	//}
	fmt.Println(buffer.String())
	//fmt.Println(string(source))
}

func trimIncompleteKind(mask string) (string, error) {
	mask = strings.Trim(mask, "()")
	ks := strings.Split(mask, "|")
	if len(ks) == 1 {
		return ks[0], nil
	}
	if len(ks) == 2 && ks[0] == "null" {
		return ks[1], nil
	}
	return "", fmt.Errorf("invalid incomplete kind: %s", mask)

}
