FROM golang:1.19

# Create the working directory.
WORKDIR /app

COPY . .

WORKDIR cmd/gidari

# Run the tests.
CMD ["go", "test", "-count", "3", "-v", "./..."]
