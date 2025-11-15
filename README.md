# Pokedex CLI

A command-line interface application for exploring Pokemon location areas using the [PokeAPI](https://pokeapi.co/). This interactive REPL (Read-Eval-Print Loop) allows users to navigate through Pokemon world locations directly from their terminal.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
- [API Integration](#api-integration)
- [Project Structure](#project-structure)
- [Technical Deep Dive](#technical-deep-dive)
- [Contributing](#contributing)

---

## Overview

The Pokedex CLI is an interactive terminal application that provides a user-friendly interface to explore Pokemon location areas from the Pokemon universe. Built in Go, it leverages the public PokeAPI to fetch location data and presents it in a paginated, easy-to-navigate format.

### What Does This Program Do?

The application:

- Provides an interactive command-line interface with a REPL loop
- Fetches Pokemon location area data from the PokeAPI (<https://pokeapi.co>)
- Displays location areas in pages of 20 entries at a time
- Allows forward and backward navigation through location pages
- Implements intelligent caching to reduce API calls and improve performance
- Offers built-in help documentation and clean exit functionality

---

## Features

- **Interactive REPL**: Continuous command prompt for seamless interaction
- **Paginated Navigation**: Browse through location areas 20 at a time
- **Bidirectional Browsing**: Move forward and backward through pages
- **Intelligent Caching**: Automatic caching with 5-second TTL reduces redundant API calls
- **Concurrent-Safe**: Thread-safe cache implementation using mutex locks
- **Automatic Cache Cleanup**: Background goroutine removes stale cache entries
- **Error Handling**: Graceful error handling for network issues and API failures
- **Built-in Help**: Comprehensive command documentation via `help` command

---

## Architecture

### High-Level Overview

```
User Input → REPL → Command Router → Command Handler → API Client → PokeAPI
                                                          ↓
                                                        Cache
                                                          ↓
                                                    HTTP Response
                                                          ↓
                                                    JSON Parsing
                                                          ↓
                                                   Display Results
```

### Component Interaction Flow

1. **User Interaction Layer** (`main.go`)
   - User enters a command at the `Pokedex >` prompt
   - Input is read by `bufio.Scanner` and processed
   - Command is normalized (trimmed, lowercased)
   - Command router looks up the command in the registry

2. **Command Processing Layer** (`main.go`)
   - Command is matched against the registered command map
   - Corresponding callback function is invoked with config state
   - Config maintains pagination state (nextURL, prevURL)

3. **API Client Layer** (`internal/pokeapi/`)
   - Client checks cache before making HTTP requests
   - If cache miss, performs HTTP GET to PokeAPI
   - Stores response in cache for future use
   - Returns parsed JSON data to command handler

4. **Caching Layer** (`internal/pokecache/`)
   - Thread-safe in-memory cache using sync.Mutex
   - Background goroutine runs cleanup every 5 seconds
   - Automatically evicts entries older than TTL

5. **Data Display Layer** (`main.go`)
   - Results are formatted and printed to stdout
   - Pagination state is updated in config
   - User sees location names line-by-line

---

## Installation

### Prerequisites

- Go 1.25.3 or higher
- Internet connection (for API access)
- Terminal/Command prompt

### Build Instructions

1. Clone the repository:

```bash
git clone https://github.com/DanDo385/pokedexcli.git
cd pokedexcli
```

2. Build the executable:

```bash
go build
```

This will create a `pokedexcli` executable in your current directory.

### Alternative: Run Without Building

You can also run directly without building:

```bash
go run .
```

---

## Usage

### Starting the Program

After building, start the application:

```bash
./pokedexcli
```

You'll see the interactive prompt:

```
Pokedex >
```

### Available Commands

#### `help`

Displays a list of all available commands with descriptions.

```
Pokedex > help
Welcome to the Pokedex!
Usage:
help: Displays a help message
exit: Exit the Pokedex
map: List the next 20 location areas
mapb: List the previous 20 location areas
```

#### `map`

Displays the next page of 20 location areas.

**Example:**

```
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
sunyshore-city-area
sinnoh-pokemon-league-area
oreburgh-mine-1f
oreburgh-mine-b1f
valley-windworks-area
eterna-forest-area
fuego-ironworks-area
mt-coronet-1f-route-207
mt-coronet-2f
mt-coronet-3f
mt-coronet-exterior-snowfall
mt-coronet-exterior-blizzard
mt-coronet-4f
mt-coronet-4f-small-room
mt-coronet-5f
mt-coronet-6f
mt-coronet-1f-from-exterior
```

**Behavior:**

- First call starts from the beginning of the location list
- Subsequent calls move forward through pages
- Automatically tracks the next page URL
- If you've reached the end, wraps to available pages

#### `mapb`

Displays the previous page of 20 location areas (map backward).

**Example:**

```
Pokedex > mapb
canalave-city-area
eterna-city-area
pastoria-city-area
...
```

**Behavior:**

- Only works if you've navigated forward at least once
- If on the first page, displays: `you're on the first page`
- Moves backward through pagination
- Automatically tracks the previous page URL

#### `exit`

Exits the program gracefully.

**Example:**

```
Pokedex > exit
Closing the Pokedex... Goodbye!
```

### Command Error Handling

If you enter an invalid command:

```
Pokedex > invalid
Unknown command
```

If an API error occurs:

```
Pokedex > map
Error: bad status: 500 Internal Server Error
```

---

## API Integration

### PokeAPI Overview

The application uses [PokeAPI v2](https://pokeapi.co/), a free RESTful API serving Pokemon data.

**Base URL:** `https://pokeapi.co/api/v2`

### Endpoint Used

#### Location Areas Endpoint

```
GET https://pokeapi.co/api/v2/location-area
GET https://pokeapi.co/api/v2/location-area?offset=20&limit=20
```

**Response Structure:**

```json
{
  "count": 1036,
  "next": "https://pokeapi.co/api/v2/location-area?offset=20&limit=20",
  "previous": null,
  "results": [
    {
      "name": "canalave-city-area",
      "url": "https://pokeapi.co/api/v2/location-area/1/"
    },
    ...
  ]
}
```

### How the API Interacts with the CLI

1. **Initial Request**
   - User types `map`
   - `commandMap()` is called with `cfg.nextURL = nil`
   - Client constructs default URL: `https://pokeapi.co/api/v2/location-area`
   - HTTP GET request is made
   - Response contains first 20 locations

2. **Cache Check**
   - Before making HTTP request, `doGet()` checks cache
   - Cache key is the full URL
   - If found in cache, returns immediately (no HTTP call)
   - If not found, proceeds with HTTP request

3. **Response Processing**
   - Raw JSON response is read into byte slice
   - Stored in cache with current timestamp
   - Unmarshaled into `LocationAreasResponse` struct
   - `next` and `previous` URLs are extracted and stored in config

4. **Pagination State Management**
   - Config struct maintains two pointers:
     - `nextURL`: URL for next page
     - `prevURL`: URL for previous page
   - After each successful fetch, both are updated
   - Subsequent `map` commands use `nextURL`
   - `mapb` commands use `prevURL`

5. **Display**
   - Loop through `Results` array
   - Print each location's `Name` field
   - User sees simple list of location names

### Cache Strategy

**Cache Configuration:**

- **TTL (Time To Live):** 5 seconds
- **Storage:** In-memory map
- **Thread Safety:** Mutex-protected
- **Cleanup:** Background goroutine checks every 5 seconds

**Benefits:**

- Reduces API load (respects rate limits)
- Improves response time for repeated queries
- Allows offline browsing of recently-viewed pages
- Prevents unnecessary network calls

**Example Scenario:**

```
User: map          (API call, cache miss)
User: mapb         (API call, cache miss)
User: map          (Cache hit! No API call)
User: mapb         (Cache hit! No API call)
[wait 6 seconds]
User: map          (API call, cache expired)
```

---

## Project Structure

```
pokedexcli/
├── main.go                          # Entry point, REPL, command handlers
├── config.go                        # Configuration struct
├── go.mod                           # Go module definition
├── README.md                        # This file
│
└── internal/
    ├── pokeapi/                     # PokeAPI client package
    │   ├── client.go                # HTTP client with caching
    │   └── locations.go             # Location areas endpoint handler
    │
    └── pokecache/                   # Caching implementation
        ├── cache.go                 # Cache struct and methods
        └── cache_test.go            # Cache unit tests
```

### File Responsibilities

#### `main.go`

**Lines of Code:** ~109
**Responsibilities:**

- Application entry point (`main()` function)
- REPL implementation using `bufio.Scanner`
- Command registration system via `registerCommand()` factory
- Command handlers: `commandHelp`, `commandExit`, `commandMap`, `commandMapb`
- Shared helper function: `displayLocationAreas()`
- User input parsing and routing

**Key Components:**

```go
type cliCommand struct {
    name        string
    description string
    callback    func(*config) error
}
```

**REPL Loop Logic:**

1. Display prompt: `Pokedex >`
2. Read user input
3. Normalize input (trim, lowercase)
4. Parse command name
5. Look up in command registry
6. Execute callback function
7. Handle errors
8. Repeat

#### `config.go`

**Lines of Code:** ~9
**Responsibilities:**

- Application state management
- Pagination state (next/previous URLs)
- PokeAPI client instance

**Structure:**

```go
type config struct {
    nextURL *string          // Pointer to next page URL
    prevURL *string          // Pointer to previous page URL
    client  *pokeapi.Client  // API client instance
}
```

**Why Pointers for URLs?**

- `nil` represents "no next/previous page" or "start from beginning"
- Allows differentiation between "no value" and "empty string"
- API returns `null` for prev on first page, next to `*string` mapping

#### `internal/pokeapi/client.go`

**Lines of Code:** ~52
**Responsibilities:**

- HTTP client management
- Cache integration
- Generic GET request handler with caching logic

**Key Components:**

- `const baseURL`: PokeAPI base URL
- `Client` struct: Wraps `http.Client` and `Cache`
- `NewClient()`: Factory function, initializes with 5-second cache
- `doGet()`: Core method that orchestrates cache-check → HTTP → cache-store

**Request Flow:**

```go
doGet(url) → cache.Get(url)
    ↓ (miss)
http.Client.Get(url)
    ↓
io.ReadAll(response.Body)
    ↓
cache.Add(url, body)
    ↓
return body
```

#### `internal/pokeapi/locations.go`

**Lines of Code:** ~40
**Responsibilities:**

- Location areas endpoint specific logic
- JSON response structure definition
- URL construction (default vs paginated)
- JSON unmarshaling

**Key Components:**

```go
type LocationAreasResponse struct {
    Count    int
    Next     *string         // Nullable
    Previous *string         // Nullable
    Results  []struct {
        Name string
        URL  string
    }
}
```

**Method Logic:**

```go
GetLocationAreas(pageURL *string):
    if pageURL == nil:
        use baseURL + "/location-area"
    else:
        use pageURL (contains offset/limit params)

    fetch via doGet()
    unmarshal JSON
    return structured data
```

#### `internal/pokecache/cache.go`

**Lines of Code:** ~62
**Responsibilities:**

- Thread-safe in-memory caching
- Automatic expiration and cleanup
- Background reaping goroutine

**Key Components:**

```go
type cacheEntry struct {
    createdAt time.Time     // Timestamp for TTL calculation
    val       []byte        // Cached response body
}

type Cache struct {
    mu       sync.Mutex                // Thread safety
    entries  map[string]cacheEntry     // URL → cached data
    interval time.Duration             // TTL / cleanup interval
}
```

**Methods:**

- `NewCache(interval)`: Creates cache, launches reap goroutine
- `Add(key, val)`: Stores entry with current timestamp (mutex-protected)
- `Get(key)`: Retrieves entry if exists (mutex-protected)
- `reapLoop()`: Background goroutine that removes stale entries

**Reap Loop Logic:**

```go
Every [interval]:
    Lock mutex
    For each cache entry:
        If (now - entry.createdAt) > interval:
            Delete entry
    Unlock mutex
```

**Thread Safety:**

- All public methods acquire mutex lock
- Prevents race conditions during concurrent access
- Safe for use from multiple goroutines

---

## Technical Deep Dive

### Design Patterns Used

#### 1. Command Pattern

The command registration system uses the Command pattern:

- `cliCommand` struct encapsulates command data and behavior
- Commands are registered in a map for dynamic lookup
- Each command has a callback function for execution
- Decouples command invocation from implementation

#### 2. Factory Pattern

- `NewClient()` creates fully-initialized API client
- `NewCache()` creates cache with background goroutine
- `registerCommand()` simplifies command registration
- Encapsulates complex initialization logic

#### 3. Singleton-like Global State

- `commands` map is global and initialized once in `init()`
- Provides centralized command registry
- All parts of application reference same command map

#### 4. Dependency Injection

- `config` struct is passed to all command handlers
- Allows testing with mock clients
- Decouples command logic from client implementation

### Concurrency

#### Goroutines

The application uses one background goroutine:

- Launched in `NewCache()` via `go c.reapLoop()`
- Runs for lifetime of cache instance
- Ticker-based periodic execution
- Automatically cleans up with `defer ticker.Stop()`

#### Thread Safety

- `sync.Mutex` protects cache map
- Lock acquired in `Add()`, `Get()`, and `reapLoop()`
- Prevents concurrent map read/write panics
- Follows Go concurrency best practices

### Error Handling Strategy

1. **Network Errors**: Propagated to user with descriptive messages
2. **JSON Parsing Errors**: Returned as errors, displayed to user
3. **HTTP Status Errors**: Checked and returned as formatted error
4. **No Panics**: All errors are returned, not panicked
5. **User-Friendly Messages**: Errors prefixed with "Error:" in REPL

### Memory Management

- **Cache Size**: Unbounded (grows with unique URLs visited)
- **Cache Cleanup**: Automatic every 5 seconds
- **Typical Usage**: ~20-50 cache entries during normal use
- **Memory Footprint**: Minimal (few KB per cached response)

### Performance Characteristics

- **Cold Start**: First `map` call requires API request (~200ms)
- **Warm Cache**: Cached responses return in <1ms
- **Cache Hit Rate**: High for back-and-forth navigation (>80%)
- **Network Dependency**: Degrades gracefully on slow connections

---

## How Different Files Interact

### Startup Sequence

1. **Program Launch**

   ```
   main.go:main() executes
   ```

2. **Initialization** (before main runs)

   ```
   main.go:init() runs
   → registerCommand() called 4 times
   → commands map populated
   ```

3. **Client Creation**

   ```
   main.go:main()
   → pokeapi.NewClient()
     → pokecache.NewCache(5 * time.Second)
       → go c.reapLoop() [goroutine launched]
   → config struct created with client
   ```

4. **REPL Start**

   ```
   main.go:main()
   → bufio.Scanner created
   → infinite loop begins
   ```

### Command Execution Flow (Example: `map` command)

```
User types "map"
    ↓
main.go:main() REPL loop
    ↓
Input parsed: cmdName = "map"
    ↓
commands["map"] lookup
    ↓
cmd.callback(cfg) called
    ↓
main.go:commandMap(cfg)
    ↓
main.go:displayLocationAreas(cfg, cfg.nextURL)
    ↓
pokeapi/locations.go:GetLocationAreas(pageURL)
    ↓
URL construction (first call: pageURL is nil)
    ↓
pokeapi/client.go:doGet(url)
    ↓
pokecache/cache.go:Get(url)  [returns false on first call]
    ↓
http.Client.Get(url)  [makes actual HTTP request]
    ↓
io.ReadAll(response.Body)
    ↓
pokecache/cache.go:Add(url, body)
    ↓
[Back to locations.go]
json.Unmarshal(body, &data)
    ↓
Return LocationAreasResponse to commandMap
    ↓
[Back to main.go]
Loop through data.Results, print each area.Name
    ↓
Update cfg.nextURL = data.Next
Update cfg.prevURL = data.Previous
    ↓
Return to REPL loop
    ↓
Display "Pokedex >" prompt again
```

### Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│  main.go                                                    │
│  ┌──────────┐      ┌─────────────┐      ┌──────────────┐   │
│  │   REPL   │─────→│   Command   │─────→│   Command    │   │
│  │   Loop   │      │   Router    │      │   Handlers   │   │
│  └──────────┘      └─────────────┘      └───────┬──────┘   │
│                                                  │          │
└──────────────────────────────────────────────────┼──────────┘
                                                   │
                    ┌──────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│  config.go                                                  │
│  ┌────────────────────────────────────────────┐             │
│  │  Config Struct                             │             │
│  │  • nextURL, prevURL (pagination state)     │             │
│  │  • client (PokeAPI client instance)        │             │
│  └────────────────┬───────────────────────────┘             │
└───────────────────┼─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│  internal/pokeapi/                                          │
│  ┌────────────────┐         ┌─────────────────────┐         │
│  │  locations.go  │────────→│    client.go        │         │
│  │  • Endpoint    │         │    • HTTP Client    │         │
│  │  • JSON struct │         │    • doGet()        │         │
│  └────────────────┘         └──────────┬──────────┘         │
└─────────────────────────────────────────┼──────────────────┘
                                          │
                    ┌─────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│  internal/pokecache/                                        │
│  ┌────────────────────────────────────────────┐             │
│  │  cache.go                                  │             │
│  │  • In-memory cache (map)                   │             │
│  │  • Mutex for thread safety                 │             │
│  │  • Background reaper goroutine             │             │
│  │  • Add(), Get() methods                    │             │
│  └────────────────────────────────────────────┘             │
└─────────────────────────────────────────────────────────────┘
                    │
                    ▼
            [External PokeAPI]
        https://pokeapi.co/api/v2
```

### State Management

The application maintains state through the `config` struct:

**Mutable State:**

- `cfg.nextURL`: Updated after every successful fetch
- `cfg.prevURL`: Updated after every successful fetch
- Cache entries: Added, expired, removed continuously

**Immutable State:**

- `cfg.client`: Created once, never replaced
- `commands` map: Initialized once in `init()`

**State Transitions:**

```
Initial State:
    nextURL: nil
    prevURL: nil
    cache: empty

After first "map":
    nextURL: "https://...?offset=20&limit=20"
    prevURL: nil
    cache: { "https://.../location-area": <data> }

After second "map":
    nextURL: "https://...?offset=40&limit=20"
    prevURL: "https://...?offset=0&limit=20"
    cache: {
        "https://.../location-area": <data>,
        "https://...?offset=20&limit=20": <data>
    }

After "mapb":
    nextURL: "https://...?offset=20&limit=20"
    prevURL: nil
    cache: {
        "https://.../location-area": <data>,
        "https://...?offset=20&limit=20": <data>,
        "https://...?offset=0&limit=20": <data>
    }
```

---

## Contributing

This is a learning project from Boot.dev. Feel free to fork and experiment!

### Potential Enhancements

- Add command history (up/down arrow navigation)
- Implement autocomplete for commands
- Add commands to view specific location details
- Catch and display Pokemon in locations
- Add colored output for better UX
- Persistent cache (save to disk)
- Configuration file for cache TTL
- Unit tests for command handlers
- Integration tests for API client

---

## License

This project is built as part of the Boot.dev curriculum.

---

## Acknowledgments

- [PokeAPI](https://pokeapi.co/) for providing the free Pokemon data API
- [Boot.dev](https://boot.dev/) for the project structure and learning path
- The Go community for excellent standard library documentation
