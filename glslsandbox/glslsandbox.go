package glslsandbox

import (
	"regexp"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/polyfloyd/shady"
)

func init() {
	glsl.RegisterEnvironmentDetector(func(shaderSource string) string {
		// Quick and dirty: run some regular expressions on the source to infer
		// the environment.
		reGLSLSandbox := regexp.MustCompile(`uniform\s+vec2\s+resolution`)
		if reGLSLSandbox.MatchString(shaderSource) {
			return "glslsandbox"
		}
		return ""
	})
}

// GLSLSandbox implements the Environment interface to simulate the canvas of glslsandbox.com.
type GLSLSandbox struct {
	// Source is the single fragment shader that should be used for rendering.
	ShaderSources []glsl.Source
}

func (gs GLSLSandbox) Sources() (map[glsl.Stage][]glsl.Source, error) {
	return map[glsl.Stage][]glsl.Source{
		glsl.StageVertex: {glsl.SourceBuf(`
			attribute vec3 vert;
			varying vec2 surfacePosition;

			void main(void) {
				surfacePosition = vert.xy;
				gl_Position = vec4(vert, 1.0);
			}
		`)},
		glsl.StageFragment: gs.ShaderSources,
	}, nil
}

func (GLSLSandbox) Setup(glsl.RenderState) error { return nil }

func (GLSLSandbox) PreRender(state glsl.RenderState) {
	if loc, ok := state.Uniforms["resolution"]; ok {
		gl.Uniform2f(loc.Location, float32(state.CanvasWidth), float32(state.CanvasHeight))
	}
	if loc, ok := state.Uniforms["time"]; ok {
		gl.Uniform1f(loc.Location, float32(state.Time)/float32(time.Second))
	}
	if loc, ok := state.Uniforms["mouse"]; ok {
		gl.Uniform2f(loc.Location, float32(state.CanvasWidth)*0.5, float32(state.CanvasHeight)*0.5)
	}
	if loc, ok := state.Uniforms["surfaceSize"]; ok {
		gl.Uniform2f(loc.Location, float32(state.CanvasWidth), float32(state.CanvasHeight))
	}
	if loc, ok := state.Uniforms["backbuffer"]; ok {
		gl.BindTexture(gl.TEXTURE_2D, state.PreviousFrameTexID())
		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(loc.Location, 0)
	}
}

func (GLSLSandbox) Close() error { return nil }
