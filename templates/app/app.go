package app

import (
	"context"
	"embed"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"{{CHIFRA}}/pkg/base"
	"{{CHIFRA}}/pkg/config"
	"{{CHIFRA}}/pkg/file"
	coreTypes "{{CHIFRA}}/pkg/types"
	"{{CHIFRA}}/pkg/utils"
	dalle "github.com/{{ORG_NAME}}/{{ORG_LOWER}}-dalle/v2"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/fileserver"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/logging"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/msgs"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/preferences"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/project"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/store"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types/exports"
	sdk "github.com/{{ORG_NAME}}/{{ORG_LOWER}}-sdk/v5"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	Assets      embed.FS
	Preferences *preferences.Preferences
	Projects    *project.Manager
	chainList   *utils.ChainList
	// ADD_ROUTE
	// QUESTION: DO WE ACTUALLY NEED THIS?
	exports *exports.ExportsCollection
	// ADD_ROUTE
	collections []types.Collection
	meta        *coreTypes.MetaData
	fileServer  *fileserver.FileServer
	prefsMu     sync.RWMutex
	ctx         context.Context
	apiKeys     map[string]string
	ensMap      map[string]base.Address
	Dalle       *dalle.Context
}

func NewApp(assets embed.FS) (*App, *menu.Menu) {
	app := &App{
		Projects: project.NewManager(),
		Preferences: &preferences.Preferences{
			Org:  preferences.OrgPreferences{},
			User: preferences.UserPreferences{},
			App:  preferences.AppPreferences{},
		},
		Assets:  assets,
		apiKeys: make(map[string]string),
		ensMap:  make(map[string]base.Address),
	}
	// ADD_ROUTE
	// QUESTION: DO WE ACTUALLY NEED THIS?
	// Note: exports created on-demand per chain/address when needed
	app.exports = nil
	// ADD_ROUTE

	app.collections = make([]types.Collection, 0, 4)

	app.chainList, _ = utils.UpdateChainList(config.PathToRootConfig())

	if file.FileExists(".env") {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		} else if app.apiKeys["openAi"] = os.Getenv("OPENAI_API_KEY"); app.apiKeys["openAi"] == "" {
			log.Fatal("No OPENAI_API_KEY key found")
		}
	}

	app.Dalle = dalle.NewContext("./output")

	appMenu := app.buildAppMenu()
	return app, appMenu
}

func (a *App) GetContext() context.Context {
	return a.ctx
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	msgs.InitializeContext(ctx)

	org, err := preferences.GetOrgPreferences()
	if err != nil {
		msgs.EmitError("Loading org preferences failed", err)
		return
	}

	user, err := preferences.GetUserPreferences()
	if err != nil {
		msgs.EmitError("Loading user preferences failed", err)
		return
	}

	appPrefs, err := preferences.GetAppPreferences()
	if err != nil {
		msgs.EmitError("Loading app preferences failed", err)
		return
	}

	a.Preferences.Org = org
	a.Preferences.User = user
	a.Preferences.App = appPrefs

	a.fileServer = fileserver.NewFileServer()
	if err := a.fileServer.Start(); err != nil {
		msgs.EmitError("Failed to start file server", err)
	}
	go a.watchImagesDir()

	if len(a.Preferences.App.RecentProjects) > 0 {
		mostRecentPath := a.Preferences.App.RecentProjects[0]
		if file.FileExists(mostRecentPath) {
			_, err := a.Projects.Open(mostRecentPath)
			if err != nil {
				msgs.EmitError("Failed to open recent project", err)
			}
		}
	}
}

func (a *App) DomReady(ctx context.Context) {
	a.ctx = ctx
	if a.IsReady() {
		if !a.Preferences.App.Bounds.IsValid() {
			// Sometimes, during development, the window size is corrupted
			// and we need to reset it to a default value. Should really
			// happen in production.
			a.Preferences.App.Bounds = preferences.NewBounds()
		}
		runtime.WindowSetSize(ctx, a.Preferences.App.Bounds.Width, a.Preferences.App.Bounds.Height)
		runtime.WindowSetPosition(ctx, a.Preferences.App.Bounds.X, a.Preferences.App.Bounds.Y)
		runtime.WindowShow(ctx)
	}
	go a.watchWindowBounds() // if the window moves or resizes, we want to know
}

func (a *App) BeforeClose(ctx context.Context) bool {
	x, y := runtime.WindowGetPosition(ctx)
	w, h := runtime.WindowGetSize(ctx)
	a.SaveBounds(x, y, w, h)

	if a.fileServer != nil {
		if err := a.fileServer.Stop(); err != nil {
			log.Printf("Error shutting down file server: %v", err)
		}
	}

	return false // allow window to close
}

func (a *App) watchWindowBounds() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var lastX, lastY, lastW, lastH int
	for range ticker.C {
		if !a.IsReady() {
			continue
		}
		x, y := runtime.WindowGetPosition(a.ctx)
		w, h := runtime.WindowGetSize(a.ctx)
		if x != lastX || y != lastY || w != lastW || h != lastH {
			a.SaveBounds(x, y, w, h)
			lastX, lastY, lastW, lastH = x, y, w, h
		}
	}
}

