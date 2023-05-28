package v1beta2

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// struct which contains the information about the namespace regex
type NamespacesRegex struct {

	// +kubebuilder:default={default}
	// +kubebuilder:validation:UniqueItems=true
	// +kubebuilder:validation:MinItems=1
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AvoidRegex []string `json:"avoidregex"`

	// +kubebuilder:default={default}
	// +kubebuilder:validation:UniqueItems=true
	// +kubebuilder:validation:MinItems=1
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	MatchRegex []string `json:"matchregex"`
}

// get two lists of namespaces
//
// the 1. list contains all namespaces
//
// the 2. list, contains all namespaces, which match with the
// list of regexpressions from the matches-array, without the namespaces,
// which match with the avoid-array
func (nsr NamespacesRegex) CalculateNamespaces(l logr.Logger, ctx context.Context, c client.Client) (all, matches []v1.Namespace, err error) {
	var namespaceList = &v1.NamespaceList{}
	l.Info("calculating namespaces for the following lists", fmt.Sprintf("NamespacesRegex:%v", nsr))
	if err = c.List(ctx, namespaceList, &client.ListOptions{}); err != nil {
		return
	} else {
		for i := range namespaceList.Items {
			// append namespace to all.List
			all = append(all, namespaceList.Items[i])
			// check the namespace regex lists and add them accordingly
			var inList bool = false
			if inList, err = stringMatchesRegExpList(namespaceList.Items[i].Name, nsr.AvoidRegex); err != nil {
				l.Error(err, "error calculating avoids", fmt.Sprintf("current namespace: %s | avoids: %v", namespaceList.Items[i].Name, nsr.AvoidRegex))
				return
			} else {
				if !inList {
					if inList, err = stringMatchesRegExpList(namespaceList.Items[i].Name, nsr.MatchRegex); err != nil {
						l.Error(err, "error calculating matches", fmt.Sprintf("current namespace: %s | matches: %v", namespaceList.Items[i].Name, nsr.MatchRegex))
						return
					} else {
						if inList {
							matches = append(matches, namespaceList.Items[i])
						}
					}
				}
			}
		}
	}
	return
}

// check whether a string exists in a list of regexpressions or not
func stringMatchesRegExpList(comp string, regexpList []string) (bool, error) {
	for i := range regexpList {
		if matched, err := regexp.MatchString(regexpList[i], comp); err != nil {
			return false, nil
		} else {
			if matched {
				return true, nil
			}
		}
	}
	return false, nil
}

// find a string in a list of string
//func stringInList(comp string, list []string) bool {
//	for i := range list {
//		if list[i] == comp {
//			return true
//		}
//	}
//	return false
//}
