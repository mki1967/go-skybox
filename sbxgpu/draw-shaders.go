package sbxgpu

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	// "math"
	// "math/rand"
	// "strconv"
)

/* shaders - see: http://learnopengl.com/#!Advanced-OpenGL/Cubemaps */

var sbx_vertexShaderSource = "" +
	"attribute vec3 position;\n" +
	"varying vec3 TexCoords;\n" +
	"uniform mat4 projection;\n" +
	"uniform mat4 view;\n" +
	"void main()\n" +
	"{\n" +
	"    vec4 pos = projection * view * vec4(position, 1.0);\n" +
	"    gl_Position = pos.xyww;\n" +
	//    "    gl_Position = vec4(pos.xy, 1.0,1.0);\n"+
	"    TexCoords = position;\n" +
	"}\n"

var sbx_fragmentShaderSource = "" +
	"precision mediump float;\n" +
	"varying vec3 TexCoords;\n" +
	"uniform samplerCube skybox;\n" +
	"void main()\n" +
	"{\n" +
	"    gl_FragColor = textureCube(skybox, TexCoords);\n" +
	//    "    gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);\n"+
	"}\n"

/* shaders */
var sbx_vertexShader uint32
var sbx_fragmentShader uint32

/* shader program */
var sbx_shaderProgram uint32

/* vertex attributes locations */
var sbx_position int32

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

func sbx_makeShaderProgram() {
	/* Parameters:
	   gl - WebGL context
	*/

	// sbx_shaderProgram= sbx_makeShaderProgramTool(gl,  sbx_vertexShaderSource,  sbx_fragmentShaderSource);
	sbx_shaderProgram, err := newProgram(sbx_vertexShaderSource, sbx_fragmentShaderSource)

	if err != nil {
		panic(err)
	}

	gl.UseProgram(sbx_shaderProgram)

	/* set vertex attributes locations */
	sbx_position = gl.GetAttribLocation(sbx_shaderProgram, gl.Str("position\x00"))

	/* set uniform variables locations */
	sbx_projection = gl.GetUniformLocation(sbx_shaderProgram, gl.Str("projection\x00"))
	sbx_view = gl.GetUniformLocation(sbx_shaderProgram, gl.Str("view\x00"))
	sbx_skybox = gl.GetUniformLocation(sbx_shaderProgram, gl.Str("skybox\x00"))

	/* load buffer data */
	// sbx_arrayBuffer= gl.CreateBuffer();
	gl.GenBuffers(1, &sbx_arrayBuffer)

	gl.BindBuffer(gl.ARRAY_BUFFER, sbx_arrayBuffer)
	// gl.BufferData(gl.ARRAY_BUFFER, sbx_Float32Array , gl.STATIC_DRAW );
	gl.BufferData(gl.ARRAY_BUFFER, len(sbx_Float32Array)*4 /* 4 bytes per flat32 */, gl.Ptr(sbx_Float32Array), gl.STATIC_DRAW)

	// SUCCESS
	// return sbx_shaderProgram;
}
