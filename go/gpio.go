package golibbeaglebone

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Go inexplicably(!) lacks enums.
// As a result we're forced to use this obtuse method to alias user inputs so
// that the user doesn't have pass strings etc.irectly into the methods.
// It's ugly and inferior to enums + pattern matching, but it works and it's
// better than the alternative.

// PinDirection is the direction of the pin, which can be either an input or output.
type PinDirection int

const (
	// In means the pin is an input, and will read logic levels.
	In PinDirection = iota
	// Out means the pin is an output, and will write logic levels.
	Out
)

// PinState is the logic level of an output GPIO pin, either high or low.
type PinState int

const (
	// High means the pin is either writing or reading the high logic level,
	// which is 3.3V on the BeagleBone.
	High PinState = iota
	// Low means the pin is either writing or reading the low logic level,
	// which is 0V.
	Low
)

// GPIO represents a pin configured as a GPIO.
type GPIO struct {
	pinNum  uint8
	dirPath string
}

// NewGPIO creates a new GPIO pin object.
func NewGPIO(pinNum uint8) *GPIO {
	gpio := new(GPIO)
	gpio.pinNum = pinNum
	gpio.dirPath = fmt.Sprintf("/sys/class/gpio/gpio%d", gpio.pinNum)
	return gpio
}

// SetDirection sets the GPIO pin as an input or an output.
func (gpio *GPIO) SetDirection(pd PinDirection) error {
	var direction []byte

	switch pd {
	case In:
		direction = []byte("in")
	case Out:
		direction = []byte("out")
	}

	file, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpio.pinNum), os.O_WRONLY|os.O_SYNC, 0666)
	_, err = file.Write(direction)

	if err != nil {
		return err
	}
	return nil
}

// SetExportState enables or disables the GPIO pin.
func (gpio *GPIO) SetExportState(es ExportState) error {
	if _, err := os.Stat(gpio.dirPath); os.IsNotExist(err) && es == Exported {
		// Try to export if the GPIO isn't already exported
		file, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY|os.O_SYNC, 0666)
		_, err = file.Write([]byte(strconv.Itoa(int(gpio.pinNum))))
		if err != nil {
			return err
		}
	} else if _, err := os.Stat(gpio.dirPath); err == nil && es == UnExported {
		// Try to unexport if the GPIO is already exported
		file, err := os.OpenFile("/sys/class/gpio/unexport", os.O_WRONLY|os.O_SYNC, 0666)
		_, err = file.Write([]byte(strconv.Itoa(int(gpio.pinNum))))
		if err != nil {
			return err
		}
	} else {
		// User either tried to export an exported pin or unexport an unexported pin.
		return errors.New("Unable to export or unexport GPIO pin")
	}
	return nil
}

// Write writes to the GPIO pin, setting it either High or Low.
func (gpio *GPIO) Write(s PinState) error {
	// Figure out what we're writing to this file.
	var state []byte
	switch s {
	case High:
		state = []byte("1")
	case Low:
		state = []byte("0")
	}

	// Write to the file and make sure it worked.
	file, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpio.pinNum), os.O_WRONLY|os.O_SYNC, 0666)
	_, err = file.Write(state)
	if err != nil {
		return err
	}
	return nil
}

// Read returns the state (logic level) of the GPIO pin, either High or Low.
func (gpio *GPIO) Read() (state PinState, err error) {
	// Read from the file.
	buf := make([]byte, 4)
	file, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpio.pinNum), os.O_RDONLY|os.O_SYNC, 0666)
	numBytes, err := file.Read(buf)

	// Check if the file was read properly.
	if numBytes != 2 {
		return 0, errors.New("Unable to read from GPIO pin")
	} else if buf[0] == '1' {
		return High, nil
	} else {
		return Low, nil
	}
}
