package mutate

import (
	"testing"

	"github.com/qiniu/logkit/sender"
	"github.com/qiniu/logkit/transforms"
	"github.com/qiniu/logkit/utils"
	"github.com/stretchr/testify/assert"
)

func TestParamTransformer(t *testing.T) {
	par := &UrlParam{
		Key: "myword",
	}
	data, err := par.Transform([]sender.Data{
		{"myword": "platform=2&vid=372&vu=caea966558&chan=android_sougou&sign=ad225ec02942c79bdb710e3ad0cf1b43&nonce_str=1510555032"},
		{"myword": "platform=2&vid=&vu=caea966558&chan=&sign=ad225ec02942c79bdb710e3ad0cf1b43&nonce_str=1510555032"},
	})
	assert.NoError(t, err)
	exp := []sender.Data{
		{
			"myword":           "platform=2&vid=372&vu=caea966558&chan=android_sougou&sign=ad225ec02942c79bdb710e3ad0cf1b43&nonce_str=1510555032",
			"myword_platform":  "2",
			"myword_vid":       "372",
			"myword_vu":        "caea966558",
			"myword_chan":      "android_sougou",
			"myword_sign":      "ad225ec02942c79bdb710e3ad0cf1b43",
			"myword_nonce_str": "1510555032",
		},
		{
			"myword":           "platform=2&vid=&vu=caea966558&chan=&sign=ad225ec02942c79bdb710e3ad0cf1b43&nonce_str=1510555032",
			"myword_platform":  "2",
			"myword_vid":       "",
			"myword_vu":        "caea966558",
			"myword_chan":      "",
			"myword_sign":      "ad225ec02942c79bdb710e3ad0cf1b43",
			"myword_nonce_str": "1510555032",
		},
	}
	assert.Equal(t, len(exp), len(data))
	for i, d := range data {
		ex := exp[i]
		for k, v := range d {
			e, exist := ex[k]
			assert.Equal(t, true, exist)
			assert.Equal(t, e, v)
		}
	}
	assert.Equal(t, par.Stage(), transforms.StageAfterParser)
	assert.Equal(t, utils.StatsInfo{Success: 2}, par.stats)
}

func TestParamTransformerError(t *testing.T) {
	par := &UrlParam{
		Key: "myword",
	}
	data, err := par.Transform([]sender.Data{
		{"myword": "platform=2=372&vu=caea966558&chan=android_sougou&sign=ad225ec02942c79bdb710e3ad0cf1b43&nonce_str=1510555032"},
		{"myword": "platform=2&vid&vu=caea966558&chan=&sign=ad225ec02942c79bdb710e3ad0cf1b43&nonce_str=1510555032"},
	})
	assert.Error(t, err)
	exp := []sender.Data{
		{"myword": "platform=2=372&vu=caea966558&chan=android_sougou&sign=ad225ec02942c79bdb710e3ad0cf1b43&nonce_str=1510555032"},
		{"myword": "platform=2&vid&vu=caea966558&chan=&sign=ad225ec02942c79bdb710e3ad0cf1b43&nonce_str=1510555032"},
	}
	assert.Equal(t, exp, data)
	assert.Equal(t, par.Stage(), transforms.StageAfterParser)

	par.stats.LastError = ""
	assert.Equal(t, utils.StatsInfo{Errors: 2}, par.stats)
}

func TestParamTransformerKeyRepeat(t *testing.T) {
	par := &UrlParam{
		Key: "myword",
	}
	data, err := par.Transform([]sender.Data{
		{"myword": "a=a&a=b&a=c&a=d"},
		{"myword": "a=a&a=b&b=c&b=d&b=e"},
	})
	assert.NoError(t, err)
	exp := []sender.Data{
		{
			"myword":    "a=a&a=b&a=c&a=d",
			"myword_a":  "a",
			"myword_a1": "b",
			"myword_a2": "c",
			"myword_a3": "d",
		},
		{
			"myword":    "a=a&a=b&b=c&b=d&b=e",
			"myword_a":  "a",
			"myword_a1": "b",
			"myword_b":  "c",
			"myword_b1": "d",
			"myword_b2": "e",
		},
	}
	assert.Equal(t, exp, data)
	assert.Equal(t, par.Stage(), transforms.StageAfterParser)
	assert.Equal(t, utils.StatsInfo{Success: 2}, par.stats)
}
