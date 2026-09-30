package models

// EndpointType defines the connection medium of a candidate route.
type EndpointType string

const (
	EndpointLAN    EndpointType = "lan"
	EndpointPublic EndpointType = "public"
	EndpointTunnel EndpointType = "tunnel"
)

// CandidateEndpoint is a single network route (LAN, Public IP, or Cloudflare Tunnel) to reach the daemon.
type CandidateEndpoint struct {
	Type  EndpointType `json:"type"`
	URL   string       `json:"url"`
	Label string       `json:"label"`
}

// PairingPayload is serialized into the terminal QR code scanned by the mobile app.
type PairingPayload struct {
	ServerName string              `json:"serverName"`
	Token      string              `json:"token"`
	Candidates []CandidateEndpoint `json:"candidates"`
	CreatedAt  int64               `json:"createdAt"`
}
