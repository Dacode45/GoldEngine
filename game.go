package goldengine

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	sf "github.com/manyminds/gosfml"
)

var gRunning = true

//Game : Controls Actual Game
type Game struct {
	name   string
	logger *log.Logger
	window *window
	scenes map[string]*Scene

	GameWidth  float32
	GameHeight float32

	PrefabsFolderName   string
	ResourcesFolderName string
	ScenesFolderName    string
	prefabsFolder       []os.FileInfo
	resourcesFolder     []os.FileInfo
	scenesFolder        []os.FileInfo
}

// GameConfig : Adjust Paramaters for an App
/*
Defaults

Name: App:NoName
LogFile: os.Stdout
*/
type GameConfig struct {
	Name    string
	LogFile io.Writer

	PrefabsFolderName   string
	ResourcesFolderName string
	ScenesFolderName    string
}

const (
	//DefaultPrefabsFolderName : prefabs
	DefaultPrefabsFolderName = "prefabs"
	//DefaultResourcesFolderName : resources
	DefaultResourcesFolderName = "resources"
	//DefaultScenesFolderName : scenes
	DefaultScenesFolderName = "scenes"
)

//NewGame : Returns an app.
func NewGame(config GameConfig, wc WindowConfig) *Game {
	//Set  NewSoundBufferFromSamples
	name := config.Name
	if name == "" {
		name = "APP:NoName"
	}
	//Create Logger
	logFile := config.LogFile
	if logFile == nil {
		logFile = os.Stdout
	}
	logger := log.New(logFile, "["+name+"]", log.Ldate)
	prefabsFolderName := config.PrefabsFolderName
	if prefabsFolderName == "" {
		prefabsFolderName = DefaultPrefabsFolderName
	}
	resourcesFolderName := config.ResourcesFolderName
	if resourcesFolderName == "" {
		resourcesFolderName = DefaultResourcesFolderName
	}
	scenesFolderName := config.ScenesFolderName
	if scenesFolderName == "" {
		scenesFolderName = DefaultScenesFolderName
	}

	prefabsFolder, err := ioutil.ReadDir(prefabsFolderName)
	if err != nil {
		panic(err)
	}
	resourcesFolder, err := ioutil.ReadDir(resourcesFolderName)
	if err != nil {
		panic(err)
	}
	scenesFolder, err := ioutil.ReadDir(scenesFolderName)
	if err != nil {
		panic(err)
	}

	app := Game{
		name:                name,
		logger:              logger,
		window:              newWindow(wc),
		PrefabsFolderName:   prefabsFolderName,
		ResourcesFolderName: resourcesFolderName,
		ScenesFolderName:    scenesFolderName,

		prefabsFolder:   prefabsFolder,
		resourcesFolder: resourcesFolder,
		scenesFolder:    scenesFolder,

		scenes: make(map[string]*Scene),
	}

	windowSize := app.window.renderWindow.GetSize()
	app.GameWidth = float32(windowSize.X)
	app.GameHeight = float32(windowSize.Y)
	return &app
}

//Init : Loads prefabs resources and scenes
func (g *Game) Init() {
	for _, prefabFile := range g.prefabsFolder {
		path := filepath.Join(g.PrefabsFolderName, prefabFile.Name())
		PrefabRegister.RegisterFromFile(path)
	}
	for _, sceneFile := range g.scenesFolder {
		path := filepath.Join(g.ScenesFolderName, sceneFile.Name())
		scene, err := g.LoadSceneFromFile(path)
		if err != nil {
			panic(err)
		}
		g.logger.Printf("Scene %q loaded", scene.Name)
	}
}

//LoadSceneFromFile : Gets Scene from File
func (g *Game) LoadSceneFromFile(path string) (*Scene, error) {
	dat, err := ioutil.ReadFile(path)
	def, err := ParseSceneDef(g, filepath.Base(path), string(dat))
	if err != nil {
		return nil, err
	}
	//TODO create scene from scenedef
	scene, err := SceneFromSceneDef(def)
	if err != nil {
		return nil, err
	}
	g.scenes[scene.Name] = scene
	return scene, nil

}

//ChangeScene : Loads a new scne
func (g *Game) ChangeScene(name string) {
	scene, ok := g.scenes[name]
	if !ok {
		panic("No Scene with that name")
	}
	g.window.scene = scene
}

//ProcessArguments : Handles CommandLine Arguments
func (g *Game) ProcessArguments() {
	g.logger.Printf("Done")
}

//Run : Runs Game
func (g *Game) Run() {
	window := g.window
	renderWindow := g.window.renderWindow
	go window.Run()
	for renderWindow.IsOpen() {
		select {
		case <-window.Ticker.C:
			//poll events
			for event := renderWindow.PollEvent(); event != nil; event = renderWindow.PollEvent() {
				switch ev := event.(type) {
				case sf.EventKeyReleased:
					switch ev.Code {
					case sf.KeyEscape:
						renderWindow.Close()
					}
				}
			}

		}
	}
}
