package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/skip2/go-qrcode"
	"gopkg.in/yaml.v3"
)

type Config struct {
	PublicKey string `yaml:"publicKey" json:"publicKey"`
	Users     []User `yaml:"users" json:"users"`
}

type Secrets struct {
	PrivateKey string `yaml:"privateKey" json:"privateKey"`
}

type User struct {
	Code  string `yaml:"code" json:"code"`
	Name  string `yaml:"name" json:"name"`
	Role  string `yaml:"role" json:"role"`
	Image string `yaml:"image" json:"image"`
}

// UpdateYAML reads a YAML file, updates a key with a new value, and saves it back
func UpdateYAML() error {

	// Generate a new Ed25519 key pair
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	secrets := Secrets{
		PrivateKey: hex.EncodeToString(privKey),
	}

	newData, err := yaml.Marshal(&secrets)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}

	// Write back to file
	if err := os.WriteFile("secrets.yaml", newData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Read file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse YAML into a map
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	// Update the value
	config.PublicKey = hex.EncodeToString(pubKey)

	// Marshal back to YAML
	newDataConfig, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}

	// Write back to file
	if err := os.WriteFile("config.yaml", newDataConfig, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func GenerateHTML() error {

	data, err := os.ReadFile("template/index.html.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	t, err := template.New("config").Parse(string(data))
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Read file
	data_, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse YAML into a map
	var config Config
	if err := yaml.Unmarshal(data_, &config); err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	// Execute template with data, writing to stdout
	var rendered bytes.Buffer
	if err := t.Execute(&rendered, config); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	output := rendered.String()

	// Write back to file
	if err := os.WriteFile("index.html", []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil

}

func GenerateQRCodes() error {

	data_, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse YAML into a map
	var config Config
	if err := yaml.Unmarshal(data_, &config); err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	dataSecrets, err := os.ReadFile("secrets.yaml")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse YAML into a map
	var secrets Secrets
	if err := yaml.Unmarshal(dataSecrets, &secrets); err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	privateKey, err := hex.DecodeString(secrets.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	for _, user := range config.Users {
		signature := ed25519.Sign(privateKey, []byte(user.Code))
		hexString := hex.EncodeToString(signature)
		err = qrcode.WriteFile(fmt.Sprintf("%s-%s", user.Code, hexString), qrcode.Medium, 256, fmt.Sprintf("qrcodes/%s.png", user.Code))
		if err != nil {
			log.Fatal("Failed to create QR code:", err)
		}
	}

	return nil
}

func main() {
	err := UpdateYAML()
	if err != nil {
		log.Fatalf("Error updating YAML: %v", err)
	}
	err = GenerateHTML()
	if err != nil {
		log.Fatalf("Error Generating HTML: %v", err)
	}
	err = GenerateQRCodes()
	if err != nil {
		log.Fatalf("Error Generating HTML: %v", err)
	}
	fmt.Println("YAML updated successfully!")
}
