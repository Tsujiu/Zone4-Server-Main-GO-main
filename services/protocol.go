package services

// Operações (op) aceitas no protocolo JSON
const (
	OpLogin = "login"
	OpMove  = "move"
	OpChat  = "chat"

	OpAck   = "ack"
	OpState = "state"
	OpMsg   = "msg"
	OpError = "error"
)

// Base para roteamento
type Base struct {
	Op string `json:"op"`
}

// Tipos de requisição
type LoginReq struct {
	Op       string `json:"op"` // "login"
	Username string `json:"username"`
	Password string `json:"password"`
	// Sala/opcional. Se não vier, caímos na sala 0.
	RoomID int `json:"room_id,omitempty"`
}

type MoveReq struct {
	Op         string `json:"op"` // "move"
	PlayerID   string `json:"player_id"`
	Dir        Vec3   `json:"dir"`
	ClientTick uint64 `json:"client_tick"`
	// Opcional: mover explicitamente em outra sala (normalmente não usa)
	RoomID int `json:"room_id,omitempty"`
}

type ChatReq struct {
	Op       string `json:"op"` // "chat"
	PlayerID string `json:"player_id"`
	Text     string `json:"text"`
	RoomID   int    `json:"room_id,omitempty"`
}

// Tipos auxiliares
type Vec3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

// Respostas
type AckResp struct {
	Op     string `json:"op"` // "ack"
	Ok     bool   `json:"ok"`
	Ts     int64  `json:"ts"`
	RoomID int    `json:"room_id,omitempty"`
}

type ErrorResp struct {
	Op    string `json:"op"` // "error"
	Error string `json:"error"`
	Ts    int64  `json:"ts"`
}

type StateResp struct {
	Op         string `json:"op"` // "state"
	PlayerID   string `json:"player_id"`
	ServerTick uint64 `json:"server_tick"`
	Pos        Vec3   `json:"pos"`
	RoomID     int    `json:"room_id"`
}

type ChatMsg struct {
	Op       string `json:"op"` // "msg"
	PlayerID string `json:"player_id"`
	Text     string `json:"text"`
	Ts       int64  `json:"ts"`
	RoomID   int    `json:"room_id"`
}
