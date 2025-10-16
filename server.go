package main

import (
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TMDBServer struct {
	client *tmdb.Client
}

func NewTMDBServer(apiKey string) *TMDBServer {
	client, err := tmdb.Init(apiKey)
	if err != nil {
		panic(err)
	}
	return &TMDBServer{
		client: client,
	}
}

func (s *TMDBServer) RegisterTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_movies",
		Description: "Search for movies by title to find their TMDB ID and basic information. Use this when you need to find a specific movie mentioned by name (e.g., 'Inception', 'The Matrix'). Returns movie ID, title, release date, overview, rating, and poster. Required when you need a movie's ID for other operations like getting details or recommendations.",
	}, s.SearchMovies)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_movie_details",
		Description: "Get comprehensive details about a specific movie using its TMDB ID. Returns full information including runtime, budget, revenue, genres, tagline, and more. Use this after searching for a movie to get complete information, or when you already have a movie ID and need detailed data.",
	}, s.GetMovieDetails)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_trending",
		Description: "Get movies that are currently trending on TMDB, either today or this week. Use this when users ask about 'trending', 'popular right now', 'what's hot', or 'current popular movies'. Does not require any movie name - provides a curated list of what's trending.",
	}, s.GetTrending)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_person",
		Description: "Search for actors, directors, producers, or other people in the film industry by name. Use when users ask about a specific person, actor, or director. Returns person ID, name, known for, and profile image. Required when you need a person's ID for getting their detailed filmography.",
	}, s.SearchPerson)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_person_details",
		Description: "Get detailed information about a person including biography, birthday, and complete filmography. Use when users want to know more about an actor or director, their background, or what movies they've been in. Requires person ID from search_person.",
	}, s.GetPersonDetails)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_recommendations",
		Description: "Get personalized movie recommendations based on a specific movie the user enjoyed. IMPORTANT: Use this when users say they 'liked', 'loved', 'enjoyed' a movie and want 'similar', 'recommendations', or 'movies like' that one. First use search_movies to get the movie ID, then use this tool with that ID. Returns movies that are similar in genre, theme, and style.",
	}, s.GetRecommendations)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "discover_movies",
		Description: "Discover movies using advanced filters for genre, year, minimum rating, and sorting. Use this when users want movies matching specific criteria like 'sci-fi movies from 2020s', 'highly rated action films', 'recent comedies', or combining multiple attributes. More flexible than recommendations when users specify genres, years, or ratings rather than a specific movie they liked.",
	}, s.DiscoverMovies)
}
