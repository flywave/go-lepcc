package lepcc

/*
#include "lepcc_c_api.h"
#cgo CFLAGS: -I ./
*/
import "C"
import (
	"errors"
	"unsafe"
)

const (
	XYZ       = 0
	RGB       = 1
	Intensity = 2
	FlagBytes = 3
)

type LepccStatus C.lepcc_status

type LepccBlobType C.lepcc_blobType

type LepccContext struct {
	ctx C.lepcc_ContextHdl
}

func NewContext() LepccContext {
	return LepccContext{ctx: C.lepcc_createContext()}
}

func (c *LepccContext) Close() {
	C.lepcc_deleteContext(&c.ctx)
}

func computeCompressedSizeXYZ(lctx LepccContext, xyz []float64, errs [3]float64) (uint32, []uint32, error) {
	if len(xyz)%3 != 0 {
		return 0, nil, errors.New("error")
	}
	ptr := (*C.double)(unsafe.Pointer(&xyz[0]))
	len := len(xyz) / 3
	var outlen uint32
	buff := make([]uint32, len)
	state := C.lepcc_computeCompressedSizeXYZ(lctx.ctx, C.uint(len), ptr,
		C.double(errs[0]), C.double(errs[1]), C.double(errs[2]),
		(*C.uint)(unsafe.Pointer(&outlen)), (*C.uint)(unsafe.Pointer(&buff[0])))
	if uint(state) != 0 {
		return 0, nil, errors.New("lepcc_computeCompressedSizeXYZ error")
	}
	return outlen, buff, nil
}

func computeCompressedSizeRGB(lctx LepccContext, rgb []byte) (uint32, error) {
	if len(rgb)%3 != 0 {
		return 0, errors.New("error")
	}
	ptr := (*C.uchar)(unsafe.Pointer(&rgb[0]))
	len := len(rgb) / 3
	var outlen uint32
	state := C.lepcc_computeCompressedSizeRGB(lctx.ctx, C.uint(len), ptr,
		(*C.uint)(unsafe.Pointer(&outlen)))
	if uint(state) != 0 {
		return 0, errors.New("lepcc_computeCompressedSizeRGB error")
	}
	return outlen, nil
}

func computeCompressedSizeIntensity(lctx LepccContext, vals []uint16) (uint32, error) {
	ptr := (*C.ushort)(unsafe.Pointer(&vals[0]))
	len := len(vals)
	var outlen uint32
	state := C.lepcc_computeCompressedSizeIntensity(lctx.ctx, C.uint(len), ptr,
		(*C.uint)(unsafe.Pointer(&outlen)))
	if uint(state) != 0 {
		return 0, errors.New("lepcc_computeCompressedSizeIntensity error")
	}
	return outlen, nil
}

func computeCompressedSizeFlagBytes(lctx LepccContext, vals []byte) (uint32, error) {
	ptr := (*C.uchar)(unsafe.Pointer(&vals[0]))
	len := len(vals)
	var outlen uint32
	state := C.lepcc_computeCompressedSizeFlagBytes(lctx.ctx, C.uint(len), ptr,
		(*C.uint)(unsafe.Pointer(&outlen)))
	if uint(state) != 0 {
		return 0, errors.New("lepcc_computeCompressedSizeFlagBytes error")
	}
	return outlen, nil
}

func EncodeXYZ(lctx LepccContext, xyz []float64, errs [3]float64) ([]byte, error) {
	size, _, err := computeCompressedSizeXYZ(lctx, xyz, errs)
	if err != nil {
		return nil, err
	}
	buff := make([]byte, size)
	cptr := (*C.uchar)(unsafe.Pointer(&buff[0]))
	state := C.lepcc_encodeXYZ(lctx.ctx, &cptr, C.int(size))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_encodeXYZ error")
	}
	return buff, nil
}

func EncodeRGB(lctx LepccContext, rgb []byte) ([]byte, error) {
	size, err := computeCompressedSizeRGB(lctx, rgb)
	if err != nil {
		return nil, err
	}
	buff := make([]byte, size)
	cptr := (*C.uchar)(unsafe.Pointer(&buff[0]))
	state := C.lepcc_encodeRGB(lctx.ctx, &cptr, C.int(size))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_encodeRGB error")
	}
	return buff, nil
}

func EncodeIntensity(lctx LepccContext, vals []uint16) ([]byte, error) {
	size, err := computeCompressedSizeIntensity(lctx, vals)
	if err != nil {
		return nil, err
	}
	buff := make([]byte, size)
	cptr := (*C.uchar)(unsafe.Pointer(&buff[0]))
	ptr := (*C.ushort)(unsafe.Pointer(&vals[0]))
	len := len(vals)
	state := C.lepcc_encodeIntensity(lctx.ctx, &cptr, C.int(size), ptr, C.uint(len))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_encodeIntensity error")
	}
	return buff, nil
}

