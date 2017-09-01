package golibbeaglebone

import (
	"fmt"
	"os"
	"strconv"
)

// ADC represents a pin configured as an ADC.
type ADC struct {
	adcNum        uint8
	scalingFactor float32
}

// NewADC creates a new ADC object.
func NewADC(adcNum uint8, scalingFactor float32) *ADC {
	adc := new(ADC)
	adc.adcNum = adcNum
	adc.scalingFactor = scalingFactor
	return adc
}

// Read reads the raw voltage from the ADC pin.
func (adc *ADC) Read() (uint32, error) {
	path := fmt.Sprintf("/sys/bus/iio/devices/iio:device0/in_voltage%d_raw", adc.adcNum)
	buf := make([]byte, 16)
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_SYNC, 0666)

	numBytes, err := file.Read(buf)

	var value uint64
	value, err = strconv.ParseUint(string(buf[0:numBytes-1]), 10, 16)
	if err != nil {
		return 0, err
	}

	return uint32(value), nil
}

// ScaledRead reads the raw voltage of the ADC and applies a scaling factor to it.
func (adc *ADC) ScaledRead() (float32, error) {
	path := fmt.Sprintf("/sys/bus/iio/devices/iio:device0/in_voltage%d_raw", adc.adcNum)
	buf := make([]byte, 16)
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_SYNC, 0666)

	numBytes, err := file.Read(buf)

	var value uint64
	value, err = strconv.ParseUint(string(buf[0:numBytes-1]), 10, 16)
	if err != nil {
		return 0, err
	}

	return float32(value) * adc.scalingFactor, nil
}
