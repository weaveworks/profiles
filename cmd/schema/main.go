package main

import (
	"fmt"
	"io/ioutil"
	"os"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/schemer/definition"
	schemapkg "github.com/weaveworks/schemer/schema"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: schema <object-name> <outfile>. Assumes path to file exists.")
		os.Exit(1)
	}
	obj := os.Args[1]
	outputFile := os.Args[2]

	// generate the schema
	schema, err := schemapkg.GenerateSchema("github.com/weaveworks/profiles/api", profilesv1.GroupVersion.Version, obj, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// edit some things
	schema = editSchema(obj, schema)

	// jsonify and write to the outputFile
	bytes, err := schemapkg.ToJSON(schema)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(outputFile, bytes, 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("schema file generated for `%s`\n", obj)
}

func editSchema(obj string, schema schemapkg.Schema) schemapkg.Schema {
	cc := schema.Definitions[obj]

	// clear the descriptions of kind and apiVersion
	if t, ok := cc.Properties["kind"]; ok {
		t.Enum = []string{obj}
		t.Description = ""
		t.HTMLDescription = ""

	}
	if t, ok := cc.Properties["apiVersion"]; ok {
		t.Enum = []string{fmt.Sprintf("%s/%s", profilesv1.GroupVersion.Group, profilesv1.GroupVersion.Version)}
		t.Description = ""
		t.HTMLDescription = ""
	}

	// change the reference of the metadata, we don't need the whole kube one here
	if t, ok := cc.Properties["metadata"]; ok {
		t.Ref = "#/definitions/Meta"
	}
	schema.Definitions["Meta"] = &definition.Definition{
		Properties: map[string]*definition.Definition{
			"name": {
				Type:            "string",
				Description:     "profile name",
				HTMLDescription: "profile name",
			},
		},
		PreferredOrder:       []string{"name"},
		AdditionalProperties: false,
	}

	// delete various things which we don't want to render
	delete(schema.Definitions, "k8s.io|apimachinery|pkg|apis|meta|v1.ObjectMeta")
	delete(schema.Definitions, "k8s.io|apimachinery|pkg|types.UID")
	delete(schema.Definitions, "ProfileDefinitionStatus")
	delete(cc.Properties, "status")
	cc.PreferredOrder = cc.PreferredOrder[:len(cc.PreferredOrder)-1]

	return schema
}
