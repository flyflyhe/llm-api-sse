package aes

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"log"
	"testing"
)

func TestGenerateKeyIv(t *testing.T) {
	key, iv := GenerateKeyIv()
	log.Println("key:", key)
	log.Println("iv:", iv)
}

func TestService_Crypt(t *testing.T) {
	service := &Service{
		Key: []byte("12345678901234567890123456789012"),
		Iv:  []byte("abcdefghijklmnop"),
	}

	content := "hi hello 18137373259"
	text, err := service.Crypt(content)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("加密串", text)

	if decode, err := service.Decrypt(text); err != nil {
		t.Error(err)
	} else {
		fmt.Println("解密串", decode)
		assert.Equal(t, content, decode)
	}
}
