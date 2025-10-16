# Testing TMDB MCP Server in HTTP Mode

This guide explains how to test the MCP server in HTTP mode locally.

## Prerequisites

- TMDB API key set in environment: `export TMDB_API_KEY=your_key`
- Server built: `make build` or `go build -o tmdb-mcp-server`

---

## Method 1: MCP Inspector (Recommended)

The **MCP Inspector** is the official tool from Anthropic for testing MCP servers.

### Installation

```bash
# Using npx (no installation needed)
npx @modelcontextprotocol/inspector

# Or install globally
npm install -g @modelcontextprotocol/inspector
mcp-inspector
```

### Usage

1. **Start your MCP server in HTTP mode:**
   ```bash
   # Default port 8080
   ./tmdb-mcp-server --mode http --port 8080

   # Or use make (default port 8080):
   make run-http

   # Custom port with make:
   make run-http HTTP_PORT=3000

   # Or with environment variable:
   HTTP_PORT=9000 make run-http
   ```

2. **Open MCP Inspector in your browser:**
   ```bash
   npx @modelcontextprotocol/inspector
   ```
   This will open `http://localhost:5173` in your browser.

3. **Connect to your server:**
   - In the inspector UI, select "HTTP/SSE" as the transport
   - Enter the URL: `http://localhost:8080`
   - Click "Connect"

4. **Test your tools:**
   - Browse available tools in the left sidebar
   - Click on any tool (e.g., `search_movies`)
   - Fill in the parameters
   - Click "Execute" to see results

### Example Test Flow

```
1. Tool: search_movies
   Input: {"query": "Inception", "year": 2010}

2. Tool: get_movie_details
   Input: {"movie_id": 27205}

3. Tool: get_recommendations
   Input: {"movie_id": 27205, "limit": 5}
```

---

## Method 2: cURL Commands

Test the server using simple HTTP requests:

### 1. Start the server
```bash
# Default port 8080
./tmdb-mcp-server --mode http --port 8080

# Custom port
./tmdb-mcp-server --mode http --port 3000

# Or with make
make run-http HTTP_PORT=3000
```

### 2. Initialize connection
```bash
curl -X POST http://localhost:8080/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "1.0.0",
      "clientInfo": {
        "name": "test-client",
        "version": "1.0.0"
      }
    }
  }'
```

### 3. List available tools
```bash
curl -X POST http://localhost:8080/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list"
  }'
```

### 4. Call a tool (search movies)
```bash
curl -X POST http://localhost:8080/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "search_movies",
      "arguments": {
        "query": "Inception",
        "year": 2010
      }
    }
  }'
```

### 5. Get movie details
```bash
curl -X POST http://localhost:8080/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/call",
    "params": {
      "name": "get_movie_details",
      "arguments": {
        "movie_id": 27205
      }
    }
  }'
```

### 6. Get recommendations
```bash
curl -X POST http://localhost:8080/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 5,
    "method": "tools/call",
    "params": {
      "name": "get_recommendations",
      "arguments": {
        "movie_id": 27205,
        "limit": 5
      }
    }
  }'
```

---

## Method 3: Quick Test Script

Create a test script for quick testing:

```bash
#!/bin/bash
# test-http.sh

echo "Testing TMDB MCP Server in HTTP mode..."
echo ""

# Start server in background
./tmdb-mcp-server --mode http --port 8080 &
SERVER_PID=$!
echo "Server started with PID: $SERVER_PID"
sleep 2

# Test 1: Search for a movie
echo ""
echo "=== Test 1: Search for 'Inception' ==="
curl -s -X POST http://localhost:8080/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "search_movies",
      "arguments": {"query": "Inception"}
    }
  }' | jq '.result.content[0].text | fromjson | .results[0]'

# Test 2: Get trending movies
echo ""
echo "=== Test 2: Get trending movies ==="
curl -s -X POST http://localhost:8080/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/call",
    "params": {
      "name": "get_trending",
      "arguments": {"time_window": "week"}
    }
  }' | jq '.result.content[0].text | fromjson | .results[0:3]'

# Cleanup
echo ""
echo "Stopping server..."
kill $SERVER_PID
echo "Done!"
```

Save as `test-http.sh`, make executable, and run:
```bash
chmod +x test-http.sh
./test-http.sh
```

---

## Method 4: Connect from claude.ai Web

Once you have the server running and exposed (e.g., via ngrok or on a public server):

1. **Start the server:**
   ```bash
   ./tmdb-mcp-server --mode http --port 8080
   ```

2. **Expose it (if testing locally):**
   ```bash
   # Using ngrok
   ngrok http 8080
   ```

3. **Add to claude.ai:**
   - Go to claude.ai
   - Navigate to Settings → Integrations
   - Click "Add custom MCP server"
   - Enter your URL (e.g., `https://your-ngrok-url.ngrok.io`)

---

## Method 5: Using Postman or Insomnia

1. Import the JSON-RPC 2.0 requests from the cURL examples above
2. Create a collection with all the tool calls
3. Test each tool interactively with a nice UI

---

## Troubleshooting

### Server won't start
```bash
# Check if port is already in use
lsof -i :8080

# Kill existing process
kill $(lsof -t -i :8080)

# Or use a different port
./tmdb-mcp-server --mode http --port 8081
```

### Can't connect from MCP Inspector
- Ensure server is running: `curl http://localhost:8080`
- Check firewall settings
- Verify TMDB_API_KEY is set

### Tools return errors
- Verify TMDB_API_KEY is valid
- Check TMDB API rate limits (40 requests per 10 seconds)
- Look at server logs for detailed error messages

---

## Quick Start

The fastest way to test:

```bash
# Terminal 1: Start server
make run-http

# Terminal 2: Run MCP Inspector
npx @modelcontextprotocol/inspector
```

Then open http://localhost:5173 in your browser and start testing! 🚀
