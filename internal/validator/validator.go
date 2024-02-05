package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const (
	gamerTagMinLength    = 3
	gamerTagMaxLength    = 20
	gamerTagRegexPattern = "^[a-zA-Z0-9_]+$"
)

func ValidateWalletAddress(walletAddress string) error {
	if !common.IsHexAddress(walletAddress) {
		return errors.New("not a valid hex address")
	}

	if !strings.HasPrefix(walletAddress, "0x") {
		return errors.New("must be prefixed with '0x'")
	}

	zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")
	inputAddress := common.HexToAddress(strings.ToLower(walletAddress))
	if inputAddress == zeroAddress {
		return errors.New("must not be a zero address")
	}

	return nil
}

func ValidateGamerTag(gamerTag string) error {
	if len(gamerTag) > gamerTagMaxLength {
		return fmt.Errorf("must be at most %d characters long", gamerTagMaxLength)
	}

	if len(gamerTag) < gamerTagMinLength {
		return fmt.Errorf("must be at least %d characters long", gamerTagMinLength)
	}

	re := regexp.MustCompile(gamerTagRegexPattern)
	if !re.MatchString(gamerTag) {
		return errors.New("can only contain alphanumeric characters and underscores")
	}

	return nil
}
