package middlewares

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Make sure the client closed connection won't trigger
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				log.WithField("detail", err).Errorf("%s", formatRecover(4))

				if brokenPipe {
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}

			}
		}()
		c.Next()
	}
}

// Format panic source in the form of runtime.Stack.
// "skip" stands for the number of the frames to skip before identifing the source of panic
func formatRecover(skip int) (buffer []byte) {
	pc := make([]uintptr, 10)
	depth := runtime.Callers(skip, pc)

	buf := bytes.Buffer{}
	buf.WriteString(getGoroutineState() + "\n")
	var lines []string
	for i := 0; i < depth; i++ {
		fn := runtime.FuncForPC(pc[i])
		if fn != nil {
			file, line := fn.FileLine(pc[i])
			lines = append(lines, fmt.Sprintf("%s()\n\t%s:%d +%#x", fn.Name(), file, line, fn.Entry()))

		}
	}
	buf.WriteString(strings.Join(lines, "\n"))
	buffer = buf.Bytes()
	return
}

func getGoroutineState() string {
	stack := make([]byte, 64)
	stack = stack[:runtime.Stack(stack, false)]
	stack = stack[:bytes.Index(stack, []byte("\n"))]

	return string(stack)
}
