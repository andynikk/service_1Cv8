package encryption

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"Service_1Cv8/internal/constants"
)

type KeyEncryption struct {
	TypeEncryption string
	PublicKey      *rsa.PublicKey
	PrivateKey     *rsa.PrivateKey
}

func (key *KeyEncryption) RsaEncrypt(msg []byte) ([]byte, error) {
	if key == nil {
		return msg, nil
	}
	encryptedBytes, err := rsa.EncryptOAEP(sha512.New512_256(), rand.Reader, key.PublicKey, msg, nil)
	return encryptedBytes, err
}

func (key *KeyEncryption) RsaDecrypt(msgByte []byte) ([]byte, error) {
	if key == nil {
		return msgByte, nil
	}
	msgByte, err := key.PrivateKey.Decrypt(nil, msgByte, &rsa.OAEPOptions{Hash: crypto.SHA512_256})
	return msgByte, err
}

func CreateCert() ([]bytes.Buffer, error) {
	var numSert int64
	var subjectKeyID string
	var lenKeyByte int

	fmt.Print("Введите уникальный номер сертификата: ")
	if _, err := fmt.Fscan(os.Stdin, &numSert); err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Print("Введите ИД ключа субъекта (пример ввода 12346): ")
	if _, err := fmt.Fscan(os.Stdin, &subjectKeyID); err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Print("Длина ключа в байтах: ")
	if _, err := fmt.Fscan(os.Stdin, &lenKeyByte); err != nil {
		log.Println(err)
		return nil, err
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(numSert),
		Subject: pkix.Name{
			Organization: []string{"TELEMATIKA"},
			Country:      []string{"RU"},
		},
		NotBefore: time.Now(),
		NotAfter: time.Now().AddDate(constants.TimeLivingCertificateYaer, constants.TimeLivingCertificateMounth,
			constants.TimeLivingCe5rtificateDay),
		SubjectKeyId: []byte(subjectKeyID),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, lenKeyByte)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	var certPEM bytes.Buffer
	_ = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var privateKeyPEM bytes.Buffer
	_ = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return []bytes.Buffer{certPEM, privateKeyPEM}, nil
}

func InitPrivateKey(data []byte) (*KeyEncryption, error) {

	pvkBlock, _ := pem.Decode(data)
	pvk, err := x509.ParsePKCS1PrivateKey(pvkBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return &KeyEncryption{TypeEncryption: constants.TypeEncryption, PrivateKey: pvk, PublicKey: &pvk.PublicKey}, nil
}

func InitPublicKey(data []byte) (*KeyEncryption, error) {

	certBlock, _ := pem.Decode(data)
	cert, _ := x509.ParseCertificate(certBlock.Bytes)
	certPublicKey := cert.PublicKey.(*rsa.PublicKey)
	return &KeyEncryption{TypeEncryption: constants.TypeEncryption, PublicKey: certPublicKey}, nil
}
