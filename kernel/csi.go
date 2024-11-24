package main

/*
#include <stdbool.h>
typedef struct CursorMoveCoordinate {
	int amount;
	bool isAbsolute;
} CursorMoveCoordinate;
typedef struct CursorMove {
	CursorMoveCoordinate x;
	CursorMoveCoordinate y;
} CursorMove;
*/
import "C"
import (
	"errors"
	"strconv"
	"strings"
)

type CursorMoveCoordinate struct {
	amount     int
	isAbsolute bool
}

type CursorMove struct {
	x CursorMoveCoordinate
	y CursorMoveCoordinate
}

func handleDirectionalCursorMove(final byte, parameters []byte) (*CursorMove, error) {
	if final == 'H' {
		before, after, found := strings.Cut(string(parameters), ";")
		if !found {
			return nil, errors.New("CSI H: invalid parameter format " + string(parameters))
		}

		x, err := strconv.Atoi(before)
		if err != nil {
			if before == "" {
				x = 1
			} else {
				return nil, err
			}
		}
		y, err := strconv.Atoi(after)
		if err != nil {
			if before == "" {
				y = 1
			} else {
				return nil, err
			}
		}

		return &CursorMove{
			x: CursorMoveCoordinate{
				amount:     x,
				isAbsolute: true,
			},
			y: CursorMoveCoordinate{
				amount:     y,
				isAbsolute: true,
			},
		}, nil
	}

	amount, err := strconv.Atoi(string(parameters))
	if err != nil {
		if string(parameters) == "" {
			amount = 1
		} else {
			return nil, err
		}
	}

	if final == 'A' {
		return &CursorMove{
			x: CursorMoveCoordinate{
				amount:     0,
				isAbsolute: false,
			},
			y: CursorMoveCoordinate{
				amount:     amount,
				isAbsolute: false,
			},
		}, nil
	}
	if final == 'B' {
		return &CursorMove{
			x: CursorMoveCoordinate{
				amount:     0,
				isAbsolute: false,
			},
			y: CursorMoveCoordinate{
				amount:     -amount,
				isAbsolute: false,
			},
		}, nil
	}
	if final == 'C' {
		return &CursorMove{
			x: CursorMoveCoordinate{
				amount:     amount,
				isAbsolute: false,
			},
			y: CursorMoveCoordinate{
				amount:     0,
				isAbsolute: false,
			},
		}, nil
	}
	if final == 'D' {
		return &CursorMove{
			x: CursorMoveCoordinate{
				amount:     -amount,
				isAbsolute: false,
			},
			y: CursorMoveCoordinate{
				amount:     0,
				isAbsolute: false,
			},
		}, nil
	}
	if final == 'E' {
		return &CursorMove{
			x: CursorMoveCoordinate{
				amount:     0,
				isAbsolute: true,
			},
			y: CursorMoveCoordinate{
				amount:     -amount,
				isAbsolute: false,
			},
		}, nil
	}
	if final == 'F' {
		return &CursorMove{
			x: CursorMoveCoordinate{
				amount:     0,
				isAbsolute: true,
			},
			y: CursorMoveCoordinate{
				amount:     amount,
				isAbsolute: false,
			},
		}, nil
	}
	if final == 'G' {
		return &CursorMove{
			x: CursorMoveCoordinate{
				amount:     amount,
				isAbsolute: true,
			},
			y: CursorMoveCoordinate{
				amount:     0,
				isAbsolute: false,
			},
		}, nil
	}
	return nil, errors.New("unknown CursorMove CSI Sequence")
}

func dispatchCSI(parameters []byte, intermediates []byte, final byte) ([]byte, *Event) {
	if final < 0x40 || final > 0x7E {
		Log("CSI", "final byte in CSI Escape Sequence out of range\n")
		return nil, &Event{
			name: "Error",
			data: append(append(parameters, intermediates...), final),
			kind: EventKindError,
		}
	}

	if final >= 'A' && final <= 'H' {
		cursorMove, err := handleDirectionalCursorMove(final, parameters)
		if err != nil {
			Log("CSI", "handleDirectionalCursorMove error: %s\n", err.Error())
			return nil, &Event{
				name: "Error",
				data: append(append(parameters, intermediates...), final),
				kind: EventKindError,
			}
		}
		return nil, &Event{
			name: "CursorMove",
			data: cursorMove,
			kind: EventKindCursorMove,
		}
	}

	Log("CSI", "Unknown CSI Escape Sequence\n")
	return nil, &Event{
		name: "Error",
		data: append(append(parameters, intermediates...), final),
		kind: EventKindError,
	}
}

func ParseCSI(bytes []byte) ([]byte, *Event) {
	idx := 0

	// read parameter bytes in range 0x30–0x3F
	var paramBytes []byte
	for ; bytes[idx] >= 0x30 && bytes[idx] <= 0x3F; idx++ {
		paramBytes = append(paramBytes, bytes[idx])
	}

	// read intermediate bytes in range 0x20–0x2F
	var interBytes []byte
	for ; bytes[idx] >= 0x30 && bytes[idx] <= 0x3F; idx++ {
		interBytes = append(interBytes, bytes[idx])
	}

	// final byte 0x40–0x7E
	if idx > len(bytes)-1 {
		Log("CSI", "CSI Escape Sequence missing final byte\n")
		return nil, &Event{
			name: "Error",
			data: bytes,
			kind: EventKindError,
		}
	}
	finalByte := bytes[idx]

	return dispatchCSI(paramBytes, interBytes, finalByte)
}
