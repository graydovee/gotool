package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	uuidNum       int
	uuidVersion   int
	uuidArg       string
	uuidSimple    bool
	uuidUppercase bool
)

func NewUuid(version int, arg string) (uuid.UUID, error) {
	switch version {
	case 1:
		return uuid.NewUUID()
	case 3:
		return uuid.NewMD5(uuid.UUID{}, []byte(arg)), nil
	case 4:
		return uuid.NewRandom()
	case 5:
		return uuid.NewSHA1(uuid.UUID{}, []byte(arg)), nil
	default:
		return uuid.UUID{}, errors.New("not support")
	}
}

func newUuidString() (string, error) {
	newUUID, err := NewUuid(uuidVersion, uuidArg)
	if err != nil {
		return "", err
	}
	var uuidStr string
	if uuidSimple {
		binary, err := newUUID.MarshalBinary()
		if err != nil {
			return "", err
		}
		uuidStr = hex.EncodeToString(binary)
	} else {
		uuidStr = newUUID.String()
	}
	if uuidUppercase {
		uuidStr = strings.ToUpper(uuidStr)
	}
	return uuidStr, nil
}

func NewUUIDCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uuid [num]",
		Short: "Print the version number",
		Long:  `All software has versions. This is mine`,
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < uuidNum; i++ {
				newUUID, err := newUuidString()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(newUUID)
			}
		},
	}
	cmd.Flags().IntVarP(&uuidVersion, "version", "v", 1, "uuid version, support")
	cmd.Flags().BoolVarP(&uuidUppercase, "uppercase", "u", false, "uuid uppercase")
	cmd.Flags().BoolVarP(&uuidSimple, "simple", "s", false, "output uuid is simple")
	cmd.Flags().IntVarP(&uuidNum, "num", "n", 1, "number of uuid generated")
	cmd.Flags().StringVarP(&uuidArg, "arg", "a", "", "arg of uuid, used in version 3, 5")
	return cmd
}
