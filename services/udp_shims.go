package services

import (
	"fmt"
	"log"
)

// Shim mínimo para compilar e manter o fluxo de jogo ativo.
// Ajuste depois para enviar realmente via UDP conforme seu protocolo.

func SendUDP(opcode int, payload string) error {
	// TODO: implementar envio UDP real (serialize + net.UDPConn)
	log.Printf("[UDP SHIM] SendUDP opcode=%d len=%d", opcode, len(payload))
	return nil
}

func SendUDPToPlayer(header interface{}, body interface{}, p interface{}) error {
	// TODO: implementar envio UDP real ao jogador (usar endereço guardado no player)
	log.Printf("[UDP SHIM] SendUDPToPlayer header=%T body=%T player=%T", header, body, p)
	return nil
}

// Opcional: helper para sinalizar que ainda é shim
func NotImplemented(feature string) error {
	return fmt.Errorf("%s: not implemented (UDP shim)", feature)
}
