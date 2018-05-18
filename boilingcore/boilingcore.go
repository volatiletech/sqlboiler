// Package boilingcore has types and methods useful for generating code that
// acts as a fully dynamic ORM might.
package boilingcore

import (
	"encoding/json"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/drivers"
	"github.com/volatiletech/sqlboiler/strmangle"
)

const (
	templatesDirectory          = "templates"
	templatesSingletonDirectory = "templates/singleton"

	templatesTestDirectory          = "templates_test"
	templatesSingletonTestDirectory = "templates_test/singleton"

	templatesTestMainDirectory = "templates_test/main_test"
)

// State holds the global data needed by most pieces to run
type State struct {
	Config *Config

	Driver  drivers.Interface
	Schema  string
	Tables  []drivers.Table
	Dialect drivers.Dialect

	Templates              *templateList
	TestTemplates          *templateList
	SingletonTemplates     *templateList
	SingletonTestTemplates *templateList
}

// New creates a new state based off of the config
func New(config *Config) (*State, error) {
	s := &State{
		Config: config,
	}

	s.Driver = drivers.GetDriver(config.DriverName)

	err := s.initDBInfo(config.DriverConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize tables")
	}

	if s.Config.Debug {
		b, err := json.Marshal(s.Tables)
		if err != nil {
			return nil, errors.Wrap(err, "unable to json marshal tables")
		}
		fmt.Printf("%s\n", b)
	}

	err = s.initOutFolder()
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize the output folder")
	}

	err = s.initTemplates()
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize templates")
	}

	err = s.initTags(config.Tags)
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize struct tags")
	}

	return s, nil
}

// Run executes the sqlboiler templates and outputs them to files based on the
// state given.
func (s *State) Run() error {
	singletonData := &templateData{
		Tables:           s.Tables,
		DriverName:       s.Config.DriverName,
		PkgName:          s.Config.PkgName,
		NoHooks:          s.Config.NoHooks,
		NoAutoTimestamps: s.Config.NoAutoTimestamps,
		StructTagCasing:  s.Config.StructTagCasing,
		Dialect:          s.Dialect,
		LQ:               strmangle.QuoteCharacter(s.Dialect.LQ),
		RQ:               strmangle.QuoteCharacter(s.Dialect.RQ),

		StringFuncs: templateStringMappers,
	}

	if err := generateSingletonOutput(s, singletonData); err != nil {
		return errors.Wrap(err, "singleton template output")
	}

	if !s.Config.NoTests {
		if err := generateSingletonTestOutput(s, singletonData); err != nil {
			return errors.Wrap(err, "unable to generate singleton test template output")
		}
	}

	for _, table := range s.Tables {
		if table.IsJoinTable {
			continue
		}

		data := &templateData{
			Tables:           s.Tables,
			Table:            table,
			DriverName:       s.Config.DriverName,
			PkgName:          s.Config.PkgName,
			NoHooks:          s.Config.NoHooks,
			NoAutoTimestamps: s.Config.NoAutoTimestamps,
			StructTagCasing:  s.Config.StructTagCasing,
			Tags:             s.Config.Tags,
			Dialect:          s.Dialect,
			LQ:               strmangle.QuoteCharacter(s.Dialect.LQ),
			RQ:               strmangle.QuoteCharacter(s.Dialect.RQ),

			StringFuncs: templateStringMappers,
		}

		// Generate the regular templates
		if err := generateOutput(s, data); err != nil {
			return errors.Wrap(err, "unable to generate output")
		}

		// Generate the test templates
		if !s.Config.NoTests {
			if err := generateTestOutput(s, data); err != nil {
				return errors.Wrap(err, "unable to generate test output")
			}
		}
	}

	return nil
}

// Cleanup closes any resources that must be closed
func (s *State) Cleanup() error {
	// Nothing here atm, used to close the driver
	return nil
}

