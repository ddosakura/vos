package vos

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
)

// RsaPub Key
func RsaPub(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, ErrKeyError
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

// RsaPriv Key
func RsaPriv(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, ErrKeyError
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

// RsaEncrypt - 加密
func RsaEncrypt(pub *rsa.PublicKey, origin []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origin)
}

// RsaDecrypt - 解密
func RsaDecrypt(priv *rsa.PrivateKey, cipher []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipher)
}

// RsaSign - 签名
func RsaSign(priv *rsa.PrivateKey, text []byte) ([]byte, error) {
	hash := sha1.New()
	hash.Write(text)
	digest := hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA1, digest)
}

// RsaVerify - 验证
func RsaVerify(pub *rsa.PublicKey, text []byte, sig []byte) error {
	hash := sha1.New()
	hash.Write(text)
	digest := hash.Sum(nil)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA1, digest, sig)
}
