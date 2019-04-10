// Copyright (c) 2019 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package taskconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uber/peloton/.gen/peloton/api/v0/peloton"
	"github.com/uber/peloton/.gen/peloton/api/v0/task"
	v1peloton "github.com/uber/peloton/.gen/peloton/api/v1alpha/peloton"
	"github.com/uber/peloton/.gen/peloton/api/v1alpha/pod"
)

// TestIsPelotonLabelChanged checks peloton Label compare
func TestIsPelotonLabelChanged(t *testing.T) {
	l1 := []*peloton.Label{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}
	l2 := []*peloton.Label{
		{Key: "key2", Value: "value2"},
		{Key: "key1", Value: "value1"},
	}
	l3 := []*peloton.Label{
		{Key: "key3", Value: "value3"},
		{Key: "key1", Value: "value1"},
	}
	l4 := []*peloton.Label{
		{Key: "key2", Value: "value3"},
		{Key: "key1", Value: "value1"},
	}
	l5 := []*peloton.Label{
		{Key: "key2", Value: "value2"},
		{Key: "key1", Value: "value1"},
		{Key: "key11", Value: "value11"},
	}

	assert.False(t, IsPelotonLabelChanged(l1, l2))
	assert.True(t, IsPelotonLabelChanged(l1, l3))
	assert.True(t, IsPelotonLabelChanged(l1, l4))
	assert.True(t, IsPelotonLabelChanged(l1, l5))
}

// TestIsPelotonV1LabelChanged checks v1 peloton Label compare
func TestIsPelotonV1LabelChanged(t *testing.T) {
	l1 := []*v1peloton.Label{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}
	l2 := []*v1peloton.Label{
		{Key: "key2", Value: "value2"},
		{Key: "key1", Value: "value1"},
	}
	l3 := []*v1peloton.Label{
		{Key: "key3", Value: "value3"},
		{Key: "key1", Value: "value1"},
	}
	l4 := []*v1peloton.Label{
		{Key: "key2", Value: "value3"},
		{Key: "key1", Value: "value1"},
	}
	l5 := []*v1peloton.Label{
		{Key: "key2", Value: "value2"},
		{Key: "key1", Value: "value1"},
		{Key: "key11", Value: "value11"},
	}

	assert.False(t, IsPelotonV1LabelChanged(l1, l2))
	assert.True(t, IsPelotonV1LabelChanged(l1, l3))
	assert.True(t, IsPelotonV1LabelChanged(l1, l4))
	assert.True(t, IsPelotonV1LabelChanged(l1, l5))
}

// TestIsPortConfigChanged checks PortConfig compare
func TestIsPortConfigChanged(t *testing.T) {
	p1 := []*task.PortConfig{
		{Name: "name1", Value: 1111, EnvName: "env1"},
		{Name: "name2", Value: 2222, EnvName: "env2"},
	}
	p2 := []*task.PortConfig{
		{Name: "name2", Value: 2222, EnvName: "env2"},
		{Name: "name1", Value: 1111, EnvName: "env1"},
	}
	p3 := []*task.PortConfig{
		{Name: "name2", Value: 2222, EnvName: "env2"},
		{Name: "name1", Value: 1111, EnvName: "env1"},
		{Name: "name1", Value: 3333, EnvName: "env1"},
	}
	p4 := []*task.PortConfig{
		{Name: "name3", Value: 2222, EnvName: "env3"},
		{Name: "name1", Value: 1111, EnvName: "env1"},
	}
	p5 := []*task.PortConfig{
		{Name: "name2", Value: 2222, EnvName: "env3"},
		{Name: "name1", Value: 1111, EnvName: "env1"},
	}

	assert.False(t, IsPortConfigChanged(p1, p2))
	assert.True(t, IsPortConfigChanged(p1, p3))
	assert.True(t, IsPortConfigChanged(p1, p4))
	assert.True(t, IsPortConfigChanged(p1, p5))
}

// TestIsPortSpecChanged checks PortSpec compare
func TestIsPortSpecChanged(t *testing.T) {
	p1 := []*pod.PortSpec{
		{Name: "name1", Value: 1111, EnvName: "env1"},
		{Name: "name2", Value: 2222, EnvName: "env2"},
	}
	p2 := []*pod.PortSpec{
		{Name: "name2", Value: 2222, EnvName: "env2"},
		{Name: "name1", Value: 1111, EnvName: "env1"},
	}
	p3 := []*pod.PortSpec{
		{Name: "name2", Value: 2222, EnvName: "env2"},
		{Name: "name1", Value: 1111, EnvName: "env1"},
		{Name: "name1", Value: 3333, EnvName: "env1"},
	}
	p4 := []*pod.PortSpec{
		{Name: "name3", Value: 2222, EnvName: "env3"},
		{Name: "name1", Value: 1111, EnvName: "env1"},
	}
	p5 := []*pod.PortSpec{
		{Name: "name2", Value: 2222, EnvName: "env3"},
		{Name: "name1", Value: 1111, EnvName: "env1"},
	}

	assert.False(t, IsPortSpecChanged(p1, p2))
	assert.True(t, IsPortSpecChanged(p1, p3))
	assert.True(t, IsPortSpecChanged(p1, p4))
	assert.True(t, IsPortSpecChanged(p1, p5))
}

