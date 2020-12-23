# Deployment
This is the setup to deploy to Digital Ocean.

## App Platform
### Steps
1. If app not created, to create run `doctl apps create --spec app.yml`
2. If app exists, to update run `doctl apps update [INSERT APP ID] --spec app.yml`
3. Add Twilio environment variables with secrets by following [this](https://www.digitalocean.com/docs/app-platform/how-to/use-environment-variables/#define-build-time-environment-variables)
```
TWILIO_ENABLED
TWILIO_ACCOUNTSID
TWILIO_AUTHTOKEN
TWILIO_NUMBERTO
TWILIO_NUMBERFROM
```