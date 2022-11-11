package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dadosjusbr/coletores/status"
	"github.com/dadosjusbr/indice"
	"github.com/dadosjusbr/proto/coleta"
	"google.golang.org/protobuf/encoding/prototext"
)

func main() {
	// Processa entrada que vem do coletor, que é um resultado de coleta.
	rcIn, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		status.ExitFromError(status.NewError(4, fmt.Errorf("error reading crawling result: %q", err)))
	}
	var rc coleta.ResultadoColeta
	if err = prototext.Unmarshal(rcIn, &rc); err != nil {
		status.ExitFromError(status.NewError(5, fmt.Errorf("error unmarshaling crawling resul from STDIN: %q\n\n %s ", err, string(rcIn))))
	}
	// Define se o formato é aberto
	extensions := []string{"PDF", "ODS", "JSON", "CSV", "HTML", "ODT"}
	for _, extensao := range extensions {
		if rc.Metadados.Extensao.String() == extensao {
			rc.Metadados.FormatoAberto = true
		}
	}
	// Calcula índice e atualiza proto.
	score := indice.CalcScore(*rc.Metadados)
	rc.Metadados.IndiceCompletude = float32(score.CompletenessScore)
	rc.Metadados.IndiceFacilidade = float32(score.EasinessScore)
	rc.Metadados.IndiceTransparencia = float32(score.Score)

	// Imprime resultado de coleta atualizado.
	b, err := prototext.Marshal(&rc)
	if err != nil {
		err = status.NewError(status.Unknown, fmt.Errorf("error marshalling execution result with score:%w", err))
		status.ExitFromError(err)
	}
	fmt.Printf("%s", b)
}
