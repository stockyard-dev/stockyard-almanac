# Stockyard Almanac

**Changelog and release notes.** Write entries, publish a public page, offer RSS and email subscribe. Every product needs this. Nobody wants to build it. Single binary, no external dependencies.

Part of the [Stockyard](https://stockyard.dev) suite of self-hosted developer tools.

## Quick Start

```bash
curl -sfL https://stockyard.dev/install/almanac | sh
almanac
```

Dashboard at [http://localhost:8940/ui](http://localhost:8940/ui)
Public changelog at [http://localhost:8940/changelog](http://localhost:8940/changelog)

## Usage

```bash
# Create a changelog entry
curl -X POST http://localhost:8940/api/entries \
  -H "Content-Type: application/json" \
  -d '{"title":"Dark mode","version":"v2.1.0","content":"Added dark mode.","tags":"feature","published":true}'

# Public changelog page (share this URL)
open http://localhost:8940/changelog

# RSS feed
open http://localhost:8940/changelog/rss

# Embed subscribe form on your site
# <form method="POST" action="http://localhost:8940/changelog/subscribe">
#   <input name="email" type="email">
#   <button>Subscribe</button>
# </form>
```

## Free vs Pro

| Feature | Free | Pro ($1.99/mo) |
|---------|------|----------------|
| Entries | 20 | Unlimited |
| Public changelog | ✓ | ✓ |
| RSS feed | ✓ | ✓ |
| Email subscribers | ✓ | ✓ |
| Email notifications | — | ✓ |
| Custom domain | — | ✓ |

## License

Apache 2.0 — see [LICENSE](LICENSE).
