package main

import (
	"log"

	"github.com/fapiko/go-learn-gl/opengl-tutorial/common"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw3/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Tutorial 07 - Model loading ported from
// http://www.opengl-tutorial.org/beginners-tutorials/tutorial-7-model-loading/
func main() {

	if err := glfw.Init(); err != nil {
		panic("Failed to initialize GLFW")
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.Samples, 4)

	// Drawing the triangle threw an error with OpenGL 3.3, downgrading to 2.1 seemed to solve it
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	// Open a window and create its OpenGL context
	window, err := glfw.CreateWindow(1024, 768, "Tutorial 07", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()

	// Initialize OpenGL - Go bindings use Glow and now Glew
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Ensure we can capture the escape key being pressed below
	window.SetInputMode(glfw.StickyKeysMode, gl.TRUE)

	// Hide the mouse and enable unlimited movement
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Dark blue background
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	// Enable depth test
	gl.Enable(gl.DEPTH_TEST)

	// Accept fragment if it is closer to the camera than the former one
	gl.DepthFunc(gl.LESS)

	// Cull triangles which normal is not towards the camera
	gl.Enable(gl.CULL_FACE)

	// Create and compile our GLSL program from the shaders
	programId := common.LoadShaders("TransformVertexShader.vertexshader", "TextureFragmentShader.fragmentshader")
	defer gl.DeleteProgram(programId)

	// Get a handle for our "MVP" uniform
	matrixId := gl.GetUniformLocation(programId, gl.Str("MVP\x00"))

	// Get a handle for our buffers
	vertexPositionModelspaceId := uint32(gl.GetAttribLocation(programId, gl.Str("vertexPosition_modelspace\x00")))

	vertices, uvs, _, err := common.LoadObj("cube.obj")
	if err != nil {
		log.Panic(err)
	}

	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4*3, gl.Ptr(vertices), gl.STATIC_DRAW)
	defer gl.DeleteBuffers(1, &vertexBuffer)

	//_, err = loadBmpCustom("uvtemplate.bmp")
	textureId, err := common.LoadBmpCustom("uvmap.bmp")
	if err != nil {
		panic(err)
	}

	var textureBuffer uint32
	gl.GenBuffers(1, &textureBuffer)
	gl.BindBuffer(gl.TEXTURE_BUFFER, textureBuffer)
	gl.BufferData(gl.TEXTURE_BUFFER, len(uvs)*4*2, gl.Ptr(uvs), gl.STATIC_DRAW)
	defer gl.DeleteBuffers(1, &textureBuffer)

	// Set the mouse at the center of the screen
	glfw.PollEvents()
	windowWidth, windowHeight := window.GetSize()
	window.SetCursorPos(float64(windowWidth/2), float64(windowHeight/2))

	for window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose() {

		// Clear the screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Use our shader
		gl.UseProgram(programId)

		// Compute the MVP matrix from keyboard and mouse input
		common.ComputeMatricesFromInputs()
		projection := common.GetProjectionMatrix()
		view := common.GetViewMatrix()
		model := mgl32.Ident4()
		mvp := projection.Mul4(view).Mul4(model)

		// Send our transformation to the currently bound shader, in the "MVP" uniform
		gl.UniformMatrix4fv(matrixId, 1, false, &mvp[0])

		// Bind our texture in Texture Unit 0
		gl.ActiveTexture(gl.TEXTURE0)

		// Set our "myTextureSampler" sampler to user Texture Unit 0
		gl.Uniform1i(textureId, 0)

		// 1st attribute buffer : vertices
		gl.EnableVertexAttribArray(vertexPositionModelspaceId)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
		gl.VertexAttribPointer(
			vertexPositionModelspaceId, // The attribute we want to configure
			3,               // size
			gl.FLOAT,        // type
			false,           // normalized?
			0,               // stride
			gl.PtrOffset(0)) // array buffer offset

		// Draw the triangle !
		gl.DrawArrays(gl.TRIANGLES, 0, 12*3) // 12*3 indices starting at 0 -> 12 triangles
		gl.DisableVertexAttribArray(vertexPositionModelspaceId)

		// 2nd attribute buffer : colors
		gl.EnableVertexAttribArray(1)
		gl.BindBuffer(gl.ARRAY_BUFFER, textureBuffer)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, nil)

		// Swap buffers
		window.SwapBuffers()
		glfw.PollEvents()

	}

}
