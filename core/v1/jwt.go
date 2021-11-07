package v1
import (
	"crypto/rsa"
)
type Jwt struct {
	PublicKey  *rsa.PublicKey
}

