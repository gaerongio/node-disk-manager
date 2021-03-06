package ndmutil

import (
	"fmt"
	"github.com/ghodss/yaml"
	sysutil "github.com/openebs/CITF/utils/system"
	"github.com/openebs/node-disk-manager/cmd/controller"
	"io/ioutil"
	"strings"
)

type ConfigMapPatch controller.NodeDiskManagerConfig

// Replace the ConfigMap portion in the configuration
// yaml with the given ConfigMapPatch struct and apply it.
func ReplaceAndApplyConfig(configMapPatch ConfigMapPatch) error {
	yamlBytes, err := ioutil.ReadFile(GetNDMOperatorFilePath())
	if err != nil {
		return err
	}

	yamlString := string(yamlBytes)
	configString, err := yaml.Marshal(configMapPatch)
	stringConfig := strings.Replace(string(configString), "\n", "\n    ", -1)
	s1 := strings.Split(yamlString, "node-disk-manager.config: |")[0]
	s2 := strings.SplitN(yamlString, "---", 2)[1]
	yamlString = s1 + "\n  node-disk-manager.config: | \n    " + stringConfig + "\n--- \n" + s2
	return sysutil.RunCommandWithGivenStdin("kubectl apply -f -", yamlString)
}

// Get the ConfigMap from the NDMConfiguration file
// ConfigMapPatch is generated by parsing the configuration yaml
func GetNDMConfig(fileName string) ConfigMapPatch {
	var configMapPatch ConfigMapPatch
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal(data, &configMapPatch)
	if err != nil {
		fmt.Println(err)
	}
	yamlString := string(data)
	yamlString = strings.Split(yamlString, "---")[0]
	yamlString = strings.Split(yamlString, "node-disk-manager.config: |")[1]
	yaml.Unmarshal([]byte(yamlString), &configMapPatch)
	return configMapPatch
}

// SetIncludePath is used to set the include section in
// path filter of NDMConfig
func (c *ConfigMapPatch) SetIncludePath(deviceList ...string) {
	for index, element := range c.FilterConfigs {
		if element.Key == "path-filter" {
			c.FilterConfigs[index].Include = strings.Join(deviceList, ",")
			c.FilterConfigs[index].Exclude = ""
		}
	}
}

// SetExcludePath is used to set the exclude section in
// path filter of NDMConfig
func (c *ConfigMapPatch) SetExcludePath(deviceList ...string) {
	for index, element := range c.FilterConfigs {
		if element.Key == "path-filter" {
			c.FilterConfigs[index].Exclude = strings.Join(deviceList, ",")
			c.FilterConfigs[index].Include = ""
		}
	}
}

// SetPathFilter is used to change the state of
// path filter in NDMConfig
func (c *ConfigMapPatch) SetPathFilter(state string) {
	for index, element := range c.FilterConfigs {
		if element.Key == "path-filter" {
			c.FilterConfigs[index].State = state
		}
	}
}
