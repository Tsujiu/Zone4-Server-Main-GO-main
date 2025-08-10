package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"maxion-zone4/config"
	"maxion-zone4/services"
)

func imprimirInfo(canalPort int) {
	fmt.Println("================================")
	fmt.Printf("🟢 Canal está iniciando\n")
	fmt.Printf("🎮 Porta do Jogo: %d\n", canalPort)
	fmt.Println("================================")
	log.Printf("🩺 Debug/Métricas em http://127.0.0.1:%s ( /metrics /healthz /debug/pprof/ )",
		config.AppConfig["PPROF_PORT"])
	log.Println("ℹ️ ", config.DebugPortsSummary())
}

func main() {
	flagPort := flag.Int("port", -1, "Porta do canal/jogo")
	flagIndex := flag.Int("index", -1, "Índice em CHANNEL_PORTS (0..n-1)")
	flag.Parse()

	config.LoadConfig()

	var canalPort int
	switch {
	case *flagPort > 0:
		canalPort = *flagPort
	case *flagIndex >= 0:
		canalPort = config.ChannelPort(*flagIndex)
	default:
		canalPort = config.ChannelPort(config.EnvIntDefault("CHANNEL_INDEX", 0))
	}

	imprimirInfo(canalPort)
	go services.StartChannelTCP(canalPort)

	for {
		time.Sleep(10 * time.Second)
	}
}
