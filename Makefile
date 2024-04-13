docker.remove:
	docker image rm anon-chat
docker.build:
	templ generate
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify
	docker build -t anon-chat:latest .
docker.push:
	docker tag anon-chat:latest $(IMAGE_URL)
	docker push $(IMAGE_URL)

air:
	air -c .air.toml
tailwind:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
templ:
	templ generate --watch --proxy="http://localhost:3000"
