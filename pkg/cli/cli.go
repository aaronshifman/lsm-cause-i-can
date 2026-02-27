package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aaronshifman/lsm-cause-i-can/pkg/lsm"
)

type Op int

const (
	GET  Op = iota
	PUT  Op = iota
	QUIT Op = iota
)

// Sentinel error to quit
var ErrQuit = errors.New("quit")

type interaction struct {
	operation Op
	input     string
}

// LSM represents the GET/PUT components of the LSM interface
type LSM interface {
	Get(key string) (string, bool)
	Put(key string, val string)
}

// Cli contains the cli application
type Cli struct {
	// MemCache is the in-memory cache of the LSM tree
	// for the CLI its a map of string:string
	MemCache LSM
}

// Parse converts the raw sting into an operation and an argument
func (c *Cli) Parse(input string) (interaction, error) {
	parts := strings.SplitN(input, " ", 2)

	// QUIT technically has no argument so 1-2 args are valid
	if len(parts) == 0 || len(parts) > 2 {
		return interaction{}, errors.New("malformed input - expected [Operation] 'key' ")
	}

	var op Op
	switch parts[0] {
	case "GET":
		op = GET
	case "PUT":
		op = PUT
	case "QUIT":
		op = QUIT
	default:
		return interaction{}, fmt.Errorf("malformed input - unknown operation %s", parts[0])
	}

	var str string
	if len(parts) == 2 {
		str = parts[1]
	} else {
		str = ""
	}

	return interaction{operation: op, input: str}, nil
}

// Execute runs the parsed interaction against the in-memory cache.
func (c *Cli) Execute(inter interaction) (string, error) {
	switch inter.operation {
	case GET:
		return c.executeGet(inter.input)
	case PUT:
		return c.executePut(inter.input)
	case QUIT:
		return "", ErrQuit
	default:
		return "", fmt.Errorf("unknown operation %d", inter.operation)
	}
}

func (c *Cli) executeGet(input string) (string, error) {
	if input == "" {
		return "", errors.New("GET requires a key")
	}

	val, ok := c.MemCache.Get(input)
	if !ok {
		return "", fmt.Errorf("key not found: %s", input)
	}

	return val, nil
}

func (c *Cli) executePut(input string) (string, error) {
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("PUT requires a key and value")
	}

	c.MemCache.Put(parts[0], parts[1])

	return "OK", nil
}

// NewCli creates the CLI object
func NewCli() *Cli {
	memCache := lsm.NewMemoryCache[string, string]()
	return &Cli{
		MemCache: memCache,
	}
}
