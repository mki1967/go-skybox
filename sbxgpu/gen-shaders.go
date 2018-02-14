package sbxgpu

import (
	"strconv"
	"math"
	"math/rand"
	"github.com/go-gl/gl/v3.3-core/gl"

)

const sbx_CUBE_SIZE= 1024;

var sbx_srcCubeSize= "const int cubeSize= " + strconv.Itoa(sbx_CUBE_SIZE) +";\n";
var sbx_srcPI= "const float PI = " +strconv.FormatFloat(math.Pi, 'f', -1,64) +";\n";
var sbx_srcFunRPrefix= "float fR(float x,float y,float z){ return ";
var sbx_srcFunGPrefix= "float fG(float x,float y,float z){ return ";
var sbx_srcFunBPrefix= "float fB(float x,float y,float z){ return ";
var sbr_srcFunSuffix= "; }\n"

// collection of color intensity functions
var sbx_srcFunStrings= [...] string{
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
};

var sbx_renderTextureVS2=""+ // prepend constant definitions fR, fG, fB
	"attribute float h;\n"+
	"uniform float v;\n"+
	"const float depth=1.0;\n"+
	"uniform mat3 xyz;\n"+
	"varying vec4 color;\n"+
	"void main()\n"+
	"{\n"+
	"  float  args[6];\n"+
	"  float h=h-float(cubeSize)/2.0;\n"+
	"  float v=v-float(cubeSize)/2.0;\n"+
	"  float d=depth*float(cubeSize)-float(cubeSize)/2.0;\n"+
	"  vec3 hvd= vec3(h,v,d);\n"+
	"  vec3 vxyz=normalize(xyz*hvd);\n"+
	"  float x=vxyz.x;\n"+
	"  float y=vxyz.y;\n"+
	"  float z=vxyz.z;\n"+
	"  color= 0.5*vec4( fR(x,y,z), fG(x,y,z), fB(x,y,z), 1.0 )+vec4(0.5,0.5,0.5,0.5);\n"+
	"  gl_Position = vec4( h/float(cubeSize), v/float(cubeSize), 0.0, 0.5 );\n"+ /// w=0.5 for perspective division
	"  gl_PointSize=1.0;\n"+ /// test it
	"}\n";

var sbx_renderTextureFS=""+
	"precision mediump float;\n"+
	"varying vec4 color;\n"+
	"void main()\n"+
	"{\n"+
	"  gl_FragColor= color;\n"+
	"}\n";

/* to be created by sbx_makeRenderTextureShaderProgram */
var sbx_renderTextureVS string
var sbx_renderTextureShaderProgram uint32
var  sbx_hBufferId uint32 // array: [0,1, ..., sbx_CUBE_SIZE-1]

var sbx_hLocation int32
var sbx_vLocation int32
var sbx_xyzLocation int32

func sbx_makeRenderTextureShaderProgram() {
	var fun=sbx_srcFunStrings;
	var r=rand.Intn( len(fun) )
	var g=rand.Intn( len(fun) )
	var b=rand.Intn( len(fun) )

	var sbx_srcFunR = sbx_srcFunRPrefix + sbx_srcFunStrings[r]+sbr_srcFunSuffix;
	var sbx_srcFunG = sbx_srcFunGPrefix + sbx_srcFunStrings[g]+sbr_srcFunSuffix;
	var sbx_srcFunB = sbx_srcFunBPrefix + sbx_srcFunStrings[b]+sbr_srcFunSuffix;

	sbx_renderTextureVS= 
		sbx_srcCubeSize + 
		sbx_srcPI +
		sbx_srcFunR +
		sbx_srcFunG +
		sbx_srcFunB +
		sbx_renderTextureVS2;
	// console.log(sbx_renderTextureVS); // tests

	gl.DeleteProgram( sbx_renderTextureShaderProgram ); // delete old program if exists

	// sbx_renderTextureShaderProgram=  sbx_makeShaderProgramTool(gl, sbx_renderTextureVS , sbx_renderTextureFS )
	sbx_renderTextureShaderProgram, err := newProgram(sbx_renderTextureVS, sbx_renderTextureFS)

	if err != nil {
		panic(err)
	}
	
	/* set vertex attributes locations */
	sbx_hLocation = gl.GetAttribLocation(sbx_renderTextureShaderProgram, gl.Str("h\x00"));
	if sbx_hLocation< 0 {
		panic("sbx_hLocation="+strconv.Itoa(int(sbx_hLocation)))
	}

	/* set uniform variables locations */
	sbx_vLocation = gl.GetUniformLocation(sbx_renderTextureShaderProgram,  gl.Str("v\x00"));
	if sbx_vLocation< 0 {
		panic("sbx_vLocation="+strconv.Itoa(int(sbx_vLocation)))
	}

	sbx_xyzLocation = gl.GetUniformLocation(sbx_renderTextureShaderProgram,  gl.Str("xyz\x00"));
	if sbx_xyzLocation< 0 {
		panic("sbx_xyzLocation="+strconv.Itoa(int(sbx_xyzLocation)))
	}

	/* load buffer data */
	if( sbx_hBufferId== 0 ) {
		gl.GenBuffers(1, &sbx_hBufferId)
		// sbx_hBufferId= gl.CreateBuffer();
		gl.BindBuffer(gl.ARRAY_BUFFER, sbx_hBufferId );

		var hIn [sbx_CUBE_SIZE+4]float32
		for i:=0; i< sbx_CUBE_SIZE+4; i++ {
			hIn[i]= float32(i-2);
		}
		// gl.BufferData(gl.ARRAY_BUFFER, new Float32Array( hIn ) , gl.STATIC_DRAW );
		gl.BufferData(gl.ARRAY_BUFFER, len(hIn)*4 /* 4 bytes per flat32 */, gl.Ptr(hIn), gl.STATIC_DRAW)
	}

};
