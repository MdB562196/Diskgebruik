package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	FileTypes map[string][]string `json:"categories"`
}

func formatSize(size int64) string {
	if size > 1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(1024*1024*1024)) // Grootte in GB
	}
	return fmt.Sprintf("%.2f MB", float64(size)/float64(1024*1024)) //Grootte in MB
}

func getFileType(file string, fileTypes map[string][]string) string {
	ext := strings.ToLower(filepath.Ext(file))
	for category, extensions := range fileTypes {
		for _, e := range extensions {
			if ext == e {
				return category
			}
		}
	}
	return "Overige"
}

func scanDisk(drive string, fileTypes map[string][]string) map[string]int64 {
	result := make(map[string]int64)
	filepath.Walk(drive, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) || err.Error() == "Toegang geweigerd" {
				fmt.Printf("Toegang geweigerd voor dit bestand: %s. Overgeslagen.\n", path)
				return nil
			}
			fmt.Println("Fout bij het openen van dit bestand:", err)
			return err
		}
		if !info.IsDir() {
			category := getFileType(path, fileTypes)
			result[category] += info.Size()
		}
		return nil
	})
	return result
}

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Fout bij het openen van de configuratiebestand", err)
		return
	}
	defer configFile.Close()

	var config Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		fmt.Println("Fout bij decoderen van configuratiebestand", err)
		return
	}

	var drive string
	for {
		fmt.Print("Voer de schijf/directory in die je wilt scannen (of typ 'verlaat' om te stoppen: ")
		fmt.Scanln(&drive)
		if strings.ToLower(drive) == "verlaat" {
			fmt.Println("Programma verlaten.")
			break
		}

		result := scanDisk(drive, config.FileTypes)
		fmt.Println("\nDiskgebruik overzicht:")
		for category, size := range result {
			fmt.Printf("%s: %s\n", category, formatSize(size))
		}
	}
}
