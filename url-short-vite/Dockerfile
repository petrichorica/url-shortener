# syntax=docker/dockerfile:1

# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Dockerfile reference guide at
# https://docs.docker.com/go/dockerfile-reference/

# Want to help us make this template better? Share your feedback here: https://forms.gle/ybq9Krt8jtBL3iCk7

ARG NODE_VERSION=18.16.0

################################################################################
# Use node image for base image for all stages.
FROM node:${NODE_VERSION}-alpine as base

# Set the working directory to /app
WORKDIR /app

# Copy the package.json and package-lock.json to the working directory
COPY ./package*.json ./

# Install the dependencies
RUN npm install

# Copy the remaining application files to the working directory
COPY . .

# Build the application
RUN npm run build

# Use production node environment by default.
ENV NODE_ENV production

# Expose port 3000 for the application
EXPOSE 5173

# Start the application
CMD [ "npm", "run", "dev" ]
