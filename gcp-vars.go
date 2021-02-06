package go_gcp_vars

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path"
)

type GCPVars struct {
	ProjectId   string
	ServiceName string
	Port        string
	Debug       bool
}

func init() {
	loadEnvFile()
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func loadEnvFile() {
	debug := len(os.Getenv("DEBUG")) > 0
	if !debug {
		return
	}
	var filename string
	if fileExists(".env") {
		filename = ".env"
	} else {
		var dir, _ = os.Getwd()
		for {
			dir = path.Join(dir, "..")
			filename = path.Join(dir, ".env")
			if fileExists(filename) {
				break
			}
			if dir == "" {
				break
			}
		}
	}
	_ = godotenv.Load(filename)
}

func GetGCPVars() *GCPVars {
	service := os.Getenv("K_SERVICE")
	if service == "" {
		service = os.Getenv("GAE_SERVICE")
		if service == "" {
			service = "???"
		}

	}

	revision := os.Getenv("K_REVISION")
	if revision == "" {
		revision = "???"
	}

	project := ""

	// Environment variable GOOGLE_CLOUD_PROJECT is only set locally.
	// On Cloud Run, strip the timestamp prefix from log entries.
	if project == "" {
		log.SetFlags(0)
	}

	if project == "" {
		project = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	port := os.Getenv("PORT")
	return &GCPVars{
		ProjectId:   project,
		ServiceName: service,
		Port:        port,
		Debug:       len(os.Getenv("DEBUG")) > 0,
	}
}
