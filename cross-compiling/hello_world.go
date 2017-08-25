package cross_compiling

import "fmt"
import "runtime"

func main() {
        fmt.Printf("Hello %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
