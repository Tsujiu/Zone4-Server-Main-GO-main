package services

import (
	"encoding/json"
	"fmt"
	"time"
)

// routeJSON faz unmarshal só do cabeçalho (op) e delega para o handler tipado.
func routeJSON(c *Client, raw []byte) error {
	var base Base
	if err := json.Unmarshal(raw, &base); err != nil {
		return fmt.Errorf("JSON inválido: %w", err)
	}

	switch base.Op {
	case OpLogin:
		var req LoginReq
		if err := json.Unmarshal(raw, &req); err != nil {
			return fmt.Errorf("login inválido: %w", err)
		}
		return handleLogin(c, req)

	case OpMove:
		var req MoveReq
		if err := json.Unmarshal(raw, &req); err != nil {
			return fmt.Errorf("move inválido: %w", err)
		}
		return handleMove(c, req)

	case OpChat:
		var req ChatReq
		if err := json.Unmarshal(raw, &req); err != nil {
			return fmt.Errorf("chat inválido: %w", err)
		}
		return handleChat(c, req)

	default:
		return fmt.Errorf("operação desconhecida: %q", base.Op)
	}
}

// ==== Handlers (integre sua lógica real aqui) ====

func handleLogin(c *Client, req LoginReq) error {
	if req.Username == "" || req.Password == "" {
		return WriteJSON(c.Writer, ErrorResp{Op: OpError, Error: "usuário ou senha vazios", Ts: time.Now().UnixMilli()})
	}

	// TODO: autenticação real; por enquanto só aceita e atribui PlayerID
	c.PlayerID = req.Username

	// Sala padrão 0, ou a informada no request
	room := req.RoomID
	if room == 0 && req.RoomID == 0 {
		room = 0
	}
	JoinRoom(c, room)

	// ACK para o cliente
	return WriteJSON(c.Writer, AckResp{
		Op:     OpAck,
		Ok:     true,
		Ts:     time.Now().UnixMilli(),
		RoomID: room,
	})
}

func handleMove(c *Client, req MoveReq) error {
	// Opcional: permitir override do room no pacote (normalmente não precisa)
	if req.RoomID != 0 && req.RoomID != c.RoomID {
		JoinRoom(c, req.RoomID)
	}

	// TODO: chamar seu sistema de simulação (posições, colisão, zonas, etc.)
	// Exemplo simples: aplica um fator pequeno ao vetor de direção
	pos := Vec3{
		X: req.Dir.X * 0.1,
		Y: req.Dir.Y * 0.1,
		Z: req.Dir.Z * 0.1,
	}
	state := StateResp{
		Op:         OpState,
		PlayerID:   firstNonEmpty(req.PlayerID, c.PlayerID),
		ServerTick: uint64(time.Now().UnixMilli()),
		Pos:        pos,
		RoomID:     c.RoomID,
	}

	// Broadcast do estado para TODOS na sala (incluindo o próprio)
	BroadcastToRoom(c.RoomID, state, "") // use c.Addr para excluir o remetente

	return nil
}

func handleChat(c *Client, req ChatReq) error {
	// Se o pacote trouxer room_id e diferir, movimenta o cliente
	if req.RoomID != 0 && req.RoomID != c.RoomID {
		JoinRoom(c, req.RoomID)
	}

	msg := ChatMsg{
		Op:       OpMsg,
		PlayerID: firstNonEmpty(req.PlayerID, c.PlayerID),
		Text:     req.Text,
		Ts:       time.Now().UnixMilli(),
		RoomID:   c.RoomID,
	}
	// Envia para todos MENOS o remetente (evita eco); troque "" para incluir o remetente
	BroadcastToRoom(c.RoomID, msg, c.Addr)
	return nil
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
