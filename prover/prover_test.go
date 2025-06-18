package prover

import (
	"encoding/json"
	"fmt"
	"github.com/yetanotherco/go-circom-prover-verifier/parsers"
	"github.com/yetanotherco/go-circom-prover-verifier/verifier"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCircuitsGenerateProof(t *testing.T) {
	testCircuitGenerateProof(t, "circuit1k") // 1000 constraints
	testCircuitGenerateProof(t, "circuit5k") // 5000 constraints
	// testCircuitGenerateProof(t, "circuit10k") // 10000 constraints
	// testCircuitGenerateProof(t, "circuit20k") // 20000 constraints
}

func testCircuitGenerateProof(t *testing.T, circuit string) {
	// Using json provingKey file:
	// provingKeyJson, err := ioutil.ReadFile("../testdata/" + circuit + "/proving_key.json")
	// require.Nil(t, err)
	// pk, err := parsers.ParsePk(provingKeyJson)
	// require.Nil(t, err)
	// witnessJson, err := ioutil.ReadFile("../testdata/" + circuit + "/witness.json")
	// require.Nil(t, err)
	// w, err := parsers.ParseWitness(witnessJson)
	// require.Nil(t, err)

	// Using bin provingKey file:
	// pkBinFile, err := os.Open("../testdata/" + circuit + "/proving_key.bin")
	// require.Nil(t, err)
	// defer pkBinFile.Close()
	// pk, err := parsers.ParsePkBin(pkBinFile)
	// require.Nil(t, err)

	// Using go bin provingKey file:
	pkGoBinFile, err := os.Open("../testdata/" + circuit + "/proving_key.go.bin")
	require.Nil(t, err)
	defer pkGoBinFile.Close()
	pk, err := parsers.ParsePkGoBin(pkGoBinFile)
	require.Nil(t, err)

	witnessBinFile, err := os.Open("../testdata/" + circuit + "/witness.bin")
	require.Nil(t, err)
	defer witnessBinFile.Close()
	w, err := parsers.ParseWitnessBin(witnessBinFile)
	require.Nil(t, err)

	beforeT := time.Now()
	proof, pubSignals, err := GenerateProof(pk, w)
	assert.Nil(t, err)
	fmt.Println("proof generation time for "+circuit+" elapsed:", time.Since(beforeT))

	proofStr, err := parsers.ProofToJson(proof)
	assert.Nil(t, err)

	err = ioutil.WriteFile("../testdata/"+circuit+"/proof.json", proofStr, 0644)
	assert.Nil(t, err)
	publicStr, err := json.Marshal(parsers.ArrayBigIntToString(pubSignals))
	assert.Nil(t, err)
	err = ioutil.WriteFile("../testdata/"+circuit+"/public.json", publicStr, 0644)
	assert.Nil(t, err)

	// verify the proof
	vkJson, err := ioutil.ReadFile("../testdata/" + circuit + "/verification_key.json")
	require.Nil(t, err)
	vk, err := parsers.ParseVk(vkJson)
	require.Nil(t, err)

	v := verifier.Verify(vk, proof, pubSignals)
	assert.True(t, v)

	// to verify the proof with snarkjs:
	// snarkjs verify --vk testdata/circuitX/verification_key.json -p testdata/circuitX/proof.json --pub testdata/circuitX/public.json
}

func BenchmarkGenerateProof(b *testing.B) {
	// benchmark with a circuit of 10000 constraints
	provingKeyJson, err := ioutil.ReadFile("../testdata/circuit5k/proving_key.json")
	require.Nil(b, err)
	pk, err := parsers.ParsePk(provingKeyJson)
	require.Nil(b, err)

	witnessJson, err := ioutil.ReadFile("../testdata/circuit5k/witness.json")
	require.Nil(b, err)
	w, err := parsers.ParseWitness(witnessJson)
	require.Nil(b, err)

	for i := 0; i < b.N; i++ {
		GenerateProof(pk, w)
	}
}
