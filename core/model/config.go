package model

// JSONData wrapper struct
type JSONData map[string]interface{}

// Config the main config structure
type Config struct {
	AuthKeys          string
	InternalAPIKey    string
	CoreBBHost        string
	ContentServiceURL string
}
