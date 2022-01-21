package model

type JsonData map[string]interface{}

// Config the main config structure
type Config struct {
	AuthKeys                string
	CoreServiceRegLoaderURL string
	ContentServiceURL       string
}
