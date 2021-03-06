// +build kinect

package shadertoy

import (
	"github.com/polyfloyd/shady/renderer"
	"github.com/polyfloyd/shady/shadertoy/kinect"
)

func init() {
	resourceBuilders["kinect"] = func(m Mapping, pwd string, texIndexEnum *uint32, state renderer.RenderState) (resource, error) {
		kin, err := kinect.Open(m.Name, *texIndexEnum)
		if err != nil {
			return nil, err
		}
		*texIndexEnum++
		return kin, nil
	}
}
