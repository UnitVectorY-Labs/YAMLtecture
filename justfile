
# List all available commands
default:
  @just --list

# Build the Go application
build:
  go build ./...

# Run the Go tests
test:
  go clean -testcache
  go test ./...

# Serve documentation locally for testing (bootstraps local Jekyll setup)
docs-serve:
  #!/usr/bin/env bash
  set -euo pipefail
  cd docs

  if [[ ! -f Gemfile ]]; then
    printf '%s\n' \
      'source "https://rubygems.org"' \
      '' \
      'gem "github-pages", group: :jekyll_plugins' \
      'gem "jekyll-remote-theme"' \
      'gem "webrick"' > Gemfile
  fi

  bundle config set --local path vendor/bundle
  bundle check >/dev/null || bundle install
  bundle exec jekyll serve --host 127.0.0.1 --port 4000 --livereload
