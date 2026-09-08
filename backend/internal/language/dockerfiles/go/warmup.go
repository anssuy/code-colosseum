package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	var v interface{}
	data, _ := io.ReadAll(os.Stdin)
	json.Unmarshal(data, &v)
	out, _ := json.Marshal(v)
	fmt.Println(string(out))
}
