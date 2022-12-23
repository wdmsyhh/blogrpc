package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/spf13/cast"
)

// key 需要是 base64 编码的字符串且解码后长度必须为 128 bits
// 采用 CBC 加密算法，会直接使用 key 做 IV，因此结果不具有随机性
func AESEncrypt(key, plaintext string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	if keyLen := len(keyBytes); keyLen != 16 {
		return "", errors.New("invalid key length")
	}
	plainBytes := []byte(plaintext)

	block, err := aes.NewCipher(keyBytes)
	plainBytes = PKCS7Padding(plainBytes, block.BlockSize())
	// 选择 CBC 加密算法，使用 keyBytes 作 IV
	blockModel := cipher.NewCBCEncrypter(block, keyBytes)
	cipherBytes := make([]byte, len(plainBytes))
	blockModel.CryptBlocks(cipherBytes, plainBytes)
	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

// key 需要是 base64 编码的字符串且解码后长度必须为 128 bits
func AESDecrypt(key, ciphertext string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	if keyLen := len(keyBytes); keyLen != 16 {
		return "", errors.New("invalid key length")
	}
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	blockModel := cipher.NewCBCDecrypter(block, keyBytes)
	plainBytes := make([]byte, len(cipherBytes))
	blockModel.CryptBlocks(plainBytes, cipherBytes)
	plainBytes = PKCS7UnPadding(plainBytes)
	return string(plainBytes), nil
}

// key 需要是 base64 编码的字符串且解码后长度必须为 128 bits
func AESDecryptIgnorePanic(key, ciphertext string) (plaintext string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(cast.ToString(r))
		}
	}()
	plaintext, err = AESDecrypt(key, ciphertext)
	return
}

func PKCS7Padding(src []byte, size int) []byte {
	padding := size - len(src)%size
	paddingBytes := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, paddingBytes...)
}

func PKCS7UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func IsEncryptedStr(s string) bool {
	// 由于加密过的字符串都是 base64 编码后的值，因此只需要检测是不是 base64 字符串即可判断是否加密
	_, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return false
	}
	return true
}
