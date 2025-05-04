package secrets

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"

	ctrl "sigs.k8s.io/controller-runtime"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/faithByte/kaas/internal/controller/utils"
)

func generateSshKey() (string, string, error) {
	// Private Key generation -------------------------
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	err = privateKey.Validate()
	if err != nil {
		return "", "", err
	}

	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}

	prvKey := string(pem.EncodeToMemory(&privBlock))
	// prvKey := string(encodePrivateKeyToPEM(privateKey))

	// Public Key generation -------------------------
	publicRsaKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	pubKey := string(ssh.MarshalAuthorizedKey(publicRsaKey))

	return prvKey, pubKey, nil
}

func CreateSshSecret(jobData *utils.JobData) error {
	name := "ssh-" + string(jobData.Job.GetUID())

	var isCreated corev1.Secret
	err := jobData.Client.Get(jobData.Context, client.ObjectKey{Namespace: utils.MY_NAMESPACE, Name: name}, &isCreated)
	if err == nil {
		return errors.NewAlreadyExists(schema.GroupResource{}, "")
	} else if !errors.IsNotFound(err) {
		return err
	}

	prvKey, pubKey, err := generateSshKey()
	if err != nil {
		return err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: utils.MY_NAMESPACE,
		},
		StringData: map[string]string{
			"ssh-prvkey": prvKey,
			"ssh-pubkey": pubKey,
		},
	}

	if err := ctrl.SetControllerReference(&jobData.Job, secret, jobData.Scheme); err != nil {
		return err
	}

	if err := jobData.Client.Create(jobData.Context, secret); err != nil {
		return err
	}

	return nil
}