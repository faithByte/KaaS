package utils

import (
	"encoding/json"
	"fmt"

	// "github.com/faithByte/kaas/internal/controller/scheduler"
	"github.com/google/go-cmp/cmp"
)

func PrintObj(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return
	}

	fmt.Println(string(bytes))
}

func PrintDiff(obj1, obj2 interface{}) {
	fmt.Println(cmp.Diff(obj1, obj2))
}
