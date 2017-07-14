/*
Copyright 2017 Devin All rights reserved.

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

package marathon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const fakePodInstanceName = "fake-pod.instance-dc6cfe60-6812-11e7-a18e-70b3d5800003"

func TestDeletePodInstance(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	podInstance, err := endpoint.Client.DeletePodInstance(fakePodName, fakePodInstanceName)
	assert.NoError(t, err)
	assert.Equal(t, podInstance.InstanceID.ID, fakePodInstanceName)
}

func TestDeletePodInstances(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	instances := []string{fakePodInstanceName}
	podInstances, err := endpoint.Client.DeletePodInstances(fakePodName, instances)
	assert.NoError(t, err)
	assert.Equal(t, podInstances[0].InstanceID.ID, fakePodInstanceName)
}
