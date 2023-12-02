package env

import (
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Addr           string
	ManagementAddr string
	WithRequestId  bool
}

func InitFromEnv() Config {
	switch strings.ToLower(stringWithDefault("color", "")) {
	case "force":
		gin.ForceConsoleColor()
	case "disable":
		gin.DisableConsoleColor()
	}

	switch strings.ToLower(stringWithDefault("mode", "")) {
	case gin.DebugMode:
		gin.SetMode(gin.DebugMode)
	case gin.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	}

	requestId, _ := boolWithDefault("requestId", true)

	addr := ":8080"
	if envAddr := os.Getenv("ADDR"); envAddr != "" {
		addr = envAddr
	}
	mgmtAddr := ":8081"
	if envMgmtAddr := os.Getenv("MGMT_ADDR"); envMgmtAddr != "" {
		mgmtAddr = envMgmtAddr
	} else {
		if host, port, err := net.SplitHostPort(addr); err == nil {
			if prt, err := strconv.Atoi(port); err == nil {
				mgmtAddr = net.JoinHostPort(host, strconv.Itoa(prt+1))
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
	return Config{Addr: addr, ManagementAddr: mgmtAddr, WithRequestId: requestId}
}

var env map[string]string

func init() {
	env = make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		env[pair[0]] = pair[1]
	}
}

func stringWithDefault(key, defaultValue string) string {
	if v, ok := env[key]; ok {
		return v
	} else {
		return defaultValue
	}
}

func boolWithDefault(key string, defaultValue bool) (bool, error) {
	if v, ok := env[key]; ok {
		if b, err := strconv.ParseBool(v); err != nil {
			return defaultValue, err
		} else {
			return b, nil
		}
	} else {
		return defaultValue, nil
	}
}
