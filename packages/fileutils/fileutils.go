package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	DefaultFileWritePermissions = 0o600
	DefaultDirWritePermissions  = 0o755
)

// Comments out lines containing the "$mcp" keyword
func CommentOutMCPLines(inputFile string) (outputFile string, patchList []map[string]interface{}, err error) {
	const (
		mcpPattern = "$mcp"
	)

	contents, err := os.ReadFile(inputFile)
	if err != nil {
		return outputFile, patchList, fmt.Errorf("unable to open file: %s, err: %s ", inputFile, err)
	}

	pattern := regexp.MustCompile(fmt.Sprintf(`.*%s.*`, regexp.QuoteMeta(mcpPattern)))

	// Split the file contents into lines
	lines := strings.Split(string(contents), "\n")

	// comment out every lines matching pattern
	var modifiedLines []string
	for _, line := range lines {
		if pattern.MatchString(line) {
			// Comment out the line by adding "//" at the beginning
			line = "# " + line
		}

		// Add the processed line to the result
		modifiedLines = append(modifiedLines, line)
	}

	// Join the modified lines
	modifiedString := strings.Join(modifiedLines, "\n")
	outputFile = strings.TrimSuffix(inputFile, ".yaml") + "-SetSelector.yaml"
	err = os.WriteFile(outputFile, []byte(modifiedString), DefaultFileWritePermissions)
	if err != nil {
		fmt.Printf("Error writing to file: %s, err: %s", inputFile, err)
		return "", patchList, err
	}
	fmt.Printf("Wrote converted ACM template: %s\n", outputFile)
	return outputFile, patchList, nil
}

// Replaces the "$mcp" keyword with the mcp string (worker or master)
func RenderMCPLines(inputFile, mcp string) (outputFile string, err error) {
	const (
		mcpPattern = "$mcp"
	)

	contents, err := os.ReadFile(inputFile)
	if err != nil {
		return outputFile, fmt.Errorf("unable to open file: %s, err: %s ", inputFile, err)
	}
	outputFile = inputFile
	if strings.Contains(string(contents), mcpPattern) {
		contents = []byte(strings.ReplaceAll(string(contents), mcpPattern, mcp))
		outputFile = strings.TrimSuffix(inputFile, ".yaml") + "-MCP-" + mcp + ".yaml"
	}

	err = os.WriteFile(outputFile, contents, DefaultFileWritePermissions)
	if err != nil {
		fmt.Printf("Error writing to file: %s, err: %s", inputFile, err)
		return "", err
	}
	fmt.Printf("Wrote converted ACM template: %s\n", outputFile)
	return outputFile, nil
}

// Gets All Yaml files in a path
func GetAllYAMLFilesInPath(path string) (files []string, err error) {
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//		if !info.IsDir() {
		if !info.IsDir() && (strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml")) {
			files = append(files, filePath)
		}
		return nil
	})

	return files, err
}

// Prefixes the file referred to by a path with a given prefix
func PrefixLastPathComponent(originalPath, prefix string) string {
	dir, file := filepath.Split(originalPath)
	if dir == "" {
		// If the originalPath has no directory component, simply prefix the file name
		return prefix + file
	}
	return filepath.Join(dir, prefix+file)
}

// Type used to parse the Kind of a K8s Manifest
type KindType struct {
	Kind string `yaml:"kind"`
}

// Gets the manifest kind from the file
func GetManifestKind(filePath string) (kindType KindType, err error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return kindType, fmt.Errorf("could not read %s: %s", filePath, err)
	}
	err = yaml.Unmarshal(yamlFile, &kindType)
	if err != nil {
		return kindType, fmt.Errorf("could not parse %s as yaml: %s", filePath, err)
	}
	return kindType, nil
}

type AnnotationsOnly struct {
	Metadata MetaData `yaml:"metadata"`
}

type MetaData struct {
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// Gets the manifest kind from the file
func GetAnnotationsOnly(filePath string) (annotations AnnotationsOnly, err error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return annotations, fmt.Errorf("could not read %s: %s", filePath, err)
	}
	err = yaml.Unmarshal(yamlFile, &annotations)
	if err != nil {
		return annotations, fmt.Errorf("could not parse %s as yaml: %s", filePath, err)
	}
	return annotations, nil
}
