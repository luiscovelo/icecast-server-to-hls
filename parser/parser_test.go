package parser_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/luiscovelo/icecast-server-to-hls/parser"
)

func TestMatch(t *testing.T) {
	strings := []string{
		"PUT /name HTTP/1.1",
		"GET /status-json.xsl HTTP/1.1",
		"GET /7.xsl HTTP/1.1",
	}

	pattern := `^(PUT|GET) /([^ ]+) HTTP/1\.1$`

	regex := regexp.MustCompile(pattern)

	for _, str := range strings {
		match := regex.FindStringSubmatch(str)
		if match != nil {
			// match[1] contém o método (PUT ou GET)
			// match[2] contém a parte variável (:name, status-json.xsl, ou 7.xsl)
			fmt.Printf("Método: %s, Parte Variável: %s\n", match[1], match[2])
		} else {
			fmt.Println("Não corresponde:", str)
		}
	}
}

func TestPar(t *testing.T) {
	data := []byte("PUT /stream HTTP/1.1\r\nAuthorization: Basic c291cmNlOjEyMw==\r\nHost: localhost:8090\r\nUser-Agent: butt 0.1.40\r\nContent-Type: audio/mpeg\r\nice-name: no name\r\nice-public: 0\r\nice-audio-info: ice-bitrate=128;ice-channels=2;ice-samplerate=44100\r\nExpect: 100-continue\r\n\r\n")

	req := parser.Parse(data)

	log.Println("Authorization", req.Header["Authorization"])
}
