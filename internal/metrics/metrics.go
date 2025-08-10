package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ConexoesAtuais = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "zone4_conexoes_atuais",
		Help: "Conexões ativas por sala",
	}, []string{"sala"})

	ConexoesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "zone4_conexoes_total",
		Help: "Total de conexões aceitas",
	})

	MensagensTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "zone4_mensagens_total",
		Help: "Total de mensagens processadas por operação e sala",
	}, []string{"op", "sala"})

	BytesLidos = promauto.NewCounter(prometheus.CounterOpts{
		Name: "zone4_bytes_lidos_total",
		Help: "Total de bytes lidos no TCP",
	})

	BytesEscritos = promauto.NewCounter(prometheus.CounterOpts{
		Name: "zone4_bytes_escritos_total",
		Help: "Total de bytes escritos no TCP",
	})
)

func IncConexao(sala int)               { ConexoesTotal.Inc(); ConexoesAtuais.WithLabelValues(labelSala(sala)).Inc() }
func DecConexao(sala int)               { ConexoesAtuais.WithLabelValues(labelSala(sala)).Dec() }
func AddBytesLidos(n int)               { BytesLidos.Add(float64(n)) }
func AddBytesEscritos(n int)            { BytesEscritos.Add(float64(n)) }
func IncMensagem(op string, sala int)   { MensagensTotal.WithLabelValues(op, labelSala(sala)).Inc() }

func labelSala(id int) string {
	if id <= 0 {
		return "desconhecida"
	}
	return strconv.Itoa(id)
}
