# TMDB MCP Server

A Model Context Protocol (MCP) server that provides Claude and other MCP clients access to The Movie Database (TMDB) API.

## Features

- 🎬 Search movies by title
- 📊 Get detailed movie information
- 🔥 Discover trending movies
- 👤 Search actors and directors
- 🎭 Get person details and filmography
- 💡 Get movie recommendations
- 🔍 Advanced movie discovery with filters

## Dual Mode Support

This server works in two modes:

- **Local (stdio)**: For Claude Desktop and local MCP clients
- **Remote (HTTP)**: For claude.ai web and remote integrations

## Prerequisites

- Go 1.21 or later
- TMDB API key ([Get one here](https://www.themoviedb.org/settings/api))

## Installation

```bash
# Clone the repository
git clone https://github.com/stormino/tmdb-mcp-server.git
cd tmdb-mcp-server

# Install dependencies
go mod tidy

# Build the server
go build -o tmdb-mcp-server
```

## Configuration

Set your TMDB API key as an environment variable:

```bash
export TMDB_API_KEY=your_api_key_here
```

## Usage

### Local Mode (Claude Desktop)

1. Build the server (see Installation above)

2. Configure Claude Desktop by editing the config file:
   - **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - **Linux**: `~/.config/Claude/claude_desktop_config.json`
   - **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

3. Add this configuration:

```json
{
  "mcpServers": {
    "tmdb": {
      "command": "/absolute/path/to/tmdb-mcp-server",
      "args": ["--mode", "stdio"],
      "env": {
        "TMDB_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

4. Restart Claude Desktop

5. Ask Claude things like:
   - "Search for movies directed by Christopher Nolan"
   - "What are the trending movies this week?"
   - "Give me details about the movie Inception"

### Remote Mode (HTTP Server)

Run the server in HTTP mode:

```bash
./tmdb-mcp-server --mode http --port 8080
# Or use make:
make run-http
```

#### Testing HTTP Mode Locally

**Option 1: MCP Inspector (Recommended)**
```bash
# Terminal 1: Start server
make run-http

# Terminal 2: Open MCP Inspector
make inspector
# Or: npx @modelcontextprotocol/inspector
```

Then connect to `http://localhost:8080` in the inspector UI at http://localhost:5173

**Option 2: cURL**
```bash
# Test with curl (in another terminal)
make test-http
```

**Option 3: Connect from claude.ai**
1. Expose your local server (e.g., using ngrok: `ngrok http 8080`)
2. Go to claude.ai → Settings → Connectors
3. Click "Add custom connector"
4. Enter your server URL: `https://your-ngrok-url.ngrok.io`

For detailed testing instructions, see [TEST_HTTP_MODE.md](TEST_HTTP_MODE.md)

## Command Line Options

```bash
./tmdb-mcp-server [options]

Options:
  --mode string
        Transport mode: stdio or http (default "stdio")
        Can also be set via MCP_MODE environment variable
  
  --port string
        HTTP port, only used in http mode (default "8080")
        Can also be set via HTTP_PORT environment variable
```

## Available Tools

| Tool Name | Description | Input |
|-----------|-------------|-------|
| `search_movies` | Search for movies by title | query, year (optional) |
| `get_movie_details` | Get detailed movie information | movie_id |
| `get_trending` | Get trending movies | time_window (day/week) |
| `search_person` | Search for people | query |
| `get_person_details` | Get person details and filmography | person_id |
| `get_recommendations` | Get movie recommendations | movie_id, limit |
| `discover_movies` | Advanced movie search with filters | genre, year, min_rating, etc. |

## Example Queries

- "Find movies similar to Interstellar"
- "Who directed The Dark Knight and what are their other movies?"
- "What are the top-rated sci-fi movies from 2023?"
- "Show me trending movies this week"
- "Search for Tom Hanks and list his most popular movies"

Following, a bunch of very complex prompts to test and explore the MCP server capabilities.

### Multi-Tool Chain Prompts

- Deep Filmmaker Analysis:
  
    "Analyze Christopher Nolan's career trajectory: find all his movies, track how his average ratings evolved over time, identify his most frequent collaborators (actors who appeared in 3+ of his films), and create a markdown report showing which genres he's explored and how his box office performance has grown. Include recommendations for similar directors I should explore."

- Franchise Deep Dive:

    "I want to marathon the Marvel Cinematic Universe. Search for all MCU movies, sort them chronologically by release date, tell me which ones are highest rated, estimate total runtime for watching them all, and then for the top 5 rated ones, give me similar movie recommendations that aren't MCU films."

- Actor Career Comparison:

    "Compare the careers of Leonardo DiCaprio and Tom Hanks: search both actors, get their complete filmographies, calculate their average movie ratings, identify their peak decade, find any movies they've both been in, and recommend 3 movies from each that I absolutely shouldn't miss based on ratings and popularity."

### Trending & Discovery

- Intelligent Trend Analysis:

    "What are the trending movies this week? For each one, tell me: why it's trending (analyze release date, rating, genre), get similar movie recommendations, and create a ranked list of which trending movies I should prioritize based on ratings above 7.5 and runtime under 2.5 hours."

- Genre Exploration Challenge:

    "I want to explore sci-fi movies I've never heard of. Use discover to find sci-fi movies from 2020-2024 with ratings above 7.0 that aren't in the top 20 most popular. For the top 5 results, get full details and tell me what makes each one unique. Then find similar hidden gems."

- Decade Time Machine:

    "Take me on a tour of the 1990s: find the highest-rated action movies from that decade, get details on the top 5, identify which actors appeared in multiple films, and create a '90s action movie starter pack' with viewing order based on release date and thematic connections."

### Creative & Complex

- Perfect Movie Night Curator:

    "I'm hosting a movie night for 4 people. We loved Inception, The Shawshank Redemption, and Parasite. Find movies that are similar to at least 2 of these, filter for runtime under 150 minutes and rating above 8.0, check if any share directors or actors, and create a ranked shortlist of 5 perfect options with explanations why each would work."

- Director's Vision Evolution:

    "Track Quentin Tarantino's evolution as a director: find all his movies chronologically, analyze how his ratings changed over time, identify recurring actors in his films, get details on his highest and lowest rated films, and create a 'Tarantino in Three Acts' viewing guide that shows his early, middle, and recent period masterpieces."

- Ultimate Heist Movie Database:

    "Build me the ultimate heist movie collection: discover all heist/crime movies with ratings above 7.0, group them by decade, find which actors appear most frequently in heist films, get recommendations for deep cuts I might have missed, and create a tiered watchlist from 'Essential' to 'Hidden Gems'."

### Data Analysis & Insights

- Rating Pattern Analysis:

    "Find all Best Picture Oscar winners from the last 20 years, get their TMDB ratings, and tell me: Are Oscar winners actually better rated by audiences? Which winner has the highest rating? Get recommendations similar to the top 3, and identify any patterns in genres or directors."

- Streaming Strategy:

    "I have one free weekend. Search for trending movies this week, get full details on all of them, calculate which combination of 5-6 movies would give me the most variety (different genres), highest average rating, and fit within 12 hours of total runtime. Create an optimal viewing schedule."

- Franchise Fatigue Check:

    "Are movie franchises getting worse over time? Pick 3 major franchises (search for Star Wars, Fast and Furious, and Mission Impossible movies), track how ratings changed from first to most recent film, identify the peak film in each franchise, and tell me if there's evidence of declining quality."

### Stress Test / Show-Off Prompts

- The TMDB Grand Tour:

    "Let's do the ultimate test: (1) Find what's trending today, (2) Get details on the #1 trending movie, (3) Find the director of that movie, (4) Get their complete filmography, (5) Find recommendations based on their highest-rated film, (6) For those recommendations, identify any shared actors, (7) Create a comprehensive viewing guide with ratings, genres, and why I should watch each one."

- Cinematic Universe Builder:

    "I want to create my own themed movie marathon. Search for movies with 'time travel' in the plot, get details on the top 10 rated ones, find which actors appear in multiple time travel movies, get recommendations for each that might also involve time/reality manipulation, and build a 'Temporal Cinema Festival' with viewing order, brief descriptions, and fun facts about connections between films."

- The Six Degrees Challenge:

    "Start with Kevin Bacon. Get his details and top 10 movies. Pick 3 of those movies and find their full cast details. For those cast members, find their other highest-rated movies. Create a connection map showing how different actors and films link together, and recommend a viewing path through this network."

### Creative Output Formats
    
- Movie Recommendation Newsletter:

    "Create a movie newsletter: Find trending movies this week, discover hidden gems (high rated but low popularity), get a classic recommendation (pre-2000, rating >8), and format everything as a markdown newsletter with sections, ratings, emoji, and personalized commentary on why each is worth watching."

- Film Festival Programmer:

    "I'm programming a mini film festival. Theme: 'Dreams and Reality'. Discover movies with high ratings that match this theme across different decades, get full details, find directorial connections, and create a festival program with opening night pick, centerpiece, and closing film, complete with descriptions of why each fits the theme."

- Algorithmic Movie Critic:

    "Act as a movie critic AI: Search for this week's trending movies, get detailed information on each, analyze their ratings vs budget (if available), identify which genres are performing well right now, and write a trend analysis piece about what's working in cinema today based on the data."

## Development

### Architecture

This server implements the Model Context Protocol (MCP) with:
- **Client-Host-Server architecture** built on JSON-RPC 2.0
- **Dual transport mode**: Same codebase works locally (stdio) and remotely (HTTP+SSE)
- **7 fully implemented tools** for movie discovery and information retrieval

### Tech Stack

- **Language**: Go 1.25+
- **MCP SDK**: `github.com/modelcontextprotocol/go-sdk` v1.0.0
- **TMDB Library**: `github.com/cyruzin/golang-tmdb` v1.8.2

### Building

```bash
# Install dependencies
go mod download
go mod tidy

# Build
go build -o tmdb-mcp-server

# Or use Makefile
make build
```

### Testing

```bash
# Test stdio mode
./tmdb-mcp-server --mode stdio

# Test HTTP mode
./tmdb-mcp-server --mode http --port 8080

# Run tests (when implemented)
go test -v ./...

# Or use Makefile
make test
```

## Docker Support (Coming Soon)

```bash
# Build Docker image
docker build -t tmdb-mcp-server .

# Run in local mode
docker run -e TMDB_API_KEY=your_key tmdb-mcp-server

# Run in HTTP mode
docker run -p 8080:8080 -e HTTP_PORT=8080 -e TMDB_API_KEY=your_key tmdb-mcp-server
```

## Troubleshooting

### "TMDB_API_KEY environment variable is required"
Make sure you've set the TMDB_API_KEY environment variable or included it in your Claude Desktop config.

### Server not connecting in Claude Desktop
1. Check the path to the executable is absolute
2. Verify the TMDB_API_KEY is valid
3. Check Claude Desktop logs for errors
4. Try running the server manually to test: `./tmdb-mcp-server --mode stdio`

## Contributing

Contributions are welcome! Areas for improvement:

**Priority Enhancements:**
- Unit and integration tests
- Caching layer for API responses
- Rate limiting implementation
- Streaming availability integration

**Additional Features:**
- MCP Resources (e.g., `tmdb://trending/movies/week`)
- MCP Prompts (e.g., "Analyze movie reviews")
- Image URL construction helpers
- Multi-language support
- Pagination for large result sets
- Authentication for HTTP mode
- Docker deployment support for HTTP mode

**Questions to Consider:**
- Should we add streaming availability data? (requires additional API)
- Should we support multiple languages for results?
- Should we add TV show support in addition to movies?
- Should we implement WebSocket transport in addition to HTTP+SSE?

## License

MIT License - feel free to use this in your own projects!

## Resources

- [Model Context Protocol Documentation](https://modelcontextprotocol.io)
- [TMDB API Documentation](https://developer.themoviedb.org/reference/getting-started)
- [golang-tmdb Library](https://github.com/cyruzin/golang-tmdb)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)

## Credits

Built using:
- [Model Context Protocol](https://modelcontextprotocol.io) by Anthropic
- [golang-tmdb](https://github.com/cyruzin/golang-tmdb) by cyruzin
- [TMDB API](https://www.themoviedb.org) for movie data