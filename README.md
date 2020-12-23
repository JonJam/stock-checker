# stock-checker

![CI](https://github.com/jonjam/stock-checker/workflows/CI/badge.svg?branch=main)
[![Quality gate](https://sonarcloud.io/api/project_badges/quality_gate?project=stock-checker)](https://sonarcloud.io/dashboard?id=stock-checker)

## Overview
Inspired by [How To Get a PlayStation 5 When It's Always Out of Stock](https://dev.to/marisayou/how-to-get-a-playstation-5-when-it-s-always-out-of-stock-5d4i).

This app checks various retailers in the United Kingdom to see if they have an Xbox Series X console in stock. 

If any have stock, it will send an SMS via Twilio detailling the results.

## Deploy to DigitalOcean
Before deploying to DigitalOcean, you will need a Twilio account which you can sign up for [here](https://www.twilio.com/try-twilio)

[![Deploy to DO](https://mp-assets1.sfo2.digitaloceanspaces.com/deploy-to-do/do-btn-blue.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/jonjam/stock-checker/tree/main)