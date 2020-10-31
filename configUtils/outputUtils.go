package configUtils

import (
	"configopt/globalUtils"
	"configopt/model"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

var OutputFile string

func RewriteConfig(config model.Configuration) {
	config = cleanConfig(config)
	FullConfig = InjectMappingRules(FullConfigBytes, config)
	if OutputFile != "" {
		jsonized, err := json.Marshal(&FullConfig.FullConfigContainer)
		if err != nil {
			fmt.Println(err)
			return
		}
		absFile, err := filepath.Abs(OutputFile)
		if err != nil {
			panic(err)
		}
		file, error := os.Create(absFile)
		if error != nil {
			panic(error)
		}
		fmt.Fprint(file, string(jsonized))
		file.Close()
	}
}

// func removeRule(proxyRules []model.MappingRule, i int) []model.MappingRule {
// 	proxyRules[i] = proxyRules[len(proxyRules)-1]
// 	return proxyRules[:len(proxyRules)-1]
// }

func cleanConfig(config model.Configuration) model.Configuration {

	for outerIndex, proxyConfigOuter := range config.ProxyConfigsOuter {
		k := 0
		rules := proxyConfigOuter.ProxyConfig.Content.Proxy.Proxy_rules
		for _, proxyRule := range rules {
			//removing rules marked for deletion here
			if !proxyRule.IsMarkedForDeletion {
				rules[k] = proxyRule
				k++
			}
		}
		config.ProxyConfigsOuter[outerIndex].ProxyConfig.Content.Proxy.Proxy_rules = rules[:k]
	}
	return config
}

func PrintIssues() {
	sort.Slice(globalUtils.Issues, func(i, j int) bool {
		return globalUtils.Issues[i].Severity < globalUtils.Issues[j].Severity
	})
	if OutputFile == "" {
		for _, issue := range globalUtils.Issues {
			fmt.Println(issue.String())
		}
		return
	}
	absFile, err := filepath.Abs(OutputFile)
	if err != nil {
		panic(err)
	}
	file, error := os.Create(absFile)
	if error != nil {
		panic(error)
	}
	for _, issue := range globalUtils.Issues {
		fmt.Fprint(file, issue.String())
	}
	file.Close()
}