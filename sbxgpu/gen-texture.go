package sbxgpu

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	// "math"
	// "math/rand"
	// "strconv"
)

/* texture parameters */
var sbx_textureId uint32
var sbx_textureIdExists = false
var sbx_textureUnit uint32 = 0 // default
// var sbx_textureSize=1024;
var sbx_frameBufferId uint32
var sbx_frameBufferIdExists = false

/* arguments permutations */
var sbx_xyzXPlus = [9]float32{
	0, 0, -1,
	0, 1, 0,
	1, 0, 0} //[z,y,-x];
var sbx_xyzXMinus = [9]float32{
	0, 0, 1,
	0, 1, 0,
	-1, 0, 0} //[-z,y,x];

var sbx_xyzYPlus = [9]float32{
	1, 0, 0,
	0, 0, 1,
	0, -1, 0} //[x,-z,y];

var sbx_xyzYMinus = [9]float32{
	1, 0, 0,
	0, 0, -1,
	0, 1, 0} //[x,z,-y];

var sbx_xyzZPlus = [9]float32{
	1, 0, 0,
	0, 1, 0,
	0, 0, 1} // [x,y,z];
var sbx_xyzZMinus = [9]float32{
	-1, 0, 0,
	0, 1, 0,
	0, 0, -1} // [-x,y,-z];

var sbx_xyzArray = [6][9]float32{
	sbx_xyzXPlus,
	sbx_xyzXMinus,
	sbx_xyzYPlus,
	sbx_xyzYMinus,
	sbx_xyzZPlus,
	sbx_xyzZMinus,
}

/* rendering random skybox in a frame */
func sbx_renderRandomCube() {
	// var i,j;
	var defaultFBO int32
	gl.GetIntegerv(gl.FRAMEBUFFER_BINDING, &defaultFBO)

	// var hth= gl.ViewportHeight;
	// 	var wth=  gl.ViewportWidth;

	var viewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &viewport[0]) // save viewport parameters

	if sbx_textureIdExists == false {
		/* create texture object and allocate image memories */
		// sbx_textureId=gl.CreateTexture();
		gl.GenTextures(1, &sbx_textureId)
		sbx_textureIdExists = true
		gl.ActiveTexture(gl.TEXTURE0 + sbx_textureUnit)
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, sbx_textureId)
		for i := 0; i < 6; i++ {
			gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.RGBA, sbx_CUBE_SIZE, sbx_CUBE_SIZE, 0, /* border */
				gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
		}
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}

	if sbx_frameBufferIdExists == false {
		/* create framebuffer object */
		// sbx_frameBufferId=gl.CreateFramebuffer();
		gl.GenFramebuffers(1, &sbx_frameBufferId)
		sbx_frameBufferIdExists = true
	}
	gl.DeleteProgram(sbx_renderTextureShaderProgram) // delete old
	sbx_makeRenderTextureShaderProgram()
	gl.UseProgram(sbx_renderTextureShaderProgram)

	gl.ActiveTexture(gl.TEXTURE0 + sbx_textureUnit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, sbx_textureId)

	gl.BindFramebuffer(gl.FRAMEBUFFER, sbx_frameBufferId)
	gl.Viewport(0, 0, sbx_CUBE_SIZE, sbx_CUBE_SIZE)

	for i := 0; i < 6; i++ {
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), sbx_textureId, 0)
		// console.log(gl.checkFramebufferStatus(gl.FRAMEBUFFER)); // test
		// console.log(gl); // test

		gl.UniformMatrix3fv(sbx_xyzLocation, 1 /* count */, false, &sbx_xyzArray[i][0])
		gl.EnableVertexAttribArray(uint32(sbx_hLocation))
		gl.BindBuffer(gl.ARRAY_BUFFER, sbx_hBufferId)
		for j := 0; j < sbx_CUBE_SIZE+4; j++ {
			gl.Uniform1f(sbx_vLocation, float32(j-2))
			gl.VertexAttribPointer(uint32(sbx_hLocation), 1, gl.FLOAT, false, 0, gl.PtrOffset(0))
			gl.DrawArrays(gl.POINTS, 0, sbx_CUBE_SIZE+4)
		}
	}
	gl.GenerateMipmap(gl.TEXTURE_CUBE_MAP)

	gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(defaultFBO)) // return to default screen FBO
	//gl.ViewportWidth = wth;
	// gl.ViewportHeight = hth;
	gl.Viewport(viewport[0], viewport[1], viewport[2], viewport[3]) // restore viewport

}
