package internal

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

func GetCacheFileNameByCtxMap(ctxMap map[string]interface{}, kubeContext string) string {

	oidcMap := ctxMap["oidcMap-"+kubeContext].(map[string][]string)
	oidcIssuerUrl := oidcMap["oidc-issuer-url"][0]
	oidcClientId := oidcMap["oidc-client-id"][0]
	oidcExtraScopes := oidcMap["oidc-extra-scope"]
	oidcClientSecret := ""
	if oidcMap["oidc-client-secret"] != nil {
		oidcClientSecret = oidcMap["oidc-client-secret"][0]
	}
	key := Key{
		IssuerURL:    oidcIssuerUrl,
		ClientID:     oidcClientId,
		ExtraScopes:  oidcExtraScopes,
		ClientSecret: oidcClientSecret,
	}
	filename, _ := ComputeFilename(key)
	return filename

}

func UnzipAndSave(r *http.Request, w http.ResponseWriter, dest string) error {
	file, _, err := r.FormFile("file")

	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.(io.Seeker).Seek(0, io.SeekEnd)

	if err != nil {
		return err
	}
	_, err = file.(io.Seeker).Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(file, stat)
	if err != nil {
		return err
	}
	dest = dest + "/" + zipReader.File[0].Name

	destExists, err := exists(dest)
	if err != nil {
		return err
	}
	if destExists {
		err = os.RemoveAll(dest)
		if err != nil {
			return err
		}
	}

	for _, f := range zipReader.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if strings.HasPrefix(f.Name, "__MACOSX/") || strings.HasSuffix(f.Name, ".DS_Store") {
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to handle errors
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
