package controller

import (
	"fmt"
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"

	corev1 "k8s.io/api/core/v1"
	kaasv1 "github.com/faithByte/KaaS/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const MY_NAMESPACE = "kaas-mpi-jobs"

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

func (r *JobReconciler) createSshSecret(ctx context.Context, job kaasv1.Job) error {
	name := fmt.Sprintf("ssh-%s", job.GetUID())

	var isCreated corev1.Secret
	err := r.Get(ctx, client.ObjectKey{Namespace: MY_NAMESPACE, Name: name}, &isCreated)
	if err == nil {
		return errors.NewAlreadyExists(schema.GroupResource{}, "");
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
			Namespace: MY_NAMESPACE,
		},
		StringData: map[string]string{
			"ssh-prvkey": prvKey,
			"ssh-pubkey": pubKey,
		},
	}

	if err := r.Create(ctx, secret); err != nil {
		return err
	}

	return nil
}
