package marathon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const key = "testKey"
const val = "testValue"

const fakePodName = "/fake-pod"

func TestPodLabels(t *testing.T) {
	pod := NewPod()
	pod.AddLabel(key, val)
	assert.Equal(t, pod.Labels[key], val)

	pod.EmptyLabels()
	assert.Equal(t, len(pod.Labels), 0)
}

func TestPodEnvironmentVars(t *testing.T) {
	pod := NewPod()
	pod.AddEnvironment(key, val)

	newVal, err := pod.GetEnvironmentVariable(key)
	assert.Equal(t, newVal, val)
	assert.Equal(t, err, nil)

	badVal, err := pod.GetEnvironmentVariable("fakeKey")
	assert.Equal(t, badVal, "")
	assert.NotNil(t, err)

	pod.EmptyEnvironment()
	assert.Equal(t, len(pod.Environment), 0)
}

func TestSecrets(t *testing.T) {
	pod := NewPod()
	pod.AddSecret(key, val)

	newVal, err := pod.GetSecretSource(key)
	assert.Equal(t, newVal, val)
	assert.Equal(t, err, nil)

	badVal, err := pod.GetSecretSource("fakeKey")
	assert.Equal(t, badVal, "")
	assert.NotNil(t, err)

	pod.EmptySecrets()
	assert.Equal(t, len(pod.Environment), 0)
}

func TestCreatePod(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	pod := NewPod().Name(fakePodName)
	pod, err := endpoint.Client.CreatePod(pod)
	assert.NoError(t, err)
	assert.NotNil(t, pod)
	assert.Equal(t, pod.ID, fakePodName)
}

// TODO

// /v2/pods/{id} PUT
//		in: Pod
//		out: Pod
// /v2/pods/{id} GET
//		Pod
// /v2/pods/{id} DELETE
//		out: code, header "Marathon-Deployment-Id" (for others too!)

// /v2/pods GET
// 		[]Pod
// /v2/pods/::status GET
//		[]PodStatus
// /v2/pods/{id}::status GET
//		PodStatus
// /v2/pods/{id}::versions GET
//		[]string
// /v2/pods/{id}::versions/{version} GET
//		Pod

// these do not give deployment id
// /v2/pods/{id}::instances DELETE
// 		[]PodInstanceStatus
// /v2/pods/{id}::instances/{instance} DELETE
//		PodInstanceStatus
