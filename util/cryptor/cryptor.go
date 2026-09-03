package cryptor

import "github.com/duke-git/lancet/v2/cryptor"

func Sha256(str string) string {
	return cryptor.Sha256(str)
}

func Md5(str string) string {
	return cryptor.Md5String(str)
}

func AesEcbEncrypt(data, key []byte) []byte {
	return cryptor.AesEcbEncrypt(data, key)
}
func AesEcbDecrypt(encrypted, key []byte) []byte {
	return cryptor.AesEcbDecrypt(encrypted, key)
}