func (a *App) SaveBounds(x, y, w, h int) {
	if !a.IsReady() {
		return
	}

	a.Preferences.App.Bounds = preferences.Bounds{
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
	}

	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

func (a *App) IsReady() bool {
	return a.ctx != nil
}

func (a *App) IsInitialized() bool {
	_, appFolder := preferences.GetConfigFolders()
	fn := filepath.Join(appFolder, ".initialized")
	return file.FileExists(fn)
}

func (a *App) SetInitialized(isInit bool) error {
	_, appFolder := preferences.GetConfigFolders()
	fn := filepath.Join(appFolder, ".initialized")
	if isInit {
		if !file.Touch(fn) {
			return fmt.Errorf("failed to create " + fn + " file")
		} else {
			return nil
		}
	} else {
		_ = os.Remove(fn)
		return nil // do not fail even if not found
	}
}

func (a *App) GetWizardReturn() string {
	if a.Preferences.App.LastViewNoWizard == "" {
		return "/"
	}
	return a.Preferences.App.LastViewNoWizard
}

func (a *App) GetAppId() preferences.Id {
	return preferences.GetAppId()
}

func (a *App) GetOpenProjects() []map[string]interface{} {
	projectIDs := a.Projects.GetOpenProjectIDs()
	result := make([]map[string]interface{}, 0, len(projectIDs))

	for _, id := range projectIDs {
		project := a.Projects.GetProjectByID(id)
		if project == nil {
			continue
		}

		projectInfo := map[string]interface{}{
			"id":         id,
			"name":       project.GetName(),
			"path":       project.GetPath(),
			"isActive":   id == a.Projects.ActiveID,
			"isDirty":    project.IsDirty(),
			"lastOpened": project.LastOpened,
			// "createdAt":  project.CreatedAt,
		}

		result = append(result, projectInfo)
	}

	return result
}

func (a *App) GetUserPreferences() *preferences.UserPreferences {
	return &a.Preferences.User
}

func (a *App) SetUserPreferences(userPrefs *preferences.UserPreferences) error {
	a.Preferences.User = *userPrefs
	return preferences.SetUserPreferences(userPrefs)
}

func (a *App) GetOrgPreferences() *preferences.OrgPreferences {
	return &a.Preferences.Org
}

func (a *App) SetOrgPreferences(orgPrefs *preferences.OrgPreferences) error {
	a.Preferences.Org = *orgPrefs
	return preferences.SetOrgPreferences(orgPrefs)
}

func (app *App) GetChainList() *utils.ChainList {
	return app.chainList
}

// GetProjectAddress returns the address of the active project
func (a *App) GetProjectAddress() base.Address {
	active := a.Projects.Active()
	if active == nil {
		return base.ZeroAddr
	}
	return active.GetAddress()
}

// SetProjectAddress sets the address of the active project
func (a *App) SetProjectAddress(addr base.Address) {
	active := a.Projects.Active()
	if active != nil {
		active.SetAddress(addr)
	}
}

func (a *App) BuildDalleDressForProject() (map[string]interface{}, error) {
	active := a.Projects.Active()
	if active == nil {
		return nil, fmt.Errorf("no active project")
	}
	addr := active.GetAddress()
	if addr == base.ZeroAddr {
		return nil, fmt.Errorf("project address is not set")
	}

	// Always resolve ENS/address using ConvertToAddress
	resolved, ok := a.ConvertToAddress(addr.Hex())
	if !ok || resolved == base.ZeroAddr {
		return nil, fmt.Errorf("invalid address or ENS name")
	}

	dress, err := a.Dalle.MakeDalleDress(resolved.Hex())
	if err != nil {
		return nil, err
	}

	imagePath := filepath.Join("generated", dress.Filename+".png")
	imageURL := ""
	if a.fileServer != nil {
		imageURL = a.fileServer.GetURL(imagePath)
	}

	return map[string]interface{}{
		"imageUrl": imageURL,
		"parts":    dress,
	}, nil
}

func (a *App) CancelFetch(dataFacet types.DataFacet) {
	storeName := ""

	for _, collection := range a.collections {
		for _, facet := range collection.GetSupportedFacets() {
			if facet == dataFacet {
				storeName = collection.GetStoreName(facet)
				break
			}
		}
		if storeName != "" {
			break
		}
	}

	if storeName != "" {
		store.CancelFetch(storeName)
	}
}

func (a *App) CancelAllFetches() int {
	return store.CancelAllFetches()
}

func (a *App) GetNodeStatus() *coreTypes.MetaData {
	defer logging.Silence()()
	a.meta, _ = sdk.GetMetaData(preferences.GetLastChain())
	return a.meta
}

func (a *App) RegisterCollection(collection types.Collection) {
	a.collections = append(a.collections, collection)
}

func getCollectionPage[T any](
	collection interface {
		GetPage(*types.Payload, int, int, sdk.SortSpec, string) (types.Page, error)
	},
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (T, error) {
	var zero T

	dataFacet := payload.DataFacet
	page, err := collection.GetPage(payload, first, pageSize, sort, filter)
	if err != nil {
		return zero, err
	}

	typedPage, ok := page.(T)
	if !ok {
		return zero, types.NewValidationError("app", dataFacet, "getCollectionPage",
			fmt.Errorf("GetPage returned unexpected type %T, expected %T", page, zero))
	}

	return typedPage, nil
}

func (a *App) ResetStore(storeName string) {
	for _, collection := range a.collections {
		for _, facet := range collection.GetSupportedFacets() {
			if collection.GetStoreName(facet) == storeName {
				collection.Reset(facet)
			}
		}
	}
}

func (a *App) Encode(fn sdk.Function, params []interface{}) (string, error) {
	packed, err := fn.Pack(params)
	if err != nil {
		return "", fmt.Errorf("failed to pack function call: %w", err)
	}
	return "0x" + hex.EncodeToString(packed), nil
}
