package radio

import "strings"

type Metadata struct {
	Title  string
	Artist string
}

func ParseMetadata(metadataStr string) Metadata {
	parts := strings.Split(metadataStr, ";")
	metadata := Metadata{}
	for _, part := range parts {
		if strings.HasPrefix(part, "StreamTitle='") {
			title := strings.TrimPrefix(part, "StreamTitle='")
			title = strings.TrimSuffix(title, "'")

			titleParts := strings.SplitN(title, " - ", 2)
			if len(titleParts) == 2 {
				metadata.Artist = titleParts[0]
				metadata.Title = titleParts[1]
			}
		}
	}
	return metadata
}
