package internal

import (
	"flag"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

func Kubeconfig() string {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	return *kubeconfig
}

func KubeconfigList(configPath string) map[string]*api.Context {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: configPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	return config.Contexts
}

func MapToString(data map[string]string) string {
	var result string
	for key, value := range data {
		result += key + ":" + value + "\n"
	}
	return result
}

// func encrypt(text string, key []byte) (string, error) {
//     plaintext := []byte(text)

//     block, err := aes.NewCipher(key)
//     if err != nil {
//         return "", err
//     }

//     aesGCM, err := cipher.NewGCM(block)
//     if err != nil {
//         return "", err
//     }

//     nonce := make([]byte, aesGCM.NonceSize())
//     if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
//         return "", err
//     }

//     ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
//     return base64.URLEncoding.EncodeToString(ciphertext), nil
// }

// func decrypt(encryptedText string, key []byte) (string, error) {
//     enc, err := base64.URLEncoding.DecodeString(encryptedText)
//     if err != nil {
//         return "", err
//     }

//     block, err := aes.NewCipher(key)
//     if err != nil {
//         return "", err
//     }

//     aesGCM, err := cipher.NewGCM(block)
//     if err != nil {
//         return "", err
//     }

//     nonceSize := aesGCM.NonceSize()
//     if len(enc) < nonceSize {
//         return "", fmt.Errorf("ciphertext too short")
//     }

//     nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
//     plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
//     if err != nil {
//         return "", err
//     }

//     return string(plaintext), nil
// }
