// bluemoon-graphics project bluemoon-graphics.go
package graphics

import (
	"strconv"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type renderer struct {
	offset       int32
	vertexBuffer *uint32
	indexBuffer  *uint32
}

func CreateRenderer() *renderer {
	r := new(renderer)

	gl.GenBuffers(1, r.vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, *r.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 2000*4*3*strconv.IntSize, nil, gl.DYNAMIC_DRAW)

	gl.GenBuffers(1, r.vertexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, *r.indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 2000*6*strconv.IntSize, nil, gl.DYNAMIC_DRAW)

	pointer := gl.MapBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.WRITE_ONLY)

	var j uint = 0
	for i := 0; i < 12000; i += 6 {
		*(*uint)(unsafe.Pointer(uintptr(pointer) + uintptr(i*strconv.IntSize))) = j
		*(*uint)(unsafe.Pointer(uintptr(pointer) + uintptr((i+1)*strconv.IntSize))) = j + 1
		*(*uint)(unsafe.Pointer(uintptr(pointer) + uintptr((i+2)*strconv.IntSize))) = j + 2
		*(*uint)(unsafe.Pointer(uintptr(pointer) + uintptr((i+3)*strconv.IntSize))) = j + 1
		*(*uint)(unsafe.Pointer(uintptr(pointer) + uintptr((i+4)*strconv.IntSize))) = j + 2
		*(*uint)(unsafe.Pointer(uintptr(pointer) + uintptr((i+5)*strconv.IntSize))) = j + 3
		j += 4
	}

	gl.UnmapBuffer(gl.ELEMENT_ARRAY_BUFFER)

	return r
}

func (r *renderer) Begin() {
	r.offset = 0
	gl.BindBuffer(gl.ARRAY_BUFFER, *r.vertexBuffer)
}

func (r *renderer) DrawRectangle(x, y, w, h int) {
	array := [12]int{
		x, y, 0,
		x + w, y, 0,
		x, y + h, 0,
		x + w, y + h, 0,
	}
	gl.BufferSubData(gl.ARRAY_BUFFER, (int)(r.offset*4*strconv.IntSize), 4*strconv.IntSize, unsafe.Pointer(&array))
	r.offset += 1
	if r.offset == 2000 {
		r.Flush()
		r.offset = 0
	}
}

func (r *renderer) Flush() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, *r.indexBuffer)
	gl.DrawElements(gl.TRIANGLES, r.offset*6, gl.UNSIGNED_INT, nil)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (r *renderer) End() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
