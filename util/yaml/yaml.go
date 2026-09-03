package yaml

import (
	"go.yaml.in/yaml/v4"
)

func Marshal(in any) ([]byte, error) {
	return yaml.Marshal(in)
}

func Unmarshal(in []byte, out any) error {
	return yaml.Unmarshal(in, out)
}
