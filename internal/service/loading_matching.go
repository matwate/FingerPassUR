package service

import (
	"context"
	"image/png"
	"log"
	"os"
	"sync"

  "github.com/matwate/corner/internal/model"
  "github.com/matwate/corner/internal/repository"

	"github.com/miqdadyyy/go-sourceafis"
	"github.com/miqdadyyy/go-sourceafis/templates"
)

var transparencyLogger = sourceafis.NewTransparencyLogger(UrTransparencyProvider{})

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
		}(img, i)

	}
	wg.Wait()
	// Now set the templates to the pointer
	*templs = ImagesTemplates
}

func MatchTemplates(
	probe *templates.SearchTemplate,
	candidates []*templates.SearchTemplate,
) []float64 {
	matcher, err := sourceafis.NewMatcher(nil, probe)
	if err != nil {
		log.Fatal(err)
	}
	// Get all similarity scores concurrently
	var matches []float64 = make([]float64, len(candidates))
	var wg sync.WaitGroup
	for i, candidate := range candidates {
		wg.Add(1)
		go func(i int, candidate *templates.SearchTemplate) {
			defer wg.Done()
			similarity := matcher.Match(context.Background(), candidate)
			matches[i] = similarity
		}(i, candidate)
	}

	return matches
}

