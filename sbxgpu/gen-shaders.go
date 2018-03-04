package sbxgpu

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"math"
	"math/rand"
	"strconv"
)

const sbx_CUBE_SIZE = 1024

var sbx_srcCubeSize = "const int cubeSize= " + strconv.Itoa(sbx_CUBE_SIZE) + ";\n"
var sbx_srcPI = "const float PI = " + strconv.FormatFloat(math.Pi, 'f', -1, 64) + ";\n"
var sbx_srcFunRPrefix = "float fR(float x,float y,float z){ return "
var sbx_srcFunGPrefix = "float fG(float x,float y,float z){ return "
var sbx_srcFunBPrefix = "float fB(float x,float y,float z){ return "
var sbx_srcFunSuffix = "; }\n"

// collection of color intensity functions
var sbx_srcFunStrings = [...]string{
	" x",
	" y",
	" z",
	" sin( x * PI * 2.0 )",
	" sin( z * PI * 2.0 )",
	" sin( y * PI * 2.0 )",
	" cos( x * PI * 2.0 )",
	" cos( z * PI * 2.0 )",
	" cos( y * PI * 2.0 )",
	" sin( x * PI * 2.0 )*cos( y * PI * 2.0 )",
	" sin( x * PI * 2.0 )*cos( z * PI * 2.0 )",
	" sin( x * PI * 4.0 )",
	" sin( z * PI * 4.0 )",
	" sin( y * PI * 4.0 )",
	" cos( x * PI * 4.0 )",
	" cos( z * PI * 4.0 )",
	" cos( y * PI * 4.0 )",
	" sin( x * PI * 4.0 )*cos( y * PI * 4.0 )",
	" sin( x * PI * 4.0 )*cos( z * PI * 4.0 )",
}

// prepend constant definitions fR, fG, fB
var sbx_renderTextureVS2 = ` 

layout (location = 0) in float h;
uniform float v;
const float depth=1.0;
uniform mat3 xyz;
out vec4 color;

void main()
{
  float  args[6];
  float h=h-float(cubeSize)/2.0;
  float v=v-float(cubeSize)/2.0;
  float d=depth*float(cubeSize)-float(cubeSize)/2.0;
  vec3 hvd= vec3(h,v,d);
  vec3 vxyz=normalize(xyz*hvd);
  float x=vxyz.x;
  float y=vxyz.y;
  float z=vxyz.z;
  color= 0.5*vec4( fR(x,y,z), fG(x,y,z), fB(x,y,z), 1.0 )+vec4(0.5,0.5,0.5,0.5);
  gl_Position = vec4( h/float(cubeSize), v/float(cubeSize), 0.0, 0.5 ); /// w=0.5 for perspective division
  gl_PointSize=1.0; /// test it
}
` + "\x00"

var sbx_renderTextureFS = `
#version 330

in vec4 color;

void main()
{
  gl_FragColor= color;
}
` + "\x00"

var sbx_hBufferId uint32 // array: [0,1, ..., sbx_CUBE_SIZE-1]
var sbx_hBufferIdExists = false

func (sbx *SbxGpu) makeRenderTextureShaderProgram() {
	var fun = sbx_srcFunStrings
	var r = rand.Intn(len(fun))
	var g = rand.Intn(len(fun))
	var b = rand.Intn(len(fun))

	var sbx_srcFunR = sbx_srcFunRPrefix + sbx_srcFunStrings[r] + sbx_srcFunSuffix
	var sbx_srcFunG = sbx_srcFunGPrefix + sbx_srcFunStrings[g] + sbx_srcFunSuffix
	var sbx_srcFunB = sbx_srcFunBPrefix + sbx_srcFunStrings[b] + sbx_srcFunSuffix

	sbx_renderTextureVS := "#version 330\n\n" +
		sbx_srcCubeSize +
		sbx_srcPI +
		sbx_srcFunR +
		sbx_srcFunG +
		sbx_srcFunB +
		sbx_renderTextureVS2

	gl.DeleteVertexArrays(1, &sbx.renderTextureVAO)  // delete old VAO
	gl.DeleteProgram(sbx.renderTextureShaderProgram) // delete old program if exists

	sbx_renderTextureShaderProgram, err := newProgram(sbx_renderTextureVS, sbx_renderTextureFS)

	if err != nil {
		panic(err)
	}

	sbx.renderTextureShaderProgram = sbx_renderTextureShaderProgram // exists

	/* set vertex attributes locations */
	sbx.hLocation = gl.GetAttribLocation(sbx.renderTextureShaderProgram, gl.Str("h\x00"))
	if sbx.hLocation < 0 {
		panic("sbx.hLocation=" + strconv.Itoa(int(sbx.hLocation)))
	}

	/* set uniform variables locations */
	sbx.vLocation = gl.GetUniformLocation(sbx.renderTextureShaderProgram, gl.Str("v\x00"))
	if sbx.vLocation < 0 {
		panic("sbx.vLocation=" + strconv.Itoa(int(sbx.vLocation)))
	}

	sbx.xyzLocation = gl.GetUniformLocation(sbx.renderTextureShaderProgram, gl.Str("xyz\x00"))
	if sbx.xyzLocation < 0 {
		panic("sbx.xyzLocation=" + strconv.Itoa(int(sbx.xyzLocation)))
	}

	/* load buffer data */
	if sbx_hBufferIdExists == false {
		gl.GenBuffers(1, &sbx_hBufferId)
		gl.BindBuffer(gl.ARRAY_BUFFER, sbx_hBufferId)

		var hIn [sbx_CUBE_SIZE + 4]float32
		for i := 0; i < sbx_CUBE_SIZE+4; i++ {
			hIn[i] = float32(i - 2)
		}
		gl.BufferData(gl.ARRAY_BUFFER, len(hIn)*4 /* 4 bytes per flat32 */, gl.Ptr(&hIn[0]), gl.STATIC_DRAW)
		sbx_hBufferIdExists = true
	}

	/* init VAO */
	gl.UseProgram(sbx.renderTextureShaderProgram)
	gl.GenVertexArrays(1, &sbx.renderTextureVAO)
	gl.BindVertexArray(sbx.renderTextureVAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, sbx_hBufferId)
	gl.EnableVertexAttribArray(uint32(sbx.hLocation))
	gl.VertexAttribPointer(uint32(sbx.hLocation), 1, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.BindVertexArray(0) // unbind VAO
}
