// This file is part of MinIO DirectPV
// Copyright (c) 2021, 2022 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Code generated by client-gen. DO NOT EDIT.

package clientset

import (
	"fmt"

	directv1alpha1 "github.com/minio/directpv/pkg/clientset/typed/direct.csi.min.io/v1alpha1"
	directv1beta1 "github.com/minio/directpv/pkg/clientset/typed/direct.csi.min.io/v1beta1"
	directv1beta2 "github.com/minio/directpv/pkg/clientset/typed/direct.csi.min.io/v1beta2"
	directv1beta3 "github.com/minio/directpv/pkg/clientset/typed/direct.csi.min.io/v1beta3"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	DirectV1alpha1() directv1alpha1.DirectV1alpha1Interface
	DirectV1beta1() directv1beta1.DirectV1beta1Interface
	DirectV1beta2() directv1beta2.DirectV1beta2Interface
	DirectV1beta3() directv1beta3.DirectV1beta3Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	directV1alpha1 *directv1alpha1.DirectV1alpha1Client
	directV1beta1  *directv1beta1.DirectV1beta1Client
	directV1beta2  *directv1beta2.DirectV1beta2Client
	directV1beta3  *directv1beta3.DirectV1beta3Client
}

// DirectV1alpha1 retrieves the DirectV1alpha1Client
func (c *Clientset) DirectV1alpha1() directv1alpha1.DirectV1alpha1Interface {
	return c.directV1alpha1
}

// DirectV1beta1 retrieves the DirectV1beta1Client
func (c *Clientset) DirectV1beta1() directv1beta1.DirectV1beta1Interface {
	return c.directV1beta1
}

// DirectV1beta2 retrieves the DirectV1beta2Client
func (c *Clientset) DirectV1beta2() directv1beta2.DirectV1beta2Interface {
	return c.directV1beta2
}

// DirectV1beta3 retrieves the DirectV1beta3Client
func (c *Clientset) DirectV1beta3() directv1beta3.DirectV1beta3Interface {
	return c.directV1beta3
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.directV1alpha1, err = directv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.directV1beta1, err = directv1beta1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.directV1beta2, err = directv1beta2.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.directV1beta3, err = directv1beta3.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.directV1alpha1 = directv1alpha1.NewForConfigOrDie(c)
	cs.directV1beta1 = directv1beta1.NewForConfigOrDie(c)
	cs.directV1beta2 = directv1beta2.NewForConfigOrDie(c)
	cs.directV1beta3 = directv1beta3.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.directV1alpha1 = directv1alpha1.New(c)
	cs.directV1beta1 = directv1beta1.New(c)
	cs.directV1beta2 = directv1beta2.New(c)
	cs.directV1beta3 = directv1beta3.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
