package crypto

import (
	"crypto/sha1"
	"encoding/base32"
	"time"

	"github.com/pquerna/otp/totp"
)

// generateOTPSecretFromClientID sinh secret TOTP cố định từ clientID
func generateOTPSecretFromClientID(clientID string) string {
	hasher := sha1.New()
	hasher.Write([]byte(clientID))
	sha1Bytes := hasher.Sum(nil) // 20 bytes
	p1 := sha1Bytes[0:4]
	p2 := sha1Bytes[4:8]
	p3 := sha1Bytes[8:12]
	p4 := sha1Bytes[12:16]
	p5 := sha1Bytes[16:20]
	var secretBytes []byte
	secretBytes = append(secretBytes, p5...)
	secretBytes = append(secretBytes, p1...)
	secretBytes = append(secretBytes, p4...)
	secretBytes = append(secretBytes, p2...)
	secretBytes = append(secretBytes, p3...)
	secret := base32.StdEncoding.EncodeToString(secretBytes)
	return secret
}

// GetTOTPByClientID sinh mã TOTP từ clientID
func GetTOTPByClientID(clientID string) (string, error) {
	secret := generateOTPSecretFromClientID(clientID)
	return totp.GenerateCode(secret, time.Now())
}

// VerifyTOTPByClientID xác thực mã TOTP với clientID
func VerifyTOTPByClientID(clientID, code string) bool {
	secret := generateOTPSecretFromClientID(clientID)
	return totp.Validate(code, secret)
}

// GetTOTPWithExpireByClientID sinh mã TOTP và trả về số giây còn lại đến khi hết hạn (theo chuẩn TOTP 30s)
func GetTOTPWithExpireByClientID(clientID string) (code string, secondsLeft int, err error) {
	secret := generateOTPSecretFromClientID(clientID)
	now := time.Now()
	code, err = totp.GenerateCode(secret, now)
	if err != nil {
		return "", 0, err
	}
	period := 30
	secondsLeft = period - int(now.Unix()%int64(period))
	return code, secondsLeft, nil
}
