package jwt

var accessSecret []byte
var refreshSecret []byte

type secretType int

const (
	access  secretType = 1
	refresh secretType = 2
)

func SetSecretKey(accessTokenKey, refreshTokenKey string) {
	accessSecret = []byte(accessTokenKey)
	refreshSecret = []byte(refreshTokenKey)
}

func getSecretKey(t secretType) ([]byte, error) {
	switch t {
	case access:
		if accessSecret == nil {
			break
		}
		return accessSecret, nil
	case refresh:
		if refreshSecret == nil {
			break
		}
		return refreshSecret, nil
	}

	return nil, ErrNoSecret
}
