package sm2

import (
    "errors"
    "crypto/rand"

    "github.com/deatil/go-cryptobin/gm/sm2"
)

// 私钥签名
func (this SM2) Sign() SM2 {
    if this.privateKey == nil {
        err := errors.New("SM2: privateKey error.")
        return this.AppendError(err)
    }

    hashed, err := this.dataHash(this.signHash, this.data)
    if err != nil {
        return this.AppendError(err)
    }

    parsedData, err := this.privateKey.Sign(rand.Reader, hashed, sm2.SignerOpts{
        Uid: this.uid,
    })
    if err != nil {
        return this.AppendError(err)
    }

    this.parsedData = parsedData

    return this
}

// 公钥验证
// 使用原始数据[data]对比签名后数据
func (this SM2) Verify(data []byte) SM2 {
    if this.publicKey == nil {
        err := errors.New("SM2: publicKey error.")
        return this.AppendError(err)
    }

    hashed, err := this.dataHash(this.signHash, data)
    if err != nil {
        return this.AppendError(err)
    }

    this.verify = this.publicKey.Verify(hashed, this.data, sm2.SignerOpts{
        Uid: this.uid,
    })

    return this
}

// ===============

// 私钥签名 ASN1
func (this SM2) SignASN1() SM2 {
    return this.Sign()
}

// 公钥验证 ASN1
// 使用原始数据[data]对比签名后数据
func (this SM2) VerifyASN1(data []byte) SM2 {
    return this.Verify(data)
}

// ===============

// 私钥签名 Bytes
// 兼容[招行]
func (this SM2) SignBytes() SM2 {
    if this.privateKey == nil {
        err := errors.New("SM2: privateKey error.")
        return this.AppendError(err)
    }

    hashed, err := this.dataHash(this.signHash, this.data)
    if err != nil {
        return this.AppendError(err)
    }

    parsedData, err := this.privateKey.SignBytes(rand.Reader, hashed, sm2.SignerOpts{
        Uid: this.uid,
    })
    if err != nil {
        return this.AppendError(err)
    }

    this.parsedData = parsedData

    return this
}

// 公钥验证 Bytes
// 兼容[招行]
// 使用原始数据[data]对比签名后数据
func (this SM2) VerifyBytes(data []byte) SM2 {
    if this.publicKey == nil {
        err := errors.New("SM2: publicKey error.")
        return this.AppendError(err)
    }

    if len(this.data) != 64 {
        err := errors.New("SM2: sig error.")
        return this.AppendError(err)
    }

    hashed, err := this.dataHash(this.signHash, data)
    if err != nil {
        return this.AppendError(err)
    }

    this.verify = this.publicKey.VerifyBytes(hashed, this.data, sm2.SignerOpts{
        Uid: this.uid,
    })

    return this
}

// ===============

// 签名后数据
func (this SM2) dataHash(fn HashFunc, data []byte) ([]byte, error) {
    if fn == nil {
        return data, nil
    }

    h := fn()
    h.Write(data)

    return h.Sum(nil), nil
}
