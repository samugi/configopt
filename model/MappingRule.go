package model

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type MappingRule struct {
	Id                     *int64            `json:"id"`
	Proxy_id               *int64            `json:"proxy_id"`
	Http_method            *string           `json:"http_method"`
	Pattern                *string           `json:"pattern"`
	Owner_id               *int64            `json:"owner_id"`
	Owner_type             *string           `json:"owner_type"`
	Querystring_parameters map[string]string `json:"querystring_parameters"`

	Metric_id          *int64   `json:"metric_id"`
	Metric_system_name *string  `json:"metric_system_name"`
	Delta              *int64   `json:"delta"`
	Tenant_id          *int64   `json:"tenant_id"`
	Created_at         *string  `json:"created_at"`
	Updated_at         *string  `json:"updated_at"`
	Redirect_url       *string  `json:"redirect_url"`
	Position           *int64   `json:"position"`
	Last               *bool    `json:"last"`
	Parameters         []string `json:"parameters"`

	//stuff I added
	//	QueryPairs          map[string]string `json:"-"`
	IsExactMatch        bool   `json:"-"`
	Host                string `json:"-"`
	IsMarkedForDeletion bool   `json:"-"`
}

func (rule *MappingRule) Initialize(host string) {
	query := rule.getQuery(*(rule.Pattern))
	rule.Host = host
	if query != "" {
		*rule.Pattern = strings.Replace(*(rule.Pattern), "?"+query, "", 1)
	}
	if strings.HasSuffix(*(rule.Pattern), "$") {

		//	*rule.Pattern = (*rule.Pattern)[0 : len(*(rule.Pattern))-1]
		rule.IsExactMatch = true
	}
	if query == "" {
		return
	}
	// pairs := strings.Split(query, "&")
	// for _, pair := range pairs {
	// 	index := strings.Index(pair, "=")
	// 	if index > 0 {
	// 		rule.QueryPairs[pair[0:index]] = pair[index+1:]
	// 	}
	// }
}

func (rule *MappingRule) SetMarkedForDeletion(marked bool) {
	(*rule).IsMarkedForDeletion = marked
}

func (rule MappingRule) String() string {
	var proxy_id int64
	var owner_type string
	if rule.Proxy_id != nil {
		proxy_id = *rule.Proxy_id
	}
	if rule.Owner_type != nil {
		owner_type = *rule.Owner_type
	}
	return *rule.Http_method + " " + rule.getRealPath() + " - ID: " + strconv.FormatInt(*(rule.Id), 10) + " - Service ID: " + fmt.Sprintf("%d", int(proxy_id)) + " Host: " + rule.Host + " Owner type: " + owner_type
}

func (rule MappingRule) BrutalMatch(mr *MappingRule) bool {
	return rule.matches(*mr) && !rule.optimizationMatch(*mr)
}

func (rule MappingRule) CanBeOptimized(mr *MappingRule) bool {
	return rule.matches(*mr) && rule.optimizationMatch(*mr)
}

func (rule MappingRule) getRealPath() string {
	//	if rule.IsExactMatch {
	//		return *rule.Pattern + "$"
	//	}
	return *rule.Pattern
}

func (rule *MappingRule) SetExactMatch(em bool) {
	rule.IsExactMatch = em
}

func (rule MappingRule) optimizationMatch(mr MappingRule) bool {
	shorterExactMatch := GetShorter(&rule, &mr).IsExactMatch
	sameSectionsLengths := rule.getPathSectionsLength() == mr.getPathSectionsLength() && rule.getLastSectionLength() == mr.getLastSectionLength()
	return !sameSectionsLengths && shorterExactMatch
}

func (rule MappingRule) matches(mr MappingRule) bool {
	matghingMethods := strings.EqualFold(*rule.Http_method, *mr.Http_method)
	matchingQP := rule.matchQP(mr)
	matchingPath := rule.matchingPath(mr)
	return matghingMethods && matchingQP && matchingPath
}

func (rule MappingRule) getPathSectionsLength() int {
	pattern := strings.Trim(*rule.Pattern, "$")
	mr1 := strings.Split(pattern, "/")
	return len(mr1)
}

func (rule MappingRule) getLastSectionLength() int {
	pattern := strings.Trim(*rule.Pattern, "$")
	mr1 := strings.Split(pattern, "/")
	if len(mr1) == 0 {
		return 0
	}
	lastSection := mr1[len(mr1)-1]
	return len(lastSection)
}

func GetShorter(mr1 *MappingRule, mr2 *MappingRule) *MappingRule {
	p1 := strings.Trim(*mr1.Pattern, "$")
	p2 := strings.Trim(*mr2.Pattern, "$")
	if len(p1) < len(p2) {
		return mr1
	}
	return mr2
}

func (rule MappingRule) matchQP(mr MappingRule) bool {
	return reflect.DeepEqual(rule.Querystring_parameters, mr.Querystring_parameters)
}

func (rule MappingRule) matchingPath(mr MappingRule) bool {
	p1 := strings.Trim(*rule.Pattern, "$")
	p2 := strings.Trim(*mr.Pattern, "$")
	if strings.HasPrefix(p1, p2) || strings.HasPrefix(p2, p1) {
		return true
	}
	if rule.matchWithParams(mr) {
		return true
	}
	return false
}

func (rule MappingRule) matchWithParams(mr MappingRule) bool {
	p1 := strings.Trim(*rule.Pattern, "$")
	p2 := strings.Trim(*mr.Pattern, "$")
	mr1 := strings.Split(p1, "/")
	mr2 := strings.Split(p2, "/")
	if len(mr1) != len(mr2) {
		return false
	}
	for i := 0; i < len(mr1); i++ {
		if mr1[i] != mr2[i] && !isParam(mr1[i]) && !isParam(mr2[i]) {
			return false
		}
	}
	last1 := mr1[len(mr1)-1]
	last2 := mr2[len(mr2)-1]
	return matchLastPartial(last1, last2)
}

func matchLastPartial(last1 string, last2 string) bool {
	return strings.HasPrefix(last1, last2) || strings.HasPrefix(last2, last1)
}

func isParam(p string) bool {
	return strings.HasPrefix(p, "{") && strings.HasPrefix(p, "}")
}

func (rule MappingRule) getQuery(pattern string) string {
	lastQuery := strings.LastIndex(pattern, "?")
	lastSlash := strings.LastIndex(pattern, "/")
	if lastQuery > lastSlash {
		return pattern[lastQuery+1:]
	}
	return ""
}
