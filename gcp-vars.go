package go_gcp_vars

import (
	"cloud.google.com/go/compute/metadata"
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

	// Only attempt to check the Cloud Run metadata server if it looks like
	// the service is deployed to Cloud Run or GOOGLE_CLOUD_PROJECT not already set.
	if project == "" || service != "???" {
		var err error
		if project, err = metadata.ProjectID(); err != nil {
			log.Printf("metadata.ProjectID: Cloud Run metadata server: %v", err)
		}
		var ip string
		if ip, err = metadata.InternalIP(); err == nil {
			log.Printf("internal IP: %s\n", ip)
		}
		if ip, err = metadata.ExternalIP(); err == nil {
			log.Printf("external IP: %s\n", ip)
		}
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
