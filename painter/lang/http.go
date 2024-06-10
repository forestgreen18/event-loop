package lang

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

// CommandHttpHandler constructs an HTTP request handler that takes data from the request and passes it to CommandProcessor,
// then sends the resulting list of operations to painter.Loop.
func CommandHttpHandler(loop *painter.EventLoop, cp *CommandProcessor) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var input io.Reader = r.Body
		if r.Method == http.MethodGet {
			input = strings.NewReader(r.URL.Query().Get("cmd"))
		}

		operations, err := cp.ProcessCommands(input)
		if err != nil {
			log.Printf("Error processing script: %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, op := range operations {
			loop.Enqueue(op) // Enqueue each operation individually
		}
		rw.WriteHeader(http.StatusOK)
	})
}
