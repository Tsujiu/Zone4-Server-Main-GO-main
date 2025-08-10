package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	readBufSize  = 64 * 1024
	writeBufSize = 64 * 1024
	maxJSONSize  = 512 * 1024
)

// StartTCPListener inicia o servidor TCP usando TCP_HOST/TCP_PORT do ambiente.
// Ex.: TCP_HOST=127.0.0.1  TCP_PORT=9090 (compatível com o cliente local)
func StartTCPListener() {
	host := getenv("TCP_HOST", "127.0.0.1")
	port := getenv("TCP_PORT", "9090")
	addr := net.JoinHostPort(host, port)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("❌ Erro ao iniciar servidor TCP:", err)
		return
	}
	fmt.Println("✅ Servidor TCP ativo em", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("❌ Erro ao aceitar conexão:", err)
			continue
		}

		if tc, ok := conn.(*net.TCPConn); ok {
			_ = tc.SetNoDelay(true)
			_ = tc.SetKeepAlive(true)
			_ = tc.SetKeepAlivePeriod(30 * time.Second)
			_ = tc.SetReadBuffer(readBufSize)
			_ = tc.SetWriteBuffer(writeBufSize)
		}

		br := bufio.NewReaderSize(conn, readBufSize)
		bw := bufio.NewWriterSize(conn, writeBufSize)
		c := &Client{
			Addr:   conn.RemoteAddr().String(),
			Conn:   conn,
			Writer: bw,
			// PlayerID vazio até login
			// RoomID padrão = 0 (lobby)
		}
		RegisterClient(c)

		fmt.Println("🔗 Novo cliente conectado:", conn.RemoteAddr())
		go handleConnectionListener(c, br)
	}
}

func handleConnectionListener(c *Client, br *bufio.Reader) {
	defer func() {
		// Limpa salas/cliente e fecha a conexão
		UnregisterClient(c)
		fmt.Println("🔌 Cliente desconectado:", c.Conn.RemoteAddr())
		_ = c.Conn.Close()
	}()

	for {
		raw, err := readOneJSONObject(br)
		if err != nil {
			if !errors.Is(err, net.ErrClosed) && !strings.Contains(err.Error(), "closed") {
				fmt.Println("⚠️ Erro ao ler dados do cliente:", err)
			}
			return
		}

		// Roteia com base no campo "op"
		if err := routeJSON(c, raw); err != nil {
			fmt.Println("⚠️ Erro ao processar mensagem:", err)
			_ = WriteJSON(c.Writer, ErrorResp{Op: OpError, Error: err.Error(), Ts: time.Now().UnixMilli()})
		}
	}
}

// readOneJSONObject lê UM objeto JSON completo do stream TCP (sem depender de '\n').
func readOneJSONObject(br *bufio.Reader) ([]byte, error) {
	var buf bytes.Buffer
	inString := false
	escaped := false
	depth := 0
	started := false

	for {
		if buf.Len() > maxJSONSize {
			return nil, fmt.Errorf("JSON muito grande")
		}
		b, err := br.ReadByte()
		if err != nil {
			return nil, err
		}

		if !started {
			if b <= ' ' {
				continue
			}
			if b != '{' {
				return nil, fmt.Errorf("esperado '{', mas recebido %q", b)
			}
			started = true
			depth = 1
			buf.WriteByte(b)
			continue
		} else {
			buf.WriteByte(b)
		}

		if inString {
			if escaped {
				escaped = false
				continue
			}
			if b == '\\' {
				escaped = true
			} else if b == '"' {
				inString = false
			}
			continue
		}

		if b == '"' {
			inString = true
			continue
		}

		switch b {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return buf.Bytes(), nil
			}
		}
	}
}

// WriteJSON exportado para outros arquivos do pacote services
func WriteJSON(bw *bufio.Writer, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	// Adiciona \n para facilitar logs e leitores NDJSON
	b = append(b, '\n')
	if _, err := bw.Write(b); err != nil {
		return err
	}
	return bw.Flush()
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
