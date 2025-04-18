# README.md for key-bored

Key-bored is a fun fast-paced game testing your speed, focus, and coordination at your desk. In the game you open up a number of browsers and try your best to click and type the letter into browsers as they appear. The more letters you type, the higher your score, but be careful, a wrong keystroke will cost you a life. The game ends when you either run out of lives or the game time ends. Each game is 60 seconds from when you press start. Can you get the highest score?

## Project information and how to run:
The project is made up of a Nextjs front end and a Golang backend. We used [Panfactum](https://panfactum.com/) to bootstrap the repo to guarantee a quick and deterministic developer experience. The only utilities needed are nix, git, and direnv. To run the project from scratch:
1. Install nix using the [Determinate Nix Installer](https://github.com/DeterminateSystems/nix-installer)
2. Run the following commands to install git and direnv through nix
    - `nix profile install nixpkgs#git` (You can skip if you already have git on your machine)
    - `nix profile install nixpkgs#direnv` (You can skip if you already have direnv on your machine)
3. Use git to clone the repo
4. cd into the repo and run `direnv allow`. This will download all the required binaries at the right versions to work in the repo. When you exit the directory they will be removed from your path, so no need to worry about any lingering side effects to your machine.
5. Open a new terminal and cd `packages/server`
6. Run `go mod download` followed by `go run .` to start the server
7. Open a new terminal and cd `packages/client`
8. Run `npm install` followed by `npm run dev`
9. Open `http://localhost:3000` and you are ready to continue on to the **How to play** section

## How to play:
1. Use the link provided on the homepage to open up a set of browsers. We recommend 4-6 for the best experience, but you're welcome to experiment and see what works well for you
2. Press start. You should see each of the non-homepage browsers start to explode with letters
3. Click on a browser to focus on it and then press a matching key on your keyboard to claim your points. But careful, wrong keystrokes or keystrokes to the wrong window will cost you a life
4. The game ends and you can appreciate your high score and how far you've come 

## Multiplayer mode:
Playing alone is fun, but what's life without a little teamwork? Grab a friend and have one of you on the mouse while the other drives the keys for a co-op experience that will test your friendship!

## Supported platforms:
- Linux
- MacOS
- Windows (via wsl)