package gost_test

import (
    "fmt"
    "testing"
    "crypto/rand"
    "encoding/pem"

    "github.com/deatil/go-cryptobin/gost"
)

func decodePEM(pubPEM string) []byte {
    block, _ := pem.Decode([]byte(pubPEM))
    if block == nil {
        panic("failed to parse PEM block containing the key")
    }

    return block.Bytes
}

func encodePEM(src []byte, typ string) string {
    keyBlock := &pem.Block{
        Type:  typ,
        Bytes: src,
    }

    keyData := pem.EncodeToMemory(keyBlock)

    return string(keyData)
}

func TestEqual(t *testing.T) {
    testOneCurve(t, gost.CurveIdGostR34102001TestParamSet())
    testOneCurve(t, gost.CurveIdGostR34102001CryptoProAParamSet())
    testOneCurve(t, gost.CurveIdGostR34102001CryptoProBParamSet())
    testOneCurve(t, gost.CurveIdGostR34102001CryptoProCParamSet())

    testOneCurve(t, gost.CurveIdGostR34102001CryptoProXchAParamSet())
    testOneCurve(t, gost.CurveIdGostR34102001CryptoProXchBParamSet())

    testOneCurve(t, gost.CurveIdtc26gost34102012256paramSetA())

    testOneCurve(t, gost.CurveIdtc26gost34102012512paramSetA())
    testOneCurve(t, gost.CurveIdtc26gost34102012512paramSetB())
    testOneCurve(t, gost.CurveIdtc26gost34102012512paramSetC())
}

func testOneCurve(t *testing.T, curue *gost.Curve) {
    t.Run(fmt.Sprintf("PKCS8 %s", curue), func(t *testing.T) {
        priv, err := gost.GenerateKey(rand.Reader, curue)
        if err != nil {
            t.Fatal(err)
        }

        pub := priv.Public().(*gost.PublicKey)

        pubDer, err := gost.MarshalPublicKey(pub)
        if err != nil {
            t.Fatal(err)
        }
        privDer, err := gost.MarshalPrivateKey(priv)
        if err != nil {
            t.Fatal(err)
        }

        if len(privDer) == 0 {
            t.Error("expected export key Der error: priv")
        }
        if len(pubDer) == 0 {
            t.Error("expected export key Der error: pub")
        }

        newPub, err := gost.ParsePublicKey(pubDer)
        if err != nil {
            t.Fatal(err)
        }
        newPriv, err := gost.ParsePrivateKey(privDer)
        if err != nil {
            t.Fatal(err)
        }

        if !newPriv.Equal(priv) {
            t.Error("Marshal privekey error")
        }
        if !newPub.Equal(pub) {
            t.Error("Marshal public error")
        }
    })
}

func Test_Pkcs8(t *testing.T) {
    curue := gost.CurveIdtc26gost34102012512paramSetA()

    priv, err := gost.GenerateKey(rand.Reader, curue)
    if err != nil {
        t.Fatal(err)
    }

    pub := priv.Public().(*gost.PublicKey)

    pubDer, err := gost.MarshalPublicKey(pub)
    if err != nil {
        t.Fatal(err)
    }
    privDer, err := gost.MarshalPrivateKey(priv)
    if err != nil {
        t.Fatal(err)
    }

    if len(privDer) == 0 {
        t.Error("expected export key Der error: priv")
    }
    if len(pubDer) == 0 {
        t.Error("expected export key Der error: pub")
    }

    pri2 := encodePEM(privDer, "PRIVATE KEY")
    pub2 := encodePEM(pubDer, "PUBLIC KEY")

    if len(pri2) == 0 {
        t.Error("expected export key PEM error: priv")
    }
    if len(pub2) == 0 {
        t.Error("expected export key PEM error: pub")
    }

    // t.Error(pri2)
    // t.Error(pub2)
}

var GostR3410_2001_CryptoPro_A_ParamSet_prikey = `
-----BEGIN PRIVATE KEY-----
MEUCAQAwHAYGKoUDAgITMBIGByqFAwICIwEGByqFAwICHgEEIgIgRhUDJ1WQASIf
nx+aUM2eagzV9dCt6mQ5wdtenr2ZS/Y=
-----END PRIVATE KEY-----
`
var GostR3410_2001_CryptoPro_A_ParamSet_pubkey = `
-----BEGIN PUBLIC KEY-----
MGMwHAYGKoUDAgITMBIGByqFAwICIwEGByqFAwICHgEDQwAEQORQaJaqv4S10bz4
jw112dGlrtD+DyGR8TqkhmOvlJB46VUIbpBsEHs8nn0pXtzsIfEwgV8Oxo/QA0Ri
Qu5j7SU=
-----END PUBLIC KEY-----
`

func Test_Check_GostR3410_2001_CryptoPro_A_ParamSet(t *testing.T) {
    pri := decodePEM(GostR3410_2001_CryptoPro_A_ParamSet_prikey)
    if len(pri) == 0 {
        t.Error("decodePEM prikey empty")
    }

    prikey, err := gost.ParsePrivateKey(pri)
    if err != nil {
        t.Fatal(err)
    }

    pubkey := &prikey.PublicKey

    pub, err := gost.MarshalPublicKey(pubkey)
    if err != nil {
        t.Fatal(err)
    }

    pubPem := encodePEM(pub, "PUBLIC KEY")
    if len(pubPem) == 0 {
        t.Error("make pub error")
    }

    // ========

    pub2 := decodePEM(GostR3410_2001_CryptoPro_A_ParamSet_pubkey)
    if len(pub2) == 0 {
        t.Error("decodePEM pub empty")
    }

    pubkey2, err := gost.ParsePublicKey(pub2)
    if err != nil {
        t.Fatal(err)
    }

    if !pubkey2.Equal(pubkey) {
        t.Error("parse pubkey fail")
    }
}

var Openssl_Gost_Prikey = `
-----BEGIN PRIVATE KEY-----
MEYCAQAwHwYIKoUDBwEBAQEwEwYHKoUDAgIjAQYIKoUDBwEBAgIEIJ3L20nIrfUo
MdMNKTx9pxh3e7Etf7abOI73mypFZToK
-----END PRIVATE KEY-----
`

func Test_Check_Openssl_Gost_Prikey(t *testing.T) {
    pri := decodePEM(Openssl_Gost_Prikey)
    if len(pri) == 0 {
        t.Error("decodePEM prikey empty")
    }

    prikey, err := gost.ParsePrivateKey(pri)
    if err != nil {
        t.Fatal(err)
    }

    pubkey := &prikey.PublicKey

    pub, err := gost.MarshalPublicKey(pubkey)
    if err != nil {
        t.Fatal(err)
    }

    pubPem := encodePEM(pub, "PUBLIC KEY")
    if len(pubPem) == 0 {
        t.Error("make pub error")
    }
}
