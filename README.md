# EVE Ran

![EVE Ran Services Design](https://github.com/tadeasf/eve-ran-monorepo/raw/main/docs/ran-services-design.png)

**EVE Ran** is a web application for tracking and analyzing character kills in **EVE Online**.

## Demo

![Frontend Demo 1](https://github.com/tadeasf/eve-ran-monorepo/raw/main/docs/ran-frontend-1.png)

![Frontend Demo 2](https://github.com/tadeasf/eve-ran-monorepo/raw/main/docs/ran-frontend-2.png)

## Deployment

To deploy **EVE Ran** using Docker Compose:

1. **Clone the repository:**

   ```bash
   git clone https://github.com/tadeasf/eve-ran-monorepo.git
   cd eve-ran-monorepo
   ```

2. Create a `.env` file in the root directory with the necessary environment variables (refer to the `.env.example` files in the frontend and backend directories).

3. **Build and start the containers:**

   ```bash
   docker-compose up -d
   ```

4. **Access the application at `http://localhost:80` (or the configured domain).**

## Architecture

EVE Ran consists of the following components:

* **Frontend:** next.js & tailwind & shadcn/ui
* **Backend:** go & gin
* **Database:** postgreSQL
* **Web server & reverse proxy:** caddy

The services work together as follows:

- The frontend, built with Next.js and Tailwind CSS, is hosted on Vercel and provides a responsive user interface.
- Caddy acts as a reverse proxy, routing requests to the appropriate services.
- The Golang Gin backend handles API requests, interacts with the database, and communicates with external APIs ([zKillboard](https://zkillboard.com/) and [EVE Online ESI](https://esi.evetech.net/)).
- PostgreSQL stores all application data, including character information and killmail data.
- Cron jobs in the backend periodically fetch and enrich killmail data.

## Backend

The backend is built with Golang using the Gin framework. It provides various endpoints for fetching and managing data related to characters, regions, systems, constellations, and items. The main functionality includes:

- **Character management** (adding, removing, and fetching character data)
- Killmail **retrieval** and **processing**
- Region and system data management
- **Cron jobs** for periodic data **updates** and **enhancement** with data from [ESI API](https://esi.evetech.net/)

Key features of the backend include:
- RESTful API endpoints for data management and retrieval
- Integration with [zKillboard](https://zkillboard.com/) and [EVE Online ESI API](https://esi.evetech.net/)
- Cron jobs for periodic data updates
- Database interactions using **GORM**

## Frontend

The frontend is a Next.js application with Tailwind CSS for styling. It provides a user-friendly interface for:

- **Visualize PVP stats and performance for for members squads and sigs inside [Goonswarm Federation](https://goonfleet.com/)**
    - opportunity to get bragging rights
- Filtering and displaying killmail data
    - Aggregate
    - View individual kills and their ISK/point value
    - Calculates trends for both ISK/points
- Managing tracked characters

Key components of the frontend include:

1. **Dashboard**: The main page displaying character statistics, kill charts, and filtering options.
2. **Character Table**: A component for displaying detailed character information and kill statistics.
3. **Filter Controls**: Allow users to select regions, date ranges, and apply filters to the displayed data.
4. **Charts**: Various charts displaying kill statistics over time and ISK destroyed.

The frontend communicates with the backend API to fetch and display data, and uses React Query for efficient data management and caching. For a detailed look at the Dashboard component, which serves as the main interface for the application, refer to the frontend directory.

## Contributing

Contributions are welcome! Please feel free to submit a **Pull Request**

## Feature Requests

Feature requests are encouraged as well. Feel free to add a new **Issue** with description of your desired new feature.

## License

This project is licensed under the GPL-3.0 License.
