package serialization

import "gopkg.in/yaml.v2"

func DeserializeFromYaml(value interface{}, buff []byte) error {
	return yaml.Unmarshal(buff, value)
}
