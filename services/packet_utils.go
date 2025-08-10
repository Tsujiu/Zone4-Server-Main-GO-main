package services

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"strings"
)

// ReadNextMessage lê do reader e:
// - retorna (json, false, nil) quando montou um JSON completo (balanceando {} com suporte a strings/escapes)
// - retorna (nil, true, nil) quando recebeu keepalive '0' (tratado e respondido)
// - retorna (nil, false, nil) quando veio lixo não-JSON: só loga prévia e continua
// - retorna erro != nil quando a conexão caiu
func ReadNextMessage(r *bufio.Reader, c *Client) ([]byte, bool, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, false, err
	}

	// Keepalive '0' (não derruba; responde opcionalmente)
	if b != '{' {
		if b == '0' {
			if c != nil && c.Writer != nil {
				_ = c.Writer.WriteByte('0')
				_ = c.Writer.WriteByte('\n')
				_ = c.Writer.Flush()
			}
			return nil, true, nil
		}
		// Log de prévia para depurar protocolos diferentes
		preview := make([]byte, 16)
		n, _ := r.Read(preview)
		fmt.Printf("⚠️ Primeiro byte != '{' (0x%02X). Prévia(hex): %s (addr=%s)\n",
			b, strings.ToUpper(hex.EncodeToString(preview[:n])), c.Addr)
		return nil, false, nil
	}

	// Montagem de JSON com balanceamento e suporte a strings/escapes
	depth := 1
	inStr := false
	esc := false
	buf := []byte{'{'} // já temos a primeira chave

	for {
		ch, e := r.ReadByte()
		if e != nil {
			return nil, false, e
		}
		buf = append(buf, ch)

		if inStr {
			if esc {
				esc = false
			} else if ch == '\\' {
				esc = true
			} else if ch == '"' {
				inStr = false
			}
			continue
		}
		if ch == '"' {
			inStr = true
			continue
		}

		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
			if depth == 0 {
				break
			}
		}
	}
	return buf, false, nil
}
