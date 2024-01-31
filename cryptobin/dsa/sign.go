package dsa

import (
    "errors"
    "strings"
    "math/big"
    "crypto/dsa"
    "crypto/rand"
    "encoding/asn1"
)

// 私钥签名
func (this DSA) Sign(separator ...string) DSA {
    if this.privateKey == nil {
        err := errors.New("dsa: privateKey error.")
        return this.AppendError(err)
    }

    hashed, err := this.dataHash(this.signHash, this.data)
    if err != nil {
        return this.AppendError(err)
    }

    r, s, err := dsa.Sign(rand.Reader, this.privateKey, hashed)
    if err != nil {
        return this.AppendError(err)
    }

    rt, err := r.MarshalText()
    if err != nil {
        return this.AppendError(err)
    }

    st, err := s.MarshalText()
    if err != nil {
        return this.AppendError(err)
    }

    sep := "+"
    if len(separator) > 0 {
        sep = separator[0]
    }

    signStr := string(rt) + sep + string(st)

    this.parsedData = []byte(signStr)

    return this
}

// 公钥验证
// 使用原始数据[data]对比签名后数据
func (this DSA) Verify(data []byte, separator ...string) DSA {
    if this.publicKey == nil {
        err := errors.New("dsa: publicKey error.")
        return this.AppendError(err)
    }

    hashed, err := this.dataHash(this.signHash, data)
    if err != nil {
        return this.AppendError(err)
    }

    sep := "+"
    if len(separator) > 0 {
        sep = separator[0]
    }

    split := strings.Split(string(this.data), sep)
    if len(split) != 2 {
        err := errors.New("dsa: sign data is error.")
        return this.AppendError(err)
    }

    rStr := split[0]
    sStr := split[1]
    rr := new(big.Int)
    ss := new(big.Int)

    err = rr.UnmarshalText([]byte(rStr))
    if err != nil {
        return this.AppendError(err)
    }

    err = ss.UnmarshalText([]byte(sStr))
    if err != nil {
        return this.AppendError(err)
    }

    this.verify = dsa.Verify(this.publicKey, hashed, rr, ss)

    return this
}

// ===============

type DSASignature struct {
    R, S *big.Int
}

// 私钥签名
func (this DSA) SignASN1() DSA {
    if this.privateKey == nil {
        err := errors.New("dsa: privateKey error.")
        return this.AppendError(err)
    }

    hashed, err := this.dataHash(this.signHash, this.data)
    if err != nil {
        return this.AppendError(err)
    }

    r, s, err := dsa.Sign(rand.Reader, this.privateKey, hashed)
    if err != nil {
        return this.AppendError(err)
    }

    parsedData, err := asn1.Marshal(DSASignature{r, s})

    this.parsedData = parsedData

    return this.AppendError(err)
}

// 公钥验证
// 使用原始数据[data]对比签名后数据
func (this DSA) VerifyASN1(data []byte) DSA {
    if this.publicKey == nil {
        err := errors.New("dsa: publicKey error.")
        return this.AppendError(err)
    }

    var dsaSign DSASignature
    _, err := asn1.Unmarshal(this.data, &dsaSign)
    if err != nil {
        return this.AppendError(err)
    }

    hashed, err := this.dataHash(this.signHash, data)
    if err != nil {
        return this.AppendError(err)
    }

    r := dsaSign.R
    s := dsaSign.S

    this.verify = dsa.Verify(this.publicKey, hashed, r, s)

    return this
}

// ===============

const (
    // 字节大小
    dsaSubgroupBytes = 32
)

// 私钥签名
func (this DSA) SignBytes() DSA {
    if this.privateKey == nil {
        err := errors.New("dsa: privateKey error.")
        return this.AppendError(err)
    }

    hashed, err := this.dataHash(this.signHash, this.data)
    if err != nil {
        return this.AppendError(err)
    }

    r, s, err := dsa.Sign(rand.Reader, this.privateKey, hashed)
    if err != nil {
        return this.AppendError(err)
    }

    rBytes := r.Bytes()
    sBytes := s.Bytes()
    if len(rBytes) > dsaSubgroupBytes || len(sBytes) > dsaSubgroupBytes {
        err := errors.New("dsa: DSA signature too large.")
        return this.AppendError(err)
    }

    out := make([]byte, 2*dsaSubgroupBytes)
    copy(out[dsaSubgroupBytes-len(rBytes):], rBytes)
    copy(out[len(out)-len(sBytes):], sBytes)

    this.parsedData = out

    return this
}

// 公钥验证
// 使用原始数据[data]对比签名后数据
func (this DSA) VerifyBytes(data []byte) DSA {
    if this.publicKey == nil {
        err := errors.New("dsa: publicKey error.")
        return this.AppendError(err)
    }

    // 签名结果数据
    sig := this.data

    if len(sig) != 2*dsaSubgroupBytes {
        err := errors.New("dsa: sig data error.")
        return this.AppendError(err)
    }

    r := new(big.Int).SetBytes(sig[:dsaSubgroupBytes])
    s := new(big.Int).SetBytes(sig[dsaSubgroupBytes:])

    hashed, err := this.dataHash(this.signHash, data)
    if err != nil {
        return this.AppendError(err)
    }

    this.verify = dsa.Verify(this.publicKey, hashed, r, s)

    return this
}

// ===============

// 签名后数据
func (this DSA) dataHash(fn HashFunc, data []byte) ([]byte, error) {
    h := fn()
    h.Write(data)

    return h.Sum(nil), nil
}
