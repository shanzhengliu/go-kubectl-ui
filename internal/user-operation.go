package internal

import (
	"context"
	"fmt"

	"k8s.io/client-go/tools/auth"
)

func GetUserInfo(ctx context.Context, namespace string) (string, bool) {
	ctxMap := ctx.Value("map").(map[string]interface{})
	path := ctxMap["configPath"].(string)
	info, _ := auth.LoadFromFile(path)

	fmt.Println(info)
	return "", true
}
