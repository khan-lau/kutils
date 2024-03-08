package logger

import "testing"

func TestXxx(t *testing.T) {
	logger := LoggerInstanceOnlyConsole(-1)

	logger.D("fuck off")
	logger.D("{} fuck off", "maybe")

	logger.D("")

	logger.D("int8 {} fuck off", int8(8))
	logger.D("uint8 {} fuck off", uint8(8))
	logger.D("int16 {} fuck off", int16(16))
	logger.D("uint16 {} fuck off", uint16(16))
	logger.D("int {} fuck off", int(10))
	logger.D("uint {} fuck off", uint(10))
	logger.D("int32 {} fuck off", int32(32))
	logger.D("uint32 {} fuck off", uint32(32))
	logger.D("int32 {} fuck off", int32(64))
	logger.D("uint32 {} fuck off", uint32(64))
	logger.D("float32 {} fuck off", float32(4.45))
	logger.D("float64 {} fuck off", float64(2.1))

	logger.D("")

	logger.D("int8 {} fuck off", []int8{0, 1, 2, 3, 4})
	logger.D("uint8 {} fuck off", []uint8{0, 1, 2, 3, 4})
	logger.D("int16 {} fuck off", []int16{0, 1, 2, 3, 4})
	logger.D("uint16 {} fuck off", []uint16{0, 1, 2, 3, 4})
	logger.D("int {} fuck off", []int{0, 1, 2, 3, 4})
	logger.D("uint {} fuck off", []uint{0, 1, 2, 3, 4})
	logger.D("int32 {} fuck off", []int32{0, 1, 2, 3, 4, 0x44, 0x38})
	logger.D("uint32 {} fuck off", []uint32{0, 1, 2, 3, 4})
	logger.D("int64 {} fuck off", []int64{0, 1, 2, 3, 4})
	logger.D("uint64 {} fuck off", []uint64{0, 1, 2, 3, 4})

	logger.D("float32 {} fuck off", []float32{0, 1, 2, 3, 4})
	logger.D("float64 {} fuck off", []float64{0, 1, 2, 3, 4})

	logger.D("")

	logger.D("string {} fuck off", []string{"0", "1", "2", "3", "4"})

	logger.D("")

	cmp := complex(4, 4)
	cmp64 := complex64(cmp)
	logger.D("complex128 {} complex64 {} fuck off", cmp, cmp64)

	logger.D("complex64 {} fuck off", []complex64{complex(4, 0), complex(4, 1), complex(4, 2), complex(4, 3), complex(4, 4)})
	logger.D("complex128 {} fuck off", []complex128{complex(4, 0), complex(4, 1), complex(4, 2), complex(4, 3), complex(4, 4)})

	logger.D("")

	type AA struct {
		A int
		B string
		C complex128
	}

	aa := AA{A: 12, B: "string", C: complex(4, -1)}

	logger.D("obj {} fuck off", aa)
	logger.D("*obj {} fuck off", &aa)

	logger.D("obj {} fuck off", []AA{aa, aa})

	logger.D("obj {} fuck off", []*AA{&aa, &aa})

	logger.D("")
}
