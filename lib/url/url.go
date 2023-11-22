package url

import (
	"fmt"
	urllib "net/url"
	"strings"
)

func toURLWithEncodedPath(rawURL string) string {
	// 1. Check for "://". If it doesn't exist, return since we wont be able to parse it
	if !strings.Contains(rawURL, "://") {
		return rawURL
	}

	// 2. Remove fragment, if any
	rawWithoutFragment, fragment, _ := strings.Cut(rawURL, "#")

	// 3. Remove queryParams, if any
	rawWithoutQuery, queryParams, _ := strings.Cut(rawWithoutFragment, "?")

	// 4. string.Cut '://' to isolate hostname+path
	scheme, hostnameAndPath, _ := strings.Cut(rawWithoutQuery, "://")

	// 5. split the path by '/'
	hostnamePathSplit := strings.Split(hostnameAndPath, "/")

	// 6. If there isn't both hostname and path, return since there is nothing to encode
	if len(hostnamePathSplit) <= 1 {
		return rawURL
	}
	hostname, pathParts := hostnamePathSplit[0], hostnamePathSplit[1:]

	encodedPathParts := make([]string, len(pathParts))
	for i, part := range pathParts {
		encodedPathParts[i] = urllib.PathEscape(part)
	}

	// 7. reensemble everything
	result := []string{
		scheme,
		"://",
		hostname,
	}

	if len(encodedPathParts) > 0 {
		result = append(result, "/", strings.Join(encodedPathParts, "/"))
	}

	if queryParams != "" {
		result = append(result, "?", queryParams)
	}

	if fragment != "" {
		result = append(result, "#", fragment)
	}

	return strings.Join(result, "")
}

// SanitizeURL takes a raw URL string and ensures it is a valid, properly formatted URL.
//
// Parameters:
// - rawURL: A string representing the URL to be sanitized.
//
// Returns:
// - A sanitized URL string if the input is valid and meets the criteria.
// - An error if the input string is not a valid URL or if the URL scheme is not supported (i.e., not HTTP or HTTPS).
// Note: This function does not modify the URL's host, path, query parameters, or fragment.
func SanitizeURL(rawURL string) (string, error) {
	url, err := urllib.Parse(toURLWithEncodedPath(rawURL))

	if err != nil {
		return "", fmt.Errorf("received string %s is not a valid URL: %s", rawURL, err)
	}

	if url.Scheme != "https" && url.Scheme != "http" {
		return "", fmt.Errorf("%s scheme is not supported", url.Scheme)
	}

	if url.Hostname() == "" {
		return "", fmt.Errorf("received string %s does not have a hostname", rawURL)
	}

	return url.String(), nil
}
