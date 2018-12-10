package store

import (
	"context"
	"encoding/base64"

	"github.com/carousell/Orion/utils/errors"
	"github.com/carousell/Orion/utils/spanutils"
	"golang.org/x/crypto/scrypt"
)

func cause(err error) error {
	if e, ok := err.(errors.ErrorExt); ok {
		return e.Cause()
	}
	return err
}

func getPasswordHash(ctx context.Context, password, salt string) string {
	span, ctx := spanutils.NewInternalSpan(ctx, "HashPasword")
	defer span.Finish()
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 65536, 8, 1, 32)
	return base64.StdEncoding.EncodeToString(dk)
}
