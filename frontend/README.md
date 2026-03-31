# VinylNet Frontend

This is the frontend application for VinylNet, built with Vue 3 and Vite.

## Prerequisites

- Node.js (v18 or higher recommended)
- npm or yarn

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   # or
   yarn install
   ```

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
VITE_API_URL=http://localhost:3000/api
VITE_DISCOGS_API_KEY=your_discogs_api_key
```

## Development

Start the development server:

```bash
npm run dev
# or
yarn dev
```

## Build

Build the application for production:

```bash
npm run build
# or
yarn build
```

## Deployment

Deploy the contents of the `dist` folder to your web server.
