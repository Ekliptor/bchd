package pingserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type PingServer struct {
	ctx  context.Context
	port int
}

func NewPingServer(ctx context.Context, port int) *PingServer {
	return &PingServer{
		ctx:  ctx,
		port: port,
	}
}

// StartServer starts the server and listens on the specified port.
// This call is blocking, so it should be called in a new goroutine.
func (srv *PingServer) StartServer() error {
	s := &http.Server{
		Addr: fmt.Sprintf(":%d", srv.port),
		//Handler:     cors(router),
		Handler:     http.NewServeMux(), // responds with: 404 page not found
		ReadTimeout: 30 * time.Second,
		//ErrorLog:    goLog.New(log.GetWriter(), "HTTP: ", 0),
	}

	done := make(chan struct{})
	go func() {
		<-srv.ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error on HTTP server shutdown: %+v\n", err)
		}
		close(done)
	}()

	fmt.Printf("Serving HTTP ping at port: %d\n", srv.port)
	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		fmt.Printf("Error starting %s server: %+v\n", "HTTP", err)
	}
	<-done

	return nil
}
