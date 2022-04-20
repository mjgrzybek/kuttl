package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestIsSubset(t *testing.T) {
	assert.Nil(t, IsSubset(map[string]interface{}{
		"hello": "world",
	}, map[string]interface{}{
		"hello": "world",
		"bye":   "moon",
	}))

	assert.NotNil(t, IsSubset(map[string]interface{}{
		"hello": "moon",
	}, map[string]interface{}{
		"hello": "world",
		"bye":   "moon",
	}))

	assert.Nil(t, IsSubset(map[string]interface{}{
		"hello": map[string]interface{}{
			"hello": "world",
		},
	}, map[string]interface{}{
		"hello": map[string]interface{}{
			"hello": "world",
			"bye":   "moon",
		},
	}))

	assert.NotNil(t, IsSubset(map[string]interface{}{
		"hello": map[string]interface{}{
			"hello": "moon",
		},
	}, map[string]interface{}{
		"hello": map[string]interface{}{
			"hello": "world",
			"bye":   "moon",
		},
	}))

	assert.NotNil(t, IsSubset(map[string]interface{}{
		"hello": map[string]interface{}{
			"hello": "moon",
		},
	}, map[string]interface{}{
		"hello": "world",
	}))

	assert.NotNil(t, IsSubset(map[string]interface{}{
		"hello": "world",
	}, map[string]interface{}{}))

	assert.Nil(t, IsSubset(map[string]interface{}{
		"hello": []int{
			1, 2, 3,
		},
	}, map[string]interface{}{
		"hello": []int{
			1, 2, 3,
		},
	}))

	assert.Nil(t, IsSubset(map[string]interface{}{
		"hello": map[string]interface{}{
			"hello": []map[string]interface{}{
				{
					"image": "hello",
				},
			},
		},
	}, map[string]interface{}{
		"hello": map[string]interface{}{
			"hello": []map[string]interface{}{
				{
					"image": "hello",
					"bye":   "moon",
				},
			},
		},
	}))

	assert.Nil(t, IsSubset(map[string]interface{}{
		"hello": map[string]interface{}{
			"hello": []map[string]interface{}{
				{
					"image": "hello",
				},
			},
		},
	}, map[string]interface{}{
		"hello": map[string]interface{}{
			"hello": []map[string]interface{}{
				{
					"image": "hello",
					"bye":   "moon",
				},
				{
					"bye": "moon",
				},
			},
		},
	}))

	assert.NotNil(t, IsSubset(map[string]interface{}{
		"hello": map[string]interface{}{
			"hello": []map[string]interface{}{
				{
					"image": "hello",
				},
			},
		},
	}, map[string]interface{}{
		"hello": map[string]interface{}{
			"hello": []map[string]interface{}{
				{
					"image": "world",
				},
			},
		},
	}))
}

func TestIsSubset2(t *testing.T) {
	createPodSpecWithEnv := func(e []corev1.EnvVar) corev1.Pod {
		return corev1.Pod{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Env: e,
					},
				},
			},
			Status: corev1.PodStatus{},
		}
	}

	x := []corev1.EnvVar{
		{
			Name:  "a",
			Value: "a",
		},
		{
			Name:  "b",
			Value: "b",
		},
	}

	y := []corev1.EnvVar{
		{
			Name:  "b",
			Value: "b",
		},
		{
			Name:  "a",
			Value: "a",
		},
	}

	z := []corev1.EnvVar{
		{
			Name:  "b",
			Value: "b",
		},
		{
			Name:  "c",
			Value: "c",
		},
		{
			Name:  "a",
			Value: "a",
		},
	}

	podX := createPodSpecWithEnv(x)
	podY := createPodSpecWithEnv(y)
	podZ := createPodSpecWithEnv(z)

	ux, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&podX)
	assert.Nil(t, err)

	uy, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&podY)
	assert.Nil(t, err)

	uz, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&podZ)
	assert.Nil(t, err)

	assert.Nil(t, IsSubset(ux, ux))
	assert.Nil(t, IsSubset(ux, uy))
	assert.Nil(t, IsSubset(uz, uz))
	assert.Nil(t, IsSubset(ux, uz))
	assert.Nil(t, IsSubset(uy, uz))

	err = IsSubset(uz, ux)
	assert.Equal(t, err, &SubsetLengthError{
		SubsetError{
			message: "expected length longer than actual: 3 > 2",
		}},
	)

	err = IsSubset(uz, uy)
	assert.Equal(t, err, &SubsetLengthError{
		SubsetError{
			message: "expected length longer than actual: 3 > 2",
		}},
	)
}
