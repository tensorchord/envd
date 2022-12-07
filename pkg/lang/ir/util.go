package ir

import "encoding/json"

func (rg *RuntimeGraph) Dump() (string, error) {
	b, err := json.Marshal(rg)
	if err != nil {
		return "", nil
	}
	runtimeGraphCode := string(b)
	return runtimeGraphCode, nil
}

func (rg *RuntimeGraph) Load(code []byte) error {
	var newrg *RuntimeGraph
	err := json.Unmarshal(code, newrg)
	if err != nil {
		return err
	}
	rg.RuntimeCommands = newrg.RuntimeCommands
	rg.RuntimeDaemon = newrg.RuntimeDaemon
	rg.RuntimeEnviron = newrg.RuntimeEnviron
	rg.RuntimeExpose = newrg.RuntimeExpose
	return nil
}
