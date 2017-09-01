package golibbeaglebone

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// PWMState is state in which the PWM is in, either exported or unexported.
type PWMState int

const (
	// Enabled means the PWM is ready to use.
	Enabled PWMState = iota
	// Disabled means the PWM is unavailable for use.
	// This is the default state when the PWM is first created.
	Disabled
)

// PWM is a representation of a PWM device.
type PWM struct {
	pwmChipNum uint8
	pwmNum     uint8
	period     uint32
	dutyCycle  uint32
	state      PWMState
}

// NewPWM creates a new PWM device.
func NewPWM(pwmChipNum uint8, pwmNum uint8) *PWM {
	pwm := new(PWM)
	pwm.pwmChipNum = pwmChipNum
	pwm.pwmNum = pwmNum
	pwm.dutyCycle = 0
	pwm.state = Disabled
	return pwm
}

// SetExportState exports the PWM.
func (pwm *PWM) SetExportState(es ExportState) error {
	path := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d", pwm.pwmChipNum, pwm.pwmNum)

	if _, err := os.Stat(path); os.IsNotExist(err) && es == Exported {
		// Try to export if the PWM isn't already exported
		file, err := os.OpenFile(fmt.Sprintf("/sys/class/pwm/pwmchip%d/export", pwm.pwmChipNum), os.O_WRONLY|os.O_SYNC, 0666)
		// strconv.Itoa(int(arg)) is probably the best way to convert a number to the string the file expects.
		_, err = file.Write([]byte(strconv.Itoa(int(pwm.pwmNum))))
		if err != nil {
			return err
		}
	} else if _, err := os.Stat(path); err == nil && es == UnExported {
		// Try to unexport if the PWM is already exported
		file, err := os.OpenFile(fmt.Sprintf("/sys/class/pwm/pwmchip%d/unexport", pwm.pwmChipNum), os.O_WRONLY|os.O_SYNC, 0666)
		_, err = file.Write([]byte(strconv.Itoa(int(pwm.pwmNum))))
		if err != nil {
			return err
		}
	} else {
		// User either tried to export an exported PWM or unexport an unexported PWM.
		return errors.New("Unable to export or unexport PWM")
	}
	return nil
}

// SetPeriod sets the period of the PWM in nanoseconds.
func (pwm *PWM) SetPeriod(period uint32) error {
	path := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/period", pwm.pwmChipNum, pwm.pwmNum)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_SYNC, 0666)
	_, err = file.Write([]byte(strconv.Itoa(int(period))))
	if err != nil {
		return err
	}

	pwm.period = period
	return nil
}

// SetState sets the PWMState and enables or disabled the PWM.
func (pwm *PWM) SetState(ps PWMState) error {
	path := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/enable", pwm.pwmChipNum, pwm.pwmNum)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_SYNC, 0666)

	// Map the PWMState to the string the file expects.
	var psString string
	if ps == Enabled {
		psString = "1"
	} else {
		psString = "0"
	}

	_, err = file.Write([]byte(psString))
	if err != nil {
		return err
	}
	pwm.state = ps
	return nil
}

// Write sets the duty cycle of the PWM as a percentage from %0-%100.
func (pwm *PWM) Write(percentage float32) error {
	path := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/duty_cycle", pwm.pwmChipNum, pwm.pwmNum)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_SYNC, 0666)

	newDutyCycle := uint32((percentage / 100.0) * float32(pwm.period))
	_, err = file.Write([]byte(strconv.Itoa(int(newDutyCycle))))
	if err != nil {
		return err
	}

	pwm.dutyCycle = newDutyCycle
	return nil
}

// SetDutyCycle sets the duty cycle as a period in nanoseconds.
func (pwm *PWM) SetDutyCycle(dutyCycle uint32) error {
	path := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/duty_cycle", pwm.pwmChipNum, pwm.pwmNum)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_SYNC, 0666)

	_, err = file.Write([]byte(strconv.Itoa(int(dutyCycle))))
	if err != nil {
		return err
	}

	pwm.dutyCycle = dutyCycle
	return nil
}
