package config

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var AppConfig map[string]string
var AESKey []byte
var AESIV []byte
var Conn net.Conn
var Addr *net.UDPAddr
var ConnUDP *net.UDPConn

// ---------- Utils ----------
func envOrDefault(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func EnvIntDefault(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func parsePorts(s string) []int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	out := []int{}
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.Contains(p, "-") {
			r := strings.SplitN(p, "-", 2)
			if len(r) != 2 {
				continue
			}
			start, err1 := strconv.Atoi(strings.TrimSpace(r[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(r[1]))
			if err1 == nil && err2 == nil && start > 0 && end >= start {
				for x := start; x <= end; x++ {
					out = append(out, x)
				}
			}
		} else {
			if n, err := strconv.Atoi(p); err == nil && n > 0 {
				out = append(out, n)
			}
		}
	}
	// dedup preservando ordem
	seen := map[int]bool{}
	final := []int{}
	for _, n := range out {
		if !seen[n] {
			seen[n] = true
			final = append(final, n)
		}
	}
	return final
}

// ---------- API de canais/portas ----------
func ChannelPorts() []int {
	ports := []int{}

	// Single
	if v := os.Getenv("CHANNEL_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			ports = append(ports, n)
		}
	}

	// Lista/ranges
	if v := os.Getenv("CHANNEL_PORTS"); v != "" {
		ports = append(ports, parsePorts(v)...)
	}

	// Alias opcional
	if v := os.Getenv("GAME_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			ports = append(ports, n)
		}
	}

	if len(ports) == 0 {
		ports = []int{9090}
	}

	// dedup final
	seen := map[int]bool{}
	final := []int{}
	for _, n := range ports {
		if !seen[n] {
			seen[n] = true
			final = append(final, n)
		}
	}
	return final
}

func ChannelPort(index int) int {
	ps := ChannelPorts()
	if index < 0 || index >= len(ps) {
		return ps[0]
	}
	return ps[index]
}

func DebugPortsSummary() string {
	ps := ChannelPorts()
	ss := make([]string, 0, len(ps))
	for _, p := range ps {
		ss = append(ss, strconv.Itoa(p))
	}
	idx := EnvIntDefault("CHANNEL_INDEX", 0)
	return fmt.Sprintf("Canais disponíveis: [%s]  |  CHANNEL_INDEX=%d  |  Em uso: %d",
		strings.Join(ss, ", "), idx, ChannelPort(idx))
}

// ---------- Loaders ----------
func LoadConfig() {
	_ = godotenv.Load()

	AppConfig = make(map[string]string)

	// DSNs
	AppConfig["REDIS_ADDR"] = os.Getenv("REDIS_ADDR")
	AppConfig["REDIS_USER"] = os.Getenv("REDIS_USER")
	AppConfig["REDIS_PASS"] = os.Getenv("REDIS_PASS")
	AppConfig["SQLSERVER_GAME"] = os.Getenv("SQLSERVER_GAME")
	AppConfig["SQLSERVER_GAME_TEST"] = os.Getenv("SQLSERVER_GAME_TEST")
	AppConfig["SQLSERVER_GAME_INVENTORY"] = os.Getenv("SQLSERVER_GAME_INVENTORY")
	AppConfig["SQLSERVER_GAME_INVENTORT_TEST"] = os.Getenv("SQLSERVER_GAME_INVENTORT_TEST")
	AppConfig["SQLSERVER_GAME_RECORD"] = os.Getenv("SQLSERVER_GAME_RECORD")

	// Portas base
	AppConfig["TCP_PORT"] = envOrDefault("TCP_PORT", "9090")
	AppConfig["UDP_PORT"] = envOrDefault("UDP_PORT", "9091")
	AppConfig["PPROF_PORT"] = envOrDefault("PPROF_PORT", "6060")

	// Canal (single + lista + índice)
	AppConfig["CHANNEL_PORT"] = envOrDefault("CHANNEL_PORT", "9090")
	AppConfig["CHANNEL_PORTS"] = strings.TrimSpace(os.Getenv("CHANNEL_PORTS"))
	AppConfig["GAME_PORT"] = strings.TrimSpace(os.Getenv("GAME_PORT"))
	AppConfig["CHANNEL_INDEX"] = envOrDefault("CHANNEL_INDEX", "0")

	// Outros
	AppConfig["WORKER_POOL"] = os.Getenv("WORKER_POOL")
	AppConfig["MAX_MATCHING_WORKER"] = os.Getenv("MAX_MATCHING_WORKER")
	AppConfig["DB_USER"] = os.Getenv("DB_USER")
	AppConfig["DB_PASSWORD"] = os.Getenv("DB_PASSWORD")
	AppConfig["DB_NAME"] = os.Getenv("DB_NAME")
	AppConfig["DB_SERVER"] = os.Getenv("DB_SERVER")
	AppConfig["DB_PORT"] = os.Getenv("DB_PORT")
	AppConfig["MAX_USER_CHANNEL"] = os.Getenv("MAX_USER_CHANNEL")

	// Cripto
	keyString := strings.TrimSpace(os.Getenv("AES_KEY"))
	if keyString == "" {
		log.Fatal("❌ AES_KEY not set in .env")
	}
	AESKey = []byte(keyString)

	aesIVHex := strings.TrimSpace(os.Getenv("AES_IV"))
	if aesIVHex == "" {
		log.Fatal("❌ AES_IV not set in .env")
	}
	if len(aesIVHex) != 32 {
		log.Fatalf("❌ Invalid AES_IV length: got %d, want 32 hex chars", len(aesIVHex))
	}
	aesIV, err := hex.DecodeString(aesIVHex)
	if err != nil {
		log.Fatalf("❌ Invalid AES_IV format: %v", err)
	}
	if len(aesIV) != aes.BlockSize {
		log.Fatalf("❌ AES_IV must be %d bytes, got %d", aes.BlockSize, len(aesIV))
	}
	AESIV = aesIV
}

func LoadConfigTestLocal() {
	AppConfig = make(map[string]string)
	AppConfig["REDIS_ADDR"] = "localhost:6379"
	AppConfig["SQLSERVER_DSN"] = "sqlserver://server01:Server01@localhost:1433?database=game"
	AppConfig["SQLSERVER_GAME"] = "sqlserver://server01:Server01@localhost:1433?database=game"

	// Portas padrão locais
	AppConfig["TCP_PORT"] = "9090"
	AppConfig["UDP_PORT"] = "9091"
	AppConfig["PPROF_PORT"] = "6060"

	// Canal local: 9090 + lista de canais
	AppConfig["CHANNEL_PORT"] = "9090"
	AppConfig["CHANNEL_PORTS"] = "29998,29996,29995,29993,29994,29992"
	AppConfig["CHANNEL_INDEX"] = "0"

	// Chaves demo (somente local)
	AESKey = []byte("p*{Ilqw<8AT_@poI2Kq3D1uVcp`*@bRh")
	ivParts := strings.Split("235,72,71,0,201,74,178,207,129,184,192,91,50,78,209,100", ",")
	var iv []byte
	for _, part := range ivParts {
		val, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			log.Fatal("Invalid AES_IV format:", err)
		}
		iv = append(iv, byte(val))
	}
	AESIV = iv
}
