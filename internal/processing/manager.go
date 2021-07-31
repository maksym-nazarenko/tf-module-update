package processing

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/maxim-nazarenko/tf-module-update/internal/module"
	"github.com/maxim-nazarenko/tf-module-update/internal/strategies"
)

// Config holds configuration for revision manager
type Config struct {
	// Write indicates if files should be updated
	Write bool
}

// RevisionManager is responsible for managing module source updates
type RevisionManager struct {
	config        Config
	strategy      strategies.Strategy
	resultFactory *ResultFactory
}

// ProcessPaths processes Terraform in the given paths
//
// The process is recursive and checks only files that are not ignored (see ignoredFile())
func (m *RevisionManager) ProcessPaths(paths []string, results *Results) {
	var absPath string
	var err error
	for _, p := range paths {
		absPath, err = filepath.Abs(p)
		if err != nil {
			results.Append(err)
			continue
		}
		info, err := os.Stat(absPath)
		if err != nil {
			results.Append(err)
			continue
		}
		if info.IsDir() {
			m.processDir(absPath, results)
			continue
		}
		if m.ignoredFile(absPath) {
			continue
		}
		results.Append(m.processFile(absPath))
	}
}

func (m *RevisionManager) ignoredFile(absPath string) bool {
	return filepath.Ext(absPath) != ".tf"
}

func (m *RevisionManager) processDir(path string, results *Results) {
	items, err := ioutil.ReadDir(path)
	if err != nil {
		results.Append(err)
		return
	}

	var itemPath string
	for _, item := range items {
		itemPath = filepath.Join(path, item.Name())
		if item.IsDir() {
			m.processDir(itemPath, results)
			continue
		}
		if m.ignoredFile(itemPath) {
			continue
		}
		results.Append(m.processFile(itemPath))
	}
}

func (m *RevisionManager) processFile(fileName string) *Results {

	results := &Results{}

	f, err := os.Open(fileName)
	if err != nil {
		results.Append(err)
		return results
	}

	defer f.Close()
	src, err := io.ReadAll(f)
	if err != nil {
		results.Append(err)
		return results
	}
	normalizedPath, err := filepath.Abs(fileName)
	if err != nil {
		results.Append(err)
		return results
	}

	infileHeader := m.resultFactory.Info("In file " + fileName + ":")
	bodyResults := &Results{}
	updatedFileBody, err := m.updateFileBody(src, normalizedPath, bodyResults)
	if string(updatedFileBody) != string(src) {
		results.Append(infileHeader)
	}
	results.Append(bodyResults)
	if err != nil {
		results.Append(err)
		return results
	}

	if m.config.Write {
		if err = ioutil.WriteFile(normalizedPath, updatedFileBody, 0644); err != nil {
			results.Append(err)
			return results
		}
	}

	return results
}

func (m *RevisionManager) updateFileBody(src []byte, normalizedPath string, results *Results) ([]byte, error) {
	parsed, diags := hclwrite.ParseConfig(src, normalizedPath, hcl.InitialPos)
	if diags.HasErrors() {
		results.Append(diags.Errs())
		return src, errors.New("parsing HCL syntax failed")
	}
	parsedBody := parsed.Body()

	blocks := parsedBody.Blocks()

	for _, b := range blocks {
		// we can process only modules
		if b.Type() != "module" {
			continue
		}
		// there must be a "source" attribute to update
		if _, ok := b.Body().Attributes()["source"]; !ok {
			continue
		}

		results.Append(m.processBlock(b))
	}

	return parsed.Bytes(), nil
}

func (m *RevisionManager) processBlock(block *hclwrite.Block) Results {
	results := Results{}
	// just a sanity check
	sourceAttr := block.Body().GetAttribute("source")
	if block.Type() != "module" || sourceAttr == nil {
		results.Append(errors.New("current block is not a module or does not have source attribute"))
		return results
	}

	exprTokens := sourceAttr.Expr().BuildTokens(nil)
	if exprTokens[1].Type != hclsyntax.TokenQuotedLit {
		return results
	}

	source, err := module.ParseSource(string(exprTokens[1].Bytes))
	if err != nil {
		results.Append(err)
		return results
	}

	if !m.strategy.Decide(source) {
		results.Append(m.resultFactory.Debug("skipping source due to updater decision: " + source.String()))
		return results
	}

	newSource := m.strategy.Apply(source)

	if source.String() == newSource.String() {
		return results
	}

	exprTokens[1] = &hclwrite.Token{
		Type:         exprTokens[1].Type,
		Bytes:        []byte(newSource.String()),
		SpacesBefore: exprTokens[1].SpacesBefore,
	}

	results.Append(
		m.resultFactory.Info("  - "+source.String()),
		m.resultFactory.Info("  + "+newSource.String()),
	)

	block.Body().SetAttributeRaw("source",
		exprTokens,
	)

	return results
}

func NewManager(config Config, strategy strategies.Strategy) *RevisionManager {
	return &RevisionManager{config: config, strategy: strategy, resultFactory: NewResultFactory()}
}
