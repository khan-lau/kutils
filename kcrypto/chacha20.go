package kcrypto

import "golang.org/x/crypto/chacha20"

// chacha20 加密, key必须32字节, nonce 12或24字节
func Chacha20Encrypt(src, key, nonce []byte) ([]byte, error) {
	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err == nil {
		cipher.SetCounter(0) // 轮转8次
		dst := make([]byte, len(src))
		cipher.XORKeyStream(dst, src)
		return dst, nil
	} else {
		return nil, err
	}
}

// chacha20 解密
func ChaCha20Decrypt(src, key, nonce []byte) ([]byte, error) {
	return Chacha20Encrypt(src, key, nonce)
}
