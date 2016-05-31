package goldengine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"sort"
	"time"

	sf "github.com/manyminds/gosfml"
)

//SceneDefEntity : Defines an Entity within the scene
type SceneDefEntity struct {
	Name               string
	Parent             string
	Prefab             string
	TransformArguments map[string]interface{}
	Position           Vector
	Scale              Vector
	Rotation           float32
}

//EntityFromSceneDefEntity : Creates an Entity from a scene definition entity
func EntityFromSceneDefEntity(def SceneDefEntity) (*Entity, error) {
	if def.Name == "" {
		return nil, fmt.Errorf("All Prefabs require a Name. %v", def)
	}
	if def.Prefab == "" {
		return nil, fmt.Errorf("%s Requires field \"Prefab\"", def.Name)
	}
	prefab, ok := PrefabRegister.Get(def.Prefab)
	if !ok {
		return nil, fmt.Errorf("Cannot find Prefab with the name: %s", def.Prefab)
	}
	for k, v := range def.TransformArguments {
		prefab.Transformer.Arguments[k] = v
	}
	entity := EntityFromEntityPrefab(prefab)
	entity.Name = def.Name
	return entity, nil
}

//SceneDef : Informaiton to Create a scene from a JSON Template
type SceneDef struct {
	Name     string
	Entities []SceneDefEntity
}

//ParseSceneDef : Get SceneDef from string definition
func ParseSceneDef(g *Game, sceneName, sceneString string) (*SceneDef, error) {
	//SceneFuncMap : Functions that can be used in a Scene Def template
	var SceneFuncMap = template.FuncMap{
		"gameHeight": func() float32 {
			pixelHeight := float32(g.GetWindow().renderWindow.GetSize().Y)
			unit := pixelHeight / scale
			return unit
		},
		"gameWidth": func() float32 {
			return 100
		},
		"divide": func(a, b float32) float32 {
			return a / b
		},
		"subtract": func(a, b float32) float32 {
			return a - b
		},
	}

	t := template.Must(template.New(sceneName).Funcs(SceneFuncMap).Parse(sceneString))
	var rawJSON bytes.Buffer
	err := t.Execute(&rawJSON, g)
	if err != nil {
		return nil, err
	}
	var sceneDef SceneDef
	err = json.Unmarshal(rawJSON.Bytes(), &sceneDef)
	if err != nil {
		return nil, err
	}
	if sceneDef.Name == "" {
		return nil, fmt.Errorf("Scenes Must have a name")
	}
	return &sceneDef, err
}

//SceneFromSceneDef : Creates a new scene given a definition
func SceneFromSceneDef(def *SceneDef) (*Scene, error) {
	entityMap := make(map[uint32]*Entity)
	entityNodeMap := make(map[string]*entityNode)
	entityDefMap := make(map[uint32]SceneDefEntity)
	scene := Scene{
		Name: def.Name,
		root: &entityNode{
			entity:   NewEntity(),
			children: make([]*entityNode, 0),
		},
		entityMap:     entityMap,
		entityNodeMap: entityNodeMap,
		entityDefMap:  entityDefMap,
	}
	entityNodeMap[RootNodeName] = scene.root
	scene.root.scene = &scene
	for i, entityDef := range def.Entities {
		//Create Entity
		if _, ok := entityNodeMap[entityDef.Name]; ok {
			return nil, fmt.Errorf("Entities in SceneDef must have a unique name: %s", entityDef.Name)
		}
		if entityDef.Name == RootNodeName {
			return nil, fmt.Errorf("Entities in SceneDef can't thave the name: %s", RootNodeName)
		}
		entity, err := EntityFromSceneDefEntity(entityDef)
		if err != nil {
			return nil, err
		}
		entity.scene = &scene
		entityMap[entity.id] = entity
		node := entityNode{
			scene:  &scene,
			entity: entity,
		}
		entityNodeMap[entityDef.Name] = &node
		entityDefMap[entity.id] = def.Entities[i]
	}
	//Node Operations
	for _, entityDef := range def.Entities {
		//Node Operations
		node := entityNodeMap[entityDef.Name]
		entity := entityMap[node.entity.id]
		//Add Parent and Children
		if entityDef.Name == "" {
			scene.root.children = append(scene.root.children, node)
		} else {
			parentNode, ok := entityNodeMap[entityDef.Parent]
			if !ok {
				parentNode = entityNodeMap[RootNodeName]
			}
			if parentNode.children == nil {
				parentNode.children = make([]*entityNode, 1)
			}
			parentNode.children = append(parentNode.children, node)
			node.parent = parentNode

			parent, ok := entityMap[parentNode.entity.id]
			if ok {
				if parent.children == nil {
					parent.children = make(map[uint32]*Entity)
				}
				parent.children[entity.id] = parent
				entity.parent = parent
			}
		}
		//Set Transform Properties
		entity.Transfrom.SetPosition(entityDef.Position.ToSFML())
		if entityDef.Scale != ZeroVector {
			entity.Transfrom.SetScale(entityDef.Scale.ToSFML())
		}
		entity.Transfrom.SetRotation(entityDef.Rotation)
	}
	//Scene Operations

	return &scene, nil
}

