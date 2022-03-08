package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dadosjusbr/coletores/status"
	"github.com/dadosjusbr/indice"
	"github.com/dadosjusbr/proto/coleta"
	"github.com/dadosjusbr/proto/pipeline"
	"google.golang.org/protobuf/encoding/prototext"
)

func main() {
	// Processa entrada.
	var er pipeline.ResultadoExecucao
	er.Rc = new(coleta.ResultadoColeta)

	erIN, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		status.ExitFromError(status.NewError(4, fmt.Errorf("error reading crawling result: %q", err)))
	}
	if err = prototext.Unmarshal(erIN, er.Rc); err != nil {
		status.ExitFromError(status.NewError(5, fmt.Errorf("error unmarshaling crawling resul from STDIN: %q\n\n %s ", err, string(erIN))))
	}

	// Calcula índice e atualiza proto.
	score := indice.CalcScore(*er.Rc.Metadados)
	er.Rc.Metadados.IndiceCompletude = float32(score.CompletenessScore)
	er.Rc.Metadados.IndiceFacilidade = float32(score.CompletenessScore)
	er.Rc.Metadados.IndiceTransparencia = float32(score.Score)

	// Imprime proto atualizado.
	b, err := prototext.Marshal(&er)
	if err != nil {
		err = status.NewError(status.Unknown, fmt.Errorf("error marshalling execution result with score:%w", err))
		status.ExitFromError(err)
	}
	fmt.Printf("%s", b)
}
