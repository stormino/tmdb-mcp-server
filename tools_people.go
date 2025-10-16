package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *TMDBServer) SearchPerson(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SearchPersonInput,
) (*mcp.CallToolResult, SearchPersonOutput, error) {
	results, err := s.client.GetSearchPeople(input.Query, nil)
	if err != nil {
		return nil, SearchPersonOutput{}, fmt.Errorf("failed to search person: %w", err)
	}

	people := make([]PersonResult, 0, len(results.Results))
	for _, person := range results.Results {
		knownForTitles := make([]string, 0, len(person.KnownFor))
		for i, item := range person.KnownFor {
			if i >= 3 {
				break
			}
			if item.Title != "" {
				knownForTitles = append(knownForTitles, item.Title)
			} else if item.Name != "" {
				knownForTitles = append(knownForTitles, item.Name)
			}
		}
		knownForStr := person.KnownForDepartment
		if len(knownForTitles) > 0 {
			knownForStr = fmt.Sprintf("%s (%s)", person.KnownForDepartment, strings.Join(knownForTitles, ", "))
		}

		people = append(people, PersonResult{
			ID:          int(person.ID),
			Name:        person.Name,
			KnownFor:    knownForStr,
			ProfilePath: person.ProfilePath,
		})
	}

	return nil, SearchPersonOutput{
		Results: people,
		Total:   len(people),
	}, nil
}

func (s *TMDBServer) GetPersonDetails(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetPersonDetailsInput,
) (*mcp.CallToolResult, GetPersonDetailsOutput, error) {
	person, err := s.client.GetPersonDetails(int(input.PersonID), nil)
	if err != nil {
		return nil, GetPersonDetailsOutput{}, fmt.Errorf("failed to get person details: %w", err)
	}

	credits, err := s.client.GetPersonMovieCredits(int(input.PersonID), nil)
	if err != nil {
		return nil, GetPersonDetailsOutput{}, fmt.Errorf("failed to get person movie credits: %w", err)
	}

	movieCredits := make([]MovieCredit, 0)

	for _, cast := range credits.Cast {
		movieCredits = append(movieCredits, MovieCredit{
			ID:          int(cast.ID),
			Title:       cast.Title,
			Character:   cast.Character,
			ReleaseDate: cast.ReleaseDate,
		})
	}

	for _, crew := range credits.Crew {
		if crew.Job == "Director" || crew.Job == "Producer" || crew.Job == "Writer" || crew.Job == "Screenplay" {
			movieCredits = append(movieCredits, MovieCredit{
				ID:          int(crew.ID),
				Title:       crew.Title,
				Job:         crew.Job,
				ReleaseDate: crew.ReleaseDate,
			})
		}
	}

	if len(movieCredits) > 20 {
		movieCredits = movieCredits[:20]
	}

	details := PersonDetails{
		ID:           int(person.ID),
		Name:         person.Name,
		Biography:    person.Biography,
		Birthday:     person.Birthday,
		Deathday:     person.Deathday,
		KnownFor:     person.KnownForDepartment,
		PlaceOfBirth: person.PlaceOfBirth,
		ProfilePath:  person.ProfilePath,
		Movies:       movieCredits,
	}

	return nil, GetPersonDetailsOutput{Person: details}, nil
}
