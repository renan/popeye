package sanitize

import (
	"testing"

	"github.com/derailed/popeye/internal/cache"
	"github.com/derailed/popeye/internal/issues"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func loadCodes(t *testing.T) *issues.Codes {
	codes, err := issues.LoadCodes()
	assert.Nil(t, err)
	return codes
}

func TestConfigMapSanitize(t *testing.T) {
	cm := NewConfigMap(issues.NewCollector(loadCodes(t)), newConfigMap())
	cm.Sanitize(nil)

	assert.Equal(t, 4, len(cm.Outcome()))

	ii := cm.Outcome()["default/cm3"]
	assert.Equal(t, 1, len(ii))
	assert.Equal(t, "[POP-400] Used? Unable to locate resource reference", ii[0].Message)
	assert.Equal(t, issues.InfoLevel, ii[0].Level)

	ii = cm.Outcome()["default/cm4"]
	assert.Equal(t, 1, len(ii))
	assert.Equal(t, `[POP-401] Key "k2" used? Unable to locate key reference`, ii[0].Message)
	assert.Equal(t, issues.InfoLevel, ii[0].Level)
}

// ----------------------------------------------------------------------------
// Helpers...

type configMap struct{}

func newConfigMap() configMap {
	return configMap{}
}

func (c configMap) PodRefs(refs cache.ObjReferences) {
	refs["cm:default/cm1"] = cache.StringSet{
		"k1": cache.Blank,
		"k2": cache.Blank,
	}
	refs["cm:default/cm2"] = cache.StringSet{
		cache.AllKeys: cache.Blank,
	}
	refs["cm:default/cm4"] = cache.StringSet{
		"k1": cache.Blank,
	}
}

func (c configMap) ListConfigMaps() map[string]*v1.ConfigMap {
	return map[string]*v1.ConfigMap{
		"default/cm1": makeConfigMap("cm1"),
		"default/cm2": makeConfigMap("cm2"),
		"default/cm3": makeConfigMap("cm3"),
		"default/cm4": makeConfigMap("cm4"),
	}
}

func makeConfigMap(n string) *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n,
			Namespace: "default",
		},
		Data: map[string]string{
			"k1": "",
			"k2": "",
		},
	}
}
