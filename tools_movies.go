package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *TMDBServer) SearchMovies(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SearchMoviesInput,
) (*mcp.CallToolResult, SearchMoviesOutput, error) {
	options := make(map[string]string)
	if input.Year > 0 {
		options["year"] = strconv.Itoa(input.Year)
	}

	results, err := s.client.GetSearchMovies(input.Query, options)
	if err != nil {
		return nil, SearchMoviesOutput{}, fmt.Errorf("failed to search movies: %w", err)
	}

	movies := make([]MovieResult, 0, len(results.Results))
	for _, movie := range results.Results {
		movies = append(movies, MovieResult{
			ID:          int(movie.ID),
			Title:       movie.Title,
			ReleaseDate: movie.ReleaseDate,
			Overview:    movie.Overview,
			Rating:      movie.VoteAverage,
			PosterPath:  movie.PosterPath,
		})
	}

	return nil, SearchMoviesOutput{
		Results: movies,
		Total:   results.TotalResults,
	}, nil
}

func (s *TMDBServer) GetMovieDetails(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetMovieDetailsInput,
) (*mcp.CallToolResult, GetMovieDetailsOutput, error) {
	movie, err := s.client.GetMovieDetails(int(input.MovieID), nil)
	if err != nil {
		return nil, GetMovieDetailsOutput{}, fmt.Errorf("failed to get movie details: %w", err)
	}

	genres := make([]string, 0, len(movie.Genres))
	for _, genre := range movie.Genres {
		genres = append(genres, genre.Name)
	}

	details := MovieDetails{
		ID:          int(movie.ID),
		Title:       movie.Title,
		ReleaseDate: movie.ReleaseDate,
		Overview:    movie.Overview,
		Rating:      movie.VoteAverage,
		Runtime:     movie.Runtime,
		Budget:      movie.Budget,
		Revenue:     movie.Revenue,
		Genres:      genres,
		PosterPath:  movie.PosterPath,
		Tagline:     movie.Tagline,
	}

	return nil, GetMovieDetailsOutput{Movie: details}, nil
}

func (s *TMDBServer) GetTrending(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetTrendingInput,
) (*mcp.CallToolResult, GetTrendingOutput, error) {
	timeWindow := input.TimeWindow
	if timeWindow != "day" && timeWindow != "week" {
		timeWindow = "week"
	}

	results, err := s.client.GetTrending("movie", timeWindow, nil)
	if err != nil {
		return nil, GetTrendingOutput{}, fmt.Errorf("failed to get trending movies: %w", err)
	}

	movies := make([]MovieResult, 0, len(results.Results))
	for _, result := range results.Results {
		movies = append(movies, MovieResult{
			ID:          int(result.ID),
			Title:       result.Title,
			ReleaseDate: result.ReleaseDate,
			Overview:    result.Overview,
			Rating:      result.VoteAverage,
			PosterPath:  result.PosterPath,
		})
	}

	return nil, GetTrendingOutput{
		Results: movies,
		Total:   len(movies),
	}, nil
}

func (s *TMDBServer) GetRecommendations(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetRecommendationsInput,
) (*mcp.CallToolResult, GetRecommendationsOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 20 {
		limit = 20
	}

	results, err := s.client.GetMovieRecommendations(input.MovieID, nil)
	if err != nil {
		return nil, GetRecommendationsOutput{}, fmt.Errorf("failed to get movie recommendations: %w", err)
	}

	movies := make([]MovieResult, 0, limit)
	for i, movie := range results.Results {
		if i >= limit {
			break
		}
		movies = append(movies, MovieResult{
			ID:          int(movie.ID),
			Title:       movie.Title,
			ReleaseDate: movie.ReleaseDate,
			Overview:    movie.Overview,
			Rating:      movie.VoteAverage,
			PosterPath:  movie.PosterPath,
		})
	}

	return nil, GetRecommendationsOutput{
		Results: movies,
		Total:   len(movies),
	}, nil
}

func (s *TMDBServer) DiscoverMovies(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DiscoverMoviesInput,
) (*mcp.CallToolResult, DiscoverMoviesOutput, error) {
	urlOptions := make(map[string]string)

	sortBy := input.SortBy
	if sortBy == "" {
		sortBy = "popularity.desc"
	}
	urlOptions["sort_by"] = sortBy

	if input.Genre != "" {
		genreMap := map[string]string{
			"action":          "28",
			"adventure":       "12",
			"animation":       "16",
			"comedy":          "35",
			"crime":           "80",
			"documentary":     "99",
			"drama":           "18",
			"family":          "10751",
			"fantasy":         "14",
			"history":         "36",
			"horror":          "27",
			"music":           "10402",
			"mystery":         "9648",
			"romance":         "10749",
			"science fiction": "878",
			"sci-fi":          "878",
			"tv movie":        "10770",
			"thriller":        "53",
			"war":             "10752",
			"western":         "37",
		}

		genreName := input.Genre
		for key, id := range genreMap {
			if key == genreName {
				urlOptions["with_genres"] = id
				break
			}
		}
		if _, ok := urlOptions["with_genres"]; !ok {
			urlOptions["with_genres"] = genreName
		}
	}

	if input.Year > 0 {
		urlOptions["primary_release_year"] = strconv.Itoa(input.Year)
	}

	if input.MinRating > 0 {
		urlOptions["vote_average.gte"] = fmt.Sprintf("%.1f", input.MinRating)
		urlOptions["vote_count.gte"] = "100"
	}

	results, err := s.client.GetDiscoverMovie(urlOptions)
	if err != nil {
		return nil, DiscoverMoviesOutput{}, fmt.Errorf("failed to discover movies: %w", err)
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 20 {
		limit = 20
	}

	movies := make([]MovieResult, 0, limit)
	for i, movie := range results.Results {
		if i >= limit {
			break
		}
		movies = append(movies, MovieResult{
			ID:          int(movie.ID),
			Title:       movie.Title,
			ReleaseDate: movie.ReleaseDate,
			Overview:    movie.Overview,
			Rating:      movie.VoteAverage,
			PosterPath:  movie.PosterPath,
		})
	}

	return nil, DiscoverMoviesOutput{
		Results: movies,
		Total:   len(movies),
	}, nil
}
