package helpers

import (
	"log"

	yaml "gopkg.in/yaml.v2"
)

// YAMLDecode reads the file located at path and unmarshals it into the input interface
func YAMLDecode(path *string, i interface{}) {
	if IsFileExists(path) {
		content := ReadFile(path)
		err := yaml.Unmarshal(*content, i)
		LogIfErrorExit(err)
	} else {
		log.Fatalf("File: %s does not exists!\n", *path)
	}
}
