package jsonschema

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"strings"
)

// MapLoader is an immutable in-memory schema resource loader.
type MapLoader struct {
	resources map[string][]byte
}

// NewMapLoader copies map entries and bytes into an immutable loader. Nil
// input creates an empty loader. Equivalent URI identities reject duplicates;
// keep inputs stable during construction. Construction/copying has no byte
// budget; Compiler applies document/aggregate limits after retrieval.
func NewMapLoader(resources map[string][]byte) (*MapLoader, error) {
	result := &MapLoader{resources: make(map[string][]byte, len(resources))}
	for identifier, raw := range resources {
		if identifier == "" {
			return nil, fmt.Errorf("%w: empty resource identifier", ErrResourceUnavailable)
		}
		normalized, err := normalizeResourceIdentifier(identifier)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: invalid resource identifier %q",
				ErrResourceUnavailable,
				safeResourceIdentifier(identifier),
			)
		}
		if _, duplicate := result.resources[normalized]; duplicate {
			return nil, fmt.Errorf(
				"%w: duplicate normalized resource identifier",
				ErrResourceUnavailable,
			)
		}
		result.resources[normalized] = append([]byte(nil), raw...)
	}
	return result, nil
}

// Load implements ResourceLoader and returns a detached caller-owned copy.
// Require a nonnil context; cancellation is checked before copying, not
// continuously during allocation. A zero MapLoader is an empty registry;
// a nil receiver returns ErrResourceUnavailable.
func (loader *MapLoader) Load(ctx context.Context, identifier string) ([]byte, error) {
	if loader == nil {
		return nil, fmt.Errorf("%w: nil map loader", ErrResourceUnavailable)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	normalized, err := normalizeResourceIdentifier(identifier)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid resource identifier", ErrResourceNotFound)
	}
	raw, exists := loader.resources[normalized]
	if !exists {
		return nil, fmt.Errorf(
			"%w: %q",
			ErrResourceNotFound,
			safeResourceIdentifier(identifier),
		)
	}
	return append([]byte(nil), raw...), nil
}

// FSLoader confines hierarchical resource identifiers to a caller-provided
// filesystem rooted at one absolute base URI. URI confinement is not a
// symlink/device sandbox: the caller owns filesystem capability, lifetime and
// concurrent safety. Reads are synchronous and cannot be canceled mid-read.
type FSLoader struct {
	base       *url.URL
	filesystem fs.FS
}

// NewFSLoader constructs a confined filesystem loader.
func NewFSLoader(baseIdentifier string, filesystem fs.FS) (*FSLoader, error) {
	if filesystem == nil {
		return nil, fmt.Errorf("%w: nil filesystem", ErrResourceUnavailable)
	}
	base, err := normalizeResourceURL(baseIdentifier)
	if err != nil || !base.IsAbs() || base.Host == "" || base.User != nil ||
		base.RawQuery != "" || base.Fragment != "" || !strings.HasSuffix(base.Path, "/") {
		return nil, fmt.Errorf(
			"%w: invalid filesystem base %q",
			ErrResourceUnavailable,
			safeResourceIdentifier(baseIdentifier),
		)
	}
	return &FSLoader{base: base, filesystem: filesystem}, nil
}

// Load implements ResourceLoader by rejecting URI authority/path escape.
// Require a nonnil context; cancellation is checked before the filesystem
// callback, not enforced while it blocks. fs.ReadFile can allocate all bytes
// before compiler admission. Filesystem panics are contained here; raw error
// causes remain inspectable despite fixed callback diagnostic text.
func (loader *FSLoader) Load(ctx context.Context, identifier string) ([]byte, error) {
	if loader == nil || loader.base == nil || loader.filesystem == nil {
		return nil, fmt.Errorf("%w: nil filesystem loader", ErrResourceUnavailable)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	requested, err := normalizeResourceURL(identifier)
	if err != nil || requested.Scheme != loader.base.Scheme ||
		requested.Host != loader.base.Host || requested.User != nil ||
		requested.RawQuery != "" || requested.Fragment != "" ||
		!strings.HasPrefix(requested.Path, loader.base.Path) {
		return nil, fmt.Errorf(
			"%w: %q",
			ErrResourceNotFound,
			safeResourceIdentifier(identifier),
		)
	}
	name := strings.TrimPrefix(requested.Path, loader.base.Path)
	if !validFilesystemResourcePath(name) {
		return nil, fmt.Errorf(
			"%w: %q",
			ErrResourceNotFound,
			safeResourceIdentifier(identifier),
		)
	}
	raw, err := callFilesystemRead(ctx, loader.filesystem, name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", ErrResourceNotFound, name)
		}
		return nil, err
	}
	return raw, nil
}

func validFilesystemResourcePath(name string) bool {
	return fs.ValidPath(name) && name != "."
}

// CompositeLoader tries loaders in order and falls through only when a loader
// classifies the resource as not found.
type CompositeLoader struct {
	loaders []ResourceLoader
}

// NewCompositeLoader copies the ordered interface slice, not loader state.
// Empty input creates an always-missing chain; nil/typed-nil entries reject.
// Caller owns nested loader lifecycle/concurrency and callback IO bounds.
func NewCompositeLoader(loaders ...ResourceLoader) (*CompositeLoader, error) {
	result := &CompositeLoader{loaders: make([]ResourceLoader, len(loaders))}
	for index, loader := range loaders {
		if resourceLoaderIsNil(loader) {
			return nil, fmt.Errorf("%w: loader %d is nil", ErrResourceUnavailable, index)
		}
		result.loaders[index] = loader
	}
	return result, nil
}

// Load tries borrowed loaders until success or a non-missing failure. Empty
// chains return ErrResourceNotFound. Callback wrappers contain nested panics
// and hide cause text but preserve causes; they do not force cancellation or
// copying. Components must honor the nonnil context and returned-byte lifetime.
func (loader *CompositeLoader) Load(ctx context.Context, identifier string) ([]byte, error) {
	if loader == nil {
		return nil, fmt.Errorf("%w: nil composite loader", ErrResourceUnavailable)
	}
	for _, candidate := range loader.loaders {
		raw, err := callResourceLoader(ctx, candidate, identifier)
		if err == nil {
			return raw, nil
		}
		if !errors.Is(err, ErrResourceNotFound) {
			return nil, err
		}
	}
	return nil, fmt.Errorf(
		"%w: %q",
		ErrResourceNotFound,
		safeResourceIdentifier(identifier),
	)
}

func safeResourceIdentifier(identifier string) string {
	parsed, err := url.Parse(identifier)
	if err != nil {
		return "<invalid>"
	}
	parsed.User = nil
	parsed.RawQuery = ""
	parsed.ForceQuery = false
	parsed.Fragment = ""
	parsed.RawFragment = ""
	return parsed.String()
}
