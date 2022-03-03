# Film Generator

A simple application that generates a random film using the TMDB API. Perfect if you're stuck on what to watch tonight.

Front end coming soon....

### Requirements
- [Go](https://golang.org/doc/install)
- [TMDB Account](https://www.themoviedb.org/signup)

### Setting up your API
- [Create a free account](https://www.themoviedb.org/signup)
- Check your e-mail to verify your account.
- Visit the [API Settings](https://www.themoviedb.org/settings/api) page in your Account Settings and request an API key.
- You should now have an API key and be ready to go!

### Setting up your Environment
#### Clone:
- `git clone git@github.com:rhc07/film-generator.git`

**OR**

- `git clone https://github.com/rhc07/film-generator.git`

#### Setting up your API Key Variable
- Create a `.env` file in the root of your project.
- Create a variable called `export API_KEY=`, and set it to your personal TMDB API Key.
- **DO NOT** push to your Github unless you have put your `.env` file in a `.gitignore` file.

#### Download Dependencies:
- `go mod tidy`
- `go mod vendor`

### Build and Run

#### Build:
- `go build`

#### Run:
- `./film-generator`
