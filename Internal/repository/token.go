package repository

import (
	"Service_1Cv8/internal/token"
	"github.com/recoilme/pudge"
	"time"

	"Service_1Cv8/internal/compression"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/encryption"
)

func GetPrivateKey() (*encryption.KeyEncryption, error) {

	cfg := &pudge.Config{StoreMode: 1}
	db, err := pudge.Open(constants.PudgeKey, cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	gzipPrivKey := []byte{}
	db.Get("public", &gzipPrivKey)

	if len(gzipPrivKey) == 0 {
		return nil, nil
	}

	jsonRSAPrivateKey, err := compression.Decompress(gzipPrivKey)
	if err != nil {
		return nil, err
	}

	pvk, err := encryption.InitPrivateKey(jsonRSAPrivateKey)
	if err != nil {
		return nil, err
	}

	return pvk, nil
}

func SetPrivateKey() (*encryption.KeyEncryption, error) {

	arrCert, err := encryption.CreateCert()
	if err != nil {
		return nil, err
	}

	privkey := &arrCert[0]

	gzipRSAPrivateKey, _ := compression.Compress([]byte(privkey.String()))

	cfg := &pudge.Config{StoreMode: 1}
	db, err := pudge.Open(constants.PudgeKey, cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	pvk, err := encryption.InitPrivateKey([]byte(privkey.String()))
	if err != nil {
		return nil, err
	}

	db.Set("public", gzipRSAPrivateKey)
	return pvk, nil
}

func GetTokens() ([]token.ClaimStore, error) {
	cfg := &pudge.Config{StoreMode: 1}
	db, err := pudge.Open(constants.PudgeTokens, cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	//arrC := []token.Claim{}
	arrCS := []token.ClaimStore{}
	keys, _ := db.Keys(0, 0, 0, true)
	for _, key := range keys {
		var cs token.ClaimStore
		db.Get(key, &cs)
		cs.Value, _ = compression.Decompress(cs.Value)
		cs.Secret, _ = compression.Decompress(cs.Secret)
		//var c token.Claim
		arrCS = append(arrCS, cs)

		//claims, ok := token.ExtractClaims(string(strToken))
		//if !ok {
		//	continue
		//}
		//
		//c.Authorized = claims["authorized"].(bool)
		//c.Key = claims["key"].(string)
		//c.Exp = claims["exp"].(float64)
		//
		//arrC = append(arrC, c)
	}

	return arrCS, nil
}

func GetToken(k string) (token.ClaimStore, error) {
	cs := token.ClaimStore{}

	cfg := &pudge.Config{StoreMode: 1}
	db, err := pudge.Open(constants.PudgeTokens, cfg)
	if err != nil {
		return cs, err
	}
	defer db.Close()

	_ = db.Get(k, &cs)
	cs.Value, _ = compression.Decompress(cs.Value)
	cs.Secret, _ = compression.Decompress(cs.Secret)

	return cs, nil
}

func SetToken(c *token.ClaimStore) error {
	cfg := &pudge.Config{StoreMode: 1}
	db, err := pudge.Open(constants.PudgeTokens, cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	db.Set(c.Key, c)

	return nil
}

func DelToken(c *token.ClaimStore) error {
	cfg := &pudge.Config{StoreMode: 1}
	db, err := pudge.Open(constants.PudgeTokens, cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_ = db.Delete(c.Key)

	return nil
}

func CheckToken(k string) (token.ClaimStore, bool) {
	cfg := &pudge.Config{StoreMode: 1}
	db, err := pudge.Open(constants.PudgeTokens, cfg)
	if err != nil {
		return token.ClaimStore{}, false
	}
	defer db.Close()

	cs := token.ClaimStore{}
	db.Get(k, &cs)

	cs.Value, _ = compression.Decompress(cs.Value)
	cs.Secret, _ = compression.Decompress(cs.Secret)

	c, ok := token.ExtractClaims(string(cs.Value), cs.Secret)
	if c == nil {
		return cs, false
	}

	t1 := time.Now()
	t2 := time.Unix(int64(c["exp"].(float64)), 0)

	dif := int(t2.Sub(t1)/(time.Hour*24)) + 1
	if ok && dif > 0 && dif < constants.TimeLiveToken {

		tc := token.NewClaims(k, time.Duration(constants.TimeLiveToken))
		sk := string(cs.Secret)
		if sk == "" {
			sk = constants.SecretKey
		}
		byteSK := []byte(sk)

		tokenString, err := tc.GenerateJWT(byteSK)
		tokenByte := []byte(tokenString)

		gzipTokenString, err := compression.Compress(tokenByte)

		okChangeToken := true
		if err != nil {
			okChangeToken = false
		}

		gziSecretString, err := compression.Compress(byteSK)
		if err != nil {
			okChangeToken = false
		}

		if okChangeToken {
			claimStore := token.ClaimStore{
				Key:    k,
				Value:  gzipTokenString,
				Secret: gziSecretString,
			}
			err = SetToken(&claimStore)
			if err == nil {
				cs.Value = tokenByte
			}
		}

	}
	return cs, ok
}

func ExtendToken(lifeTime int64, k string) error {

	//t1 := time.Now()
	//t2 := time.Unix(lifeTime, 0)
	//
	//dif := int(t2.Sub(t1)/(time.Hour*24)) + 1
	//if dif > 0 && dif < constants.TimeLiveToken {
	//
	//	tc := token.NewClaims(k, time.Duration(constants.TimeLiveToken))
	//	sk := string(cs.Secret)
	//	if sk == "" {
	//		sk = constants.SecretKey
	//	}
	//	byteSK := []byte(sk)
	//
	//	tokenString, err := tc.GenerateJWT(byteSK)
	//	tokenByte := []byte(tokenString)
	//
	//	gzipTokenString, err := compression.Compress(tokenByte)
	//
	//	okChangeToken := true
	//	if err != nil {
	//		okChangeToken = false
	//	}
	//
	//	gziSecretString, err := compression.Compress(byteSK)
	//	if err != nil {
	//		okChangeToken = false
	//	}
	//
	//	if okChangeToken {
	//		claimStore := token.ClaimStore{
	//			Key:    k,
	//			Value:  gzipTokenString,
	//			Secret: gziSecretString,
	//		}
	//		err = SetToken(&claimStore)
	//		if err == nil {
	//			cs.Value = tokenByte
	//		}
	//	}
	//}

	return nil
}
