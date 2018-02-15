package sbxgpu

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	// "math"
	// "math/rand"
	// "strconv"
)

func sbx_drawSkybox(view, projection mgl32.Mat4) {
	/* use after drawing the scene */
	/* Parameters:
	   gl - WebGL context
	   view, projection - gl matrices 4x4 (column major)
	   textureUnit - integer from [0 ... gl.MAX_TEXTURE_IMAGE_UNITS]
	*/
	gl.DepthFunc(gl.LEQUAL)

	gl.UseProgram(sbx_shaderProgram)

	// gl.uniformMatrix4fv(sbx_view, false, view);
	gl.UniformMatrix4fv(sbx_view, 1, false, &(view[0]))
	// gl.uniformMatrix4fv(sbx_projection, false, projection);
	gl.UniformMatrix4fv(sbx_projection, 1, false, &(projection[0]))

	gl.EnableVertexAttribArray(uint32(sbx_position))
	gl.BindBuffer(gl.ARRAY_BUFFER, sbx_arrayBuffer)
	gl.VertexAttribPointer(uint32(sbx_position), 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.ActiveTexture(gl.TEXTURE0 + sbx_textureUnit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, sbx_textureId)
	gl.Uniform1i(sbx_skybox, int32(sbx_textureUnit))
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, sbx_textureId)

	//  gl.drawArrays(gl.TRIANGLES, 0, sbx_Float32Array.length/3 );
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.DepthFunc(gl.LESS)
}
