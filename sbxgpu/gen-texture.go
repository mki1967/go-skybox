package sbxgpu

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

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

// RenderRandomCube creates new random skybox.
// Use RenderRandomCube before using DrawSkybox.
func (sbx *sbxGpu) RenderRandomCube() {

	var defaultFBO int32
	gl.GetIntegerv(gl.FRAMEBUFFER_BINDING, &defaultFBO)

	var viewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &viewport[0]) // save viewport parameters

	if sbx.textureIdExists == false {
		/* create texture object and allocate image memories */
		// sbx_textureId=gl.CreateTexture();
		gl.GenTextures(1, &sbx.textureId)
		sbx.textureIdExists = true
		gl.ActiveTexture(gl.TEXTURE0 + sbx.TextureUnit)
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, sbx.textureId)
		for i := 0; i < 6; i++ {
			gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.RGBA, sbx_CUBE_SIZE, sbx_CUBE_SIZE, 0, /* border */
				gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))
		}
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}

	if sbx.frameBufferIdExists == false {
		/* create framebuffer object */
		// sbx_frameBufferId=gl.CreateFramebuffer();
		gl.GenFramebuffers(1, &sbx.frameBufferId)
		sbx.frameBufferIdExists = true
	}
	gl.DeleteProgram(sbx_renderTextureShaderProgram) // delete old
	sbx_makeRenderTextureShaderProgram()
	gl.UseProgram(sbx_renderTextureShaderProgram)

	gl.ActiveTexture(gl.TEXTURE0 + sbx.TextureUnit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, sbx.textureId)

	gl.BindFramebuffer(gl.FRAMEBUFFER, sbx.frameBufferId)
	gl.Viewport(0, 0, sbx_CUBE_SIZE, sbx_CUBE_SIZE)

	for i := 0; i < 6; i++ {
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), sbx.textureId, 0)
		// console.log(gl.checkFramebufferStatus(gl.FRAMEBUFFER)); // test
		// console.log(gl); // test

		gl.UniformMatrix3fv(sbx_xyzLocation, 1 /* count */, false, &sbx_xyzArray[i][0])

		/* // replaced with gl.BindVertexArray(sbx_renderTextureVAO )
		gl.EnableVertexAttribArray(uint32(sbx_hLocation))
		gl.BindBuffer(gl.ARRAY_BUFFER, sbx_hBufferId)
		*/

		gl.BindVertexArray(sbx.renderTextureVAO)

		for j := 0; j < sbx_CUBE_SIZE+4; j++ {
			gl.Uniform1f(sbx_vLocation, float32(j-2))
			gl.VertexAttribPointer(uint32(sbx_hLocation), 1, gl.FLOAT, false, 0, gl.PtrOffset(0))
			gl.DrawArrays(gl.POINTS, 0, sbx_CUBE_SIZE+4)
		}
	}
	gl.GenerateMipmap(gl.TEXTURE_CUBE_MAP)

	gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(defaultFBO))          // return to default screen FBO
	gl.Viewport(viewport[0], viewport[1], viewport[2], viewport[3]) // restore viewport

}
