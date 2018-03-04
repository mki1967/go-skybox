package sbxgpu

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// DrawSkybox draws the generated skybox in the background.
// Use DrawSkybox after drawing the scene.
//
// Parameters:
//
// view, projection - gl matrices 4x4 (column major)
//
// textureUnit - integer from [0 ... gl.MAX_TEXTURE_IMAGE_UNITS]
func (sbx *SbxGpu) DrawSkybox(view, projection mgl32.Mat4) {

	var depthTest bool // previous depth test
	gl.GetBooleanv(gl.DEPTH_TEST, &depthTest)
	gl.Enable(gl.DEPTH_TEST)
	var depthFunc int32 // previous depth function
	gl.GetIntegerv(gl.DEPTH_FUNC, &depthFunc)
	gl.DepthFunc(gl.LEQUAL)

	gl.UseProgram(sbx.shaderProgram)

	gl.UniformMatrix4fv(sbx.view, 1, false, &(view[0]))
	gl.UniformMatrix4fv(sbx.projection, 1, false, &(projection[0]))

	gl.BindVertexArray(sbx.vao)

	gl.ActiveTexture(gl.TEXTURE0 + sbx.TextureUnit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, sbx.textureId)
	gl.Uniform1i(sbx.skybox, int32(sbx.TextureUnit))
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, sbx.textureId)

	//  gl.drawArrays(gl.TRIANGLES, 0, sbx_Float32Array.length/3 );
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.DepthFunc(uint32(depthFunc))
	if !depthTest {
		gl.Disable(gl.DEPTH_TEST)
	}
	gl.BindVertexArray(0) // unbind vao

}
