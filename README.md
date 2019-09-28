# Lotto; a Toy microsite for looking at UK Lottery data
lotto is a toy project for playing with vuejs/golang and utilising Plotly to display
data.

# Build
Building locally requires the following build tools:
 - Go
 - NPM
 - Quasar

Use your package manager to install Go and NPM/Yarn and then install @quasar/cli for the node build.

Or if you have Docker you can run the most recent stable build with:

```
docker run -d -p 8000:8000 ncboughton/lotto
```

And open a web browser to localhost:8000