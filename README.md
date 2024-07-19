## Introduction
Hey there! If you're as frustrated with the housing situation in the Netherlands as I am, you're in the right place. Those websites charge you a subscription just to send you emails when new houses get posted. This is scammy, and I hate scammers. They don't even provide an API! So, I had to create a scraper that checks the websites for new houses. It connects to a Telegram bot and notifies you directly. This project looks rough (I made it in 2 hours), and it’s not clean or efficient. I’ll update it… unless I find a house, then I probably won’t update it at all. <br><br>

## What This Program Does
- Scrapes housing websites for new listings in specified cities.
- Filters listings based on maximum price.
- Sends new listing notifications to a specified Telegram user via a bot.
## Prerequisites
- Docker
- Docker Compose
- A Telegram bot token and user ID (follow this guide to create a bot and get the token)
## How to Run
1. Clone the repository

2. Navigate to the docker directory and run the Docker Compose:
   ```bash
   cd docker
   docker-compose up
   ```

## Docker Compose Configuration
The docker-compose.yml file is already provided in the docker directory. Here’s what it looks like and what each part does:
   ```yaml
   services:
    find-house-almere:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    environment:
      - CITY=almere
      - CSV_FILE=/csv/listings-almere.csv
      - MAX_PRICE=1500
      - TELEGRAM_BOT_TOKEN=<TOKEN>
      - TELEGRAM_USER=<TELEGRAM_USERNAME>
    volumes:
      - ../csv:/csv
   ```


Environment Variables:

- `CITY=almere`: Specifies the city to scrape.
- `CSV_FILE=/csv/listings-almere.csv`: Path to the CSV file where listings are stored.
- `MAX_PRICE=1500`: Maximum price for filtering listings(euros).
- `TELEGRAM_BOT_TOKEN=<TOKEN>`: Token for the Telegram bot.
- `TELEGRAM_USER=<TELEGRAM_USERNAME>`: Username of the Telegram user to notify.
- volumes: Mounts a volume to persist data. Maps the ../csv directory on the host to /csv in the container.
<br>

*Note: Replace `<TOKEN>` and `<TELEGRAM_USERNAME>` with your actual Telegram bot token and username.*
 <br><br>
## Future Improvements
- Improve error handling and logging.
- Optimize the scraper for efficiency.
- Add more configurable parameters.
- Create a web interface for easier configuration and monitoring.
- Use AI to create messages to the agenct on Telegram so you don't even have to open those stupid websites
 
### **Disclaimer**
This project is a quick hack and may require frequent updates if the target websites change their structure. Use at your own risk.

Happy house hunting!

