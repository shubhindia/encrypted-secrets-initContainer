package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/shubhindia/encrypted-secrets/pkg/providers"
	"github.com/shubhindia/encrypted-secrets/pkg/providers/utils"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"

	secretsv1alpha1 "github.com/shubhindia/encrypted-secrets/api/v1alpha1"
)

func init() {

	_ = secretsv1alpha1.AddToScheme(scheme.Scheme)
}
func main() {
	k8sClient, err := utils.GetKubeClient()

	if err != nil {
		fmt.Printf("failed to get kubeclient %v", err)
	}

	log.Print("getting encrypted secret")
	raw, err := k8sClient.RESTClient().
		Get().
		AbsPath("/apis/secrets.shubhindia.xyz/v1alpha1/namespaces/dev/encryptedsecrets/encryptedsecret-sample").
		DoRaw(context.TODO())

	if err != nil {
		fmt.Printf("failed to get encrypted secret %v", err)
	}

	log.Print("decoding encrypted secret")

	fmt.Printf("raw data %s", raw)

	codecs := serializer.NewCodecFactory(scheme.Scheme, serializer.EnableStrict)
	obj, _, err := codecs.UniversalDeserializer().Decode(raw, &schema.GroupVersionKind{
		Group:   secretsv1alpha1.GroupVersion.Group,
		Version: "v1alpha1",
		Kind:    "EncryptedSecret",
	}, nil)
	if err != nil {
		if ok, _ := regexp.MatchString("no kind(.*)is registered for version", err.Error()); ok {
			log.Printf("no kind is registered for version")
		}
		panic(err)
	}
	encryptedSecret, ok := obj.(*secretsv1alpha1.EncryptedSecret)
	if !ok {
		// should never happen
		panic("failed to convert runtimeObject to encryptedSecret")
	}

	log.Print("decrypting encrypted secret")
	decryptedObj, err := providers.DecodeAndDecrypt(encryptedSecret)
	if err != nil {
		fmt.Printf("failed to decrypt value for %s", err.Error())
	}
	// write data to file with key as filename
	var basePath = "/tmp/opt/conf/"

	err = os.MkdirAll(basePath, 0755)
	if err != nil {
		fmt.Printf("failed to create directory %s", err.Error())
	}

	log.Println("writing decrypted secret to file")
	for key, value := range decryptedObj.Data {
		log.Println("writing file for key", key)
		err = os.WriteFile(basePath+key, []byte(value), 0644)
		if err != nil {
			fmt.Printf("failed to write file %s", err.Error())
		}
	}
}
