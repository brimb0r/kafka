package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Accessor interface {
	Get(key string) (interface{}, error)
}

func (config *configurator) yamlConfig(yamlFile []byte, conf interface{}) error {
	configMap := make(map[string]interface{})
	err := yaml.Unmarshal(yamlFile, &configMap)
	if err != nil {
		return fmt.Errorf("unmarshaling yaml to configMap: %w", err)
	}
	//traverse all map items and accessorInterpolate
	return config.mapConfig(configMap, err, conf)
}

func (config *configurator) setAccessor(id string, accessor Accessor) {
	if config.accessors == nil {
		config.accessors = make(map[string]*Accessor)
	}
	config.accessors[id] = &accessor
}

func (config *configurator) mapConfig(configMap map[string]interface{}, err error, conf interface{}) error {
	configMap, err = config.applyAccessorsAndEnvironment(configMap, "")
	if err != nil {
		return fmt.Errorf("apply config accessors: %w", err)
	}

	//marshal into yaml format
	marshal, err := yaml.Marshal(configMap)
	if err != nil {
		return fmt.Errorf("marshal configMap to yaml: %w", err)
	}

	//unmarshal
	err = yaml.Unmarshal(marshal, conf)
	if err != nil {
		return err
	}
	return nil
}

func (config *configurator) applyAccessorsAndEnvironment(configMap map[string]interface{}, keyPath string) (map[string]interface{}, error) {
	var environmentVariable string

	for k, v := range configMap {
		environmentVariable = bashEnvString(fmt.Sprintf("%s%s", keyPath, k))
		switch v.(type) {
		case map[string]interface{}:
			var err error
			configMap[k], err = config.applyAccessorsAndEnvironment(v.(map[string]interface{}), fmt.Sprintf("%s%s_", keyPath, k))
			if err != nil {
				return configMap, fmt.Errorf("key[%s] %w", k, err)
			}
		case []interface{}:
			replacement := make([]interface{}, len(v.([]interface{})))
			for i, str := range v.([]interface{}) {
				if envVar := envReplacement(fmt.Sprintf("%s_%d", environmentVariable, i)); envVar != "" {
					replacement[i] = envVar
				} else {
					switch str.(type) {
					case string:
						interpolation, err := config.accessorInterpolate(str.(string))
						if err != nil {
							return configMap, fmt.Errorf("accessorInterpolate key[%s] value[%s]; %w", k, v, err)
						}
						replacement[i] = interpolation
					default:
						replacement[i] = str
					}
				}
			}
			configMap[k] = replacement
		case string:
			if envVar := envReplacement(environmentVariable); envVar != "" {
				configMap[k] = envVar
			} else {
				interpolation, err := config.accessorInterpolate(v.(string))
				if err != nil {
					return configMap, fmt.Errorf("accessorInterpolate key[%s] value[%s]; %w", k, v, err)
				}
				configMap[k] = interpolation
			}
		case nil:
			if envVar := envReplacement(environmentVariable); envVar != "" {
				configMap[k] = envVar // will be set as a string...
			} else {
				log.Printf("%s not set\n", environmentVariable)
			}
		case int:
			if envVar := envReplacement(environmentVariable); envVar != "" {
				intVersion, err := strconv.Atoi(envVar)
				if err != nil {
					return configMap, fmt.Errorf("atoi key[%s] %w", k, err)
				}
				configMap[k] = intVersion
			} else {
				configMap[k] = v
			}
		case bool:
			if envVar := envReplacement(environmentVariable); envVar != "" {
				boolVersion, err := strconv.ParseBool(envVar)

				if err != nil {
					return configMap, fmt.Errorf("boolean key[%s] %w", k, err)
				}
				configMap[k] = boolVersion
			} else {
				configMap[k] = v
			}
		default:
			return configMap, fmt.Errorf("key[%s]: we do not support type %T", k, v)
		}
	}
	return configMap, nil
}

func (config *configurator) accessorInterpolate(str string) (interface{}, error) {
	rx := regexp.MustCompile(`^\{(.*?)\}$`)
	if rx.MatchString(str) {
		tempStr := str[1 : len(str)-1]
		split := strings.Split(tempStr, ":")
		if len(split) == 2 {
			accessorId := split[0]
			accessorKey := split[1]

			if accessor, ok := config.accessors[accessorId]; ok {
				return (*accessor).Get(accessorKey)
			}
		}
	}
	return str, nil
}

func envReplacement(str string) string {
	environmentVariable := os.Getenv(str)
	if len(environmentVariable) != 0 {
		return environmentVariable
	}
	return ""
}

func bashEnvString(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}
