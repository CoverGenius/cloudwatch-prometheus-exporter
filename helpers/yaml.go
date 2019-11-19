package helpers

import (
	"gopkg.in/yaml.v2"
	"log"
)

func YAMLDecode(f *string, i interface{}) {
	if IsFileExists(f) {
		content := ReadFile(f)
		err := yaml.Unmarshal(*content, i)
		LogErrorExit(err)
	} else {
		log.Fatalf("File: %s does not exists!\n", *f)
	}
}
