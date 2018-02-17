package sbxgpu

import (
// "github.com/go-gl/gl/v3.3-core/gl"
// "github.com/go-gl/mathgl/mgl32"
// "math"
// "math/rand"
// "strconv"
)

type SbxGpu struct {
	/* shader program for drawing the skybox */
	shaderProgram uint32
	/* vertex attributes locations */
	position int32
	VAO      uint32

	renderTextureVAO uint32
}
