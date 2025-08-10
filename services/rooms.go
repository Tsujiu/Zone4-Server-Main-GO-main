package services

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"maxion-zone4/internal/metrics"
)


// Client representa um jogador/conexão ativa
type Client struct {
	net.Conn            // <— EMBED: agora Client implementa net.Conn
	Addr     string
	Writer   *bufio.Writer
	PlayerID string
	RoomID   int
}

// Tabelas globais (com locks)
var (
	Clients     = make(map[string]*Client)      // addr -> *Client
	ClientMutex sync.RWMutex

	Rooms      = make(map[int]map[string]*Client) // roomID -> (addr -> *Client)
	RoomsMutex sync.RWMutex
)

// --- Helpers internos ---

// ensureClientIO garante Addr e Writer válidos (idempotente).
func ensureClientIO(c *Client) {
	if c == nil {
		return
	}
	// se já tem Conn, garante Addr/Writer
	if c.Conn != nil {
		if c.Addr == "" {
			c.Addr = c.RemoteAddr().String() // método promovido
		}
		if c.Writer == nil {
			c.Writer = bufio.NewWriter(c.Conn)
		}
	}
}

// roomHasClient verifica se o cliente está mapeado na sala atual.
func roomHasClient(c *Client) bool {
	if c == nil {
		return false
	}
	RoomsMutex.RLock()
	defer RoomsMutex.RUnlock()
	room, ok := Rooms[c.RoomID]
	if !ok {
		return false
	}
	_, ok = room[c.Addr]
	return ok
}

// --- API pública ---

// RegisterClient registra o cliente recém-conectado
func RegisterClient(c *Client) {
	if c == nil {
		return
	}
	ensureClientIO(c)

	ClientMutex.Lock()
	Clients[c.Addr] = c
	ClientMutex.Unlock()
}

// UnregisterClient remove o cliente e o tira da sala atual
func UnregisterClient(c *Client) {
	if c == nil {
		return
	}

	// remove da sala, se estiver em alguma
	if c.RoomID != 0 && roomHasClient(c) {
		LeaveRoom(c)
	}

	ClientMutex.Lock()
	delete(Clients, c.Addr)
	ClientMutex.Unlock()
}

// JoinRoom coloca o cliente em uma sala (remove da anterior se necessário)
func JoinRoom(c *Client, roomID int) {
	if c == nil {
		return
	}
	ensureClientIO(c)

	// Se já está na mesma sala e mapeado, não faz nada
	if c.RoomID == roomID && roomHasClient(c) {
		return
	}

	// Sair da sala anterior (se estiver em alguma)
	oldRoom := c.RoomID
	if oldRoom != 0 && roomHasClient(c) {
		LeaveRoom(c) // LeaveRoom já ajusta métricas da sala antiga
	}

	// Entrar na nova sala
	RoomsMutex.Lock()
	if _, ok := Rooms[roomID]; !ok {
		Rooms[roomID] = make(map[string]*Client)
	}
	Rooms[roomID][c.Addr] = c
	RoomsMutex.Unlock()

	c.RoomID = roomID
	metrics.IncConexao(roomID)
}

// LeaveRoom remove o cliente da sala atual
func LeaveRoom(c *Client) {
	if c == nil {
		return
	}

	RoomsMutex.Lock()
	oldRoom := c.RoomID
	if oldRoom != 0 {
		if room, ok := Rooms[oldRoom]; ok {
			delete(room, c.Addr)
			if len(room) == 0 {
				delete(Rooms, oldRoom) // limpa sala vazia
			}
		}
	}
	RoomsMutex.Unlock()

	if oldRoom != 0 {
		metrics.DecConexao(oldRoom)
	}

	c.RoomID = 0
}

// BroadcastToRoom envia uma mensagem JSON para todos da sala.
// Se excludeAddr != "", não envia para esse endereço (ex.: evitar eco).
// Retorna a quantidade de destinatários que receberam tentativa de envio.
func BroadcastToRoom(roomID int, payload any, excludeAddr string) int {
	// Snapshot dos clientes da sala para evitar iterar o mapa com lock durante I/O
	RoomsMutex.RLock()
	room := Rooms[roomID]
	if room == nil {
		RoomsMutex.RUnlock()
		return 0
	}

	targets := make([]*Client, 0, len(room))
	for addr, cli := range room {
		if excludeAddr != "" && addr == excludeAddr {
			continue
		}
		targets = append(targets, cli)
	}
	RoomsMutex.RUnlock()

	if len(targets) == 0 {
		return 0
	}

	metrics.IncMensagem("broadcast", roomID)

	count := 0
	for _, cli := range targets {
		if cli == nil || cli.Writer == nil {
			continue
		}
		if err := WriteJSON(cli.Writer, payload); err != nil {
			fmt.Println("⚠️ Erro ao enviar para", cli.Addr, ":", err)
			continue
		}
		count++
	}
	return count
}

// GetClientByConn busca o *Client pelo net.Conn
func GetClientByConn(conn net.Conn) *Client {
	if conn == nil {
		return nil
	}
	addr := conn.RemoteAddr().String()

	ClientMutex.RLock()
	defer ClientMutex.RUnlock()
	return Clients[addr]
}

// GetClientsInRoom retorna uma cópia (snapshot) dos clientes na sala.
func GetClientsInRoom(roomID int) []*Client {
	RoomsMutex.RLock()
	defer RoomsMutex.RUnlock()

	room := Rooms[roomID]
	if room == nil {
		return nil
	}
	out := make([]*Client, 0, len(room))
	for _, c := range room {
		out = append(out, c)
	}
	return out
}
