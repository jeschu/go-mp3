package engine

import "time"

type serverBuilder struct {
	addr            string
	managementAddr  string
	shutdownTimeout time.Duration
	withRequestId   bool
}

func (serverBuilder *serverBuilder) Addr(addr string) *serverBuilder {
	serverBuilder.addr = addr
	return serverBuilder
}
func (serverBuilder *serverBuilder) ManagementAddr(managementAddr string) *serverBuilder {
	serverBuilder.managementAddr = managementAddr
	return serverBuilder
}
func (serverBuilder *serverBuilder) ShutdownTimeout(shutdownTimeout time.Duration) *serverBuilder {
	serverBuilder.shutdownTimeout = shutdownTimeout
	return serverBuilder
}
func (serverBuilder *serverBuilder) WithRequestId(withRequestId bool) *serverBuilder {
	serverBuilder.withRequestId = withRequestId
	return serverBuilder
}
func (serverBuilder *serverBuilder) Build() (*Server, error) { return newServer(*serverBuilder) }