// TestIsTaskConfigChanged checks TaskConfig compare
func TestIsTaskConfigChanged(t *testing.T) {
	t1 := &task.TaskConfig{
		Name: "task-1",
		Labels: []*peloton.Label{
			{Key: "k1", Value: "v1"},
			{Key: "k2", Value: "v2"},
		},
		Ports: []*task.PortConfig{
			{Name: "port-1", Value: 10000},
			{Name: "port-2", Value: 10001},
		},
	}
	t2 := &task.TaskConfig{
		Name: "task-2",
		Labels: []*peloton.Label{
			{Key: "k2", Value: "v2"},
			{Key: "k1", Value: "v1"},
		},
		Ports: []*task.PortConfig{
			{Name: "port-2", Value: 10001},
			{Name: "port-1", Value: 10000},
		},
	}
	t3 := &task.TaskConfig{
		Name: "task-3",
		Labels: []*peloton.Label{
			{Key: "k2", Value: "v2"},
			{Key: "k1", Value: "v1"},
		},
		Ports: []*task.PortConfig{
			{Name: "port-2", Value: 10002},
			{Name: "port-1", Value: 10000},
		},
	}

	assert.False(t, IsTaskConfigChanged(t1, t2))
	assert.True(t, IsTaskConfigChanged(t1, t3))
}

// TestIsContainerSpecChanged checks ContainerSpec compare
func TestIsContainerSpecChanged(t *testing.T) {
	oldContainer := &pod.ContainerSpec{
		Name: "container",
		Ports: []*pod.PortSpec{
			{
				Name:    "name1",
				Value:   1111,
				EnvName: "env1",
			},
			{
				Name:    "name2",
				Value:   2222,
				EnvName: "env2",
			},
		},
	}
	newContainer := &pod.ContainerSpec{
		Name: "container",
		Ports: []*pod.PortSpec{
			{
				Name:    "name2",
				Value:   2222,
				EnvName: "env2",
			},
			{
				Name:    "name1",
				Value:   1111,
				EnvName: "env1",
			},
		},
	}

	assert.False(t, IsContainerSpecChanged(nil, nil))
	assert.True(t, IsContainerSpecChanged(oldContainer, nil))
	assert.True(t, IsContainerSpecChanged(nil, newContainer))
	assert.False(t, IsContainerSpecChanged(oldContainer, newContainer))
}

// TestIsPodSpecChanged checks PodSpec compare
func TestIsPodSpecChanged(t *testing.T) {
	p1 := &pod.PodSpec{
		PodName: &v1peloton.PodName{Value: "pod-1"},
		Labels: []*v1peloton.Label{
			{Key: "k1", Value: "v1"},
			{Key: "k2", Value: "v2"},
		},
		Containers: []*pod.ContainerSpec{
			{Name: "container-1"},
			{Name: "container-2"},
		},
	}
	p2 := &pod.PodSpec{
		PodName: &v1peloton.PodName{Value: "pod-2"},
		Labels: []*v1peloton.Label{
			{Key: "k2", Value: "v2"},
			{Key: "k1", Value: "v1"},
		},
		Containers: []*pod.ContainerSpec{
			{Name: "container-1"},
			{Name: "container-2"},
		},
	}
	p3 := &pod.PodSpec{
		PodName: &v1peloton.PodName{Value: "pod-3"},
		Labels: []*v1peloton.Label{
			{Key: "k2", Value: "v2"},
			{Key: "k1", Value: "v1"},
		},
		Containers: []*pod.ContainerSpec{
			{Name: "container-2"},
			{Name: "container-1"},
		},
	}
	p4 := &pod.PodSpec{
		PodName: &v1peloton.PodName{Value: "pod-3"},
		Labels: []*v1peloton.Label{
			{Key: "k2", Value: "v3"},
			{Key: "k1", Value: "v1"},
		},
		Containers: []*pod.ContainerSpec{
			{Name: "container-2"},
			{Name: "container-1"},
		},
	}
	p5 := &pod.PodSpec{
		PodName: &v1peloton.PodName{Value: "pod-3"},
		Labels: []*v1peloton.Label{
			{Key: "k2", Value: "v2"},
			{Key: "k1", Value: "v1"},
		},
		Containers: []*pod.ContainerSpec{
			{Name: "container-1"},
		},
	}

	assert.False(t, IsPodSpecChanged(p1, p2))
	assert.True(t, IsPodSpecChanged(p1, p3))
	assert.True(t, IsPodSpecChanged(p1, p4))
	assert.True(t, IsPodSpecChanged(p1, p5))

	p6 := &pod.PodSpec{
		PodName: &v1peloton.PodName{Value: "pod-1"},
		Labels: []*v1peloton.Label{
			{Key: "k1", Value: "v1"},
			{Key: "k2", Value: "v2"},
		},
		InitContainers: []*pod.ContainerSpec{
			{Name: "container-1"},
			{Name: "container-2"},
		},
	}
	p7 := &pod.PodSpec{
		PodName: &v1peloton.PodName{Value: "pod-3"},
		Labels: []*v1peloton.Label{
			{Key: "k2", Value: "v2"},
			{Key: "k1", Value: "v1"},
		},
		InitContainers: []*pod.ContainerSpec{
			{Name: "container-2"},
			{Name: "container-1"},
		},
	}
	p8 := &pod.PodSpec{
		PodName: &v1peloton.PodName{Value: "pod-3"},
		Labels: []*v1peloton.Label{
			{Key: "k2", Value: "v2"},
			{Key: "k1", Value: "v1"},
		},
		InitContainers: []*pod.ContainerSpec{
			{Name: "container-1"},
		},
	}

	assert.True(t, IsPodSpecChanged(p6, p7))
	assert.True(t, IsPodSpecChanged(p6, p8))
}
