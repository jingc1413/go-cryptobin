package hc

import (
    "bytes"
    "testing"
    "math/rand"
    "encoding/hex"
)

func Test_Cipher128(t *testing.T) {
    random := rand.New(rand.NewSource(99))
    max := 100

    var encrypted [16]byte
    var decrypted [16]byte

    for i := 0; i < max; i++ {
        key := make([]byte, 16)
        random.Read(key)
        iv := make([]byte, 16)
        random.Read(iv)
        value := make([]byte, 16)
        random.Read(value)

        cipher1, err := NewCipher128(key, iv)
        if err != nil {
            t.Fatal(err.Error())
        }

        cipher1.XORKeyStream(encrypted[:], value)

        cipher2, err := NewCipher128(key, iv)
        if err != nil {
            t.Fatal(err.Error())
        }

        cipher2.XORKeyStream(decrypted[:], encrypted[:])

        if !bytes.Equal(decrypted[:], value[:]) {
            t.Errorf("encryption/decryption failed: % 02x != % 02x\n", decrypted, value)
        }
    }
}

func Test_Cipher256(t *testing.T) {
    random := rand.New(rand.NewSource(99))
    max := 100

    var encrypted [16]byte
    var decrypted [16]byte

    for i := 0; i < max; i++ {
        key := make([]byte, 32)
        random.Read(key)
        iv := make([]byte, 32)
        random.Read(iv)
        value := make([]byte, 16)
        random.Read(value)

        cipher1, err := NewCipher256(key, iv)
        if err != nil {
            t.Fatal(err.Error())
        }

        cipher1.XORKeyStream(encrypted[:], value)

        cipher2, err := NewCipher256(key, iv)
        if err != nil {
            t.Fatal(err.Error())
        }

        cipher2.XORKeyStream(decrypted[:], encrypted[:])

        if !bytes.Equal(decrypted[:], value[:]) {
            t.Errorf("encryption/decryption failed: % 02x != % 02x\n", decrypted, value)
        }
    }
}

func fromHex(s string) []byte {
    h, _ := hex.DecodeString(s)
    return h
}

type testData struct {
    key []byte
    iv []byte
    pt []byte
    ct []byte
}

