package sbxgpu

import (
// "github.com/go-gl/gl/v3.3-core/gl"
// "github.com/go-gl/mathgl/mgl32"
// "math"
// "math/rand"
// "strconv"
)

// NewSbxGpu creates structure for handling  skybox
func NewSbxGpu() SbxGpu {
	sbx := SbxGpu{}

	// initial values
	sbx.textureIdExists = false
	sbx.frameBufferIdExists = false

	// default values of the public fields
	sbx.TextureUnit = 0

	sbx.makeShaderProgram()

	return sbx
}

type SbxGpu struct {
	/* shader program for drawing the skybox */
	shaderProgram uint32

	/* location of vertex attributes locations */
	position int32

	/* vertex array object  for skybox drawing */
	vao uint32

	/* location of projection matrix unifrom for skybox drawing */
	projection int32

	/* location of view matrix unifrom for skybox drawing */
	view int32

	/* location of uniform samplerCube skybox */
	skybox int32

	/* id of the shader program for rendernig the texture walls */
	renderTextureShaderProgram uint32

	/* attribute location for sequence of points of horizontal line of a texture wall in rendering skybox */
	hLocation int32

	/* uniform location for vertical coordinate of the line */
	vLocation int32

	/* uniform location of the permutation for the wall */
	xyzLocation int32

	/*  vertex array object  for rendering new skybox */
	renderTextureVAO uint32

	/* texture id of skybox texture */
	textureId       uint32
	textureIdExists bool

	/* texture unit to be used by skybox */
	TextureUnit uint32

	/* frame buffer for drawing the walls of the skybox cube */
	frameBufferId       uint32
	frameBufferIdExists bool
}
