package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"
	"os"
	"sync"

	"github.com/miqdadyyy/go-sourceafis"
	"github.com/miqdadyyy/go-sourceafis/templates"

	"github.com/matwate/corner/internal/model"
	"github.com/matwate/corner/internal/repository"
)

var (
	transparencyLogger = sourceafis.NewTransparencyLogger(UrTransparencyProvider{})
	Templates          []*templates.SearchTemplate
	UserMap            map[int]int = make(map[int]int) // From index to the array to the user id
)

func LoadTemplates(templs *[]*templates.SearchTemplate) {
	// Get all png images from the database.
	var Images []model.Image = repository.ListALLImages()
	// Open all .png files and load them in memory
	var ImagesTemplates []*templates.SearchTemplate = make([]*templates.SearchTemplate, len(Images))

	c := sourceafis.NewTemplateCreator(transparencyLogger)
	var wg sync.WaitGroup
	for i, img := range Images {
		wg.Add(1)
		go func(img model.Image, i int) {
			defer wg.Done()
			// Open the image file
			file, err := os.Open(img.Template)
			if err != nil {
				log.Fatal(err)
				panic("Something went wrong")
			}
			defer file.Close()
			// Decode the image
			image, err := png.Decode(file)
			if err != nil {
				log.Fatal(err)
				panic("Something went wrong")
			}

			// Create a template from the image
			sImage, err := sourceafis.NewFromImage(image)
			if err != nil {
				log.Fatal(err)
				panic("Something went wrong!")
			}
			// Add the template to the list
			template, err := c.Template(sImage)
			if err != nil {
				log.Fatal(err)
				panic("Something went wrong")
			}
			ImagesTemplates[i] = template
			UserMap[i] = img.User_id
		}(img, i)

	}
	wg.Wait()
	// Now set the templates to the pointer
	*templs = ImagesTemplates
}

func Base64ToTemplate(b64 string) *templates.SearchTemplate {
	// Decode the base64 string
	unbased, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewReader(unbased)
	// Decode the image
	image, err := png.Decode(r)
	if err != nil {
		log.Fatal(err)
	}

	sImage, err := sourceafis.NewFromImage(image)
	if err != nil {
		log.Fatal(err)
	}

	c := sourceafis.NewTemplateCreator(transparencyLogger)
	template, err := c.Template(sImage)
	if err != nil {
		log.Fatal(err)
	}

	// Create a template from the image
	return template
}

func MatchTemplates(
	probe *templates.SearchTemplate,
) int {
	candidates := Templates
	fmt.Println(len(candidates))
	matcher, err := sourceafis.NewMatcher(transparencyLogger, probe)
	if err != nil {
		log.Fatal(err)
	}
	// Get all similarity scores concurrently
	matches := make(map[int]float64, len(candidates))
	fmt.Println(matches)
	var wg sync.WaitGroup
	for i, candidate := range candidates {
		wg.Add(1)
		go func(i int, candidate *templates.SearchTemplate) {
			defer wg.Done()
			similarity := matcher.Match(context.Background(), candidate)
			matches[i] = similarity
		}(i, candidate)
	}
	wg.Wait()
	fmt.Println(matches)

	// Now find the highest similarity score.
	var max float64 = 0
	var maxIndex int = -1
	for i, similarity := range matches {
		if similarity > max {
			max = similarity
			maxIndex = i
		}
	}
	// Now maxIndex is the index of the user in the UserMap which has the highest similarity score
	if matches[maxIndex] < 40 {
		return -1
	}
	return UserMap[maxIndex]
}
