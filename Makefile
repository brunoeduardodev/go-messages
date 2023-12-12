tailwind: 
	npx tailwindcss -i ./templates/input.css -o ./static/output.css
tailwind-watch: 
	npx tailwindcss -i ./templates/input.css -o ./static/output.css --watch
dev:
	gow -e go,mod,html,css run .