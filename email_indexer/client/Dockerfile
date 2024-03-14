# Use the official Node.js image as the base image
FROM node:latest

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy package.json and package-lock.json files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the entire project directory into the container
COPY . .

# Build the Vue application for production
RUN npm run build

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the frontend application
CMD ["npm", "run", "dev"]
