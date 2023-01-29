package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func GetCAPath(file string) string {
	home, e := os.UserHomeDir()
	if e != nil {
		panic(e)
	}
	sooproxy := home + "/.sooproxy"

	if _, err := os.Stat(sooproxy); os.IsNotExist(err) {
		if os.MkdirAll(sooproxy, 0755) != nil {
			panic(errors.New("could not create .sooproxy directory"))
		}
	}

	return sooproxy + "/" + file
}

func GenerateCA() error {
	log.Printf("Generating CA")
	// set up our CA certificate
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			CommonName:    "SooProxy CA",
			Organization:  []string{"SooProxy"},
			Country:       []string{"Maldives"},
			Province:      []string{""},
			Locality:      []string{"Dhievhi"},
			StreetAddress: []string{"Majeedhee Magu"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	// pem encode
	caPEM, err := os.Create(GetCAPath("ca.pem"))
	if err != nil {
		return fmt.Errorf("failed to open ca.pem for writing: %s", err)
	}
	log.Printf("Writing CA to ca.pem")
	defer caPEM.Close()

	if err := pem.Encode(caPEM, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes}); err != nil {
		return fmt.Errorf("failed to write ca.pem: %s", err)
	}
	log.Printf("Generating CA key")
	caPrivKeyPEM, err := os.Create(GetCAPath("ca.key"))
	if err != nil {
		return fmt.Errorf("failed to open ca.key for writing: %s", err)
	}
	log.Printf("Writing CA key to ca.key")
	defer caPrivKeyPEM.Close()

	if err := pem.Encode(caPrivKeyPEM, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey)}); err != nil {
		return fmt.Errorf("failed to write ca.key: %s", err)
	}
	log.Printf("CA generated")
	TrustCA()
	return nil
}

func TrustCA() {

	certFile := GetCAPath("ca.pem")
	// Get the current operating system
	os := runtime.GOOS

	// Use a switch statement to determine the command to run
	// based on the current operating system
	var cmd *exec.Cmd
	var cmd2 *exec.Cmd
	switch os {
	case "windows":
		cmd = exec.Command("certutil", "-addstore", "-f", "Root", certFile)
	case "linux":
		cmd = exec.Command("sudo", "cp", certFile, "/usr/local/share/ca-certificates/")
		cmd2 = exec.Command("sudo", "update-ca-trust")
		cmd2.Stdout = cmd.Stdout
		cmd2.Stderr = cmd.Stderr
	case "darwin":
		cmd = exec.Command("security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", certFile)
	default:
		log.Printf("Unsupported operating system")
		return
	}

	// Run the command
	err := cmd.Run()
	if err != nil {
		log.Printf("Could not trust the certificate. You should run me as root perhaps.")
		return
	}

	if os == "linux" {
		err = cmd2.Run()
		if err != nil {
			log.Printf("Could not trust the certificate. You should run me as root perhaps.")
			return
		}
	}

	log.Printf("Certificate added to trusted root store successfully!")
}

func MakeTlsConfig(caCertPath string, caKeyPath string, commonName string) (*tls.Config, error) {
	// Read the CA cert file
	caCertPEM, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	// Decode the PEM-encoded CA cert
	block, _ := pem.Decode(caCertPEM)
	if block == nil {
		return nil, errors.New("failed to parse CA cert PEM")
	}

	// Parse the CA cert
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Read the CA key file
	caKeyPEM, err := ioutil.ReadFile(caKeyPath)
	if err != nil {
		return nil, err
	}

	// Decode the PEM-encoded CA key
	block, _ = pem.Decode(caKeyPEM)
	if block == nil {
		return nil, errors.New("failed to parse CA key PEM")
	}

	// Parse the CA key
	caKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	log.Printf("Generating certificate for %s", commonName)

	// Set up the new certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			CommonName:    commonName,
			Names:         []pkix.AttributeTypeAndValue{{Type: []int{2, 5, 4, 3}, Value: commonName}},
			Organization:  []string{"SooProxy"},
			Country:       []string{"Maldives"},
			Province:      []string{""},
			Locality:      []string{"Dhievhi"},
			StreetAddress: []string{"Majeedhee Magu"},
		},
		DNSNames:    []string{commonName},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(0, 3, 0),

		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	// Create the new cert
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, err
	}

	// Create the tls.Certificate
	tlsCert := tls.Certificate{
		Certificate: [][]byte{certBytes},
		PrivateKey:  caKey,
	}

	// Create the tls.Config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	return tlsConfig, nil
}
