package image

import (
	"errors"
	"image"
)

// Backend is the interface for various image
// backends
type Backend interface {
	// Available should return nil if an image
	// backend successfully runs. This could be
	// used for intializing backends.
	Available() error

	Spawn(image.Image, int, int) error

	Delete() error
}

// PossibleBackends is a list of possible image
// backends
var PossibleBackends = []Backend{
	&Ueberzug{},
	&W3M{},
}

var (
	// ErrNoBackend is returned if no backends
	// are available
	ErrNoBackend = errors.New("no backend")
)

// New finds a backend and spawns the image
func New(i image.Image) (backend Backend, err error) {
	if t == nil {
		if err := Listen(); err != nil {
			return nil, err
		}
	}

	if PixelW < 1 || PixelH < 1 {
		return nil, ErrNoBackend
	}

	for _, b := range PossibleBackends {
		if b.Available() == nil {
			backend = b
			break
		}
	}

	if backend == nil {
		return nil, ErrNoBackend
	}

	w := i.Bounds().Dx()

	if err := backend.Spawn(i, PixelW/2-w/2, 0); err != nil {
		return nil, err
	}

	return backend, nil
}

// Close has to be ran on exit!
func Close() {
	if t != nil {
		t.Close()
	}
}
