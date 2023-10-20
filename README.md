# ContextGO
### Context switcher - Simplify Your Kubernetes Workflow

This is a command-line tool to switch between Kubernetes contexts. It lists all available contexts and allows the user to select one using arrow keys and enter. The selected context is then set as the current context.

## Installation

1. Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) if you haven't already.
2. Clone this repository.
3. Build the binary using `go` or download the binary from the releases.
4. Run the binary using `./contextgo`.

## Usage

1. Run the binary using `./contextgo`.
2. Use arrow keys to select a context and press enter to confirm.
3. The selected context will be set as the current context.

![Example usage](example.gif)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
