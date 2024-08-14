# Dalamud Plugin Listing

Dalamud Plugin Listing is a small web application that combines all the plugin repositories I use when playing FFXIV, and combines them into a single list. This removes all the clutter of having 15+ different repositories saved, and instead provides a single URL that provides all the plugins I want to use.

The project is my first attempt at using Go, and is in no way perfect. I'm sure there are many things that could be done better, and I'm open to suggestions.

## Prerequisites

 - Go version 1.20 or later
 - NodeJS version 22.0.0 or later

## Installation

Start by cloning the repository:

```bash
git clone https://github.com/Senither/dalamud-plugin-listing.git
```

Next, navigate into the project directory and run the following command to install the required dependencies:

```bash
# Installs TailwindCSS so we can build the CSS
npm install

# Installs the Go dependencies
go mod download
```

Now we can build the Tailwind CSS file by using `npx`.

```bash
npx tailwindcss -i ./styles.css -o ./assets/styles.css
```

We're not finally ready to launch the service.

```bash
go run main.go
```

## Building with Docker

Building the application is made easy with Docker, the project comes with a `docker-compose.yml` file that sets up the necessary names and image tags up, so to build the image you can run the following command.

```bash
docker compose build plugin-listing
```

From here you can just start up the container with `docker compose up -d`, alternatively you can directly re-build the image with the `up` command as well.

```bash
docker compose up -d --build
```

## Testing

The project comes with some simple tests using the built in Go testing package, the tests can be ran using the following command.

```bash
go test ./... -v
```
