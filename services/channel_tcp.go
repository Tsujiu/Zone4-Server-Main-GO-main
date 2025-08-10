package services

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func StartChannelTCP(port int) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("❌ Falha ao abrir porta do canal %d: %v", port, err)
	}
	log.Printf("✅ Canal TCP ouvindo em %s", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("⚠️ Accept: %v", err)
			continue
		}
		go handleChannelConn(conn)
	}
}

func handleChannelConn(c net.Conn) {
	defer c.Close()
	_ = c.SetDeadline(time.Now().Add(30 * time.Second))
	remote := c.RemoteAddr().String()
	log.Printf("👤 Cliente conectado: %s", remote)

	r := bufio.NewReader(c)

	// Consome pequeno handshake (2 bytes) sem ecoar
	_, _ = r.Peek(2)

	// Responde lista de servidores em JSON + \n
	resp := BuildServerListJSON()
	resp = append(resp, '\n')
	if _, err := c.Write(resp); err != nil {
		log.Printf("⚠️ Write server list: %v", err)
		return
	}
	log.Printf("📤 Enviado server list (%d bytes) para %s", len(resp), remote)
}
