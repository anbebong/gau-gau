package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net"
)

// MessageType defines the type of message exchanged between client and server.
type MessageType string

const (
	HELLO MessageType = "HELLO"
	MSG   MessageType = "MSG"
	CMD   MessageType = "CMD"
	AUTH  MessageType = "AUTH"
	// Có thể mở rộng thêm các loại bản tin khác ở đây
)

// Message represents a structured message for communication.
type Message struct {
	Type     MessageType `json:"type"`
	Content  string      `json:"content"`
	Key      string      `json:"key,omitempty"`
	ClientID string      `json:"clientID,omitempty"`
	// Có thể bổ sung thêm các trường khác như Sender, Timestamp, v.v.
}

// SendMessage sends a message to the specified address using the given network type (TCP/UDP).
func SendMessage(network, address string, msg Message) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

// ReceiveMessage listens for incoming messages on the specified network and address.
func ReceiveMessage(network, address string) (Message, error) {
	var msg Message

	listener, err := net.Listen(network, address)
	if err != nil {
		return msg, err
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		return msg, err
	}
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	err = decoder.Decode(&msg)
	return msg, err
}

// SendMessageConn sends a Message qua kết nối TCP đã mở sẵn
func SendMessageConn(conn net.Conn, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

// ReceiveMessageConn nhận Message từ kết nối TCP đã mở sẵn
func ReceiveMessageConn(conn net.Conn) (Message, error) {
	var msg Message
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return msg, err
	}
	err = json.Unmarshal(buffer[:n], &msg)
	return msg, err
}

// Encrypt encrypts plain text string with AES and returns base64 string
func Encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	b := []byte(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], b)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64 string with AES and returns plain text
func Decrypt(key []byte, cryptoText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

// GenerateRandomKey sinh key ngẫu nhiên (16 bytes, hex)
func GenerateRandomKey() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
