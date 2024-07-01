package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	outputFileBaseName := filepath.Join("knowledge")
	i := 1
	outputFileName := fmt.Sprintf("%s%d.md", outputFileBaseName, i)

	fmt.Println("Creating", outputFileName)
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile.Close()

	// Enumerate the cheatsheets
	directoryPath := "./OWASP-CheatSheetSeries/cheatsheets/"
	cheatSheets, err := filepath.Glob(filepath.Join(directoryPath, "*.md"))
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	// Manually populate a list of the whitepapers
	whitePapers := []string{
		"./tag-security/supply-chain-security/supply-chain-security-paper/sscsp.md",
		"./tag-security/supply-chain-security/secure-software-factory/secure-software-factory.md",
		"./tag-security/security-fuzzing-handbook/fuzzing-handbook.md",
		"./tag-security/security-whitepaper/v2/cloud-native-security-whitepaper.md",
		"./sig-security/sig-security-docs/papers/policy_grc/kubernetes-grc.md",
		"./sig-security/sig-security-docs/papers/policy/kubernetes-policy-management.md",
		"./tag-runtime/wgs/cnai/whitepapers/cloudnativeai.md",
	}

	// Create a single list of all files
	files := append(cheatSheets, whitePapers...)

	// Concatenate files into as few files as possible, creating a new file only once it exceeds 1.5MB
	for _, file := range files {
		// Add a per-file header for sanity
		_, err := outputFile.WriteString(fmt.Sprintf("START OF %s\n\n", filepath.Base(file)))
		if err != nil {
			fmt.Printf("Error writing header: %v\n", err)
			return
		}

		// I don't know the exact limitations so this is guesswork; > 1.5MB seems to work
		fileInfo, err := os.Stat(outputFileName)
		if err != nil {
			fmt.Printf("Error getting file information: %v\n", err)
			return
		}
		fileSize := fileInfo.Size()
		fileSizeInMB := float64(fileSize) / (1024 * 1024)

		if fileSizeInMB > 1.5 {
			fmt.Printf("%s is > 1.5 MB, creating a new file...\n", outputFileName)
			i++
			outputFileName = fmt.Sprintf("%s%d.md", outputFileBaseName, i)

			outputFile, err = os.Create(outputFileName)
			if err != nil {
				fmt.Printf("Error creating output file: %v\n", err)
				return
			}
			fmt.Printf("Created %s\n", outputFileName)
			defer outputFile.Close()
		}

		fileContents, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			return
		}

		_, err = outputFile.Write(fileContents)
		if err != nil {
			fmt.Printf("Error writing file contents: %v\n", err)
			return
		}

		// Add a per-file footer for sanity
		_, err = outputFile.WriteString(fmt.Sprintf("\n\nEND OF %s\n\n", filepath.Base(file)))
		if err != nil {
			fmt.Printf("Error writing separator: %v\n", err)
			return
		}
	}
	fmt.Printf("Done!\n")
}
