package notable

import (
	"bytes"
	"fmt"
	sp "github.com/SparkPost/gosparkpost"
	"log"
	"regexp"
	"text/template"
	"time"
)

type Variables struct {
	Today           string
	NotesByCategory []CategoryNotes
}

type CategoryNotes struct {
	Name  string
	Notes []Note
}

func (categoryNotes *CategoryNotes) Title() string {
	count := len(categoryNotes.Notes)
	announcements := pluralize(count, "Announcement")

	return fmt.Sprintf("#%s &mdash; %s", categoryNotes.Name, announcements)
}

func Email() string {
	var html bytes.Buffer

	notesTemplate, err := template.ParseFiles("template.html")
	check(err)

	today := time.Now().Add(-8 * time.Hour).Format("Monday, January 2, 2006")
	variables := Variables{today, notesByCategory()}
	err = notesTemplate.Execute(&html, variables)
	check(err)

	autolinkRegexp := regexp.MustCompile(`([^"])(\b([\w-]+://?|www[.])[^\s()<>]+(?:\([\w\d]+\)|([^[:punct:]\s]|/)))`)
	return autolinkRegexp.ReplaceAllString(html.String(), "$1<a href=\"$2\">$2</a>")
}

func SendEmail(apiKey string, toEmail string, fromEmail string) {
	var sparky sp.Client
  err := sparky.Init(&sp.Config{ApiKey: apiKey})

  if err != nil {
    log.Fatalf("SparkPost client init failed: %s\n", err)
  }

  tx := &sp.Transmission{
    Recipients: []string{toEmail},
    Content: sp.Content{
      HTML:    Email(),
      From:    fromEmail,
      Subject: pluralize(len(Notes()), "Notable Announcement"),
    },
  }
  id, _, err := sparky.Send(tx)
  if err != nil {
    log.Fatal(err)
  }

  log.Printf("Transmission sent with id [%s]\n", id)
}

func notesByCategory() []CategoryNotes {
	var category string
	grouped := make(map[string]*CategoryNotes, 0)

	for _, note := range Notes() {
		category = note.Category

		if _, found := grouped[category]; !found {
			grouped[category] = &CategoryNotes{Name: category, Notes: make([]Note, 0)}
		}

		grouped[category].Notes = append(grouped[category].Notes, note)
	}

	categoryNotes := make([]CategoryNotes, 0)

	for _, value := range grouped {
		categoryNotes = append(categoryNotes, *value)
	}

	return categoryNotes
}

func pluralize(count int, singularForm string) string {
	pluralForm := fmt.Sprintf("%s%s", singularForm, "s")

	if count == 1 {
		return fmt.Sprintf("1 %s", singularForm)
	} else {
		return fmt.Sprintf("%d %s", count, pluralForm)
	}
}
