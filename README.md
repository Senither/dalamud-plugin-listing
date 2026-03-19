# Dalamud Plugin List

Dalamud Plugin List is a small web application that combines a bunch of plugin repositories for Final Fantasy 14, and then serves them as a single list that can be used to install and receive updates for all the available plugins. This removes all the clutter of having 20+ different repositories saved, and instead provides a single URL that provides all the plugins you might want to use.

## Prerequisites

 - Go version 1.25 or later
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

Now we can build the Tailwind CSS file by using the `build` or `dev` scripts.

```bash
# Builds, minimizes, and optimizes the styles.css file
npm run build
# Starts a watching to rebuild the styles anytime any of the files changes
npm run dev
```

We're now finally ready to launch the service.

```bash
go run main.go
```

> It's recommend to use [air](https://github.com/air-verse/air) during development for quickly reloading the application on file changes.

## Starting with Docker

Starting the application is made easy with Docker, the project comes with a `docker-compose.yml` file that sets up the necessary names and image tags, so to start the application you can run the following command.

```bash
docker compose up -d
```

The image will be built automatically the first time you start the application, and on subsequent runs it will boot up instantly. You can also pass the `--build` flag to rebuild the application image if you already have the image locally.

## Testing

The project comes with some simple tests using the built in Go testing package, the tests can be ran using the following command.

```bash
go test ./... -v
```
