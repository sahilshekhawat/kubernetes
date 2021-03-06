/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package fake

import (
	"k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_1"
	"k8s.io/kubernetes/pkg/client/testing/core"
	extensions_unversioned "k8s.io/kubernetes/pkg/client/typed/generated/extensions/unversioned"
	extensions_unversioned_fake "k8s.io/kubernetes/pkg/client/typed/generated/extensions/unversioned/fake"
	legacy_unversioned "k8s.io/kubernetes/pkg/client/typed/generated/legacy/unversioned"
	legacy_unversioned_fake "k8s.io/kubernetes/pkg/client/typed/generated/legacy/unversioned/fake"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/watch"
)

// Clientset returns a clientset that will respond with the provided objects
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := core.NewObjects(api.Scheme, api.Codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	fakePtr := core.Fake{}
	fakePtr.AddReactor("*", "*", core.ObjectReaction(o, api.RESTMapper))

	fakePtr.AddWatchReactor("*", core.DefaultWatchReactor(watch.NewFake(), nil))

	return &Clientset{fakePtr}
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	core.Fake
}

var _ clientset.Interface = &Clientset{}

func (c *Clientset) Legacy() legacy_unversioned.LegacyInterface {
	return &legacy_unversioned_fake.FakeLegacy{&c.Fake}
}

func (c *Clientset) Extensions() extensions_unversioned.ExtensionsInterface {
	return &extensions_unversioned_fake.FakeExtensions{&c.Fake}
}

func (c *Clientset) Discovery() unversioned.DiscoveryInterface {
	return &FakeDiscovery{&c.Fake}
}