type entityNode struct {
	scene    *Scene
	parent   *entityNode
	zIndex   int
	entity   *Entity
	children []*entityNode
}

func (node *entityNode) AddChild(child *entityNode) {
	if node.children == nil {
		node.children = make([]*entityNode, 0)
	}
	node.children = append(node.children, child)
}

func (node *entityNode) Start() {
	if node == nil {
		return
	}
	if node.entity != nil {

		node.entity.Start()
	}
	for _, child := range node.children {
		c := child
		c.Start()

	}
}

func (node *entityNode) Awake() {
	if node == nil {
		return
	}
	if node.entity != nil {

		node.entity.Awake()
	}
	fmt.Println("Waking Up entitynode	")
	for _, child := range node.children {
		c := child
		c.Awake()
	}
}

func (node *entityNode) Update(dur time.Duration) {
	if node == nil {
		return
	}
	if node.entity != nil {

		node.entity.Update(dur)
	}
	for _, child := range node.children {
		c := child
		c.Update(dur)
	}
}

func (node *entityNode) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	if node.entity != nil && node.entity.Transfrom != nil {
		// transform := node.entity.Transfrom.GetTransform()
		// combinedTransform := renderStates.Transform.Combine(&transform)
		// renderStates.Transform = *combinedTransform
		//fmt.Println(node.entity.Transfrom.GetPosition())
		target.Draw(node.entity.Transfrom, renderStates)
	}

	for _, child := range node.children {
		child.Draw(target, renderStates)
	}
}

//SceneLoadedMSG : Scene Was loaded returns pointer to scene
const SceneLoadedMSG = MessageType("SceneLoaded")

//SceneChangedMSG : Game has changed scene
const SceneChangedMSG = MessageType("SceneChanged")

//SceneAddedEntityMSG : Added Entity to scene. Sends pointer to entity
const SceneAddedEntityMSG = MessageType("SceneAddedEntity")

//Scene : Everything that is being rendered
type Scene struct {
	Name          string
	root          *entityNode
	game          *Game
	entityMap     map[uint32]*Entity
	entityNodeMap map[string]*entityNode
	entityDefMap  map[uint32]SceneDefEntity

	Awake  func()
	Start  func()
	Update func(time.Duration)

	BasicMailBox
}

//AddEntity : AddEntity to the Scene
func (s *Scene) AddEntity(e *Entity) {
	if _, ok := s.entityMap[e.id]; ok {
		return
	}
	node := &entityNode{
		entity: e,
	}
	s.root.children = append(s.root.children, node)
	s.entityNodeMap[e.Name] = node
	s.entityMap[e.id] = e
	if e.parent != nil {
		parent, ok := s.entityNodeMap[e.parent.Name]
		if ok {
			parent.AddChild(node)
		} else {
			s.root.AddChild(node)
		}
	} else {
		s.root.AddChild(node)
	}
	for _, child := range e.children {
		s.AddEntity(child)
	}
	s.PostMessage(Message{
		Message: SceneAddedEntityMSG,
		Content: e,
	})
}

//GetEntities : Gets all entities in scene
func (s *Scene) GetEntities() []*Entity {
	list := make([]*Entity, len(s.entityMap))
	counter := 0
	for _, e := range s.entityMap {
		list[counter] = e
		counter = counter + 1
	}
	return list
}

//RecalculateScale : Changes the of every entity
func (s *Scene) RecalculateScale() {
	for _, e := range s.entityMap {
		e.RecalculateScale()
	}
}

//GetEntityByName : Returns an entity in the scene with that name
func (s *Scene) GetEntityByName(name string) (*Entity, bool) {
	node, found := s.entityNodeMap[name]
	return node.entity, found
}

//SetZIndex : SetsZIndex of a node
func (s *Scene) SetZIndex(name string, index int) {
	node, found := s.entityNodeMap[name]
	if found {
		node.zIndex = index
	}
}

//Draw : Draws the scene to a render target
func (s *Scene) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	entities := make([]*entityNode, len(s.entityNodeMap))
	counter := 0
	for _, e := range s.entityNodeMap {
		entities[counter] = e
		counter = counter + 1
	}
	sort.Sort(byZIndex(entities))
	for _, e := range entities {
		if e.entity != nil {

			target.Draw(e.entity.Transfrom, renderStates)
		}
	}
}

//Start : Starts all entities in the scene
func (s *Scene) start() {
	if s.Start != nil {
		s.Start()
	}
	s.root.Start()
}

//Awake : Wakes up all entities in the scene
func (s *Scene) awake() {
	if s.Awake != nil {
		s.Awake()
	}
	s.root.Awake()
}

//Update : Updates all entities in node
func (s *Scene) update(dur time.Duration) {
	if s.Update != nil {
		s.Update(dur)
	}
	s.root.Update(dur)
}

//RecieveMessage : Handles Message
func (s *Scene) RecieveMessage(msg Message) {

}

//RootNodeName : Default Name for The rootnode of all scenes
const RootNodeName = "ROOT"

type byZIndex []*entityNode

func (s byZIndex) Len() int {
	return len(s)
}
func (s byZIndex) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byZIndex) Less(i, j int) bool {
	return s[i].zIndex < s[j].zIndex
}
