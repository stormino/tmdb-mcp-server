package main

type SearchMoviesInput struct {
	Query string `json:"query" jsonschema:"Movie title to search for"`
	Year  int    `json:"year,omitempty" jsonschema:"Optional release year to filter results"`
}

type MovieResult struct {
	ID          int     `json:"id" jsonschema:"TMDB movie ID"`
	Title       string  `json:"title" jsonschema:"Movie title"`
	ReleaseDate string  `json:"release_date" jsonschema:"Release date (YYYY-MM-DD)"`
	Overview    string  `json:"overview" jsonschema:"Movie plot summary"`
	Rating      float32 `json:"rating" jsonschema:"Average user rating (0-10)"`
	PosterPath  string  `json:"poster_path,omitempty" jsonschema:"Path to movie poster image"`
}

type SearchMoviesOutput struct {
	Results []MovieResult `json:"results" jsonschema:"List of matching movies"`
	Total   int64         `json:"total" jsonschema:"Total number of results found"`
}

type GetMovieDetailsInput struct {
	MovieID int `json:"movie_id" jsonschema:"TMDB movie ID"`
}

type MovieDetails struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	ReleaseDate string   `json:"release_date"`
	Overview    string   `json:"overview"`
	Rating      float32  `json:"rating"`
	Runtime     int      `json:"runtime"`
	Budget      int64    `json:"budget"`
	Revenue     int64    `json:"revenue"`
	Genres      []string `json:"genres"`
	PosterPath  string   `json:"poster_path,omitempty"`
	Tagline     string   `json:"tagline,omitempty"`
}

type GetMovieDetailsOutput struct {
	Movie MovieDetails `json:"movie"`
}

type GetTrendingInput struct {
	TimeWindow string `json:"time_window" jsonschema:"Trending time window: 'day' for today's trending movies or 'week' for this week's trending movies,enum=day,enum=week"`
}

type GetTrendingOutput struct {
	Results []MovieResult `json:"results"`
	Total   int           `json:"total"`
}

type GetRecommendationsInput struct {
	MovieID int `json:"movie_id" jsonschema:"TMDB movie ID to get recommendations for"`
	Limit   int `json:"limit,omitempty" jsonschema:"Maximum number of recommendations to return (default 10, max 20)"`
}

type GetRecommendationsOutput struct {
	Results []MovieResult `json:"results"`
	Total   int           `json:"total"`
}

type DiscoverMoviesInput struct {
	Genre     string  `json:"genre,omitempty" jsonschema:"Genre name (e.g., 'Science Fiction', 'Action', 'Drama')"`
	Year      int     `json:"year,omitempty" jsonschema:"Release year (e.g., 2020)"`
	MinRating float32 `json:"min_rating,omitempty" jsonschema:"Minimum average rating (0-10, e.g., 7.5)"`
	SortBy    string  `json:"sort_by,omitempty" jsonschema:"Sort order: 'popularity.desc', 'vote_average.desc', or 'release_date.desc',enum=popularity.desc,enum=vote_average.desc,enum=release_date.desc"`
	Limit     int     `json:"limit,omitempty" jsonschema:"Maximum number of results to return (default 20, max 20)"`
}

type DiscoverMoviesOutput struct {
	Results []MovieResult `json:"results"`
	Total   int           `json:"total"`
}

type SearchPersonInput struct {
	Query string `json:"query" jsonschema:"Person name to search for (actor, director, producer, etc.)"`
}

type PersonResult struct {
	ID          int    `json:"id" jsonschema:"TMDB person ID"`
	Name        string `json:"name" jsonschema:"Person's name"`
	KnownFor    string `json:"known_for" jsonschema:"What they're known for (e.g., 'Acting', 'Directing')"`
	ProfilePath string `json:"profile_path,omitempty" jsonschema:"Path to profile image"`
}

type SearchPersonOutput struct {
	Results []PersonResult `json:"results"`
	Total   int            `json:"total"`
}

type GetPersonDetailsInput struct {
	PersonID int `json:"person_id" jsonschema:"TMDB person ID"`
}

type MovieCredit struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Character   string `json:"character,omitempty"`
	Job         string `json:"job,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
}

type PersonDetails struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Biography    string        `json:"biography"`
	Birthday     string        `json:"birthday,omitempty"`
	Deathday     string        `json:"deathday,omitempty"`
	KnownFor     string        `json:"known_for_department"`
	PlaceOfBirth string        `json:"place_of_birth,omitempty"`
	ProfilePath  string        `json:"profile_path,omitempty"`
	Movies       []MovieCredit `json:"movies"`
}

type GetPersonDetailsOutput struct {
	Person PersonDetails `json:"person"`
}
