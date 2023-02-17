package terrier_translator

import (
	"awesomeProject/internal/repo/terrier"
	"log"
)

type TerrierTranslator struct {
	Terrier *terrier.Terrier
	Repo    terrier.ITerrierRepo
}

func (t *TerrierTranslator) SendSuccessCallback() error {
	return t.Repo.UpdateTerrierPublished(t.Terrier)
}

func (t *TerrierTranslator) Translate() {
	log.Printf("%v is %v", t.Terrier.DogName, t.Terrier.Activity)
}
