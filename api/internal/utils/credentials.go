package utils

import (
	"checkmate/api/internal/model"
)

// converts a PlatformCredential to SafeCredential
func ConvertToSafeCredential(cred *model.PlatformCredential) model.SafeCredential {
	return model.SafeCredential{
		ID:        cred.ID,
		UserID:    cred.UserID,
		Platform:  cred.Platform,
		CreatedAt: cred.CreatedAt,
	}
}

// converts a slice of PlatformCredentials to SafeCredential
func ConvertToSafeCredentials(creds []model.PlatformCredential) []model.SafeCredential {
	safeCreds := make([]model.SafeCredential, len(creds))
	for i := range creds {
		safeCreds[i] = ConvertToSafeCredential(&creds[i])
	}
	return safeCreds
}