// initTemplates loads all template folders into the state object.
func (s *State) initTemplates() error {
	var err error

	basePath, err := getBasePath(s.Config.BaseDir)
	if err != nil {
		return err
	}

	templates, err := findTemplates(basePath, templatesDirectory)
	if err != nil {
		return err
	}
	testTemplates, err := findTemplates(basePath, templatesTestDirectory)
	if err != nil {
		return err
	}

	for k, v := range testTemplates {
		templates[k] = v
	}

	driverTemplates, err := s.Driver.Templates()
	if err != nil {
		return err
	}

	for template, contents := range driverTemplates {
		templates[template] = base64Loader(contents)
	}

	for _, replace := range s.Config.Replacements {
		splits := strings.Split(replace, ":")
		if len(splits) != 2 {
			return errors.Errorf("replace parameters must have 2 arguments, given: %s", replace)
		}

		original, replacement := splits[0], splits[1]

		_, ok := templates[original]
		if !ok {
			return errors.Errorf("replace can only replace existing templates, %s does not exist", original)
		}

		templates[original] = fileLoader(replacement)
	}

	// For stability, sort keys to traverse the map and turn it into a slice
	keys := make([]string, 0, len(templates))
	for k := range templates {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lazyTemplates := make([]lazyTemplate, 0, len(templates))
	for _, k := range keys {
		lazyTemplates = append(lazyTemplates, lazyTemplate{
			Name:   k,
			Loader: templates[k],
		})
	}

	if s.Config.Debug {
		b, err := json.Marshal(lazyTemplates)
		if err != nil {
			return errors.Wrap(err, "unable to json marshal template list")
		}

		fmt.Printf("%s\n", b)
	}

	s.Templates, err = loadTemplates(lazyTemplates, templatesDirectory)
	if err != nil {
		return err
	}

	s.SingletonTemplates, err = loadTemplates(lazyTemplates, templatesSingletonDirectory)
	if err != nil {
		return err
	}

	if !s.Config.NoTests {
		s.TestTemplates, err = loadTemplates(lazyTemplates, templatesTestDirectory)
		if err != nil {
			return err
		}

		s.SingletonTestTemplates, err = loadTemplates(lazyTemplates, templatesSingletonTestDirectory)
		if err != nil {
			return err
		}
	}

	return nil
}

// findTemplates uses a base path: (/home/user/gopath/src/../sqlboiler/)
// and a root path: /templates
// to create a bunch of file loaders of the form:
// templates/00_struct.tpl -> /absolute/path/to/that/file
// Note the missing leading slash, this is important for the --replace argument
func findTemplates(base, root string) (map[string]templateLoader, error) {
	templates := make(map[string]templateLoader)
	baseRoot := filepath.Join(base, root)
	err := filepath.Walk(baseRoot, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		if ext := filepath.Ext(path); ext == ".tpl" {
			relative, err := filepath.Rel(base, path)
			if err != nil {
				return errors.Wrapf(err, "could not find relative path to base root: %s", baseRoot)
			}
			relative = strings.TrimLeft(relative, string(os.PathSeparator))
			templates[relative] = fileLoader(path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return templates, nil
}

var basePackage = "github.com/volatiletech/sqlboiler"

func getBasePath(baseDirConfig string) (string, error) {
	if len(baseDirConfig) > 0 {
		return baseDirConfig, nil
	}

	p, _ := build.Default.Import(basePackage, "", build.FindOnly)
	if p != nil && len(p.Dir) > 0 {
		return p.Dir, nil
	}

	return os.Getwd()
}

// initDBInfo retrieves information about the database
func (s *State) initDBInfo(config map[string]interface{}) error {
	dbInfo, err := s.Driver.Assemble(config)
	if err != nil {
		return errors.Wrap(err, "unable to fetch table data")
	}

	if len(dbInfo.Tables) == 0 {
		return errors.New("no tables found in database")
	}

	if err := checkPKeys(dbInfo.Tables); err != nil {
		return err
	}

	s.Schema = dbInfo.Schema
	s.Tables = dbInfo.Tables
	s.Dialect = dbInfo.Dialect

	return nil
}

// Tags must be in a format like: json, xml, etc.
var rgxValidTag = regexp.MustCompile(`[a-zA-Z_\.]+`)

// initTags removes duplicate tags and validates the format
// of all user tags are simple strings without quotes: [a-zA-Z_\.]+
func (s *State) initTags(tags []string) error {
	s.Config.Tags = strmangle.RemoveDuplicates(s.Config.Tags)
	for _, v := range s.Config.Tags {
		if !rgxValidTag.MatchString(v) {
			return errors.New("Invalid tag format %q supplied, only specify name, eg: xml")
		}
	}

	return nil
}

// initOutFolder creates the folder that will hold the generated output.
func (s *State) initOutFolder() error {
	if s.Config.Wipe {
		if err := os.RemoveAll(s.Config.OutFolder); err != nil {
			return err
		}
	}

	return os.MkdirAll(s.Config.OutFolder, os.ModePerm)
}

// checkPKeys ensures every table has a primary key column
func checkPKeys(tables []drivers.Table) error {
	var missingPkey []string
	for _, t := range tables {
		if t.PKey == nil {
			missingPkey = append(missingPkey, t.Name)
		}
	}

	if len(missingPkey) != 0 {
		return errors.Errorf("primary key missing in tables (%s)", strings.Join(missingPkey, ", "))
	}

	return nil
}
