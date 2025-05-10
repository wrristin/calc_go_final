package agent

import (
	"os"
	"strconv"
)

var ComputingPower = getComputingPower()

func getComputingPower() int {
	power := os.Getenv("COMPUTING_POWER")
	if power == "" {
		return 4 // Значение по умолчанию
	}
	value, err := strconv.Atoi(power)
	if err != nil {
		return 4 // Значение по умолчанию в случае ошибки
	}
	return value
}
