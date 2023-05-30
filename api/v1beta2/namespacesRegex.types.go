package v1beta2

import (
	"context"
	"regexp"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// struct which contains the information about the namespace regex
type NamespacesRegex struct {

	// +kubebuilder:default={default}
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AvoidRegex []string `json:"avoidregex"`

	// +kubebuilder:default={default}
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
func (nsr NamespacesRegex) CalculateNamespaces(l logr.Logger, ctx context.Context, c client.Client) (mustMatch, mustAvoid []v1.Namespace, err error) {

	var namespaceList = &v1.NamespaceList{}

	l.Info("calculating namespaces for the following lists", "NamespacesRegex", nsr)

	if err = c.List(ctx, namespaceList, &client.ListOptions{}); err == nil {

		// parse through all registered namespaces
		for i := range namespaceList.Items {

			// check, if the namespace has to be avoided during deployment
			var inList bool = false
			if inList, err = stringMatchesRegExpList(namespaceList.Items[i].Name, nsr.AvoidRegex); err != nil {

				l.Error(err, "error calculating avoids", "current namespace", namespaceList.Items[i].Name, "AvoidRegex", nsr.AvoidRegex)

				return

			} else {
				if inList {
					// if the namespace is in the list [AvoidRegex], append the namespace to the list [mustAvoid]
					// and calculate the next namespace
					mustAvoid = append(mustAvoid, namespaceList.Items[i])
					continue
				}
			}

			// if the namespace is not in the list [AvoidRegex], check if the namespace is in the list [MatchRegex]
			if inList, err = stringMatchesRegExpList(namespaceList.Items[i].Name, nsr.MatchRegex); err != nil {

				l.Error(err, "error calculating matches", "current namespace", namespaceList.Items[i].Name, "MatchRegex", nsr.MatchRegex)

				return

			} else {

				if inList {
					// if the namespace is in the list [MatchRegex], then append the namespace to
					// the namespaces [mustMatch]
					mustMatch = append(mustMatch, namespaceList.Items[i])

				} else {
					// if the namespace also is not in the list [MatchRegex], then append the namespace to
					// the namespaces [mustAvoid]
					mustAvoid = append(mustAvoid, namespaceList.Items[i])
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
