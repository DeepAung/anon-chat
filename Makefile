air:
	air -c .air.toml
tailwind:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
templ:
	templ generate --watch --proxy="http://localhost:3000"
