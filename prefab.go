package goldengine

import (
	"encoding/json"
	"io/ioutil"
)

type prefabRegister struct {
	register map[string]EntityPrefab
}

//PrefabRegister : Has all the EntityPrefab
var PrefabRegister = prefabRegister{
	register: make(map[string]EntityPrefab),
}

func (register *prefabRegister) RegisterFromFile(location string) error {
	dat, err := ioutil.ReadFile(location)
	if err != nil {
		return err
	}
	return register.RegisterFromData(dat)
}

func (register *prefabRegister) RegisterFromData(dat []byte) error {
	var prefab EntityPrefab
	err := json.Unmarshal(dat, &prefab)
	if err != nil {
		return err
	}
	name := prefab.Name
	register.register[name] = prefab
	return nil
}

func (register *prefabRegister) Get(name string) (EntityPrefab, bool) {
	p, ok := register.register[name]
	return p, ok
}
