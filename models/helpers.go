package models

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/chuxorg/chux-models/errors"
)

func ExtractCompanyName(urlString string) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", errors.NewChuxModelsError("ExtractCompanyName() Unable to parse the url", err)
	}

	host := parsedURL.Host
	// Split the host into parts
	parts := strings.Split(host, ".")

	// If there are at least two parts (subdomain(s) and domain)
	if len(parts) >= 2 {
		// Return the second last part, which is the domain without the extension
		return parts[len(parts)-2], nil
	}
	msg := fmt.Sprintf("Could not extract company name from url: %s", urlString)
	return "", errors.NewChuxModelsError(msg, nil)
}
