package main

import "strings"

func CleanInput(src string) string {
	profanes := [3]string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	srcSlice := strings.Split(src, " ")

	// To-Do use map instead of slice to avoid O(n^2) TC
	for i, word := range srcSlice {
		for _, profane := range profanes {
			if strings.ToLower(word) == profane {
				srcSlice[i] = "****"
			}
		}
	}
	return strings.Join(srcSlice, " ")
}
