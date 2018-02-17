package sbxgpu

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

/* shaders - see: http://learnopengl.com/#!Advanced-OpenGL/Cubemaps */

var sbx_vertexShaderSource = `
#version 330
in vec3 position;
out vec3 TexCoords;
uniform mat4 projection;
uniform mat4 view;
void main()
{
    vec4 pos = projection * view * vec4(position, 1.0);
    gl_Position = pos.xyww;
    TexCoords = position;
}
` + "\x00"

var sbx_fragmentShaderSource = `
#version 330
in vec3 TexCoords;
out vec4 gl_FragColor
uniform samplerCube skybox;
void main()
{
    gl_FragColor = textureCube(skybox, TexCoords);
}
` + "\x00"

/* shaders */
var sbx_vertexShader uint32
var sbx_fragmentShader uint32

/* shader program */
var sbx_shaderProgram uint32

/* VAO */
var sbx_VAO uint32

/* vertex attributes locations */
// var sbx_position int32

/* uniform variables locations */
var sbx_projection int32
var sbx_view int32
var sbx_skybox int32

/* input vertices of cube triangles */
var sbx_Float32Array = [...]float32{
	-1, 1, -1,
	-1, -1, -1,
	+1, -1, -1,
	+1, -1, -1,
	+1, 1, -1,
	-1, 1, -1,
	-1, -1, 1,
	-1, -1, -1,
	-1, 1, -1,
	-1, 1, -1,
	-1, 1, 1,
	-1, -1, 1,
	+1, -1, -1,
	+1, -1, 1,
	+1, 1, 1,
	+1, 1, 1,
	+1, 1, -1,
	+1, -1, -1,
	-1, -1, 1,
	-1, 1, 1,
	+1, 1, 1,
	+1, 1, 1,
	+1, -1, 1,
	-1, -1, 1,
	-1, 1, -1,
	+1, 1, -1,
	+1, 1, 1,
	+1, 1, 1,
	-1, 1, 1,
	-1, 1, -1,
	-1, -1, -1,
	-1, -1, 1,
	+1, -1, -1,
	+1, -1, -1,
	-1, -1, 1,
	+1, -1, 1,
}

var sbx_arrayBuffer uint32

func (sbx *SbxGpu) makeShaderProgram() {
	/* Parameters:
	   gl - WebGL context
	*/

	sbx_shaderProgram, err := newProgram(sbx_vertexShaderSource, sbx_fragmentShaderSource)

	if err != nil {
		panic(err)
	}

	gl.UseProgram(sbx_shaderProgram)

	/* set vertex attributes locations */
	sbx.position = gl.GetAttribLocation(sbx_shaderProgram, gl.Str("position\x00"))

	/* set uniform variables locations */
	sbx_projection = gl.GetUniformLocation(sbx_shaderProgram, gl.Str("projection\x00"))
	sbx_view = gl.GetUniformLocation(sbx_shaderProgram, gl.Str("view\x00"))
	sbx_skybox = gl.GetUniformLocation(sbx_shaderProgram, gl.Str("skybox\x00"))

	/* load buffer data */
	gl.GenBuffers(1, &sbx_arrayBuffer)

	gl.BindBuffer(gl.ARRAY_BUFFER, sbx_arrayBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(sbx_Float32Array)*4 /* 4 bytes per flat32 */, gl.Ptr(sbx_Float32Array), gl.STATIC_DRAW)

	/* init VAO */
	gl.GenVertexArrays(1, &sbx.VAO)
	gl.BindVertexArray(sbx.VAO)
	gl.EnableVertexAttribArray(uint32(sbx.position))
	gl.BindBuffer(gl.ARRAY_BUFFER, sbx_arrayBuffer)
	gl.VertexAttribPointer(uint32(sbx.position), 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.BindVertexArray(0) // unbind VAO
	sbx.shaderProgram = sbx_shaderProgram

}
