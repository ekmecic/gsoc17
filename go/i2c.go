package golibbeaglebone

import (
	"fmt"
	"os"
	"syscall"
)

// I2C_SLAVE is the Linux magic number for setting I2C slave addresses.
const I2C_SLAVE = 0x0703

// I2C represents an I2C device.
type I2C struct {
	i2cNum  uint8
	i2cFile *os.File
}

// NewI2C creates a new I2C device.
func NewI2C(i2cNum uint8) (*I2C, error) {
	i2c := new(I2C)
	i2c.i2cNum = i2cNum

	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", i2c.i2cNum), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	i2c.i2cFile = f
	return i2c, nil
}

// SetSlaveAddress sets the address of the I2C slave device.
func (i2c *I2C) SetSlaveAddress(slaveAddress uint16) error {
	if err := ioctl(i2c.i2cFile.Fd(), I2C_SLAVE, uintptr(slaveAddress)); err != nil {
		return err
	}
	return nil
}

// Write writes a single byte to an I2C slave.
func (i2c *I2C) Write(data byte) (int, error) {
	var buf [1]byte
	buf[0] = data
	return i2c.i2cFile.Write(buf[:])
}

// Read reads a single byte from an I2C slave and returns it.
func (i2c *I2C) Read(buf []byte) (int, error) {
	return i2c.i2cFile.Read(buf)
}

// Pulled from: https://github.com/davecheney/i2c/blob/master/i2c.go#L58
func ioctl(fd, cmd, arg uintptr) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}
