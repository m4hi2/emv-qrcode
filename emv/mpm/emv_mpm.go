package mpm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dongri/emv-qrcode/crc16"
)

// Encode ...
func Encode(emvqr *EMVQR) (string, error) {
	if err := emvqr.Validate(); err != nil {
		return "", err
	}
	return emvqr.GeneratePayload(), nil
}

// Decode ...
func Decode(payload string) (*EMVQR, error) {
	emvqr, err := ParseEMVQR(payload)
	if err != nil {
		return nil, err
	}
	if err := checkCRC(payload, emvqr.CRC.Value); err != nil {
		return nil, err
	}
	if err := emvqr.Validate(); err != nil {
		return emvqr, err
	}
	return emvqr, nil
}

func checkCRC(payload string, crc string) error {
	if crc == "" {
		return errors.New("CRC is mandatory")
	}

	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	computed := crc16.Checksum([]byte(payload[:len(payload)-4]), table)
	computedS := strings.ToUpper(strconv.FormatUint(uint64(computed), 16))

	crc = strings.ToUpper(crc)

	if computedS != crc {
		return fmt.Errorf("invalid CRC: computed %s, got %s", computedS, crc)
	}

	return nil
}
