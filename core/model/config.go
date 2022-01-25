package model

type JsonData map[string]interface{}

// Config the main config structure
type Config struct {
	AuthKeys                string
	InternalApiKey          string
	CoreServiceRegLoaderURL string
	ContentServiceURL       string
}