func Test_Check256(t *testing.T) {
   tests := []testData{
        {
           fromHex("0000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("0000000008000000000000000000000000000000000000000000000000000000"),
           fromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("B89065FE0B458C64FD6EDC6A893C8C8183578E7D37BE97E6FF82E45110A2596049A817CDE859B67B56CB80768D6DD2756EC368FBABC35C8B51C62AC92F913281"),
        },
        {
           fromHex("0000000000000000000000000000000000000000200000000000000000000000"),
           fromHex("0000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("689E48A85A40BD161BEC710F9B2457FD276F1156EBC10BB851A8517AFDBD692DE4827BAAFF218AF886439ED976147EBBBB1074BD599A80F6324C87BAC987B8C5"),
        },
        {
           fromHex("0000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("0000000000000000000000000000000000000000200000000000000000000000"),
           fromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("B090CC267B29A95ADFAF6BE3E147D64721ECACBF6B7D0C4061D17FB7DE0A66626D6F9FC167FB3FFF237C240AA03FAD5513B6DA848F22796DB501A8FB89F2B85D"),
        },
        {
           fromHex("3F404142434445464748494A4B4C4D4E4F505152535455565758595A5B5C5D5E"),
           fromHex("0000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("804025A410EFFBA58647A9F4B443BFC61CDDC30CA04DA8DAB3EC6A098A830D682683B59B76C60C09938E67CB41385315E2504B024DB808923B0909EFC25F0927"),
        },
        {
           fromHex("0F62B5085BAE0154A7FA4DA0F34699EC3F92E5388BDE3184D72A7DD02376C91C"),
           fromHex("288FF65DC42B92F960C72E95FC63CA3198FF66CD349B0269D0379E056CD33AA1"),
           fromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("C8632038DA61679C4685288B37D3E2327BC2D28C266B041FE0CA0D3CFEED8FD5753259BAB757168F85EA96ADABD823CA4684E918423E091565713FEDDE2CCFE0"),
        },
    }

    for i, test := range tests {
        c, err := NewCipher256(test.key, test.iv)
        if err != nil {
            t.Fatal(err.Error())
        }

        tmp := make([]byte, len(test.pt))
        c.XORKeyStream(tmp, test.pt)

        if !bytes.Equal(tmp, test.ct) {
            t.Errorf("[%d] Check error: got %x, want %x", i, tmp, test.ct)
        }

        // ===========

        c2, err := NewCipher256(test.key, test.iv)
        if err != nil {
            t.Fatal(err.Error())
        }

        tmp2 := make([]byte, len(test.ct))
        c2.XORKeyStream(tmp2, test.ct)

        if !bytes.Equal(tmp2, test.pt) {
            t.Errorf("[%d] Check Decrypt error: got %x, want %x", i, tmp2, test.pt)
        }
    }
}

func Test_Check128(t *testing.T) {
   tests := []testData{
        {
           fromHex("00002000000000000000000000000000"),
           fromHex("00000000000000000000000000000000"),
           fromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("94DAB13AE0D2F9A65283C6AE987331101C4EE45EC812AD67DDF3D1F026B51B172D366C7E3B2D55E5AE7010A279D35B0383B77E96C6B2434C3E6DDC2401D64AEC"),
        },
        {
           fromHex("00000000000000000000000000080000"),
           fromHex("00000000000000000000000000000000"),
           fromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("2EABC4033A51B3901B6340BE32F808EEA319582F21A7CF6633570E82AC879B603E438847D9E3719EAB71F8E3247FEFA5C07B2282AA2FA80CEFFA8E076304FEBA"),
        },
        {
           fromHex("ABACADAEAFB0B1B2B3B4B5B6B7B8B9BA"),
           fromHex("00000000000000000000000000000000"),
           fromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("505F4C9084D6F5C640C214EFED9E2DF08EEF8241ACAE98072B5B3EDB72F1687D586B2569DC7F58DED2C2BCD134CB6CF3D80A7A879D7878C080A5BAD5ABA1DCCF"),
        },
        {
           fromHex("0A5DB00356A9FC4FA2F5489BEE4194E7"),
           fromHex("00000000000000000000000000000000"),
           fromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("203C7A9050F5E4F98CB72D913B8E7FB9BB2635F8ECCDBFCD231B4EDCA96A24A99F71BDD76CE42B982228ADCF9385C702C2A767488DF42D5DBD8DF2884225367B"),
        },
        {
           fromHex("0F62B5085BAE0154A7FA4DA0F34699EC"),
           fromHex("288FF65DC42B92F960C72E95FC63CA31"),
           fromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
           fromHex("1CD8AEDDFE52E217E835D0B7E84E2922D04B1ADBCA53C4522B1AA604C42856A90AF83E2614BCE65C0AECABDD8975B55700D6A26D52FFF0888DA38F1DE20B77B7"),
        },
    }

    for i, test := range tests {
        c, err := NewCipher128(test.key, test.iv)
        if err != nil {
            t.Fatal(err.Error())
        }

        tmp := make([]byte, len(test.pt))
        c.XORKeyStream(tmp, test.pt)

        if !bytes.Equal(tmp, test.ct) {
            t.Errorf("[%d] Check error: got %x, want %x", i, tmp, test.ct)
        }

        // ===========

        c2, err := NewCipher128(test.key, test.iv)
        if err != nil {
            t.Fatal(err.Error())
        }

        tmp2 := make([]byte, len(test.ct))
        c2.XORKeyStream(tmp2, test.ct)

        if !bytes.Equal(tmp2, test.pt) {
            t.Errorf("[%d] Check Decrypt error: got %x, want %x", i, tmp2, test.pt)
        }
    }
}
