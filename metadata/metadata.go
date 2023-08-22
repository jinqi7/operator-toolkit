/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metadata

import (
	"errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// HasAnnotation checks if a given annotation exists
func HasAnnotation(object v1.Object, key string) bool {
	annotations, err := GetAnnotationsWithPrefix(object, key)
	if err != nil || len(annotations) == 0 {
		return false
	}
	return true
}

// HasAnnotationWithValue checks if an annotation exists by searching the key/ value 
func HasAnnotationWithValue(object v1.Object, key, value string) bool {
	if annotations, err := GetAnnotationsWithPrefix(object, key); err == nil {
		for _, val := range annotations {
			if val == value {
				return true
			}
		}
	}
	return false
}

// HasLabel checks if a given Label exists
func HasLabel(object v1.Object, key string) bool {
	labels, err := GetLabelsWithPrefix(object, key)
	if err != nil || len(labels) == 0 {
		return false
	}
	return true
}

// HasLabelWithValue checks if a Label exists by searching the key/ value
func HasLabelWithValue(object v1.Object, key, value string) bool {
	if labels, err := GetLabelsWithPrefix(object, key); err == nil {
		for _, val := range labels {
			if val == value {
				return true
			}
		}
	}
	return false
}

// AddAnnotation adds an annotation key/value to an object
func AddAnnotation(object v1.Object, key string, value string) {
	annotations := map[string]string{}
	annotations[key] = value

	AddAnnotations(object, annotations)
}

// AddLabel adds a Label key/value to an object
func AddLabel(object v1.Object, key string, value string) {
	labels := map[string]string{}
	labels[key] = value

	AddLabels(object, labels)
}

// AddAnnotations copies the map into the resource's Annotations map.
// When the destination map is nil, then the map will be created.
// The unexported function addEntries is called with args passed.
func AddAnnotations(obj v1.Object, entries map[string]string) error {
	if obj == nil {
		return errors.New("object cannot be nil")
	}

	if obj.GetAnnotations() == nil {
		obj.SetAnnotations(map[string]string{})
	}
	addEntries(entries, obj.GetAnnotations())

	return nil
}

// AddLabels copies the map into the resource's Labels map.
// When the destination map is nil, then the map will be created.
// The unexported function addEntries is called with args passed.
func AddLabels(obj v1.Object, entries map[string]string) error {
	if obj == nil {
		return errors.New("object cannot be nil")
	}

	if obj.GetLabels() == nil {
		obj.SetLabels(map[string]string{})
	}
	addEntries(entries, obj.GetLabels())

	return nil
}

// GetAnnotationsWithPrefix is a method that returns a map of key/value pairs matching a prefix string.
// The unexported function filterByPrefix is called with args passed.
func GetAnnotationsWithPrefix(obj v1.Object, prefix string) (map[string]string, error) {
	if obj == nil {
		return map[string]string{}, errors.New("object cannot be nil")
	}

	return filterByPrefix(obj.GetAnnotations(), prefix), nil
}

// GetLabelsWithPrefix is a method that returns a map of key/value pairs matching a prefix string.
// The unexported function filterByPrefix is called with args passed.
func GetLabelsWithPrefix(obj v1.Object, prefix string) (map[string]string, error) {
	if obj == nil {
		return map[string]string{}, errors.New("object cannot be nil")
	}

	return filterByPrefix(obj.GetLabels(), prefix), nil
}

// addEntries copies key/value pairs in the source map adding them into the destination map.
// The unexported function safeCopy is used to copy, and avoids clobbering existing keys in the destination map.
func addEntries(source, destination map[string]string) {
	for key, val := range source {
		safeCopy(destination, key, val)
	}
}

// filterByPrefix returns a map of key/value pairs contained in src that matches the prefix.
// When the prefix is empty/nil, the source map is returned.
// When source key does not contain the prefix string, no copy happens.
func filterByPrefix(entries map[string]string, prefix string) map[string]string {
	if len(prefix) == 0 {
		return entries
	}
	dst := map[string]string{}
	for key, val := range entries {
		if strings.HasPrefix(key, prefix) {
			dst[key] = val
		}
	}
	return dst
}

// safeCopy conditionally copies a given key/value pair into a map.
// When a key is already present in the map, no copy happens.
func safeCopy(destination map[string]string, key, val string) {
	if _, err := destination[key]; !err {
		destination[key] = val
	}
}

// copyWithNewPrefix copies key/value pairs from a source map to a destination map where the key matches the specified prefix.
// If replacementPrefix is different from prefix, the prefix will be replaced while performing the copy.
func copyWithNewPrefix(src, dest map[string]string, prefix, replacementPrefix string) {
	for key, value := range src {
		if strings.HasPrefix(key, prefix) {
			newKey := key
			if prefix != replacementPrefix {
				newKey = strings.Replace(key, prefix, replacementPrefix, 1)
			}
			dest[newKey] = value
		}
	}
}

// CopyLabelsByPrefix copies all labels from a source object to a destination object where the key matches the specified prefix.
// If replacementPrefix is different from prefix, the prefix will be replaced while performing the copy.
func CopyLabelsByPrefix(src, dest v1.Object, prefix, replacementPrefix string) {
	if src.GetLabels() == nil {
		return
	}
	if dest.GetLabels() == nil {
		dest.SetLabels(make(map[string]string))
	}
	copyWithNewPrefix(src.GetLabels(), dest.GetLabels(), prefix, replacementPrefix)
}

// CopyAnnotationsByPrefix copies all annotations from a source object to a destination object where the key matches the specified prefix.
// If replacementPrefix is different from prefix, the prefix will be replaced while performing the copy.
func CopyAnnotationsByPrefix(src, dest v1.Object, prefix, replacementPrefix string) {
	if src.GetAnnotations() == nil {
		return
	}
	if dest.GetAnnotations() == nil {
		dest.SetAnnotations(make(map[string]string))
	}
	copyWithNewPrefix(src.GetAnnotations(), dest.GetAnnotations(), prefix, replacementPrefix)
}
