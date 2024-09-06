# CLI Search Application (In Development)

This CLI app allows you to search the web directly from your terminal. It uses a configuration file where you can add or remove search engines and set defaults. The app also provides options to open URLs directly.

## Available Commands

### `config` [flag]
- **--showconfig**: Display the current configuration.
- **--showpath**: Show the path to the configuration file.
- **--set-default-browser**: Set the default browser.
- **--set-default-engine**: Set the default search engine.

### `s` [flag]
- **--url / -u**: Specify that the input is a URL rather than a search query.
- **-e**: Temporarily change the search engine for the current query without modifying the default.

# TODO
- [ ] Add support for other OS, (now only supports linux)
- [ ] Add bookmarks
- [ ] Search shell stdout
- [ ] history db

Developed in GO using Cobra cli and Viper