func EncodeFlagBytes(lctx LepccContext, vals []byte) ([]byte, error) {
	size, err := computeCompressedSizeFlagBytes(lctx, vals)
	if err != nil {
		return nil, err
	}
	buff := make([]byte, size)
	cptr := (*C.uchar)(unsafe.Pointer(&buff[0]))
	ptr := (*C.uchar)(unsafe.Pointer(&vals[0]))
	len := len(vals)
	state := C.lepcc_encodeFlagBytes(lctx.ctx, &cptr, C.int(size), ptr, C.uint(len))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_encodeFlagBytes error")
	}
	return buff, nil
}

func DecodeXYZ(lctx LepccContext, buff []byte) ([]float64, error) {
	nInfo := C.lepcc_getBlobInfoSize()

	bt := uint32(0)
	blobSize := uint32(0)
	ptr := (*C.uchar)(C.CBytes(buff))

	state := C.lepcc_getBlobInfo(lctx.ctx, ptr, nInfo, (*C.uint)(unsafe.Pointer(&bt)), (*C.uint)(unsafe.Pointer(&blobSize)))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_getBlobInfo error")
	}
	var pointcount uint32
	state = C.lepcc_getPointCount(lctx.ctx, ptr, C.int(blobSize), (*C.uint)(unsafe.Pointer(&pointcount)))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_getPointCount error")
	}
	xyz := make([]float64, pointcount*3)
	state = C.lepcc_decodeXYZ(lctx.ctx, &ptr, C.int(blobSize), (*C.uint)(unsafe.Pointer(&pointcount)), (*C.double)(unsafe.Pointer(&xyz[0])))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_decodeXYZ error")
	}
	return xyz, nil
}

func DecodeRGB(lctx LepccContext, buff []byte) ([]byte, error) {
	nInfo := C.lepcc_getBlobInfoSize()

	bt := uint32(0)
	blobSize := uint32(0)
	ptr := (*C.uchar)(C.CBytes(buff))

	state := C.lepcc_getBlobInfo(lctx.ctx, ptr, nInfo, (*C.uint)(unsafe.Pointer(&bt)), (*C.uint)(unsafe.Pointer(&blobSize)))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_getBlobInfo error")
	}
	var rcount uint32
	state = C.lepcc_getRGBCount(lctx.ctx, ptr, C.int(blobSize), (*C.uint)(unsafe.Pointer(&rcount)))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_getRGBCount error")
	}
	rgb := make([]byte, rcount*3)
	state = C.lepcc_decodeRGB(lctx.ctx, &ptr, C.int(blobSize), (*C.uint)(unsafe.Pointer(&rcount)), (*C.uchar)(unsafe.Pointer(&rgb[0])))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_decodeRGB error")
	}
	return rgb, nil
}

func DecodeIntensity(lctx LepccContext, buff []byte) ([]uint16, error) {
	nInfo := C.lepcc_getBlobInfoSize()

	bt := uint32(0)
	blobSize := uint32(0)
	ptr := (*C.uchar)(C.CBytes(buff))

	state := C.lepcc_getBlobInfo(lctx.ctx, ptr, nInfo, (*C.uint)(unsafe.Pointer(&bt)), (*C.uint)(unsafe.Pointer(&blobSize)))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_getBlobInfo error")
	}
	var count uint32
	state = C.lepcc_getIntensityCount(lctx.ctx, ptr, C.int(blobSize), (*C.uint)(unsafe.Pointer(&count)))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_getIntensityCount error")
	}
	ints := make([]uint16, count)
	state = C.lepcc_decodeIntensity(lctx.ctx, &ptr, C.int(blobSize), (*C.uint)(unsafe.Pointer(&count)), (*C.ushort)(unsafe.Pointer(&ints[0])))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_decodeIntensity error")
	}
	return ints, nil
}

func DecodeFlagByte(lctx LepccContext, buff []byte) ([]byte, error) {
	nInfo := C.lepcc_getBlobInfoSize()

	bt := uint32(0)
	blobSize := uint32(0)
	ptr := (*C.uchar)(C.CBytes(buff))

	state := C.lepcc_getBlobInfo(lctx.ctx, ptr, nInfo, (*C.uint)(unsafe.Pointer(&bt)), (*C.uint)(unsafe.Pointer(&blobSize)))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_getBlobInfo error")
	}
	var count uint32
	state = C.lepcc_getFlagByteCount(lctx.ctx, ptr, C.int(blobSize), (*C.uint)(unsafe.Pointer(&count)))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_getFlagByteCount error")
	}
	flags := make([]byte, count)
	state = C.lepcc_decodeFlagBytes(lctx.ctx, &ptr, C.int(blobSize), (*C.uint)(unsafe.Pointer(&count)), (*C.uchar)(unsafe.Pointer(&flags[0])))
	if uint(state) != 0 {
		return nil, errors.New("lepcc_decodeFlagBytes error")
	}
	return flags, nil
}
