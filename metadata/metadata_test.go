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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Metadata", func() {

	When("AddAnnotation is called", func() {
		It("should add the annotation and its value to the object", func() {
			annotation := "foo"
			value := "bar"
			pod := &corev1.Pod{}

			AddAnnotation(pod, annotation, value)
			Expect(pod.Annotations).To(HaveLen(1))
			Expect(pod.Annotations[annotation]).To(Equal(value))
		})
	})

	When("AddAnnotations is called", func() {
		It("should add the annotations to the object", func() {
			annotations := map[string]string{"foo": "bar", "baz": "qux"}
			pod := &corev1.Pod{}

			Expect(pod.Annotations).To(HaveLen(0))
			Expect(AddAnnotations(pod, annotations)).To(Succeed())
			Expect(pod.Annotations).To(Equal(annotations))
		})

		It("should add the annotations to the existing ones", func() {
			annotations := map[string]string{"foo": "bar", "baz": "qux"}
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(pod.Annotations).To(HaveLen(1))
			Expect(AddAnnotations(pod, annotations)).To(Succeed())
			Expect(pod.Annotations).To(HaveLen(3))
			for key, value := range annotations {
				Expect(pod.Annotations).To(HaveKeyWithValue(key, value))
			}
		})

		It("should error if the object is nil", func() {
			err := AddAnnotations(nil, map[string]string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})
	})

	When("AddLabel is called", func() {
		It("should add the label and its value to the object", func() {
			label := "foo"
			value := "bar"
			pod := &corev1.Pod{}

			AddAnnotation(pod, label, value)
			Expect(pod.Labels).To(HaveLen(1))
			Expect(pod.Labels[label]).To(Equal(value))
		})

	When("AddLabels is called", func() {
		It("should add the labels to the object", func() {
			labels := map[string]string{"foo": "bar", "baz": "qux"}
			pod := &corev1.Pod{}

			Expect(pod.Labels).To(HaveLen(0))
			Expect(AddLabels(pod, labels)).To(Succeed())
			Expect(pod.Labels).To(Equal(labels))
		})

		It("should add the labels to the existing ones", func() {
			labels := map[string]string{"foo": "bar", "baz": "qux"}
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(pod.Labels).To(HaveLen(1))
			Expect(AddLabels(pod, labels)).To(Succeed())
			Expect(pod.Labels).To(HaveLen(3))
			for key, value := range labels {
				Expect(pod.Labels).To(HaveKeyWithValue(key, value))
			}
		})

		It("should error if the object is nil", func() {
			err := AddLabels(nil, map[string]string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})
	})

	When("HasAnnotationWithValue is called", func() {
		It("should return true when the annotation/ with the given value is found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			found := HasAnnotation(pod, "bar")
			Expect(found).To(BeTrue())

			found = HasAnnotationWithValue(pod, "bar", "foo")
			Expect(found).To(BeTrue())
		})

		It("should return false when the annotation/ with the given value is not found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			found := HasAnnotation(pod, "nobar")
			Expect(found).To(BeFalse())

			found = HasAnnotationWithValue(pod, "bar", "nofoo")
			Expect(found).To(BeFalse())
		})
	})

	When("HasLabelWithValue is called", func() {
		It("should return true when the label/ with the given value is found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			found := HasLabel(pod, "bar")
			Expect(found).To(BeTrue())

			found = HasLabelWithValue(pod, "bar", "foo")
			Expect(found).To(BeTrue())
		})

		It("should return false when the label/ with the given value is not found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			found := HasLabel(pod, "nobar")
			Expect(found).To(BeFalse())

			found = HasLabelWithValue(pod, "bar", "nofoo")
			Expect(found).To(BeFalse())
		})
	})

	When("GetAnnotationsWithPrefix is called", func() {
		It("should return all the annotations matching the given prefix", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			annotations, err := GetAnnotationsWithPrefix(pod, "ba")
			Expect(annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bar": Equal("foo"),
				"baz": Equal("qux"),
			}))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an empty map if there are no matching annotations", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			annotations, err := GetAnnotationsWithPrefix(pod, "nil")
			Expect(annotations).To(HaveLen(0))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should error if the object is nil", func() {
			_, err := GetAnnotationsWithPrefix(nil, "prefix")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})

	})

	When("GetLabelsWithPrefix is called", func() {
		It("should return all the labels matching the given prefix", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			labels, err := GetLabelsWithPrefix(pod, "ba")
			Expect(labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bar": Equal("foo"),
				"baz": Equal("qux"),
			}))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an empty map if there are no matching labels", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			labels, err := GetLabelsWithPrefix(pod, "nil")
			Expect(labels).To(HaveLen(0))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should error if the object is nil", func() {
			_, err := GetLabelsWithPrefix(nil, "prefix")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})
	})

	When("addEntries is called", func() {
		It("should merge both maps", func() {
			source := map[string]string{"foo": "bar", "baz": "qux"}
			destination := map[string]string{"quux": "corge"}

			Expect(destination).To(HaveLen(1))
			addEntries(source, destination)
			Expect(destination).To(HaveLen(3))
			for key, value := range source {
				Expect(destination).To(HaveKeyWithValue(key, value))
			}
		})
	})

	When("filterByPrefix is called", func() {
		It("should return the matching elements", func() {
			entries := map[string]string{"bar": "foo", "baz": "qux"}

			matchingEntries := filterByPrefix(entries, "ba")
			Expect(matchingEntries).To(HaveLen(2))
			Expect(matchingEntries).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bar": Equal("foo"),
				"baz": Equal("qux"),
			}))
		})
	})

	When("safeCopy is called", func() {
		It("should copy the given key into the destination map", func() {
			destination := map[string]string{}
			safeCopy(destination, "foo", "bar")
			Expect(destination).To(HaveLen(1))
			Expect(destination).To(HaveKeyWithValue("foo", "bar"))
		})

		It("should preserve the existing key/values", func() {
			destination := map[string]string{"baz": "qux"}
			safeCopy(destination, "foo", "bar")
			Expect(destination).To(HaveLen(2))
			Expect(destination).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo": Equal("bar"),
				"baz": Equal("qux"),
			}))
		})
	})

	When("copyWithNewPrefix is called", func() {
		It("should copy key/value pairs with the prefix to destination with new prefix", func() {
			source := map[string]string{"foo": "bar", "foz": "qux"}
			destination := map[string]string{"quux": "corge"}
			copyWithNewPrefix(source, destination, "fo", "ba")
			Expect(destination).To(HaveLen(3))
			Expect(destination).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bao":  Equal("bar"),
				"baz":  Equal("qux"),
				"quux": Equal("corge"),
			}))
		})
	})

	When("CopyAnnotationsByPrefix is called", func() {
		It("should return all the annotations matching the given prefix", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"foo": "bar", "baz": "qux"},
				},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}
			CopyAnnotationsByPrefix(podSource, podDest, "fo", "ba")
			Expect(podDest.Annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bao":  Equal("bar"),
				"quux": Equal("corge"),
			}))
		})
	})

	When("CopyLabelsByPrefix is called", func() {
		It("should return all the labels matching the given prefix", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"foo": "bar", "baz": "qux"},
				},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}
			CopyLabelsByPrefix(podSource, podDest, "fo", "ba")
			Expect(podDest.Labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bao":  Equal("bar"),
				"quux": Equal("corge"),
			}))
		})
	})
})
